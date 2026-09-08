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

// Package policy holds the domain logic shared by all policy kinds: entity
// rendering, the command shapes (JSON form by default, special forms for the
// few policies without one), parsing of .show output and the "only compare
// what the spec sets" comparison.
package policy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/normalize"
)

// Entity is the target of a policy command.
type Entity struct {
	Kind     common.EntityKind
	Database string
	Name     string
	// Column narrows a table entity to one column (encoding policy).
	Column string
}

// Render returns the entity clause used in commands, e.g. "table ['T']".
func (e Entity) Render() (string, error) {
	if e.Column != "" {
		if e.Kind != common.EntityKindTable {
			return "", fmt.Errorf("column policies are only supported on tables")
		}
		return "column " + cmd.Qualified(e.Name, e.Column), nil
	}
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
		return "", fmt.Errorf("functions do not have policies")
	}
	return "", fmt.Errorf("unknown entity kind %q", e.Kind)
}

// Display returns the entity in the notation Kusto uses in .show output.
func (e Entity) Display() string {
	switch {
	case e.Column != "":
		return cmd.Qualified(e.Database, e.Name) + "." + cmd.Ident(e.Column)
	case e.Kind == common.EntityKindDatabase:
		return cmd.Ident(e.Database)
	}
	return cmd.Qualified(e.Database, e.Name)
}

// Def describes one Kusto policy type.
type Def struct {
	// Name is the policy name used in commands, e.g. "retention".
	Name string
	// Alter renders the alter command. nil means the JSON form
	// ".alter <entity> policy <name> @'<json>'".
	Alter func(e Entity, desired any) (cmd.Command, error)
	// KQLFields are JSON keys whose string values are KQL text and therefore
	// compared with the two stage normalization instead of exact equality.
	KQLFields []string
	// SetFields are JSON keys holding comma separated sets ("A, B") that are
	// compared order-insensitively.
	SetFields []string
	// NoBatch disables the ".show table * policy <name>" batch observe.
	NoBatch bool
	// HashOnly compares desired and observed only through the stage 2 hashes
	// (used for policies whose .show JSON shape is not verified).
	HashOnly bool
}

// ShowCmd is ".show <entity> policy <name>".
func (d Def) ShowCmd(e Entity) (cmd.Command, error) {
	ent, err := e.Render()
	if err != nil {
		return cmd.Command{}, err
	}
	return cmd.New(".show ", ent, " policy ", d.Name), nil
}

// ShowAllTablesCmd is ".show table * policy <name>" (batch observe).
func (d Def) ShowAllTablesCmd() cmd.Command {
	return cmd.New(".show table * policy ", d.Name)
}

// AlterCmd renders the alter command for desired.
func (d Def) AlterCmd(e Entity, desired any) (cmd.Command, error) {
	if d.Alter != nil {
		return d.Alter(e, desired)
	}
	ent, err := e.Render()
	if err != nil {
		return cmd.Command{}, err
	}
	js, err := cmd.JSON(desired)
	if err != nil {
		return cmd.Command{}, err
	}
	return cmd.New(".alter ", ent, " policy ", d.Name, " ", js), nil
}

// DeleteCmd is ".delete <entity> policy <name>".
func (d Def) DeleteCmd(e Entity) (cmd.Command, error) {
	ent, err := e.Render()
	if err != nil {
		return cmd.Command{}, err
	}
	return cmd.New(".delete ", ent, " policy ", d.Name), nil
}

// CompareWith compares desired against the observed policy JSON.
func (d Def) CompareWith(desired any, observed json.RawMessage) (Result, error) {
	if d.HashOnly {
		db, err := json.Marshal(desired)
		if err != nil {
			return Result{}, err
		}
		// Structurally vacuous; equality is decided by the hash annotations alone.
		return Result{Equal: true, Diff: "compared via hash only", DesiredTexts: []string{string(db)}, ObservedTexts: []string{Compact(observed)}}, nil
	}
	return Compare(desired, observed, d.KQLFields, d.SetFields)
}

// ParseShow returns the Policy JSON of the first row of a ".show ... policy"
// result. ok is false when there is no row or the policy is null, i.e. the
// entity has no policy of its own (it inherits).
func ParseShow(res *kusto.Result) (json.RawMessage, bool) {
	rows := res.Rows()
	if len(rows) == 0 {
		return nil, false
	}
	raw := rows[0].Dynamic("Policy")
	return raw, !IsNull(raw)
}

// ParseShowAll parses ".show table * policy <name>" into table name -> policy
// JSON. Tables without their own policy map to nil.
func ParseShowAll(res *kusto.Result) map[string]json.RawMessage {
	out := map[string]json.RawMessage{}
	for _, r := range res.Rows() {
		name := normalize.UnqualifiedTable(r.String("EntityName"))
		if name == "" {
			continue
		}
		raw := r.Dynamic("Policy")
		if IsNull(raw) {
			out[name] = nil
			continue
		}
		out[name] = raw
	}
	return out
}

// IsNull reports whether raw is empty or the JSON null literal.
func IsNull(raw json.RawMessage) bool {
	s := strings.TrimSpace(string(raw))
	return s == "" || s == "null"
}

// Compact returns raw without insignificant whitespace (for status).
func Compact(raw json.RawMessage) string {
	if IsNull(raw) {
		return ""
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return string(raw)
	}
	return string(b)
}
