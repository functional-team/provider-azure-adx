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

package table

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema/schematest"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
)

func ptr(s string) *string { return &s }

func params(cols ...common.Column) v1alpha1.TableParameters {
	return v1alpha1.TableParameters{Database: "DB", Columns: cols, SchemaUpdateMode: v1alpha1.SchemaUpdateModeMerge}
}

func observed(cols ...Column) Observed {
	return Observed{Name: "T", Columns: cols}
}

func kinds(p Plan) []StepKind {
	out := make([]StepKind, 0, len(p.Steps))
	for _, s := range p.Steps {
		out = append(out, s.Kind)
	}
	return out
}

func TestDiff(t *testing.T) {
	ts := common.Column{Name: "Ts", Type: "datetime"}
	payload := common.Column{Name: "Payload", Type: "dynamic"}
	cases := map[string]struct {
		desired  Desired
		observed Observed
		steps    []StepKind
		drift    []string
		blocked  bool
		cmds     []string
	}{
		"upToDate": {
			desired:  FromParams("T", params(ts, payload)),
			observed: observed(Column{Name: "Ts", Type: "datetime"}, Column{Name: "Payload", Type: "dynamic"}),
		},
		"aliasTypesEqual": {
			desired:  FromParams("T", params(common.Column{Name: "A", Type: "boolean"}, common.Column{Name: "B", Type: "double"})),
			observed: observed(Column{Name: "A", Type: "bool"}, Column{Name: "B", Type: "real"}),
		},
		"mergeAddsColumn": {
			desired:  FromParams("T", params(ts, payload)),
			observed: observed(Column{Name: "Ts", Type: "datetime"}),
			steps:    []StepKind{StepAlterMergeSchema},
			cmds:     []string{".alter-merge table ['T'] (['Ts']:datetime, ['Payload']:dynamic)"},
		},
		"mergeIgnoresExtraColumnsButReportsDrift": {
			desired:  FromParams("T", params(ts)),
			observed: observed(Column{Name: "Ts", Type: "datetime"}, Column{Name: "Extra", Type: "string"}, Column{Name: "Another", Type: "int"}),
			drift:    []string{"Another", "Extra"},
		},
		"mergeIgnoresOrder": {
			desired:  FromParams("T", params(payload, ts)),
			observed: observed(Column{Name: "Ts", Type: "datetime"}, Column{Name: "Payload", Type: "dynamic"}),
		},
		"replaceRemovesColumn": {
			desired: func() Desired {
				p := params(ts)
				p.SchemaUpdateMode = v1alpha1.SchemaUpdateModeReplace
				return FromParams("T", p)
			}(),
			observed: observed(Column{Name: "Ts", Type: "datetime"}, Column{Name: "Extra", Type: "string"}),
			steps:    []StepKind{StepAlterSchema},
			cmds:     []string{".alter table ['T'] (['Ts']:datetime)"},
		},
		"replaceReorders": {
			desired: func() Desired {
				p := params(payload, ts)
				p.SchemaUpdateMode = v1alpha1.SchemaUpdateModeReplace
				return FromParams("T", p)
			}(),
			observed: observed(Column{Name: "Ts", Type: "datetime"}, Column{Name: "Payload", Type: "dynamic"}),
			steps:    []StepKind{StepAlterSchema},
		},
		"typeChangeBlocked": {
			desired:  FromParams("T", params(common.Column{Name: "Ts", Type: "long"})),
			observed: observed(Column{Name: "Ts", Type: "datetime"}),
			blocked:  true,
		},
		"typeChangeBlockedEvenInReplace": {
			desired: func() Desired {
				p := params(common.Column{Name: "Ts", Type: "string"})
				p.SchemaUpdateMode = v1alpha1.SchemaUpdateModeReplace
				return FromParams("T", p)
			}(),
			observed: observed(Column{Name: "Ts", Type: "datetime"}),
			blocked:  true,
		},
		"metadata": {
			desired: func() Desired {
				p := params(ts)
				p.Folder, p.Docstring = ptr("Raw"), ptr("Landing")
				return FromParams("T", p)
			}(),
			observed: Observed{Name: "T", Columns: []Column{{Name: "Ts", Type: "datetime"}}, Folder: "Old", Docstring: ""},
			steps:    []StepKind{StepSetFolder, StepSetDocstring},
			cmds:     []string{".alter table ['T'] folder \"Raw\"", ".alter table ['T'] docstring \"Landing\""},
		},
		"metadataUnsetIgnored": {
			desired:  FromParams("T", params(ts)),
			observed: Observed{Name: "T", Columns: []Column{{Name: "Ts", Type: "datetime"}}, Folder: "Whatever", Docstring: "server default"},
		},
		"columnDocstrings": {
			desired:  FromParams("T", params(common.Column{Name: "Ts", Type: "datetime", Docstring: ptr("when")}, common.Column{Name: "New", Type: "string", Docstring: ptr("fresh")})),
			observed: observed(Column{Name: "Ts", Type: "datetime", Docstring: ptr("old")}),
			steps:    []StepKind{StepAlterMergeSchema, StepSetColumnDocstrings},
			cmds:     []string{".alter-merge table ['T'] (['Ts']:datetime, ['New']:string)", ".alter table ['T'] column-docstrings (['Ts']:\"when\", ['New']:\"fresh\")"},
		},
		"columnDocstringEqual": {
			desired:  FromParams("T", params(common.Column{Name: "Ts", Type: "datetime", Docstring: ptr("when")})),
			observed: observed(Column{Name: "Ts", Type: "datetime", Docstring: ptr("when")}),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			p, err := Diff(tc.desired, tc.observed)
			if tc.blocked {
				var b *kerrors.Blocked
				if !errors.As(err, &b) || b.Reason != kerrors.ReasonUnsupportedColumnTypeChange {
					t.Fatalf("expected Blocked(UnsupportedColumnTypeChange), got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.steps, kinds(p), cmp.Transformer("nilToEmpty", func(s []StepKind) []StepKind {
				if s == nil {
					return []StepKind{}
				}
				return s
			})); diff != "" {
				t.Errorf("steps: %s", diff)
			}
			if diff := cmp.Diff(tc.drift, p.DriftColumns); diff != "" && (len(tc.drift) != 0 || len(p.DriftColumns) != 0) {
				t.Errorf("drift: %s", diff)
			}
			for i, want := range tc.cmds {
				if got := p.Steps[i].Command.String(); got != want {
					t.Errorf("cmd[%d]\n got %s\nwant %s", i, got, want)
				}
			}
			if p.Empty() != (len(tc.steps) == 0) {
				t.Error("Empty mismatch")
			}
		})
	}
}

func TestBuildCreateDelete(t *testing.T) {
	p := params(common.Column{Name: "Ts", Type: "date"}, common.Column{Name: "P", Type: "dynamic", Docstring: ptr("body")})
	p.Folder, p.Docstring = ptr("Raw"), ptr("Landing")
	cmds := BuildCreate(FromParams("Raw Events", p))
	want := []string{
		".create table ['Raw Events'] (['Ts']:datetime, ['P']:dynamic) with (docstring=\"Landing\", folder=\"Raw\")",
		".alter table ['Raw Events'] column-docstrings (['P']:\"body\")",
	}
	if len(cmds) != len(want) {
		t.Fatalf("got %d commands", len(cmds))
	}
	for i := range want {
		if cmds[i].String() != want[i] {
			t.Errorf("create[%d]\n got %s\nwant %s", i, cmds[i].String(), want[i])
		}
	}
	if got := BuildCreate(FromParams("T", params(common.Column{Name: "A", Type: "string"})))[0].String(); got != ".create table ['T'] (['A']:string)" {
		t.Errorf("minimal create: %s", got)
	}
	if got := BuildDelete("T").String(); got != ".drop table ['T'] ifexists" {
		t.Errorf("delete: %s", got)
	}
	if got := ShowCslSchema("T").String(); got != ".show table ['T'] cslschema" {
		t.Errorf("show: %s", got)
	}
}

func TestParse(t *testing.T) {
	d, err := schema.Parse([]byte(schematest.DatabaseJSON), "Telemetry")
	if err != nil {
		t.Fatal(err)
	}
	o := FromSchema(d.Tables["RawEvents"])
	if len(o.Columns) != 2 || o.Columns[0].Type != "datetime" || o.Columns[1].Type != "dynamic" || *o.Columns[1].Docstring != "Raw JSON body" || o.Folder != "Raw" || o.Docstring != "Landing table" {
		t.Errorf("FromSchema: %+v", o)
	}
	res := kusto.NewResult(kusto.NewTable("Table_0", []string{"TableName", "Schema", "DatabaseName", "Folder", "DocString"},
		[]any{"T", "Ts:datetime,['My Col']:string,Payload:dynamic", "DB", "F", "D"}))
	c, ok := ParseCslSchema(res)
	if !ok || len(c.Columns) != 3 || c.Columns[1].Name != "My Col" || c.Columns[1].Type != "string" || c.Folder != "F" {
		t.Errorf("ParseCslSchema: %+v", c)
	}
	if _, ok := ParseCslSchema(kusto.NewResult()); ok {
		t.Error("empty result must not parse")
	}
	obs := Observation(o, []string{"X"})
	if len(obs.Columns) != 2 || obs.Columns[1].Docstring != "Raw JSON body" || obs.DriftColumns[0] != "X" {
		t.Errorf("Observation: %+v", obs)
	}
	if DescribeDrift(nil) != "" || DescribeDrift([]string{"A"}) == "" {
		t.Error("DescribeDrift")
	}
}
