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
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	mv "github.com/functional-team/provider-azure-adx/internal/adx/materializedview"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

func ptr[T any](v T) *T { return &v }

var (
	showCols = []string{"Name", "SourceTable", "Query", "IsHealthy", "IsEnabled", "Folder", "DocString", "AutoUpdateSchema", "Lookback", "LookbackColumn", "LastRunResult"}
	opCols   = []string{"OperationId", "Operation", "State", "Status"}
	opID     = "6bb29fd5-6f2e-4a83-9b2f-1e7f3c2a1b00"
)

func view(name, query string) *v1alpha1.MaterializedView {
	cr := &v1alpha1.MaterializedView{ObjectMeta: metav1.ObjectMeta{Name: "mv", Namespace: "ns"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider = v1alpha1.MaterializedViewParameters{Database: "DB", SourceTable: ptr("Raw"), Query: query}
	return cr
}

// cluster simulates the relevant Kusto behaviour: it stores views, echoes
// queries with collapsed whitespace, and tracks one async operation.
type cluster struct {
	views   map[string][]any
	opState string
	opCalls int
}

func newCluster() *cluster { return &cluster{views: map[string][]any{}} }

func nameOf(text string) string {
	start := strings.Index(text, "['") + 2
	end := strings.Index(text[start:], "']") + start
	return text[start:end]
}

func (c *cluster) handle(_ string, command cmd.Command) (*kusto.Result, error) {
	text := command.String()
	switch {
	case text == ".show materialized-views":
		rows := make([][]any, 0, len(c.views))
		for _, r := range c.views {
			rows = append(rows, r)
		}
		return kusto.NewResult(kusto.NewTable("Table_0", showCols, rows...)), nil
	case strings.HasPrefix(text, ".show materialized-view "):
		r, ok := c.views[nameOf(text)]
		if !ok {
			return nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"not found","@permanent":true}}`)
		}
		return kusto.NewResult(kusto.NewTable("Table_0", showCols, r)), nil
	case strings.HasPrefix(text, ".create ") || strings.HasPrefix(text, ".alter materialized-view"):
		name := nameOf(text)
		brace := strings.Index(text, "{")
		body := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(text[brace:], "{"), "}"))
		query := strings.Join(strings.Fields(body), " ")
		enabled := true
		if old, ok := c.views[name]; ok {
			enabled, _ = old[4].(bool)
		}
		c.views[name] = []any{name, "Raw", query, true, enabled, "", "", false, nil, "", "Success"}
		if strings.HasPrefix(text, ".create async") {
			c.opState = mv.StateInProgress
			return kusto.NewResult(kusto.NewTable("Table_0", []string{"OperationId"}, []any{opID})), nil
		}
		return kusto.NewResult(kusto.NewTable("Table_0", showCols, c.views[name])), nil
	case strings.HasPrefix(text, ".enable materialized-view"), strings.HasPrefix(text, ".disable materialized-view"):
		r := c.views[nameOf(text)]
		r[4] = strings.HasPrefix(text, ".enable")
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".drop materialized-view"):
		delete(c.views, nameOf(text))
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".show operations "+opID):
		c.opCalls++
		if c.opState == "" || c.opState == "historic" {
			// Unknown to the direct lookup (e.g. older than 6h).
			return kusto.NewResult(kusto.NewTable("Table_0", opCols)), nil
		}
		return kusto.NewResult(kusto.NewTable("Table_0", opCols, []any{opID, "MaterializedViewCreateOrAlter", c.opState, "status text"})), nil
	case strings.HasPrefix(text, ".show operations | where"):
		c.opCalls++
		if c.opState == "historic" {
			return kusto.NewResult(kusto.NewTable("Table_0", opCols, []any{opID, "MaterializedViewCreateOrAlter", mv.StateCompleted, ""})), nil
		}
		return kusto.NewResult(kusto.NewTable("Table_0", opCols)), nil
	}
	return nil, errors.New("unexpected " + text)
}

func newExternal(kc kusto.Client) *external {
	return &external{kc: kc, cache: snapshot.New(time.Minute, false)}
}

func TestLifecycleWithQueryTolerance(t *testing.T) {
	cl := newCluster()
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := newExternal(kc)
	cr := view("Latest", "Raw\n|  summarize   arg_max(Ts, *) by Id")

	obs, err := e.Observe(context.Background(), cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("initial observe: %+v %v", obs, err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if a, o := base.Hashes(cr); a == "" || o == "" {
		t.Fatal("hashes must be recorded after create")
	}
	if _, ok := cr.GetAnnotations()[base.AnnotationOperationID]; ok {
		t.Error("sync create must not record an operation id")
	}
	writes := kc.Count(".alter")
	for i := 0; i < 5; i++ {
		obs, err := e.Observe(context.Background(), cr)
		if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
			t.Fatalf("observe %d: %+v %v", i, obs, err)
		}
	}
	if kc.Count(".alter") != writes || cr.GetCondition(xpv2.TypeReady).Reason != xpv2.ReasonAvailable {
		t.Error("stable view must be Available without writes")
	}
	if cr.Status.AtProvider.Query == "" || cr.Status.AtProvider.SourceTable != "Raw" {
		t.Errorf("status: %+v", cr.Status.AtProvider)
	}

	// Spec change -> alter, then stable.
	cr.Spec.ForProvider.Query = "Raw | summarize arg_max(Ts, *) by Id, Region"
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || obs.Diff != "Alter" {
		t.Fatalf("spec change must be detected: %+v", obs)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if obs, _ = e.Observe(context.Background(), cr); !obs.ResourceUpToDate {
		t.Fatal("after update the view must be up to date")
	}

	// Cluster drift -> detected and repaired.
	cl.views["Latest"][2] = "Raw | summarize arg_max(Ts, *) by Other"
	e.cache.Invalidate("http://e", "DB")
	if obs, _ = e.Observe(context.Background(), cr); obs.ResourceUpToDate {
		t.Fatal("cluster drift must be detected")
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if obs, _ = e.Observe(context.Background(), cr); !obs.ResourceUpToDate {
		t.Fatal("drift repaired")
	}

	// Disable / enable are separate steps and do not alter the query.
	cr.Spec.ForProvider.Enabled = ptr(false)
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || obs.Diff != "Disable" {
		t.Fatalf("disable must be planned: %+v", obs)
	}
	alters := kc.Count(".alter")
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count(".disable materialized-view ['Latest']") != 1 || kc.Count(".alter") != alters {
		t.Errorf("disable step: %v", kc.Commands())
	}
	if obs, _ = e.Observe(context.Background(), cr); !obs.ResourceUpToDate || *cr.Status.AtProvider.IsEnabled {
		t.Fatalf("after disable: %+v %+v", obs, cr.Status.AtProvider)
	}

	// Delete.
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if cr.GetCondition(xpv2.TypeReady).Reason != xpv2.ReasonDeleting {
		t.Error("Delete must set Deleting")
	}
	if obs, _ = e.Observe(context.Background(), cr); obs.ResourceExists {
		t.Fatal("deleted view must not exist")
	}
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Errorf("deleting a missing view is fine: %v", err)
	}
}

func TestAsyncBackfill(t *testing.T) {
	cl := newCluster()
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := newExternal(kc)
	cr := view("Backfilled", "Raw | summarize count() by Id")
	cr.Spec.ForProvider.Backfill = ptr(true)

	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if got := cr.GetAnnotations()[base.AnnotationOperationID]; got != opID {
		t.Fatalf("operation id annotation = %q", got)
	}
	if kc.Count(".create async ifnotexists materialized-view with (backfill=true)") != 1 {
		t.Errorf("async create expected: %v", kc.Commands())
	}

	// While running: Creating, no writes, status carries the operation.
	for i := 0; i < 3; i++ {
		obs, err := e.Observe(context.Background(), cr)
		if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
			t.Fatalf("observe while running: %+v %v", obs, err)
		}
	}
	if cr.GetCondition(xpv2.TypeReady).Reason != xpv2.ReasonCreating || cr.Status.AtProvider.OperationState != mv.StateInProgress || cr.Status.AtProvider.OperationID != opID {
		t.Errorf("Creating expected: %+v", cr.Status)
	}
	if kc.Count(".show materialized-view") != 1 { // only the read-back after create
		t.Errorf("no view observe while the operation runs: %v", kc.Commands())
	}

	// Completed: annotation cleared, regular observe, Available.
	cl.opState = mv.StateCompleted
	obs, err := e.Observe(context.Background(), cr)
	if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("observe after completion: %+v %v", obs, err)
	}
	if cr.GetAnnotations()[base.AnnotationOperationID] != "" || cr.GetCondition(xpv2.TypeReady).Reason != xpv2.ReasonAvailable {
		t.Errorf("operation must be forgotten and the view Available: %v %v", cr.GetAnnotations(), cr.GetCondition(xpv2.TypeReady))
	}
	if cr.Status.AtProvider.OperationID != "" {
		t.Error("status must not keep the finished operation")
	}
	calls := cl.opCalls
	if _, err := e.Observe(context.Background(), cr); err != nil || cl.opCalls != calls {
		t.Error("no more operation lookups after completion")
	}
}

func TestAsyncFailureAndHistoricLookup(t *testing.T) {
	cl := newCluster()
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := newExternal(kc)
	cr := view("Failing", "Raw | count")
	meta.AddAnnotations(cr, map[string]string{base.AnnotationOperationID: opID})
	cl.views["Failing"] = []any{"Failing", "Raw", "Raw | count", true, true, "", "", false, nil, "", ""}

	cl.opState = mv.StateFailed
	_, err := e.Observe(context.Background(), cr)
	if err == nil || !strings.Contains(err.Error(), "Failed") || !strings.Contains(err.Error(), "status text") {
		t.Fatalf("failed operation must surface as error: %v", err)
	}
	if cr.GetAnnotations()[base.AnnotationOperationID] != "" {
		t.Error("failed operation must be forgotten so the next reconcile retries the create")
	}

	// Historic lookup: direct lookup returns nothing, the log form does.
	meta.AddAnnotations(cr, map[string]string{base.AnnotationOperationID: opID})
	cl.opState = "historic"
	obs, err := e.Observe(context.Background(), cr)
	if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("historic completed: %+v %v", obs, err)
	}
	if kc.Count(".show operations | where") != 1 {
		t.Errorf("historic form expected once: %v", kc.Commands())
	}

	// Unknown everywhere: forget it and observe normally.
	meta.AddAnnotations(cr, map[string]string{base.AnnotationOperationID: opID})
	cl.opState = ""
	if obs, err := e.Observe(context.Background(), cr); err != nil || !obs.ResourceExists {
		t.Fatalf("unknown operation: %+v %v", obs, err)
	}

	// Invalid annotation content is rejected before any command is sent.
	meta.AddAnnotations(cr, map[string]string{base.AnnotationOperationID: "x; .drop database"})
	if _, err := e.Observe(context.Background(), cr); err == nil {
		t.Error("invalid operation id must fail")
	}
}

func TestErrorsAndGuardrails(t *testing.T) {
	// Invalid spec is Blocked.
	bad := view("MV", "q")
	bad.Spec.ForProvider.Lookback = ptr(common.Timespan("soon"))
	if _, err := newExternal(fake.New("http://e")).Observe(context.Background(), bad); !kerrors.IsBlocked(err) {
		t.Errorf("bad lookback must be Blocked: %v", err)
	}
	// Alter rejected by Kusto carries the reason.
	existing := kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{"MV", "Raw", "Raw | summarize count() by A", true, true, "", "", false, nil, "", ""}))
	kc := fake.New("http://e").
		On(".show materialized-views", existing, nil).
		On(".alter materialized-view", nil, fake.HTTPError(400, `{"error":{"code":"General_BadRequest","@message":"group by change not allowed","@permanent":true}}`))
	cr := view("MV", "Raw | summarize count() by B")
	_, err := newExternal(kc).Update(context.Background(), cr)
	if err == nil || !strings.Contains(err.Error(), kerrors.ReasonMaterializedViewAlterRejected) {
		t.Errorf("rejected alter must carry the reason: %v", err)
	}
	// Transient alter failure keeps the step kind, no reason.
	kc2 := fake.New("http://e").On(".show materialized-views", existing, nil).On(".alter materialized-view", nil, fake.HTTPError(500, `{"error":{"@message":"boom"}}`))
	if _, err := newExternal(kc2).Update(context.Background(), cr); err == nil || strings.Contains(err.Error(), kerrors.ReasonMaterializedViewAlterRejected) {
		t.Errorf("transient alter failure: %v", err)
	}
	// Observe errors and not found.
	broken := fake.New("http://e").On("", nil, fake.HTTPError(500, `{"error":{"@message":"boom"}}`))
	if _, err := newExternal(broken).Observe(context.Background(), cr); err == nil {
		t.Error("observe error expected")
	}
	gone := fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"nope","@permanent":true}}`))
	if obs, err := newExternal(gone).Observe(context.Background(), cr); err != nil || obs.ResourceExists {
		t.Errorf("not found: %+v %v", obs, err)
	}
	if _, err := newExternal(gone).Delete(context.Background(), cr); err != nil {
		t.Errorf("delete not found is fine: %v", err)
	}
	// AlreadyExists on create is tolerated.
	exists := fake.New("http://e").On(".create", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityAlreadyExists","@message":"exists","@permanent":true}}`)).On(".show materialized-view", existing, nil)
	if _, err := newExternal(exists).Create(context.Background(), cr); err != nil {
		t.Errorf("already exists must be tolerated: %v", err)
	}
	// Cache disabled uses the single show.
	single := fake.New("http://e").On(".show materialized-view ['MV']", existing, nil)
	e := &external{kc: single, cache: snapshot.New(0, true)}
	up := view("MV", "Raw | summarize count() by A")
	if obs, err := e.Observe(context.Background(), up); err != nil || !obs.ResourceUpToDate {
		t.Errorf("single show: %+v %v", obs, err)
	}
	if err := e.Disconnect(context.Background()); err != nil {
		t.Error(err)
	}
}
