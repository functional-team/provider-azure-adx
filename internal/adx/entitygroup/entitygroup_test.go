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

package entitygroup

import (
	"encoding/json"
	"testing"

	"github.com/functional-team/provider-azure-adx/internal/adx/schema"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema/schematest"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
)

func TestBuild(t *testing.T) {
	c := BuildCreateOrAlter("EG", []string{" database('d').table('t') ", "cluster('c').database('d')"})
	if c.String() != ".create-or-alter entity_group ['EG'] (cluster('c').database('d'), database('d').table('t'))" {
		t.Error(c)
	}
	if BuildDelete("EG").String() != ".drop entity_group ['EG']" || ShowOne("EG").String() != ".show entity_group ['EG']" {
		t.Error("simple commands")
	}
}

func TestParse(t *testing.T) {
	d, err := schema.Parse([]byte(schematest.DatabaseJSON), "Telemetry")
	if err != nil {
		t.Fatal(err)
	}
	o := FromSchema(d.EntityGroups["EG"])
	if o.Name != "EG" || len(o.Entities) != 1 || o.Entities[0] != "cluster('c').database('d')" {
		t.Errorf("FromSchema: %+v", o)
	}
	arr := kusto.NewResult(kusto.NewTable("Table_0", []string{"Name", "Entities"}, []any{"EG", json.RawMessage(`["b","a"]`)}))
	if o, ok := ParseShow(arr, "EG"); !ok || len(o.Entities) != 2 || o.Entities[0] != "a" {
		t.Errorf("ParseShow array: %+v %v", o, ok)
	}
	str := kusto.NewResult(kusto.NewTable("Table_0", []string{"Name", "Entities"}, []any{"", "database('d, x').table('t'), cluster('c').database('d')"}))
	if o, ok := ParseShow(str, "EG"); !ok || o.Name != "EG" || len(o.Entities) != 2 || o.Entities[1] != "database('d, x').table('t')" {
		t.Errorf("ParseShow string: %+v %v", o, ok)
	}
	if _, ok := ParseShow(kusto.NewResult(), "EG"); ok {
		t.Error("empty result")
	}
	if !Equal([]string{"b", " a "}, []string{"a", "b"}) || Equal([]string{"a"}, []string{"a", "b"}) || Equal([]string{"a"}, []string{"c"}) {
		t.Error("Equal")
	}
	if len(Texts([]string{"", " x "})) != 1 {
		t.Error("Texts")
	}
	if len(Observation(o).Entities) != 1 {
		t.Error("Observation")
	}
}
