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
	"fmt"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/policy/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/policy"
)

func retention(entity common.EntityKind, name, period string) *v1alpha1.RetentionPolicy {
	cr := &v1alpha1.RetentionPolicy{ObjectMeta: metav1.ObjectMeta{Name: "r-" + name, Namespace: "it"}}
	target := v1alpha1.PolicyTarget{Database: db, Entity: common.EntityReference{Kind: entity}}
	if name != "" {
		target.Entity.Name = &name
	}
	ts := common.Timespan(period)
	rec := "Enabled"
	cr.Spec.ForProvider = v1alpha1.RetentionPolicyParameters{PolicyTarget: target, SoftDeletePeriod: &ts, Recoverability: &rec}
	return cr
}

func TestRetentionPolicyLifecycle(t *testing.T) {
	ctx := context.Background()
	run(t, ".create-merge table "+cmd.Ident("ItRet")+" (['A']:string)")
	cache := snapshot.New(time.Minute, false)
	e := policy.NewExternal(kc, cache, nil, policy.Retention())
	cr := retention(common.EntityKindTable, "ItRet", "365d")

	obs, err := e.Observe(ctx, cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("a fresh table inherits, it has no own policy: %+v %v", obs, err)
	}
	if _, err := e.Create(ctx, cr); err != nil {
		t.Fatal(err)
	}
	obs, err = e.Observe(ctx, cr)
	if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("after create: %+v %v (policy %s)", obs, err, cr.Status.AtProvider.Policy)
	}
	cr.Spec.ForProvider.SoftDeletePeriod = tsPtr("30d")
	if obs, _ = e.Observe(ctx, cr); obs.ResourceUpToDate {
		t.Fatal("changed period must be detected")
	}
	if _, err := e.Update(ctx, cr); err != nil {
		t.Fatal(err)
	}
	if obs, _ = e.Observe(ctx, cr); !obs.ResourceUpToDate {
		t.Fatalf("after update: %+v %s", obs, cr.Status.AtProvider.Policy)
	}
	if _, err := e.Delete(ctx, cr); err != nil {
		t.Fatal(err)
	}
	if obs, _ = e.Observe(ctx, cr); obs.ResourceExists {
		t.Fatal("after delete the table inherits again")
	}

	// Database level.
	dbcr := retention(common.EntityKindDatabase, "", "3650d")
	if _, err := e.Create(ctx, dbcr); err != nil {
		t.Fatal(err)
	}
	if obs, err := e.Observe(ctx, dbcr); err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("database policy: %+v %v", obs, err)
	}
	if _, err := e.Delete(ctx, dbcr); err != nil {
		t.Fatal(err)
	}
}

// TestPolicyBatchObserve answers spike S3 for the retention policy and checks
// the M1 scale criterion: 50 table policies in one database cost one batch
// command per cache TTL (or fall back to single shows when the wildcard form
// is rejected).
func TestPolicyBatchObserve(t *testing.T) {
	ctx := context.Background()
	const n = 50
	for i := 0; i < n; i++ {
		run(t, ".create-merge table "+cmd.Ident(fmt.Sprintf("ItBatch%02d", i))+" (['A']:string)")
	}
	cache := snapshot.New(time.Minute, false)
	e := policy.NewExternal(kc, cache, nil, policy.Retention())
	kc.reset()
	for i := 0; i < n; i++ {
		cr := retention(common.EntityKindTable, fmt.Sprintf("ItBatch%02d", i), "1d")
		if _, err := e.Observe(ctx, cr); err != nil {
			t.Fatalf("observe %d: %v", i, err)
		}
	}
	wild, single := kc.count("table * policy retention"), kc.count("] policy retention")
	switch {
	case wild == 1 && single == 0:
		t.Log("S3: '.show table * policy retention' is supported; one batch command for 50 resources")
	case wild == 1 && single == n:
		t.Log("S3: wildcard rejected, fell back to single .show per resource")
	default:
		t.Errorf("unexpected command pattern: wildcard=%d single=%d", wild, single)
	}
}

func TestUpdatePolicyLifecycle(t *testing.T) {
	ctx := context.Background()
	run(t, ".create-merge table "+cmd.Ident("ItUpdSrc")+" (['A']:string)")
	run(t, ".create-merge table "+cmd.Ident("ItUpdDst")+" (['A']:string)")
	cache := snapshot.New(time.Minute, false)
	e := policy.NewExternal(kc, cache, nil, policy.Update())
	name := "ItUpdDst"
	cr := &v1alpha1.UpdatePolicy{ObjectMeta: metav1.ObjectMeta{Name: "u", Namespace: "it"}}
	cr.Spec.ForProvider = v1alpha1.UpdatePolicyParameters{
		PolicyTarget: v1alpha1.PolicyTarget{Database: db, Entity: common.EntityReference{Kind: common.EntityKindTable, Name: &name}},
		Updates:      []v1alpha1.UpdatePolicyEntry{{Source: "ItUpdSrc", Query: "ItUpdSrc\n|   project   A"}},
	}
	if _, err := e.Create(ctx, cr); err != nil {
		t.Fatal(err)
	}
	kc.reset()
	for i := 0; i < 5; i++ {
		cache.Invalidate(kc.Endpoint(), db)
		obs, err := e.Observe(ctx, cr)
		if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
			t.Fatalf("observe %d: %+v %v (%s)", i, obs, err, cr.Status.AtProvider.Policy)
		}
	}
	if kc.count(".alter table") != 0 {
		t.Fatal("update policy must not be rewritten")
	}
	if _, err := e.Delete(ctx, cr); err != nil {
		t.Fatal(err)
	}
}

func tsPtr(s string) *common.Timespan { t := common.Timespan(s); return &t }
