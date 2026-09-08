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

package continuousexport

import (
	"strings"
	"testing"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

func ptr[T any](v T) *T { return &v }

var showCols = []string{"Name", "ExternalTableName", "Query", "ForcedLatency", "IntervalBetweenRuns", "CursorScopedTables", "ExportProperties", "LastRunTime", "IsDisabled", "LastRunResult", "ExportedTo", "IsRunning"}

func params() v1alpha1.ContinuousExportParameters {
	return v1alpha1.ContinuousExportParameters{
		Database: "DB", ExternalTable: "Exports", OverTables: []string{"RawEvents"},
		Query:               "RawEvents\n| where Level == \"Error\"",
		IntervalBetweenRuns: "1h", ForcedLatency: ptr(common.Timespan("10m")),
		SizeLimit: ptr(int64(104857600)), Distributed: ptr(true), Distribution: ptr("per_node"), ParquetRowGroupSize: ptr(int64(100000)), ManagedIdentity: ptr("system"),
	}
}

func TestFromParamsAndBuild(t *testing.T) {
	d, err := FromParams("CE", params())
	if err != nil {
		t.Fatal(err)
	}
	if d.IntervalBetweenRuns != timespan.Hour || *d.ForcedLatency != 10*timespan.Minute || d.OverTables[0] != "['DB'].['RawEvents']" || !d.Enabled {
		t.Errorf("FromParams: %+v", d)
	}
	want := ".create-or-alter continuous-export ['CE'] over (['RawEvents']) to table ['Exports'] with (distributed=true, distribution=\"per_node\", forcedLatency=10m, intervalBetweenRuns=1h, managedIdentity=\"system\", parquetRowGroupSize=100000, sizeLimit=104857600)\n<| RawEvents\n| where Level == \"Error\""
	if got := BuildCreateOrAlter(d).String(); got != want {
		t.Errorf("\n got %s\nwant %s", got, want)
	}
	minimal, _ := FromParams("M", v1alpha1.ContinuousExportParameters{Database: "DB", ExternalTable: "E", Query: "T", IntervalBetweenRuns: "30m", Enabled: ptr(false)})
	if got := BuildCreateOrAlter(minimal).String(); got != ".create-or-alter continuous-export ['M'] to table ['E'] with (intervalBetweenRuns=30m)\n<| T" || minimal.Enabled {
		t.Errorf("minimal: %s", got)
	}
	if BuildToggle("CE", true).String() != ".enable continuous-export ['CE']" || BuildToggle("CE", false).String() != ".disable continuous-export ['CE']" || BuildDelete("CE").String() != ".drop continuous-export ['CE']" {
		t.Error("toggle/delete")
	}
	if Show("CE").String() != ".show continuous-export ['CE']" || ShowAll().String() != ".show continuous-exports" {
		t.Error("show")
	}
	if Qualify("DB", "['Other'].['T']") != "['Other'].['T']" || Qualify("DB", "['T']") != "['DB'].['T']" || Qualify("DB", "T") != "['DB'].['T']" {
		t.Error("Qualify")
	}
	for name, p := range map[string]v1alpha1.ContinuousExportParameters{
		"noTarget":         {Database: "DB", Query: "T", IntervalBetweenRuns: "1h"},
		"badInterval":      {Database: "DB", ExternalTable: "E", Query: "T", IntervalBetweenRuns: "soon"},
		"badForcedLatency": {Database: "DB", ExternalTable: "E", Query: "T", IntervalBetweenRuns: "1h", ForcedLatency: ptr(common.Timespan("x"))},
	} {
		if _, err := FromParams("CE", p); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestParse(t *testing.T) {
	res := kusto.NewResult(kusto.NewTable("Table_0", showCols,
		[]any{"CE", "Exports", "RawEvents | where Level == \"Error\"", "00:10:00", timespan.Hour, `["['DB'].['RawEvents']"]`, `{"SizeLimit":104857600,"Distributed":true,"Distribution":"per_node","ManagedIdentity":"system"}`, "2026-09-07T10:00:00Z", false, "Completed", "2026-09-07T09:00:00Z", true},
		[]any{"", "", "", nil, nil, nil, nil, "", false, "", "", false}))
	m := ParseRows(res)
	if len(m) != 1 {
		t.Fatalf("ParseRows: %d", len(m))
	}
	o := m["CE"]
	if o.ExternalTable != "Exports" || *o.ForcedLatency != 10*timespan.Minute || *o.IntervalBetweenRuns != timespan.Hour || len(o.CursorScopedTables) != 1 || o.CursorScopedTables[0] != "['DB'].['RawEvents']" || o.IsDisabled || !o.IsRunning || o.LastRunResult != "Completed" {
		t.Errorf("Observed: %+v", o)
	}
	if v, ok := o.Properties["SizeLimit"]; !ok || v.(float64) != 104857600 {
		t.Error("properties")
	}
	if one, ok := ParseOne(res, "CE"); !ok || one.Name != "CE" {
		t.Error("ParseOne")
	}
	if _, ok := ParseOne(kusto.NewResult(), "CE"); ok {
		t.Error("ParseOne empty")
	}
	obs := Observation(o)
	if obs.ExternalTable != "Exports" || obs.LastRunResult != "Completed" || !obs.IsRunning || obs.ExportedTo == "" {
		t.Errorf("Observation: %+v", obs)
	}
}

func TestDiff(t *testing.T) {
	base := Observed{Name: "CE", ExternalTable: "Exports", Query: "RawEvents | where Level == \"Error\"", ForcedLatency: ptr(10 * timespan.Minute), IntervalBetweenRuns: ptr(timespan.Hour),
		CursorScopedTables: []string{"['DB'].['RawEvents']"}, Properties: map[string]any{"sizeLimit": float64(104857600), "Distributed": true, "Distribution": "per_node", "ParquetRowGroupSize": float64(100000), "ManagedIdentity": "system"}}
	cases := map[string]struct {
		mutate func(*v1alpha1.ContinuousExportParameters, *Observed)
		alter  bool
		toggle *bool
		reason string
	}{
		"equal":         {func(*v1alpha1.ContinuousExportParameters, *Observed) {}, false, nil, ""},
		"externalTable": {func(_ *v1alpha1.ContinuousExportParameters, o *Observed) { o.ExternalTable = "Other" }, true, nil, "externalTable"},
		"interval":      {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.IntervalBetweenRuns = "2h" }, true, nil, "intervalBetweenRuns"},
		"intervalAlias": {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.IntervalBetweenRuns = "01:00:00" }, false, nil, ""},
		"forcedLatency": {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) {
			p.ForcedLatency = ptr(common.Timespan("5m"))
		}, true, nil, "forcedLatency"},
		"forcedLatencyOff": {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.ForcedLatency = nil }, false, nil, ""},
		"overTables":       {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.OverTables = []string{"RawEvents", "Dim"} }, true, nil, "overTables"},
		"overTablesUnset":  {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.OverTables = nil }, false, nil, ""},
		"sizeLimit":        {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.SizeLimit = ptr(int64(1)) }, true, nil, "sizeLimit"},
		"distributed":      {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.Distributed = ptr(false) }, true, nil, "distributed"},
		"distribution":     {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.Distribution = ptr("single") }, true, nil, "distribution"},
		"propMissingOK":    {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.DistributionKind = ptr("default") }, false, nil, ""},
		"managedIdentity":  {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.ManagedIdentity = ptr("0000") }, true, nil, "managedIdentity"},
		"disable":          {func(p *v1alpha1.ContinuousExportParameters, _ *Observed) { p.Enabled = ptr(false) }, false, ptr(false), "disable"},
		"enable":           {func(_ *v1alpha1.ContinuousExportParameters, o *Observed) { o.IsDisabled = true }, false, ptr(true), "enable"},
		"both": {func(p *v1alpha1.ContinuousExportParameters, o *Observed) {
			o.IsDisabled = true
			p.SizeLimit = ptr(int64(2))
		}, true, ptr(true), "sizeLimit"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			p := params()
			o := base
			tc.mutate(&p, &o)
			d, err := FromParams("CE", p)
			if err != nil {
				t.Fatal(err)
			}
			plan := Diff(d, o)
			if plan.Alter != tc.alter {
				t.Errorf("Alter = %v (%s), want %v", plan.Alter, plan.String(), tc.alter)
			}
			if (plan.Toggle == nil) != (tc.toggle == nil) || (tc.toggle != nil && *plan.Toggle != *tc.toggle) {
				t.Errorf("Toggle = %v, want %v", plan.Toggle, tc.toggle)
			}
			if plan.Empty() != (!tc.alter && tc.toggle == nil) {
				t.Error("Empty mismatch")
			}
			if tc.reason != "" && !strings.Contains(plan.String(), tc.reason) {
				t.Errorf("reason %q missing in %q", tc.reason, plan.String())
			}
		})
	}
	d, _ := FromParams("CE", params())
	if DesiredTexts(d)[0] != "RawEvents\n| where Level == \"Error\"" || ObservedTexts(Observed{Query: "T \r\n"})[0] != "T" {
		t.Error("texts")
	}
}
