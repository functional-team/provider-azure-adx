/*
Copyright 2026 The provider-azure-adx Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package securityrole holds the domain logic for the SecurityRole kind:
// entity rendering, the .show/.set/.add/.drop commands, parsing of
// ".show ... principals" and the principal matching that maps what the user
// wrote (UPN, app id, group name) onto the object ids Kusto reports.
package securityrole

import (
	"fmt"
	"sort"
	"strings"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/security/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/normalize"
)

// Entity is the target of a security role.
type Entity struct {
	Kind     common.EntityKind
	Database string
	Name     string
}

// Render returns the entity clause used in commands, e.g. "table ['T']".
func (e Entity) Render() (string, error) {
	switch e.Kind {
	case common.EntityKindDatabase:
		return "database " + cmd.Ident(e.Database), nil
	case common.EntityKindTable:
		return "table " + cmd.Ident(e.Name), nil
	case common.EntityKindMaterializedView:
		return "materialized-view " + cmd.Ident(e.Name), nil
	case common.EntityKindExternalTable:
		return "external table " + cmd.Ident(e.Name), nil
	case common.EntityKindFunction:
		return "function " + cmd.Ident(e.Name), nil
	}
	return "", fmt.Errorf("unknown entity kind %q", e.Kind)
}

// Display returns the entity in Kusto notation for the status.
func (e Entity) Display() string {
	if e.Kind == common.EntityKindDatabase {
		return cmd.Ident(e.Database)
	}
	return cmd.Qualified(e.Database, e.Name)
}

// roleWord maps the role name used in commands (plural) to the word Kusto
// prints as the last token of the Role column in ".show ... principals",
// e.g. "Database Admin", "Table RawEvents Ingestor", "Database
// UnrestrictedViewer". Spike S6: the exact strings are not verified against a
// cluster; matching uses the last whitespace separated token case
// insensitively, which is robust against the entity prefix.
var roleWord = map[string]string{
	"admins":              "admin",
	"users":               "user",
	"viewers":             "viewer",
	"unrestrictedviewers": "unrestrictedviewer",
	"ingestors":           "ingestor",
	"monitors":            "monitor",
}

// Row is one principal assignment as reported by ".show ... principals".
type Row struct {
	Role        string
	Type        string
	DisplayName string
	ObjectID    string
	FQN         string
}

// ID returns the identity key of the row (object id, or FQN when missing).
func (r Row) ID() string {
	if r.ObjectID != "" {
		return strings.ToLower(r.ObjectID)
	}
	return strings.ToLower(r.FQN)
}

// ParsePrincipals parses ".show ... principals" (columns Role, PrincipalType,
// PrincipalDisplayName, PrincipalObjectId, PrincipalFQN).
func ParsePrincipals(res *kusto.Result) []Row {
	var rows []Row
	if res == nil {
		return rows
	}
	for _, r := range res.Rows() {
		row := Row{Role: r.String("Role"), Type: r.String("PrincipalType"), DisplayName: r.String("PrincipalDisplayName"), ObjectID: r.String("PrincipalObjectId"), FQN: r.String("PrincipalFQN")}
		if row.ObjectID == "" && row.FQN == "" {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

// FilterRole keeps the rows whose Role column denotes role.
func FilterRole(rows []Row, role string) []Row {
	want := roleWord[strings.ToLower(role)]
	if want == "" {
		want = strings.TrimSuffix(strings.ToLower(role), "s")
	}
	var out []Row
	for _, r := range rows {
		fields := strings.Fields(r.Role)
		if len(fields) == 0 {
			continue
		}
		if strings.EqualFold(fields[len(fields)-1], want) {
			out = append(out, r)
		}
	}
	return out
}

// Principal is a parsed principal string "type=id[;tenant]".
type Principal struct {
	Type   string
	ID     string
	Tenant string
}

// ParsePrincipal parses "aaduser=alice@contoso.com", "aadapp=<id>;<tenant>" ...
func ParsePrincipal(s string) (Principal, error) {
	s = strings.TrimSpace(s)
	typ, rest, ok := strings.Cut(s, "=")
	if !ok || strings.TrimSpace(typ) == "" || strings.TrimSpace(rest) == "" {
		return Principal{}, fmt.Errorf("principal %q must have the form type=id[;tenant]", s)
	}
	id, tenant, _ := strings.Cut(rest, ";")
	return Principal{Type: strings.ToLower(strings.TrimSpace(typ)), ID: strings.TrimSpace(id), Tenant: strings.TrimSpace(tenant)}, nil
}

// typeMatches reports whether a row's PrincipalType is compatible with the
// principal type prefix (aaduser/aadapp/aadgroup/msauser). Unknown types match.
func typeMatches(p Principal, rowType string) bool {
	rt := strings.ToLower(rowType)
	if rt == "" {
		return true
	}
	switch {
	case strings.HasSuffix(p.Type, "app"):
		return strings.Contains(rt, "app")
	case strings.HasSuffix(p.Type, "group"):
		return strings.Contains(rt, "group")
	case strings.HasSuffix(p.Type, "user"):
		return strings.Contains(rt, "user")
	}
	return true
}

// Match finds the row that corresponds to the spec principal. Rules, in
// order: FQN equality (with and without the tenant part), object id equality,
// display name equality or containment (Kusto renders "Alice (upn: alice@x)"
// and "MyApp (app id: <appId>)"). Type compatibility is required whenever
// the row carries a type.
func Match(spec string, rows []Row) (Row, bool) { //nolint:gocyclo // ordered matching rules are easier to audit inline.
	p, err := ParsePrincipal(spec)
	if err != nil {
		return Row{}, false
	}
	specLower := strings.ToLower(strings.TrimSpace(spec))
	specNoTenant := p.Type + "=" + strings.ToLower(p.ID)
	for _, r := range rows {
		fqn := strings.ToLower(strings.TrimSpace(r.FQN))
		fqnNoTenant, _, _ := strings.Cut(fqn, ";")
		if fqn == specLower || (fqnNoTenant == specNoTenant && typeMatches(p, r.Type)) {
			return r, true
		}
	}
	for _, r := range rows {
		if strings.EqualFold(r.ObjectID, p.ID) && typeMatches(p, r.Type) {
			return r, true
		}
	}
	for _, r := range rows {
		if !typeMatches(p, r.Type) {
			continue
		}
		dn := strings.ToLower(r.DisplayName)
		id := strings.ToLower(p.ID)
		if dn == id || (len(id) >= 5 && strings.Contains(dn, id)) {
			return r, true
		}
	}
	return Row{}, false
}

// Resolved maps a spec entry to the principal Kusto reports.
type Resolved struct {
	Spec        string
	ObjectID    string
	FQN         string
	DisplayName string
	Type        string
}

func fromRow(spec string, r Row) Resolved {
	return Resolved{Spec: spec, ObjectID: r.ObjectID, FQN: r.FQN, DisplayName: r.DisplayName, Type: r.Type}
}

// Input to Evaluate.
type Input struct {
	// Additive selects the mode.
	Additive bool
	// Spec principals.
	Spec []string
	// Resolved is the mapping persisted in the status.
	Resolved []Resolved
	// Rows of the role as reported by the cluster.
	Rows []Row
	// BeforeIDs are the object ids of the role before the last write
	// (annotation); used to attribute a single newly appeared row.
	BeforeIDs []string
	// SpecUnchanged is true when the spec principals equal what was last
	// written (hash annotation). It lets unresolvable entries settle.
	SpecUnchanged bool
	// Owned is true when this resource created the role assignment
	// (annotation); an empty authoritative role then still exists.
	Owned bool
}

// Evaluation is the result of comparing spec and cluster.
type Evaluation struct {
	Resolved   []Resolved
	Present    []Resolved
	Absent     []Resolved
	Stale      []Resolved
	Unresolved []string
	Extra      []Row
	Exists     bool
	UpToDate   bool
	Diff       string
	// ToAdd / ToDrop are the additive write plan (principal literals / FQNs).
	ToAdd  []string
	ToDrop []string
}

// Evaluate matches spec principals to cluster rows and decides existence and
// drift. It is a pure function; the controller persists Resolved in the
// status and writes ToAdd/ToDrop (Additive) or the whole spec (Authoritative).
func Evaluate(in Input) Evaluation { //nolint:gocyclo // A single decision table is easier to audit than helpers.
	ev := Evaluation{}
	rowsByID := make(map[string]Row, len(in.Rows))
	for _, r := range in.Rows {
		rowsByID[r.ID()] = r
	}
	claimed := map[string]bool{}
	inSpec := map[string]bool{}
	prev := make(map[string]Resolved, len(in.Resolved))
	for _, r := range in.Resolved {
		prev[strings.ToLower(strings.TrimSpace(r.Spec))] = r
	}
	unclaimed := func() []Row {
		out := make([]Row, 0, len(in.Rows))
		for _, r := range in.Rows {
			if !claimed[r.ID()] {
				out = append(out, r)
			}
		}
		return out
	}
	var spec []string
	for _, s := range in.Spec {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		spec = append(spec, s)
		key := strings.ToLower(s)
		inSpec[key] = true
		if r, ok := prev[key]; ok && r.ObjectID != "" {
			r.Spec = s
			ev.Resolved = append(ev.Resolved, r)
			if _, present := rowsByID[strings.ToLower(r.ObjectID)]; present {
				claimed[strings.ToLower(r.ObjectID)] = true
				ev.Present = append(ev.Present, r)
			} else {
				ev.Absent = append(ev.Absent, r)
			}
			continue
		}
		if row, ok := Match(s, unclaimed()); ok {
			r := fromRow(s, row)
			claimed[row.ID()] = true
			ev.Resolved = append(ev.Resolved, r)
			ev.Present = append(ev.Present, r)
			continue
		}
		ev.Unresolved = append(ev.Unresolved, s)
	}
	for _, r := range in.Resolved {
		if inSpec[strings.ToLower(strings.TrimSpace(r.Spec))] || r.ObjectID == "" {
			continue
		}
		if _, present := rowsByID[strings.ToLower(r.ObjectID)]; present && !claimed[strings.ToLower(r.ObjectID)] {
			claimed[strings.ToLower(r.ObjectID)] = true
			ev.Stale = append(ev.Stale, r)
		}
	}
	// A single unresolved entry is attributed to the single row that appeared
	// since the last write.
	if len(ev.Unresolved) == 1 {
		before := make(map[string]bool, len(in.BeforeIDs))
		for _, id := range in.BeforeIDs {
			before[strings.ToLower(strings.TrimSpace(id))] = true
		}
		var appeared []Row
		for _, r := range unclaimed() {
			if !before[r.ID()] {
				appeared = append(appeared, r)
			}
		}
		if len(appeared) == 1 {
			r := fromRow(ev.Unresolved[0], appeared[0])
			claimed[appeared[0].ID()] = true
			ev.Resolved = append(ev.Resolved, r)
			ev.Present = append(ev.Present, r)
			ev.Unresolved = nil
		}
	}
	ev.Extra = unclaimed()

	switch {
	case in.Additive:
		ev.Exists = len(spec) == 0 || len(ev.Present) > 0 || (in.SpecUnchanged && len(ev.Unresolved) > 0 && len(in.Rows) > 0)
		ev.UpToDate = len(ev.Absent) == 0 && len(ev.Stale) == 0 && (len(ev.Unresolved) == 0 || in.SpecUnchanged)
		for _, r := range ev.Absent {
			ev.ToAdd = append(ev.ToAdd, r.Spec)
		}
		ev.ToAdd = append(ev.ToAdd, ev.Unresolved...)
		for _, r := range ev.Stale {
			ev.ToDrop = append(ev.ToDrop, r.FQN)
		}
		switch {
		case len(ev.Absent) > 0:
			ev.Diff = "principals missing from the role: " + specs(ev.Absent)
		case len(ev.Stale) > 0:
			ev.Diff = "principals removed from the spec but still assigned: " + specs(ev.Stale)
		case len(ev.Unresolved) > 0 && !in.SpecUnchanged:
			ev.Diff = "principals not yet assigned: " + strings.Join(ev.Unresolved, ", ")
		}
	default:
		ev.Exists = len(in.Rows) > 0 || (len(spec) == 0 && in.Owned)
		switch {
		case len(spec) == 0:
			ev.UpToDate = len(in.Rows) == 0
			if !ev.UpToDate {
				ev.Diff = fmt.Sprintf("role has %d principals, spec wants none", len(in.Rows))
			}
		case len(ev.Absent) > 0:
			ev.Diff = "principals missing from the role: " + specs(ev.Absent)
		case len(ev.Unresolved) == 0:
			ev.UpToDate = len(ev.Extra) == 0 && len(ev.Stale) == 0
			if !ev.UpToDate {
				ev.Diff = fmt.Sprintf("%d principals in the role are not in the spec", len(ev.Extra)+len(ev.Stale))
			}
		default:
			// Some spec entries cannot be mapped to rows. Once the spec has been
			// written, accept the state when the unmatched rows account exactly
			// for the unmatched spec entries; otherwise write again.
			ev.UpToDate = in.SpecUnchanged && len(ev.Extra)+len(ev.Stale) == len(ev.Unresolved)
			if !ev.UpToDate {
				ev.Diff = "principals not yet assigned or not resolvable: " + strings.Join(ev.Unresolved, ", ")
			}
		}
	}
	return ev
}

func specs(rs []Resolved) string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.Spec)
	}
	return strings.Join(out, ", ")
}

// SpecHash hashes the set of spec principals (order and case insensitive).
func SpecHash(spec []string) string {
	norm := make([]string, 0, len(spec))
	for _, s := range spec {
		if s = strings.ToLower(strings.TrimSpace(s)); s != "" {
			norm = append(norm, s)
		}
	}
	sort.Strings(norm)
	return normalize.Hash(norm...)
}

// IDs returns the identity keys of rows.
func IDs(rows []Row) []string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.ID())
	}
	sort.Strings(out)
	return out
}

// ShowPrincipals is ".show <entity> principals".
func ShowPrincipals(e Entity) (cmd.Command, error) {
	ent, err := e.Render()
	if err != nil {
		return cmd.Command{}, err
	}
	return cmd.New(".show ", ent, " principals"), nil
}

func literals(principals []string) string {
	items := make([]string, 0, len(principals))
	for _, p := range principals {
		if p = strings.TrimSpace(p); p != "" {
			items = append(items, cmd.Str(p))
		}
	}
	return cmd.List(items)
}

// BuildSet is ".set <entity> <role> ("p1", "p2") skip-results "description""
// or ".set <entity> <role> none skip-results" for an empty list.
func BuildSet(e Entity, role string, principals []string, description *string) (cmd.Command, error) {
	ent, err := e.Render()
	if err != nil {
		return cmd.Command{}, err
	}
	if len(principals) == 0 {
		return cmd.New(".set ", ent, " ", role, " none skip-results"), nil
	}
	parts := []string{".set ", ent, " ", role, " ", literals(principals), " skip-results"}
	if description != nil && strings.TrimSpace(*description) != "" {
		parts = append(parts, " ", cmd.Str(*description))
	}
	return cmd.New(parts...), nil
}

// BuildAdd is ".add <entity> <role> ("p1") skip-results "description"".
func BuildAdd(e Entity, role string, principals []string, description *string) (cmd.Command, error) {
	ent, err := e.Render()
	if err != nil {
		return cmd.Command{}, err
	}
	parts := []string{".add ", ent, " ", role, " ", literals(principals), " skip-results"}
	if description != nil && strings.TrimSpace(*description) != "" {
		parts = append(parts, " ", cmd.Str(*description))
	}
	return cmd.New(parts...), nil
}

// BuildDrop is ".drop <entity> <role> ("p1") skip-results".
func BuildDrop(e Entity, role string, principals []string) (cmd.Command, error) {
	ent, err := e.Render()
	if err != nil {
		return cmd.Command{}, err
	}
	return cmd.New(".drop ", ent, " ", role, " ", literals(principals), " skip-results"), nil
}

// Observation converts the evaluation into the status representation. Stale
// mappings (removed from the spec but still assigned) are kept until the
// principal has actually been dropped, otherwise Update could not find them.
func Observation(ev Evaluation, rows []Row) v1alpha1.SecurityRoleObservation {
	obs := v1alpha1.SecurityRoleObservation{Unresolved: ev.Unresolved}
	for _, r := range append(append([]Resolved(nil), ev.Resolved...), ev.Stale...) {
		obs.ResolvedPrincipals = append(obs.ResolvedPrincipals, v1alpha1.ResolvedPrincipal{Spec: r.Spec, ObjectID: r.ObjectID, FQN: r.FQN, DisplayName: r.DisplayName, Type: r.Type})
	}
	for _, r := range rows {
		if r.FQN != "" {
			obs.Principals = append(obs.Principals, r.FQN)
		} else {
			obs.Principals = append(obs.Principals, r.ObjectID)
		}
	}
	sort.Strings(obs.Principals)
	return obs
}

// FromObservation restores the persisted mapping.
func FromObservation(obs v1alpha1.SecurityRoleObservation) []Resolved {
	out := make([]Resolved, 0, len(obs.ResolvedPrincipals))
	for _, r := range obs.ResolvedPrincipals {
		out = append(out, Resolved{Spec: r.Spec, ObjectID: r.ObjectID, FQN: r.FQN, DisplayName: r.DisplayName, Type: r.Type})
	}
	return out
}
