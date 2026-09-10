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

// Package table holds the pure domain logic for the Table kind: rendering the
// commands, parsing what the cluster reports and diffing desired against
// observed state with the schema guardrails (Merge vs Replace, no column
// type changes).
package table

import (
	"fmt"
	"sort"
	"strings"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/normalize"
)

// Column is a normalized column (canonical type).
type Column struct {
	Name      string
	Type      string
	Docstring *string
}

// Observed is the state of a table as reported by the cluster.
type Observed struct {
	Name      string
	Columns   []Column
	Folder    string
	Docstring string
}

// Desired is the normalized desired state.
type Desired struct {
	Name      string
	Columns   []Column
	Folder    *string
	Docstring *string
	Replace   bool
}

// FromSchema converts a schema JSON table.
func FromSchema(t schema.Table) Observed {
	o := Observed{Name: t.Name, Folder: t.Folder, Docstring: t.DocString}
	for _, c := range t.OrderedColumns {
		typ := c.CslType
		if typ == "" {
			typ = c.Type
		}
		doc := c.DocString
		o.Columns = append(o.Columns, Column{Name: c.Name, Type: normalize.ColumnType(typ), Docstring: &doc})
	}
	return o
}

// ParseCslSchema parses the result of ".show table T cslschema" (columns
// TableName, Schema, DatabaseName, Folder, DocString).
func ParseCslSchema(res *kusto.Result) (Observed, bool) {
	rows := res.Rows()
	if len(rows) == 0 {
		return Observed{}, false
	}
	r := rows[0]
	o := Observed{Name: r.String("TableName"), Folder: r.String("Folder"), Docstring: r.String("DocString")}
	for _, part := range splitSchema(r.String("Schema")) {
		i := strings.LastIndex(part, ":")
		if i < 0 {
			continue
		}
		o.Columns = append(o.Columns, Column{Name: normalize.UnqualifiedTable(part[:i]), Type: normalize.ColumnType(part[i+1:])})
	}
	return o, true
}

// splitSchema splits "A:string,['B c']:long" at commas outside brackets.
func splitSchema(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, c := range s {
		switch c {
		case '[', '(':
			depth++
		case ']', ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	if rest := strings.TrimSpace(s[start:]); rest != "" {
		parts = append(parts, rest)
	}
	return parts
}

// FromParams normalizes the spec.
func FromParams(name string, p v1alpha1.TableParameters) Desired {
	d := Desired{Name: name, Folder: p.Folder, Docstring: p.Docstring, Replace: p.SchemaUpdateMode == v1alpha1.SchemaUpdateModeReplace}
	for _, c := range p.Columns {
		d.Columns = append(d.Columns, Column{Name: c.Name, Type: normalize.ColumnType(string(c.Type)), Docstring: c.Docstring})
	}
	return d
}

// StepKind names an update step.
type StepKind string

// Update steps.
const (
	StepAlterMergeSchema    StepKind = "AlterMergeSchema"
	StepAlterSchema         StepKind = "AlterSchema"
	StepSetFolder           StepKind = "SetFolder"
	StepSetDocstring        StepKind = "SetDocstring"
	StepSetColumnDocstrings StepKind = "SetColumnDocstrings"
)

// Step is one command of an update plan.
type Step struct {
	Kind    StepKind
	Command cmd.Command
}

// Plan is the ordered list of commands needed to reach the desired state plus
// informational drift (columns that exist in the cluster but not in the spec).
type Plan struct {
	Steps        []Step
	DriftColumns []string
}

// Empty reports whether nothing needs to be sent.
func (p Plan) Empty() bool { return len(p.Steps) == 0 }

// String summarizes the plan for the Diff field / logs.
func (p Plan) String() string {
	if p.Empty() {
		return ""
	}
	kinds := make([]string, 0, len(p.Steps))
	for _, s := range p.Steps {
		kinds = append(kinds, string(s.Kind))
	}
	return strings.Join(kinds, ",")
}

// Diff compares desired and observed state. A column type change returns a
// kerrors.Blocked error (UnsupportedColumnTypeChange) and no plan.
func Diff(d Desired, o Observed) (Plan, error) {
	obs := make(map[string]Column, len(o.Columns))
	for _, c := range o.Columns {
		obs[c.Name] = c
	}
	missing, drift, err := compareColumns(d, o, obs)
	if err != nil {
		return Plan{}, err
	}
	p := Plan{DriftColumns: drift}
	switch {
	case d.Replace && (len(missing) > 0 || len(drift) > 0 || !sameOrder(d.Columns, o.Columns)):
		p.Steps = append(p.Steps, Step{Kind: StepAlterSchema, Command: cmd.New(".alter table ", cmd.Ident(d.Name), " ", schemaOf(d.Columns))})
		p.DriftColumns = nil // Replace removes them.
	case len(missing) > 0:
		p.Steps = append(p.Steps, Step{Kind: StepAlterMergeSchema, Command: cmd.New(".alter-merge table ", cmd.Ident(d.Name), " ", schemaOf(d.Columns))})
	}
	p.Steps = append(p.Steps, metadataSteps(d, o, obs)...)
	return p, nil
}

// compareColumns returns the desired columns missing in the cluster and the
// cluster columns missing in the spec, or a Blocked error on a type change.
func compareColumns(d Desired, o Observed, obs map[string]Column) (missing, drift []string, err error) {
	des := make(map[string]Column, len(d.Columns))
	for _, c := range d.Columns {
		des[c.Name] = c
		oc, ok := obs[c.Name]
		if !ok {
			missing = append(missing, c.Name)
			continue
		}
		if oc.Type != c.Type {
			return nil, nil, kerrors.NewBlocked(kerrors.ReasonUnsupportedColumnTypeChange,
				"column %s type change %s->%s is not supported by Kusto without data loss; keep the type or recreate the table manually", c.Name, oc.Type, c.Type)
		}
	}
	for _, c := range o.Columns {
		if _, ok := des[c.Name]; !ok {
			drift = append(drift, c.Name)
		}
	}
	sort.Strings(drift)
	return missing, drift, nil
}

// metadataSteps renders folder, docstring and column docstring updates.
func metadataSteps(d Desired, o Observed, obs map[string]Column) []Step {
	var steps []Step
	if d.Folder != nil && strings.TrimSpace(*d.Folder) != strings.TrimSpace(o.Folder) {
		steps = append(steps, Step{Kind: StepSetFolder, Command: cmd.New(".alter table ", cmd.Ident(d.Name), " folder ", cmd.Str(*d.Folder))})
	}
	if d.Docstring != nil && strings.TrimSpace(*d.Docstring) != strings.TrimSpace(o.Docstring) {
		steps = append(steps, Step{Kind: StepSetDocstring, Command: cmd.New(".alter table ", cmd.Ident(d.Name), " docstring ", cmd.Str(*d.Docstring))})
	}
	if c := columnDocstrings(d, obs); c != "" {
		steps = append(steps, Step{Kind: StepSetColumnDocstrings, Command: cmd.New(alterMergeColumnDocstrings, cmd.Ident(d.Name), " column-docstrings ", c)})
	}
	return steps
}

func sameOrder(d, o []Column) bool {
	if len(d) != len(o) {
		return false
	}
	for i := range d {
		if d[i].Name != o[i].Name {
			return false
		}
	}
	return true
}

// alterMergeColumnDocstrings has to be the merging verb. ".alter table T
// column-docstrings (...)" replaces the whole set: "columns not explicitly set
// will have this property removed". Since the list below holds only the
// columns whose docstring changed, ".alter" would strip every other column's
// docstring, the next Observe would see that as drift, write the previous set
// back, and drop the new one again -- the table oscillates between two states
// forever while reporting Synced=True, because each write succeeds. Reported
// as issue #1 against v0.1.0-rc.3, with three docstringed columns taking turns
// on every reconcile. ".alter-merge" leaves columns it does not name alone.
const alterMergeColumnDocstrings = ".alter-merge table "

// columnDocstrings renders the column-docstrings list for desired columns
// whose docstring is set and differs from the observed one. Returns "" if
// nothing differs.
//
// Only changed columns are listed, which is what makes the merging verb
// mandatory; see alterMergeColumnDocstrings. A docstring the spec no longer
// sets is left alone rather than cleared -- same "only add" stance as
// schemaUpdateMode: Merge.
func columnDocstrings(d Desired, obs map[string]Column) string {
	var items []string
	for _, c := range d.Columns {
		if c.Docstring == nil {
			continue
		}
		cur := ""
		if oc, ok := obs[c.Name]; ok && oc.Docstring != nil {
			cur = *oc.Docstring
		}
		if strings.TrimSpace(*c.Docstring) != strings.TrimSpace(cur) {
			items = append(items, cmd.Ident(c.Name)+":"+cmd.Str(*c.Docstring))
		}
	}
	if len(items) == 0 {
		return ""
	}
	return cmd.List(items)
}

func schemaOf(cols []Column) string {
	defs := make([]cmd.ColumnDef, 0, len(cols))
	for _, c := range cols {
		defs = append(defs, cmd.ColumnDef{Name: c.Name, Type: c.Type})
	}
	return cmd.Schema(defs)
}

// BuildCreate renders ".create table" plus, if needed, the column docstrings.
// .create table is idempotent when the table already exists (S10).
func BuildCreate(d Desired) []cmd.Command {
	props := map[string]string{}
	if d.Folder != nil {
		props["folder"] = cmd.Str(*d.Folder)
	}
	if d.Docstring != nil {
		props["docstring"] = cmd.Str(*d.Docstring)
	}
	cmds := []cmd.Command{cmd.New(".create table ", cmd.Ident(d.Name), " ", schemaOf(d.Columns), cmd.With(props))}
	if c := columnDocstrings(d, map[string]Column{}); c != "" {
		cmds = append(cmds, cmd.New(alterMergeColumnDocstrings, cmd.Ident(d.Name), " column-docstrings ", c))
	}
	return cmds
}

// BuildDelete renders ".drop table T ifexists".
func BuildDelete(name string) cmd.Command {
	return cmd.New(".drop table ", cmd.Ident(name), " ifexists")
}

// ShowCslSchema renders ".show table T cslschema" (single observe).
func ShowCslSchema(name string) cmd.Command {
	return cmd.New(".show table ", cmd.Ident(name), " cslschema")
}

// Observation converts observed state into the status representation.
func Observation(o Observed, drift []string) v1alpha1.TableObservation {
	obs := v1alpha1.TableObservation{Folder: o.Folder, Docstring: o.Docstring, DriftColumns: drift}
	for _, c := range o.Columns {
		oc := common.ObservedColumn{Name: c.Name, Type: c.Type}
		if c.Docstring != nil {
			oc.Docstring = *c.Docstring
		}
		obs.Columns = append(obs.Columns, oc)
	}
	return obs
}

// DescribeDrift renders the informational message for columns missing in the spec.
func DescribeDrift(drift []string) string {
	if len(drift) == 0 {
		return ""
	}
	return fmt.Sprintf("%s: columns exist in the cluster but not in the spec (%s); schemaUpdateMode Merge keeps them", kerrors.ReasonColumnsMissingInSpec, strings.Join(drift, ", "))
}
