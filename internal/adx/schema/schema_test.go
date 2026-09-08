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

package schema

import (
	"context"
	"testing"

	"github.com/functional-team/provider-azure-adx/internal/adx/schema/schematest"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
)

// The service returns an entity group's entities as a bare array under the
// group name, not as an object with Name and Entities. Assuming the latter made
// every Table observe fail against a real cluster with "cannot unmarshal array
// into Go struct field Database.Databases.EntityGroups", because one bad entity
// group takes down the whole database schema document.
func TestParseEntityGroupShapes(t *testing.T) {
	for _, tc := range []struct {
		name string
		json string
	}{
		{"array (service)", `{"Databases":{"DB":{"EntityGroups":{"EG":["cluster('c').database('d')"]}}}}`},
		{"object with Name and Entities", `{"Databases":{"DB":{"EntityGroups":{"EG":{"Name":"EG","Entities":["cluster('c').database('d')"]}}}}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, err := Parse([]byte(tc.json), "DB")
			if err != nil {
				t.Fatal(err)
			}
			g := d.EntityGroups["EG"]
			if g.Name != "EG" || len(g.Entities) != 1 || g.Entities[0] != "cluster('c').database('d')" {
				t.Errorf("entity group parse: %+v", g)
			}
		})
	}
	// A table must still parse when an entity group sits next to it.
	d, err := Parse([]byte(`{"Databases":{"DB":{"Tables":{"T":{"Name":"T"}},"EntityGroups":{"EG":["cluster('c').database('d')"]}}}}`), "DB")
	if err != nil {
		t.Fatalf("table alongside entity group: %v", err)
	}
	if _, ok := d.Tables["T"]; !ok {
		t.Error("table missing")
	}
}

func TestParse(t *testing.T) {
	d, err := Parse([]byte(schematest.DatabaseJSON), "Telemetry")
	if err != nil {
		t.Fatal(err)
	}
	tbl, ok := d.Tables["RawEvents"]
	if !ok || len(tbl.OrderedColumns) != 2 || tbl.OrderedColumns[1].CslType != "dynamic" || tbl.OrderedColumns[1].DocString != "Raw JSON body" || tbl.Folder != "Raw" {
		t.Errorf("table parse: %+v", tbl)
	}
	f := d.Functions["TakeSome"]
	if f.Body != "{ RawEvents | take limit }" || len(f.InputParameters) != 2 || f.InputParameters[0].CslDefaultValue != "100" || len(f.InputParameters[1].Columns) != 1 {
		t.Errorf("function parse: %+v", f)
	}
	if d.MaterializedViews["LatestEvents"].SourceTable != "RawEvents" || d.ExternalTables["Exports"].Folder != "External" || len(d.EntityGroups["EG"].Entities) != 1 {
		t.Error("mv/external/entitygroup parse")
	}
	// Single database with a different key still resolves.
	if _, err := Parse([]byte(schematest.DatabaseJSON), "Other"); err != nil {
		t.Errorf("single database fallback: %v", err)
	}
	if _, err := Parse([]byte(`{"Databases":{"A":{},"B":{}}}`), "C"); err == nil {
		t.Error("ambiguous database must error")
	}
	if _, err := Parse([]byte("nope"), "X"); err == nil {
		t.Error("invalid JSON must error")
	}
	empty, err := Parse([]byte(`{"Databases":{"E":{"Name":"E"}}}`), "E")
	if err != nil || empty.Tables == nil || empty.Functions == nil {
		t.Error("maps must be non-nil")
	}
}

func TestLoad(t *testing.T) {
	kc := fake.New("http://e", fake.Handler{Match: ".show database ['Telemetry'] schema as json",
		Result: kusto.NewResult(kusto.NewTable("Table_0", []string{"DatabaseSchema"}, []any{schematest.DatabaseJSON}))})
	d, err := Load(context.Background(), kc, "Telemetry")
	if err != nil || len(d.Tables) != 2 {
		t.Fatalf("Load: %v %v", d, err)
	}
	empty := fake.New("http://e", fake.Handler{Result: kusto.NewResult(kusto.NewTable("Table_0", []string{"DatabaseSchema"}))})
	if _, err := Load(context.Background(), empty, "Telemetry"); err == nil {
		t.Error("no rows must error")
	}
}
