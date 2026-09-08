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
	"context"
	"strings"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
)

var showCols = []string{"Name", "ExternalTableName", "Query", "ForcedLatency", "IntervalBetweenRuns", "CursorScopedTables", "ExportProperties", "LastRunTime", "IsDisabled", "LastRunResult", "ExportedTo", "IsRunning"}

func ptr[T any](v T) *T { return &v }

func export(name string) *v1alpha1.ContinuousExport {
	cr := &v1alpha1.ContinuousExport{ObjectMeta: metav1.ObjectMeta{Name: "ce", Namespace: "ns"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider = v1alpha1.ContinuousExportParameters{
		Database: "DB", ExternalTable: "Exports", OverTables: []string{"Raw"},
		Query: "Raw\n| where Level == \"Error\"", IntervalBetweenRuns: common.Timespan("1h"), ForcedLatency: ptr(common.Timespan("10m")),
	}
	return cr
}

// cluster stores exports; it echoes the query on one line (reformatted) and
// tracks the enabled flag.
type cluster struct {
	rows map[string][]any
}

func (c *cluster) handle(_ string, command cmd.Command) (*kusto.Result, error) {
	text := command.String()
	name := func() string {
		i := strings.Index(text, "['") + 2
		return text[i : i+strings.Index(text[i:], "']")]
	}
	switch {
	case strings.HasPrefix(text, ".show continuous-exports"):
		rows := make([][]any, 0, len(c.rows))
		for _, r := range c.rows {
			rows = append(rows, r)
		}
		return kusto.NewResult(kusto.NewTable("Table_0", showCols, rows...)), nil
	case strings.HasPrefix(text, ".show continuous-export "):
		if r, ok := c.rows[name()]; ok {
			return kusto.NewResult(kusto.NewTable("Table_0", showCols, r)), nil
		}
		return kusto.NewResult(kusto.NewTable("Table_0", showCols)), nil
	case strings.HasPrefix(text, ".create-or-alter continuous-export "):
		q := strings.TrimSpace(text[strings.Index(text, "<|")+2:])
		q = strings.Join(strings.Fields(q), " ")
		c.rows[name()] = []any{name(), "Exports", q, "00:10:00", "01:00:00", `["['DB'].['Raw']"]`, `{"SizeLimit":104857600}`, "", false, "", "", false}
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".disable continuous-export "), strings.HasPrefix(text, ".enable continuous-export "):
		if r, ok := c.rows[name()]; ok {
			r[8] = strings.HasPrefix(text, ".disable")
		}
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".drop continuous-export "):
		delete(c.rows, name())
		return kusto.NewResult(), nil
	}
	return nil, fake.HTTPError(400, `{"error":{"code":"General_BadRequest","@message":"unexpected: `+text+`"}}`)
}

func TestLifecycle(t *testing.T) {
	cl := &cluster{rows: map[string][]any{}}
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := &external{kc: kc, cache: snapshot.New(time.Minute, false)}
	cr := export("CE")

	obs, err := e.Observe(context.Background(), cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("initial: %+v %v", obs, err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	create := kc.Commands()[len(kc.Commands())-2]
	want := ".create-or-alter continuous-export ['CE'] over (['Raw']) to table ['Exports'] with (forcedLatency=10m, intervalBetweenRuns=1h)\n<| Raw\n| where Level == \"Error\""
	if create != want {
		t.Errorf("create:\n got %q\nwant %q", create, want)
	}
	for i := 0; i < 3; i++ {
		obs, err = e.Observe(context.Background(), cr)
		if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
			t.Fatalf("observe %d (reformatted query must be tolerated): %+v %v", i, obs, err)
		}
	}
	if cr.Status.AtProvider.ExternalTable != "Exports" || cr.Status.AtProvider.IsDisabled {
		t.Errorf("status: %+v", cr.Status.AtProvider)
	}

	// Disable through the spec -> toggle only, no alter.
	kc.Reset()
	cr.Spec.ForProvider.Enabled = ptr(false)
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || !strings.Contains(obs.Diff, "disable") {
		t.Fatalf("disable must be detected: %+v", obs)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count(".disable continuous-export ['CE']") != 1 || kc.Count(".create-or-alter") != 0 {
		t.Errorf("update commands: %v", kc.Commands())
	}
	obs, _ = e.Observe(context.Background(), cr)
	if !obs.ResourceUpToDate {
		t.Fatalf("after disable: %+v", obs)
	}

	// Interval change -> alter.
	kc.Reset()
	cr.Spec.ForProvider.IntervalBetweenRuns = "2h"
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || !strings.Contains(obs.Diff, "intervalBetweenRuns") {
		t.Fatalf("interval change must be detected: %+v", obs)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count(".create-or-alter") != 1 {
		t.Errorf("alter expected: %v", kc.Commands())
	}

	// Query drift in the cluster.
	cl.rows["CE"][2] = "Raw | where Level == \"Warning\""
	e.cache.Invalidate("http://e", "DB")
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || !strings.Contains(obs.Diff, "query") {
		t.Fatalf("query drift must be detected: %+v", obs)
	}

	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if _, ok := cl.rows["CE"]; ok {
		t.Error("delete must drop the export")
	}
	if err := e.Disconnect(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestErrors(t *testing.T) {
	bad := export("CE")
	bad.Spec.ForProvider.IntervalBetweenRuns = "soon"
	if _, err := (&external{kc: fake.New("http://e"), cache: snapshot.New(time.Minute, false)}).Observe(context.Background(), bad); !kerrors.IsBlocked(err) {
		t.Errorf("bad timespan must be Blocked: %v", err)
	}
	unresolved := export("CE")
	unresolved.Spec.ForProvider.ExternalTable = ""
	if _, err := (&external{kc: fake.New("http://e"), cache: snapshot.New(time.Minute, false)}).Observe(context.Background(), unresolved); !kerrors.IsBlocked(err) {
		t.Errorf("unresolved external table must be Blocked: %v", err)
	}
	gone := fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"not found","@permanent":true}}`))
	e := &external{kc: gone, cache: snapshot.New(0, true)}
	if obs, err := e.Observe(context.Background(), export("CE")); err != nil || obs.ResourceExists {
		t.Errorf("not found: %+v %v", obs, err)
	}
	if _, err := e.Delete(context.Background(), export("CE")); err != nil {
		t.Errorf("delete not found must be ok: %v", err)
	}
	// Create with enabled=false disables right after creating.
	cl := &cluster{rows: map[string][]any{}}
	kc := fake.New("http://e").OnFn("", cl.handle)
	cr := export("Off")
	cr.Spec.ForProvider.Enabled = ptr(false)
	if _, err := (&external{kc: kc, cache: snapshot.New(time.Minute, false)}).Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count(".disable continuous-export ['Off']") != 1 || cl.rows["Off"][8] != true {
		t.Errorf("create with enabled=false: %v", kc.Commands())
	}
}
