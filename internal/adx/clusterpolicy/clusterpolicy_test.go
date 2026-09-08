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

package clusterpolicy

import (
	"encoding/json"
	"testing"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
)

func TestCommands(t *testing.T) {
	d := Def{Name: "callout"}
	if d.ShowCmd().String() != ".show cluster policy callout" {
		t.Error(d.ShowCmd())
	}
	alter, err := d.AlterCmd([]map[string]any{{"CalloutType": "sql", "CalloutUriRegex": ".*", "CanCall": true}})
	if err != nil || alter.String() != ".alter cluster policy callout @'[{\"CalloutType\":\"sql\",\"CalloutUriRegex\":\".*\",\"CanCall\":true}]'" {
		t.Errorf("alter: %s %v", alter, err)
	}
	del, ok := d.DeleteCmd()
	if !ok || del.String() != ".delete cluster policy callout" {
		t.Error(del)
	}
	if _, ok := (Def{Name: "capacity", NoDelete: true}).DeleteCmd(); ok {
		t.Error("NoDelete must suppress the delete command")
	}
	// Kusto rejects ".alter cluster policy capacity": "The '.alter' command is
	// not supported for CapacityPolicy. Please use '.alter-merge' instead."
	merge, err := (Def{Name: "capacity", NoDelete: true, Merge: true}).AlterCmd(map[string]any{"IsEnabled": true})
	if err != nil || merge.String() != ".alter-merge cluster policy capacity @'{\"IsEnabled\":true}'" {
		t.Errorf("merge alter: %q %v", merge.String(), err)
	}
	custom := Def{Name: "request_classification", Alter: func(any) (cmd.Command, error) {
		return cmd.New(".alter cluster policy request_classification @'{\"IsEnabled\":true}'", cmd.Pipe("T")), nil
	}}
	c, _ := custom.AlterCmd(nil)
	if c.String() != ".alter cluster policy request_classification @'{\"IsEnabled\":true}'\n<| T" {
		t.Errorf("custom alter: %q", c.String())
	}
	if _, err := d.AlterCmd(make(chan int)); err == nil {
		t.Error("unmarshalable desired must fail")
	}
}

func TestCompareAndParse(t *testing.T) {
	d := Def{Name: "managed_identity", SetFields: []string{"AllowedUsages"}}
	r, err := d.CompareWith([]map[string]any{{"ObjectId": "system", "AllowedUsages": "NativeIngestion, ExternalTable"}}, json.RawMessage(`[{"ObjectId":"system","AllowedUsages":"ExternalTable, NativeIngestion"}]`))
	if err != nil || !r.Equal {
		t.Errorf("set comparison: %+v %v", r, err)
	}
	h := Def{Name: "multidatabaseadmins", HashOnly: true}
	r, err = h.CompareWith(map[string]any{"a": 1}, json.RawMessage(`{"b":2}`))
	if err != nil || !r.Equal || len(r.DesiredTexts) != 1 || r.ObservedTexts[0] != `{"b":2}` {
		t.Errorf("hash only: %+v %v", r, err)
	}
	cols := []string{"PolicyName", "EntityName", "Policy", "ChildEntities", "EntityType"}
	raw, ok := ParseShow(kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"CalloutPolicy", "", `[{"CalloutType":"sql"}]`, nil, "Cluster"})))
	if !ok || string(raw) != `[{"CalloutType":"sql"}]` {
		t.Errorf("ParseShow: %s %v", raw, ok)
	}
	if _, ok := ParseShow(kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"CalloutPolicy", "", "null", nil, "Cluster"}))); ok {
		t.Error("null must not exist")
	}
}
