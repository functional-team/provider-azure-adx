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

package ingestionmapping

import (
	"testing"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
)

func ct(s string) *common.ColumnType { c := common.ColumnType(s); return &c }

var cols = []string{"Name", "Kind", "Mapping", "LastUpdatedOn", "Database", "Table"}

func params() v1alpha1.IngestionMappingParameters {
	return v1alpha1.IngestionMappingParameters{Database: "DB", Table: "RawEvents", Kind: "json", Mapping: []v1alpha1.IngestionMappingColumn{
		{Column: "Timestamp", DataType: ct("date"), Properties: map[string]string{"path": "$.ts"}},
		{Column: "Payload", Properties: map[string]string{"Path": "$", "transform": "SourceLocation"}},
	}}
}

func TestBuild(t *testing.T) {
	d := FromParams("Raw JSON", params())
	if d.Kind != "json" || d.Mapping[0].DataType != "datetime" || d.Mapping[0].Properties["Path"] != "$.ts" || d.Mapping[1].Properties["Transform"] != "SourceLocation" {
		t.Errorf("FromParams: %+v", d)
	}
	c, err := BuildCreateOrAlter(d)
	if err != nil {
		t.Fatal(err)
	}
	want := ".create-or-alter table ['RawEvents'] ingestion json mapping \"Raw JSON\" @'[{\"Column\":\"Timestamp\",\"DataType\":\"datetime\",\"Properties\":{\"Path\":\"$.ts\"}},{\"Column\":\"Payload\",\"Properties\":{\"Path\":\"$\",\"Transform\":\"SourceLocation\"}}]'"
	if c.String() != want {
		t.Errorf("\n got %s\nwant %s", c.String(), want)
	}
	if BuildDelete("T", "Json", "M").String() != ".drop table ['T'] ingestion json mapping \"M\"" {
		t.Error("delete")
	}
	if ShowAll("T").String() != ".show table ['T'] ingestion mappings" || ShowOne("T", "CSV", "M").String() != ".show table ['T'] ingestion csv mapping \"M\"" {
		t.Error("show")
	}
	if CanonicalKey("constValue") != "ConstValue" || CanonicalKey("") != "" || CanonicalKey(" path") != "Path" {
		t.Error("CanonicalKey")
	}
	if Key("Json", "M") != "json/M" {
		t.Error("Key")
	}
}

func TestParse(t *testing.T) {
	raw := `[{"Column":"Timestamp","DataType":"System.DateTime","Properties":{"Path":"$.ts"},"CsvDataType":null},{"column":"Payload","datatype":"dynamic","properties":{"path":"$","Transform":"SourceLocation","Ordinal":null}}]`
	res := kusto.NewResult(kusto.NewTable("Table_0", cols,
		[]any{"RawJson", "Json", raw, "2026-09-07T10:00:00Z", "DB", "RawEvents"},
		[]any{"RawCsv", "Csv", `[{"Column":"A","Properties":{"Ordinal":"0"}}]`, "", "DB", "RawEvents"},
		[]any{"", "Json", "", "", "DB", "RawEvents"}))
	m, err := ParseRows(res)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatalf("rows: %v", m)
	}
	o := m[Key("json", "RawJson")]
	if o.Kind != "json" || len(o.Mapping) != 2 || o.Mapping[0].DataType != "datetime" || o.Mapping[1].Properties["Path"] != "$" || o.Mapping[1].Properties["Transform"] != "SourceLocation" || o.LastUpdatedOn == "" {
		t.Errorf("parsed: %+v", o)
	}
	if _, ok := o.Mapping[1].Properties["Ordinal"]; ok {
		t.Error("null properties must be dropped")
	}
	if o.Raw == "" || o.Raw[0] != '[' {
		t.Errorf("raw: %q", o.Raw)
	}
	one, ok, err := ParseOne(kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"RawJson", "", raw, "", "DB", "RawEvents"})), "json", "RawJson")
	if err != nil || !ok || one.Name != "RawJson" {
		t.Errorf("ParseOne fallback: %+v %v %v", one, ok, err)
	}
	if _, ok, _ := ParseOne(kusto.NewResult(), "json", "X"); ok {
		t.Error("empty result")
	}
	if _, err := ParseRows(kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"Bad", "Json", "not json", "", "DB", "T"}))); err == nil {
		t.Error("invalid JSON must error")
	}
	if c, err := ParseMapping("null"); err != nil || c != nil {
		t.Error("null mapping")
	}
	obs := Observation(o)
	if obs.Kind != "json" || obs.Table != "RawEvents" || obs.Mapping == "" {
		t.Errorf("Observation: %+v", obs)
	}
}

func TestEqual(t *testing.T) {
	d := FromParams("M", params())
	observed := Observed{Mapping: []Column{
		{Column: "Timestamp", DataType: "datetime", Properties: map[string]string{"path": "$.ts", "Extra": "x"}},
		{Column: "Payload", DataType: "dynamic", Properties: map[string]string{"Path": "$", "transform": "SourceLocation"}},
	}}
	if ok, diff := Equal(d, observed); !ok {
		t.Errorf("expected equal: %s", diff)
	}
	cases := map[string]func(o *Observed){
		"length":    func(o *Observed) { o.Mapping = o.Mapping[:1] },
		"column":    func(o *Observed) { o.Mapping[0].Column = "Ts" },
		"dataType":  func(o *Observed) { o.Mapping[0].DataType = "string" },
		"propValue": func(o *Observed) { o.Mapping[0].Properties = map[string]string{"Path": "$.time"} },
		"propMiss":  func(o *Observed) { o.Mapping[1].Properties = map[string]string{"Path": "$"} },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			o := Observed{Mapping: []Column{
				{Column: "Timestamp", DataType: "datetime", Properties: map[string]string{"Path": "$.ts"}},
				{Column: "Payload", DataType: "dynamic", Properties: map[string]string{"Path": "$", "Transform": "SourceLocation"}},
			}}
			mutate(&o)
			if ok, diff := Equal(d, o); ok || diff == "" {
				t.Errorf("expected difference for %s", name)
			}
		})
	}
	// Unset dataType in the spec is not compared.
	if ok, _ := Equal(d, Observed{Mapping: []Column{observed.Mapping[0], {Column: "Payload", DataType: "string", Properties: observed.Mapping[1].Properties}}}); !ok {
		t.Error("unset dataType must be ignored")
	}
}
