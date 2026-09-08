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

package materializedview

import (
	"strings"
	"testing"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

func ptr[T any](v T) *T { return &v }

var showCols = []string{"Name", "SourceTable", "Query", "MaterializedTo", "LastRun", "LastRunResult", "IsHealthy", "IsEnabled", "Folder", "DocString", "AutoUpdateSchema", "EffectiveDateTime", "Lookback", "LookbackColumn"}

func row(name, src, query string, enabled bool, lookback any) []any {
	return []any{name, src, query, nil, nil, "Success", true, enabled, "Views", "doc", false, "", lookback, "Timestamp"}
}

func TestFromParams(t *testing.T) {
	p := v1alpha1.MaterializedViewParameters{Database: "DB", SourceTable: ptr("Raw"), Query: "Raw | summarize count() by Id", Lookback: ptr(common.Timespan("6h")), EffectiveDateTime: ptr("2026-01-01T00:00:00Z")}
	d, err := FromParams("MV", p)
	if err != nil || d.SourceTable != "Raw" || *d.Lookback != 6*timespan.Hour || !d.Enabled || d.Backfill {
		t.Fatalf("FromParams: %+v %v", d, err)
	}
	p.Lookback = ptr(common.Timespan("soon"))
	if _, err := FromParams("MV", p); err == nil || !strings.Contains(err.Error(), "lookback") {
		t.Errorf("bad lookback must fail: %v", err)
	}
	p.Lookback = nil
	p.EffectiveDateTime = ptr("yesterday")
	if _, err := FromParams("MV", p); err == nil || !strings.Contains(err.Error(), "effectiveDateTime") {
		t.Errorf("bad datetime must fail: %v", err)
	}
	if _, err := FromParams("MV", v1alpha1.MaterializedViewParameters{Query: "x"}); err == nil || !strings.Contains(err.Error(), "sourceTable is empty") {
		t.Errorf("missing source must fail: %v", err)
	}
	d, err = FromParams("MV", v1alpha1.MaterializedViewParameters{Query: "x", SourceMaterializedView: ptr("Other"), Enabled: ptr(false), Backfill: ptr(true)})
	if err != nil || d.SourceMaterializedView != "Other" || d.Enabled || !d.Backfill {
		t.Errorf("mv source: %+v %v", d, err)
	}
}

func TestBuild(t *testing.T) {
	d, _ := FromParams("Latest Events", v1alpha1.MaterializedViewParameters{
		SourceTable: ptr("RawEvents"), Query: "RawEvents\n| summarize arg_max(Timestamp, *) by DeviceId",
		Backfill: ptr(true), EffectiveDateTime: ptr("2026-01-01T00:00:00Z"), UpdateExtentsCreationTime: ptr(false),
		Lookback: ptr(common.Timespan("6h")), LookbackColumn: ptr("Timestamp"), AutoUpdateSchema: ptr(true),
		DimensionTables: []string{"DimDevices", "Dim'Sites"}, Folder: ptr("Views"), Docstring: ptr("latest"),
		AllowMaterializedViewsWithoutRowLevelSecurity: ptr(true), MaxSourceRecordsForSingleIngest: ptr(int64(1000)), Concurrency: ptr(int64(2)),
	})
	want := `.create async ifnotexists materialized-view with (Concurrency=2, MaxSourceRecordsForSingleIngest=1000, allowMaterializedViewsWithoutRowLevelSecurity=true, autoUpdateSchema=true, backfill=true, dimensionTables=dynamic(["DimDevices", "Dim'Sites"]), docString="latest", effectiveDateTime=datetime(2026-01-01T00:00:00Z), folder="Views", lookback=6h, lookback_column="Timestamp", updateExtentsCreationTime=false) ['Latest Events'] on table ['RawEvents'] {
RawEvents
| summarize arg_max(Timestamp, *) by DeviceId
}`
	if got := BuildCreate(d).String(); got != want {
		t.Errorf("create\n got %s\nwant %s", got, want)
	}
	wantAlter := `.alter materialized-view with (autoUpdateSchema=true, dimensionTables=dynamic(["DimDevices", "Dim'Sites"]), docString="latest", folder="Views", lookback=6h, lookback_column="Timestamp") ['Latest Events'] on table ['RawEvents'] {
RawEvents
| summarize arg_max(Timestamp, *) by DeviceId
}`
	if got := BuildAlter(d).String(); got != wantAlter {
		t.Errorf("alter\n got %s\nwant %s", got, wantAlter)
	}
	minimal, _ := FromParams("MV", v1alpha1.MaterializedViewParameters{SourceMaterializedView: ptr("Src"), Query: "Src | count"})
	if got := BuildCreate(minimal).String(); got != ".create ifnotexists materialized-view ['MV'] on materialized-view ['Src'] {\nSrc | count\n}" {
		t.Errorf("minimal create: %s", got)
	}
	if BuildEnable("MV", true).String() != ".enable materialized-view ['MV']" || BuildEnable("MV", false).String() != ".disable materialized-view ['MV']" {
		t.Error("enable/disable")
	}
	if BuildDelete("MV").String() != ".drop materialized-view ['MV'] ifexists" || ShowOne("MV").String() != ".show materialized-view ['MV']" || ShowAll().String() != ".show materialized-views" {
		t.Error("simple commands")
	}
	id := "6bb29fd5-6f2e-4a83-9b2f-1e7f3c2a1b00"
	c, err := ShowOperation(id)
	if err != nil || c.String() != ".show operations "+id {
		t.Errorf("ShowOperation: %s %v", c, err)
	}
	h, err := ShowOperationHistoric(id)
	if err != nil || h.String() != `.show operations | where tostring(OperationId) == "`+id+`"` {
		t.Errorf("ShowOperationHistoric: %s %v", h, err)
	}
	if _, err := ShowOperation("x; .drop table T"); err == nil {
		t.Error("invalid operation id must be rejected")
	}
	if _, err := ShowOperationHistoric("nope"); err == nil {
		t.Error("invalid operation id must be rejected")
	}
}

func TestParse(t *testing.T) {
	res := kusto.NewResult(kusto.NewTable("Table_0", showCols,
		row("A", "Raw", "Raw | count", true, 6*timespan.Hour),
		row("B", "Raw", "Raw | take 1", false, nil),
		[]any{"", "", "", nil, nil, "", true, true, "", "", false, "", nil, ""}))
	m := ParseRows(res)
	if len(m) != 2 || m["A"].SourceTable != "Raw" || *m["A"].Lookback != 6*timespan.Hour || m["A"].Folder != "Views" || !m["A"].IsHealthy || *m["A"].AutoUpdateSchema {
		t.Errorf("ParseRows: %+v", m["A"])
	}
	if m["B"].IsEnabled || m["B"].Lookback != nil {
		t.Errorf("B: %+v", m["B"])
	}
	if o, ok := ParseOne(res, "B"); !ok || o.Name != "B" {
		t.Error("ParseOne by name")
	}
	if o, ok := ParseOne(kusto.NewResult(kusto.NewTable("Table_0", showCols, row("Z", "Raw", "q", true, nil))), "A"); !ok || o.Name != "Z" {
		t.Error("ParseOne single row fallback")
	}
	if _, ok := ParseOne(kusto.NewResult(), "A"); ok {
		t.Error("ParseOne empty")
	}
	// Minimal column set (tolerant parsing).
	sparse := kusto.NewResult(kusto.NewTable("Table_0", []string{"Name", "Query"}, []any{"S", "q"}))
	if o, ok := ParseOne(sparse, "S"); !ok || !o.IsEnabled || o.Lookback != nil {
		t.Errorf("sparse: %+v %v", o, ok)
	}
	if ParseRows(nil) == nil {
		t.Error("nil result")
	}
	obs := Observation(m["A"])
	if obs.Lookback != "06:00:00" || !*obs.IsEnabled || obs.LastRunResult != "Success" {
		t.Errorf("Observation: %+v", obs)
	}

	opCols := []string{"OperationId", "Operation", "State", "Status"}
	op, ok := ParseOperation(kusto.NewResult(kusto.NewTable("Table_0", opCols, []any{"6bb29fd5-6f2e-4a83-9b2f-1e7f3c2a1b00", "MaterializedViewCreateOrAlter", "InProgress", ""})))
	if !ok || !op.Running() || op.ID == "" {
		t.Errorf("ParseOperation: %+v %v", op, ok)
	}
	op, _ = ParseOperation(kusto.NewResult(kusto.NewTable("Table_0", opCols, []any{"id", "x", "Failed", "boom"})))
	if op.Running() || op.Status != "boom" {
		t.Errorf("failed op: %+v", op)
	}
	if _, ok := ParseOperation(kusto.NewResult()); ok {
		t.Error("empty operation result")
	}
	if ParseOperationID(kusto.NewResult(kusto.NewTable("Table_0", []string{"OperationId"}, []any{" abc "}))) != "abc" || ParseOperationID(kusto.NewResult()) != "" {
		t.Error("ParseOperationID")
	}
}

func TestDiff(t *testing.T) {
	base := Observed{Name: "MV", SourceTable: "Raw", Query: "Raw | count", IsEnabled: true, Folder: "Views", Docstring: "doc", Lookback: ptr(6 * timespan.Hour), LookbackColumn: "Ts", AutoUpdateSchema: ptr(false)}
	cases := map[string]struct {
		d          v1alpha1.MaterializedViewParameters
		o          Observed
		queryEqual bool
		want       string
	}{
		"upToDate":         {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "Raw | count"}, base, true, ""},
		"queryDrift":       {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "Raw | count"}, base, false, "Alter"},
		"lookbackEqual":    {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", Lookback: ptr(common.Timespan("06:00:00"))}, base, true, ""},
		"lookbackDiffers":  {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", Lookback: ptr(common.Timespan("1d"))}, base, true, "Alter"},
		"lookbackMissing":  {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", Lookback: ptr(common.Timespan("1d"))}, Observed{IsEnabled: true}, true, "Alter"},
		"folderDiffers":    {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", Folder: ptr("Other")}, base, true, "Alter"},
		"docstringEqual":   {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", Docstring: ptr(" doc ")}, base, true, ""},
		"autoUpdateDiff":   {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", AutoUpdateSchema: ptr(true)}, base, true, "Alter"},
		"lookbackColumn":   {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", LookbackColumn: ptr("Other")}, base, true, "Alter"},
		"disable":          {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", Enabled: ptr(false)}, base, true, "Disable"},
		"enable":           {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q"}, Observed{IsEnabled: false}, true, "Enable"},
		"alterAndDisable":  {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", Enabled: ptr(false)}, base, false, "Alter,Disable"},
		"writeOnlyIgnored": {v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "q", DimensionTables: []string{"D"}, Backfill: ptr(true), Concurrency: ptr(int64(3))}, base, true, ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			d, err := FromParams("MV", tc.d)
			if err != nil {
				t.Fatal(err)
			}
			p := Diff(d, tc.o, tc.queryEqual)
			if p.String() != tc.want {
				t.Errorf("plan = %q, want %q", p.String(), tc.want)
			}
			if p.Empty() != (tc.want == "") {
				t.Error("Empty mismatch")
			}
		})
	}
}

func TestTexts(t *testing.T) {
	d, _ := FromParams("MV", v1alpha1.MaterializedViewParameters{SourceTable: ptr("Raw"), Query: "Raw\r\n| count  \r\n"})
	if DesiredTexts(d)[0] != ObservedTexts(Observed{Query: "Raw\n| count"})[0] {
		t.Error("stage 1 normalization should equalize CRLF and trailing whitespace")
	}
}
