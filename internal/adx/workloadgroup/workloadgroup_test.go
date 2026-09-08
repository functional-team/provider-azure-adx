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

package workloadgroup

import (
	"encoding/json"
	"testing"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
)

func TestCommands(t *testing.T) {
	if ShowCmd("My WG").String() != ".show workload_group ['My WG']" || ShowAllCmd().String() != ".show workload_groups" || DropCmd("x").String() != ".drop workload_group ['x']" {
		t.Error("simple commands")
	}
	c, err := CreateOrAlterCmd("wg", map[string]any{"RequestLimitsPolicy": map[string]any{"MaxMemoryPerQueryPerNode": map[string]any{"IsRelaxable": true, "Value": 6442450944}}})
	if err != nil || c.String() != ".create-or-alter workload_group ['wg'] @'{\"RequestLimitsPolicy\":{\"MaxMemoryPerQueryPerNode\":{\"IsRelaxable\":true,\"Value\":6442450944}}}'" {
		t.Errorf("create: %s %v", c, err)
	}
	if !IsBuiltIn("default") || !IsBuiltIn(" Internal ") || IsBuiltIn("mine") {
		t.Error("IsBuiltIn")
	}
}

func TestParse(t *testing.T) {
	cols := []string{"WorkloadGroupName", "WorkloadGroup"}
	res := kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"default", `{"RequestLimitsPolicy":{}}`}, []any{"wg", `{"RequestQueuingPolicy":{"IsEnabled":true}}`}, []any{"empty", nil}))
	raw, ok := ParseShow(res, "wg")
	if !ok || string(raw) != `{"RequestQueuingPolicy":{"IsEnabled":true}}` {
		t.Errorf("ParseShow: %s %v", raw, ok)
	}
	if _, ok := ParseShow(res, "empty"); ok {
		t.Error("null group must not exist")
	}
	if _, ok := ParseShow(res, "nope"); ok {
		t.Error("unknown group must not exist")
	}
	all := ParseShowAll(res)
	if len(all) != 3 || all["default"] == nil {
		t.Errorf("ParseShowAll: %v", all)
	}
	r, err := Compare(map[string]any{"RequestQueuingPolicy": map[string]any{"IsEnabled": true}}, json.RawMessage(`{"RequestQueuingPolicy":{"IsEnabled":true},"RequestLimitsPolicy":{}}`))
	if err != nil || !r.Equal {
		t.Errorf("Compare: %+v %v", r, err)
	}
	r, _ = Compare(map[string]any{"RequestQueuingPolicy": map[string]any{"IsEnabled": false}}, json.RawMessage(`{"RequestQueuingPolicy":{"IsEnabled":true}}`))
	if r.Equal {
		t.Error("drift must be detected")
	}
}
