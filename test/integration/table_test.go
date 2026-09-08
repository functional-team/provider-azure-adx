//go:build integration

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

package integration

import (
	"context"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/table"
)

func newTable(name string, mode v1alpha1.SchemaUpdateMode, cols ...common.Column) *v1alpha1.Table {
	cr := &v1alpha1.Table{ObjectMeta: metav1.ObjectMeta{Name: "t", Namespace: "it"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider = v1alpha1.TableParameters{Database: db, Columns: cols, SchemaUpdateMode: mode}
	return cr
}

func TestTableLifecycle(t *testing.T) {
	ctx := context.Background()
	cache := snapshot.New(time.Minute, false)
	e := table.NewExternal(kc, cache, nil)
	cr := newTable("ItEvents", v1alpha1.SchemaUpdateModeMerge, common.Column{Name: "Timestamp", Type: "datetime"}, common.Column{Name: "Payload", Type: "dynamic"})
	cr.Spec.ForProvider.Folder = strPtr("Raw")

	obs, err := e.Observe(ctx, cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("initial observe: %+v %v", obs, err)
	}
	if _, err := e.Create(ctx, cr); err != nil {
		t.Fatal(err)
	}
	obs, err = e.Observe(ctx, cr)
	if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("after create: %+v %v (status %+v)", obs, err, cr.Status.AtProvider)
	}
	if len(cr.Status.AtProvider.Columns) != 2 || cr.Status.AtProvider.Folder != "Raw" {
		t.Errorf("status: %+v (spike S4: folder/columns must come through the schema JSON)", cr.Status.AtProvider)
	}

	// Drift: a column added in the cluster is reported, not removed (Merge).
	run(t, ".alter-merge table "+cmd.Ident("ItEvents")+" (['Extra']:string)")
	cache.Invalidate(kc.Endpoint(), db)
	obs, err = e.Observe(ctx, cr)
	if err != nil || !obs.ResourceUpToDate || len(cr.Status.AtProvider.DriftColumns) != 1 {
		t.Fatalf("merge drift: %+v %v drift=%v", obs, err, cr.Status.AtProvider.DriftColumns)
	}

	// Spec adds a column -> update.
	cr.Spec.ForProvider.Columns = append(cr.Spec.ForProvider.Columns, common.Column{Name: "Level", Type: "string"})
	obs, _ = e.Observe(ctx, cr)
	if obs.ResourceUpToDate {
		t.Fatal("added column must need an update")
	}
	if _, err := e.Update(ctx, cr); err != nil {
		t.Fatal(err)
	}
	obs, _ = e.Observe(ctx, cr)
	if !obs.ResourceUpToDate {
		t.Fatalf("after update: %+v", obs)
	}

	// Type change is blocked and sends no command.
	kc.reset()
	cr.Spec.ForProvider.Columns[2].Type = "long"
	if _, err := e.Observe(ctx, cr); !kerrors.IsBlocked(err) {
		t.Fatalf("type change must be blocked: %v", err)
	}
	if kc.count(".alter") != 0 {
		t.Error("blocked observe must not send commands")
	}
	cr.Spec.ForProvider.Columns[2].Type = "string"

	// Replace drops the extra column.
	cr.Spec.ForProvider.SchemaUpdateMode = v1alpha1.SchemaUpdateModeReplace
	obs, _ = e.Observe(ctx, cr)
	if obs.ResourceUpToDate {
		t.Fatal("replace with extra column must need an update")
	}
	if _, err := e.Update(ctx, cr); err != nil {
		t.Fatal(err)
	}
	obs, _ = e.Observe(ctx, cr)
	if !obs.ResourceUpToDate || len(cr.Status.AtProvider.Columns) != 3 {
		t.Fatalf("after replace: %+v %+v", obs, cr.Status.AtProvider)
	}

	if _, err := e.Delete(ctx, cr); err != nil {
		t.Fatal(err)
	}
	obs, err = e.Observe(ctx, cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("after delete: %+v %v", obs, err)
	}
}

func strPtr(s string) *string { return &s }
