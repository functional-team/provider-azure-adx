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
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/function"
)

// TestFunctionNoUpdateLoop is the M1 acceptance criterion: a function with
// ugly formatting is written once and then observed as up to date although
// Kusto echoes a reformatted body.
func TestFunctionNoUpdateLoop(t *testing.T) {
	ctx := context.Background()
	run(t, ".create-merge table "+cmd.Ident("ItFnSource")+" (['A']:string, ['B']:long)")
	cache := snapshot.New(time.Minute, false)
	e := function.NewExternal(kc, cache, nil)
	cr := &v1alpha1.Function{ObjectMeta: metav1.ObjectMeta{Name: "f", Namespace: "it"}}
	meta.SetExternalName(cr, "ItUgly")
	cr.Spec.ForProvider = v1alpha1.FunctionParameters{
		Database:   db,
		Parameters: "(limit:long   =  10,   prefix : string = \"x\")",
		Body:       "ItFnSource\r\n|  where   A startswith prefix   \r\n| take limit\r\n",
		Folder:     strPtr("Parsing"),
	}

	obs, err := e.Observe(ctx, cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("initial: %+v %v", obs, err)
	}
	if _, err := e.Create(ctx, cr); err != nil {
		t.Fatal(err)
	}
	kc.reset()
	for i := 0; i < 20; i++ {
		cache.Invalidate(kc.Endpoint(), db)
		obs, err = e.Observe(ctx, cr)
		if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
			t.Fatalf("observe %d: %+v %v (cluster body %q)", i, obs, err, cr.Status.AtProvider.Body)
		}
	}
	if n := kc.count(".create-or-alter"); n != 0 {
		t.Fatalf("no writes expected after the first one, got %d", n)
	}
	t.Logf("cluster echoes parameters %q body %q", cr.Status.AtProvider.Parameters, cr.Status.AtProvider.Body)

	// A real change is applied exactly once.
	cr.Spec.ForProvider.Body = "ItFnSource | take 1"
	obs, _ = e.Observe(ctx, cr)
	if obs.ResourceUpToDate {
		t.Fatal("body change must be detected")
	}
	if _, err := e.Update(ctx, cr); err != nil {
		t.Fatal(err)
	}
	kc.reset()
	for i := 0; i < 5; i++ {
		cache.Invalidate(kc.Endpoint(), db)
		if obs, _ = e.Observe(ctx, cr); !obs.ResourceUpToDate {
			t.Fatalf("after update observe %d: %+v", i, obs)
		}
	}
	if kc.count(".create-or-alter") != 0 {
		t.Fatal("update loop after a change")
	}

	// Drift in the cluster is detected.
	run(t, ".create-or-alter function "+cmd.Ident("ItUgly")+"(limit:long=10,prefix:string=\"x\") { ItFnSource | take 2 }")
	cache.Invalidate(kc.Endpoint(), db)
	if obs, _ = e.Observe(ctx, cr); obs.ResourceUpToDate {
		t.Fatal("cluster drift must be detected")
	}
	if _, err := e.Delete(ctx, cr); err != nil {
		t.Fatal(err)
	}
}
