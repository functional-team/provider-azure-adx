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

package policy

import (
	"encoding/json"
	"testing"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
)

func TestEntityRender(t *testing.T) {
	cases := map[string]struct {
		e       Entity
		want    string
		display string
		err     bool
	}{
		"table":    {Entity{Kind: common.EntityKindTable, Database: "DB", Name: "T"}, "table ['T']", "['DB'].['T']", false},
		"database": {Entity{Kind: common.EntityKindDatabase, Database: "DB"}, "database ['DB']", "['DB']", false},
		"mv":       {Entity{Kind: common.EntityKindMaterializedView, Database: "DB", Name: "MV"}, "materialized-view ['MV']", "['DB'].['MV']", false},
		"external": {Entity{Kind: common.EntityKindExternalTable, Database: "DB", Name: "E"}, "external table ['E']", "['DB'].['E']", false},
		"column":   {Entity{Kind: common.EntityKindTable, Database: "DB", Name: "T", Column: "C"}, "column ['T'].['C']", "['DB'].['T'].['C']", false},
		"function": {Entity{Kind: common.EntityKindFunction, Database: "DB", Name: "F"}, "", "", true},
		"colOnDB":  {Entity{Kind: common.EntityKindDatabase, Database: "DB", Column: "C"}, "", "", true},
		"unknown":  {Entity{Kind: "Weird"}, "", "", true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := tc.e.Render()
			if tc.err {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil || got != tc.want {
				t.Errorf("Render() = %q, %v; want %q", got, err, tc.want)
			}
			if d := tc.e.Display(); d != tc.display {
				t.Errorf("Display() = %q, want %q", d, tc.display)
			}
		})
	}
}

func TestCommands(t *testing.T) {
	d := Def{Name: "retention"}
	e := Entity{Kind: common.EntityKindTable, Database: "DB", Name: "T"}
	show, _ := d.ShowCmd(e)
	if show.String() != ".show table ['T'] policy retention" {
		t.Error(show)
	}
	if d.ShowAllTablesCmd().String() != ".show table * policy retention" {
		t.Error("show all")
	}
	alter, err := d.AlterCmd(e, map[string]any{"SoftDeletePeriod": "365.00:00:00"})
	if err != nil || alter.String() != ".alter table ['T'] policy retention @'{\"SoftDeletePeriod\":\"365.00:00:00\"}'" {
		t.Errorf("alter: %s %v", alter, err)
	}
	del, _ := d.DeleteCmd(e)
	if del.String() != ".delete table ['T'] policy retention" {
		t.Error(del)
	}
	custom := Def{Name: "ingestiontime", Alter: func(e Entity, desired any) (cmd.Command, error) {
		ent, _ := e.Render()
		return cmd.New(".alter ", ent, " policy ingestiontime true"), nil
	}}
	c, _ := custom.AlterCmd(e, nil)
	if c.String() != ".alter table ['T'] policy ingestiontime true" {
		t.Error(c)
	}
	if _, err := d.AlterCmd(Entity{Kind: common.EntityKindFunction}, nil); err == nil {
		t.Error("function entity must fail")
	}
	if _, err := d.AlterCmd(e, make(chan int)); err == nil {
		t.Error("unmarshalable desired must fail")
	}
}

func TestParseShow(t *testing.T) {
	cols := []string{"PolicyName", "EntityName", "Policy", "ChildEntities", "EntityType"}
	res := kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"RetentionPolicy", "[DB].[T]", `{"SoftDeletePeriod":"365.00:00:00"}`, nil, "Table"}))
	raw, ok := ParseShow(res)
	if !ok || string(raw) != `{"SoftDeletePeriod":"365.00:00:00"}` {
		t.Errorf("ParseShow: %s %v", raw, ok)
	}
	if _, ok := ParseShow(kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"RetentionPolicy", "[DB].[T]", "null", nil, "Table"}))); ok {
		t.Error("null policy must not exist")
	}
	if _, ok := ParseShow(kusto.NewResult()); ok {
		t.Error("no rows must not exist")
	}
	all := ParseShowAll(kusto.NewResult(kusto.NewTable("Table_0", cols,
		[]any{"RetentionPolicy", "[DB].[A]", `{"SoftDeletePeriod":"1.00:00:00"}`, nil, "Table"},
		[]any{"RetentionPolicy", "[DB].[B]", nil, nil, "Table"},
		[]any{"RetentionPolicy", "['DB'].['C d']", `{}`, nil, "Table"},
		[]any{"RetentionPolicy", "", nil, nil, "Table"})))
	if len(all) != 3 || all["A"] == nil || all["B"] != nil || all["C d"] == nil {
		t.Errorf("ParseShowAll: %v", all)
	}
	if Compact(json.RawMessage(" {\"a\" : 1} ")) != `{"a":1}` || Compact(nil) != "" || Compact(json.RawMessage("not json")) != "not json" {
		t.Error("Compact")
	}
}

func TestCompare(t *testing.T) {
	type retention struct {
		SoftDeletePeriod *string `json:"SoftDeletePeriod,omitempty"`
		Recoverability   *string `json:"Recoverability,omitempty"`
	}
	s := func(v string) *string { return &v }
	cases := map[string]struct {
		desired  any
		observed string
		kql      []string
		sets     []string
		equal    bool
		texts    int
		diffHas  string
	}{
		"subsetEqual":       {retention{SoftDeletePeriod: s("365d")}, `{"SoftDeletePeriod":"365.00:00:00","Recoverability":"Enabled"}`, nil, nil, true, 0, ""},
		"valueDiffers":      {retention{SoftDeletePeriod: s("30d")}, `{"SoftDeletePeriod":"365.00:00:00"}`, nil, nil, false, 0, "SoftDeletePeriod"},
		"missingKey":        {retention{Recoverability: s("Disabled")}, `{"SoftDeletePeriod":"365.00:00:00"}`, nil, nil, false, 0, "missing"},
		"observedNull":      {retention{Recoverability: s("Disabled")}, "null", nil, nil, false, 0, ""},
		"numbers":           {map[string]any{"MaximumNumberOfItems": 500, "Rate": 2.5}, `{"MaximumNumberOfItems":500,"Rate":"2.5"}`, nil, nil, true, 0, ""},
		"numberDiff":        {map[string]any{"MaximumNumberOfItems": 500}, `{"MaximumNumberOfItems":400}`, nil, nil, false, 0, "want 500"},
		"boolFromString":    {map[string]any{"IsEnabled": true}, `{"IsEnabled":"True"}`, nil, nil, true, 0, ""},
		"boolDiff":          {map[string]any{"IsEnabled": true}, `{"IsEnabled":false}`, nil, nil, false, 0, "IsEnabled"},
		"valueWrapper":      {map[string]any{"DataHotSpan": "31d"}, `{"DataHotSpan":{"Value":"31.00:00:00"}}`, nil, nil, true, 0, ""},
		"datetime":          {map[string]any{"EffectiveDateTime": "2026-01-01T00:00:00Z"}, `{"EffectiveDateTime":"2026-01-01T00:00:00.0000000"}`, nil, nil, true, 0, ""},
		"arrayEqual":        {[]map[string]any{{"TagPrefix": "drop-by:", "RetentionPeriod": "3d"}}, `[{"TagPrefix":"drop-by:","RetentionPeriod":"3.00:00:00"}]`, nil, nil, true, 0, ""},
		"arrayLength":       {[]map[string]any{{"TagPrefix": "a"}}, `[{"TagPrefix":"a"},{"TagPrefix":"b"}]`, nil, nil, false, 0, "elements"},
		"emptyArrayMissing": {map[string]any{"HotWindows": []any{}}, `{"DataHotSpan":"1d"}`, nil, nil, true, 0, ""},
		"emptyArrayVsFull":  {map[string]any{"HotWindows": []any{}}, `{"HotWindows":[{"MinValue":"x"}]}`, nil, nil, false, 0, "elements"},
		"kqlCollected":      {[]map[string]any{{"Source": "S", "Query": "S | where a == 1  \r\n"}}, `[{"Source":"S","Query":"S | where a == 1","IsEnabled":true}]`, []string{"Query"}, nil, true, 1, ""},
		"kqlNotCompared":    {map[string]any{"Query": "A"}, `{"Query":"B"}`, []string{"Query"}, nil, true, 1, ""},
		"setField":          {map[string]any{"AllowedUsages": "NativeIngestion, ExternalTable"}, `{"AllowedUsages":"externaltable,nativeingestion"}`, nil, []string{"AllowedUsages"}, true, 0, ""},
		"setFieldDiff":      {map[string]any{"AllowedUsages": "NativeIngestion"}, `{"AllowedUsages":"ExternalTable"}`, nil, []string{"AllowedUsages"}, false, 0, "AllowedUsages"},
		"typeMismatch":      {map[string]any{"A": map[string]any{"B": 1}}, `{"A":"x"}`, nil, nil, false, 0, "object"},
		"nested":            {map[string]any{"Lookback": map[string]any{"Kind": "Custom", "CustomPeriod": "1d"}}, `{"Lookback":{"Kind":"Custom","CustomPeriod":"1.00:00:00"},"AllowMerge":true}`, nil, nil, true, 0, ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r, err := Compare(tc.desired, json.RawMessage(tc.observed), tc.kql, tc.sets)
			if err != nil {
				t.Fatal(err)
			}
			if r.Equal != tc.equal {
				t.Errorf("Equal = %v, want %v (diff %q)", r.Equal, tc.equal, r.Diff)
			}
			if len(r.DesiredTexts) != tc.texts || len(r.ObservedTexts) != tc.texts {
				t.Errorf("texts = %d/%d, want %d", len(r.DesiredTexts), len(r.ObservedTexts), tc.texts)
			}
			if tc.diffHas != "" && !contains(r.Diff, tc.diffHas) {
				t.Errorf("Diff %q should contain %q", r.Diff, tc.diffHas)
			}
			if tc.texts == 1 && tc.equal && name == "kqlCollected" && r.DesiredTexts[0] != r.ObservedTexts[0] {
				t.Errorf("stage 1 should equalize the texts: %q vs %q", r.DesiredTexts[0], r.ObservedTexts[0])
			}
		})
	}
	if _, err := Compare(make(chan int), nil, nil, nil); err == nil {
		t.Error("unmarshalable desired must error")
	}
	if _, err := Compare(map[string]any{}, json.RawMessage("not json"), nil, nil); err == nil {
		t.Error("invalid observed JSON must error")
	}
	hashOnly := Def{Name: "roworder", HashOnly: true}
	r, err := hashOnly.CompareWith(map[string]any{"a": 1}, json.RawMessage(`{"b": 2}`))
	if err != nil || !r.Equal || len(r.DesiredTexts) != 1 || r.ObservedTexts[0] != `{"b":2}` {
		t.Errorf("HashOnly: %+v %v", r, err)
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && (s == sub || indexOf(s, sub) >= 0))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
