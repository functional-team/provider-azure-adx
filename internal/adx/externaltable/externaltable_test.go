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

package externaltable

import (
	"strings"
	"testing"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
)

func ptr[T any](v T) *T { return &v }

var showCols = []string{"TableName", "TableType", "Folder", "DocString", "Properties", "ConnectionStrings", "Partitions", "PathFormat"}

func storageParams() v1alpha1.ExternalTableParameters {
	return v1alpha1.ExternalTableParameters{
		Database:    "DB",
		Kind:        v1alpha1.ExternalTableKindStorage,
		Columns:     []common.Column{{Name: "Timestamp", Type: "datetime"}, {Name: "Payload", Type: "dynamic"}},
		PartitionBy: ptr("Date:datetime = bin(Timestamp, 1d)"),
		PathFormat:  ptr(`datetime_pattern("yyyy/MM/dd", Date)`),
		DataFormat:  ptr("parquet"),
		Properties:  &v1alpha1.ExternalTableProperties{Folder: ptr("External"), Compressed: ptr(true), CompressionType: ptr("snappy")},
	}
}

func TestBuildCreateOrAlter(t *testing.T) {
	d := FromParams("Exports", storageParams(), []string{"https://acct.blob.core.windows.net/exports;managed_identity=system", "https://acct2.blob.core.windows.net/c;sig=SECRET"})
	got := BuildCreateOrAlter(d).String()
	want := `.create-or-alter external table ['Exports'] (['Timestamp']:datetime, ['Payload']:dynamic) kind=storage partition by (Date:datetime = bin(Timestamp, 1d)) pathformat=(datetime_pattern("yyyy/MM/dd", Date)) dataformat=parquet (h@'https://acct.blob.core.windows.net/exports;managed_identity=system', h@'https://acct2.blob.core.windows.net/c;sig=SECRET') with (compressed=true, compressionType="snappy", folder="External")`
	if got != want {
		t.Errorf("\n got %s\nwant %s", got, want)
	}
	if red := BuildCreateOrAlter(d).Redacted(); strings.Contains(red, "SECRET") || strings.Contains(red, "managed_identity") {
		t.Errorf("redacted command leaks secrets: %s", red)
	}
	delta := FromParams("D", v1alpha1.ExternalTableParameters{Database: "DB", Kind: v1alpha1.ExternalTableKindDelta, PartitionBy: ptr("ignored")}, []string{"https://a/b;managed_identity=system"})
	if got := BuildCreateOrAlter(delta).String(); got != ".create-or-alter external table ['D'] kind=delta (h@'https://a/b;managed_identity=system')" {
		t.Errorf("delta: %s", got)
	}
	if BuildDelete("E").String() != ".drop external table ['E']" || Show("E").String() != ".show external table ['E']" || ShowAll().String() != ".show external tables" || ShowCslSchema("E").String() != ".show external table ['E'] cslschema" {
		t.Error("simple commands")
	}
	if FromParams("X", v1alpha1.ExternalTableParameters{}, nil).Kind != "storage" {
		t.Error("kind defaults to storage")
	}
}

func TestParse(t *testing.T) {
	res := kusto.NewResult(kusto.NewTable("Table_0", showCols,
		[]any{"Exports", "Blob", "External", "doc", `{"Format":"Parquet","Compressed":true,"CompressionType":"snappy"}`, `["https://acct.blob.core.windows.net/exports;*******","https://acct2.blob.core.windows.net/c;*******"]`, `[{"ColumnName":"Date","Kind":"DateTime"}]`, `datetime_pattern("yyyy/MM/dd", Date)`},
		[]any{"D", "Delta", "", "", nil, `https://a/b;*******`, nil, ""},
		[]any{"", "Blob", "", "", nil, nil, nil, ""}))
	m := ParseRows(res)
	if len(m) != 2 {
		t.Fatalf("ParseRows: %d", len(m))
	}
	e := m["Exports"]
	if e.Kind != "storage" || e.Folder != "External" || len(e.ConnectionStringURIs) != 2 || e.ConnectionStringURIs[1] != "https://acct2.blob.core.windows.net/c" || e.Partitions != `[{"ColumnName":"Date","Kind":"DateTime"}]` || propString(e.Properties, "format") != "Parquet" {
		t.Errorf("Exports: %+v", e)
	}
	if m["D"].Kind != "delta" || len(m["D"].ConnectionStringURIs) != 1 || m["D"].ConnectionStringURIs[0] != "https://a/b" || m["D"].Partitions != "" {
		t.Errorf("D: %+v", m["D"])
	}
	if o, ok := ParseOne(res, "Exports"); !ok || o.Name != "Exports" {
		t.Error("ParseOne")
	}
	if _, ok := ParseOne(kusto.NewResult(), "X"); ok {
		t.Error("ParseOne empty")
	}
	if o, ok := ParseOne(kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{"Only", "Blob", "", "", nil, nil, nil, ""})), "Other"); !ok || o.Name != "Only" {
		t.Error("ParseOne single fallback")
	}
	cols, ok := ParseCslSchema(kusto.NewResult(kusto.NewTable("Table_0", []string{"TableName", "Schema", "DatabaseName", "Folder", "DocString"}, []any{"Exports", "Timestamp:datetime,Payload:dynamic", "DB", "", ""})))
	if !ok || len(cols) != 2 || cols[1].Type != "dynamic" {
		t.Errorf("ParseCslSchema: %+v", cols)
	}
	if URIOf(" https://a/b/;sig=x ") != "https://a/b" || URIOf("plain") != "plain" {
		t.Error("URIOf")
	}
	// A SAS token lives in the query string and the service masks it, so a
	// desired connection string and what comes back only agree on the location.
	sas := "https://acct.blob.core.windows.net/exports?se=2027-09-09&sig=abc%3D"
	masked := "https://acct.blob.core.windows.net/exports?******"
	if URIOf(sas) != "https://acct.blob.core.windows.net/exports" || URIOf(sas) != URIOf(masked) {
		t.Errorf("URIOf must drop the SAS: %q vs %q", URIOf(sas), URIOf(masked))
	}
	if !sameURIs([]string{sas}, []string{masked}) {
		t.Error("a masked SAS must compare equal to the one that was sent")
	}
	obs := Observation(Observed{Kind: "storage", Columns: []Column{{"A", "string"}}, Properties: map[string]any{"Format": "Csv"}, ConnectionStringURIs: []string{"u"}})
	if obs.DataFormat != "Csv" || len(obs.Columns) != 1 || obs.ConnectionStringURIs[0] != "u" {
		t.Errorf("Observation: %+v", obs)
	}
}

func TestCompare(t *testing.T) {
	cs := []string{"https://acct.blob.core.windows.net/exports;managed_identity=system", "https://acct2.blob.core.windows.net/c;sig=SECRET"}
	observed := Observed{
		Name: "Exports", Kind: "storage", ColumnsLoaded: true,
		Columns:              []Column{{"Timestamp", "datetime"}, {"Payload", "dynamic"}},
		Folder:               "External",
		Properties:           map[string]any{"Format": "Parquet", "Compressed": true, "CompressionType": "snappy", "IncludeHeaders": "None"},
		ConnectionStringURIs: []string{"https://acct2.blob.core.windows.net/c", "https://acct.blob.core.windows.net/exports"},
	}
	cases := map[string]struct {
		mutate func(*v1alpha1.ExternalTableParameters, *Observed, *[]string)
		equal  bool
		reason string
	}{
		"equal":       {func(*v1alpha1.ExternalTableParameters, *Observed, *[]string) {}, true, ""},
		"kind":        {func(_ *v1alpha1.ExternalTableParameters, o *Observed, _ *[]string) { o.Kind = "delta" }, false, "kind"},
		"columnsType": {func(_ *v1alpha1.ExternalTableParameters, o *Observed, _ *[]string) { o.Columns[1].Type = "string" }, false, "columns"},
		"columnsOrder": {func(_ *v1alpha1.ExternalTableParameters, o *Observed, _ *[]string) {
			o.Columns[0], o.Columns[1] = o.Columns[1], o.Columns[0]
		}, false, "columns"},
		"columnsNotLoaded": {func(_ *v1alpha1.ExternalTableParameters, o *Observed, _ *[]string) {
			o.ColumnsLoaded = false
			o.Columns = nil
		}, true, ""},
		"dataFormat":        {func(p *v1alpha1.ExternalTableParameters, _ *Observed, _ *[]string) { p.DataFormat = ptr("csv") }, false, "dataFormat"},
		"dataFormatUnknown": {func(_ *v1alpha1.ExternalTableParameters, o *Observed, _ *[]string) { delete(o.Properties, "Format") }, true, ""},
		"csCount":           {func(_ *v1alpha1.ExternalTableParameters, _ *Observed, cs *[]string) { *cs = (*cs)[:1] }, false, "connection string"},
		"csURI": {func(_ *v1alpha1.ExternalTableParameters, _ *Observed, cs *[]string) {
			(*cs)[0] = "https://other/x;managed_identity=system"
		}, false, "connection string"},
		"csSecretOnlyChange": {func(_ *v1alpha1.ExternalTableParameters, _ *Observed, cs *[]string) {
			(*cs)[1] = "https://acct2.blob.core.windows.net/c;sig=ROTATED"
		}, true, ""},
		"folder": {func(p *v1alpha1.ExternalTableParameters, _ *Observed, _ *[]string) {
			p.Properties.Folder = ptr("Other")
		}, false, "folder"},
		"compressed": {func(p *v1alpha1.ExternalTableParameters, _ *Observed, _ *[]string) {
			p.Properties.Compressed = ptr(false)
		}, false, "compressed"},
		"compressionType": {func(p *v1alpha1.ExternalTableParameters, _ *Observed, _ *[]string) {
			p.Properties.CompressionType = ptr("gzip")
		}, false, "CompressionType"},
		"unsetPropIgnored": {func(p *v1alpha1.ExternalTableParameters, _ *Observed, _ *[]string) { p.Properties.IncludeHeaders = nil }, true, ""},
		"propCaseInsensitive": {func(p *v1alpha1.ExternalTableParameters, o *Observed, _ *[]string) {
			p.Properties.IncludeHeaders = ptr("None")
			o.Properties["includeheaders"] = "none"
			delete(o.Properties, "IncludeHeaders")
		}, true, ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			p := storageParams()
			o := observed
			o.Columns = append([]Column(nil), observed.Columns...)
			o.Properties = map[string]any{}
			for k, v := range observed.Properties {
				o.Properties[k] = v
			}
			c := append([]string(nil), cs...)
			tc.mutate(&p, &o, &c)
			d := FromParams("Exports", p, c)
			diff := Compare(d, o)
			if diff.Equal != tc.equal {
				t.Errorf("Equal = %v (%s), want %v", diff.Equal, diff.String(), tc.equal)
			}
			if tc.reason != "" && !strings.Contains(diff.String(), tc.reason) {
				t.Errorf("reason %q missing in %q", tc.reason, diff.String())
			}
		})
	}
}

func TestTextsAndSecretHash(t *testing.T) {
	d := FromParams("E", storageParams(), []string{"a;sig=1"})
	dt := DesiredTexts(d)
	if dt[0] != "Date:datetime = bin(Timestamp, 1d)" || dt[1] != `datetime_pattern("yyyy/MM/dd", Date)` {
		t.Errorf("DesiredTexts: %q", dt)
	}
	ot := ObservedTexts(Observed{Partitions: `[{"ColumnName":"Date"}]`, PathFormat: `datetime_pattern("yyyy/MM/dd", Date)` + "\r\n"})
	if ot[0] != `[{"ColumnName":"Date"}]` || ot[1] != dt[1] {
		t.Errorf("ObservedTexts: %q", ot)
	}
	first := SecretHash([]string{"x"})
	if SecretHash([]string{"a;sig=1"}) == SecretHash([]string{"a;sig=2"}) || first != SecretHash([]string{"x"}) {
		t.Error("SecretHash")
	}
}
