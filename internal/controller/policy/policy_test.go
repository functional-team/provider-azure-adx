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

package policy

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/policy/v1alpha1"
	adxpolicy "github.com/functional-team/provider-azure-adx/internal/adx/policy"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

var showCols = []string{"PolicyName", "EntityName", "Policy", "ChildEntities", "EntityType"}

func ptr[T any](v T) *T { return &v }

func target(kind common.EntityKind, name string) v1alpha1.PolicyTarget {
	t := v1alpha1.PolicyTarget{Database: "DB", Entity: common.EntityReference{Kind: kind}}
	if name != "" {
		t.Entity.Name = &name
	}
	return t
}

func ext[T Policy](kc kusto.Client, def Def[T], disabled bool) *external[T] {
	return &external[T]{kc: kc, cache: snapshot.New(time.Minute, disabled), def: def, st: &state{}}
}

func retention(entity common.EntityKind, name, period string) *v1alpha1.RetentionPolicy {
	cr := &v1alpha1.RetentionPolicy{ObjectMeta: metav1.ObjectMeta{Name: "r", Namespace: "ns"}}
	cr.Spec.ForProvider = v1alpha1.RetentionPolicyParameters{PolicyTarget: target(entity, name)}
	if period != "" {
		cr.Spec.ForProvider.SoftDeletePeriod = ptr(common.Timespan(period))
	}
	return cr
}

func TestRetentionObserve(t *testing.T) {
	batch := kusto.NewResult(kusto.NewTable("Table_0", showCols,
		[]any{"RetentionPolicy", "[DB].[T]", `{"SoftDeletePeriod":"365.00:00:00","Recoverability":"Enabled"}`, nil, "Table"},
		[]any{"RetentionPolicy", "[DB].[Inherits]", nil, nil, "Table"}))
	cases := map[string]struct {
		kc        *fake.Client
		cr        *v1alpha1.RetentionPolicy
		exists    bool
		upToDate  bool
		errHas    string
		blocked   string
		cmdsCount int
	}{
		"batchUpToDate":  {kc: fake.New("http://e").On("table * policy retention", batch, nil), cr: retention(common.EntityKindTable, "T", "365d"), exists: true, upToDate: true, cmdsCount: 1},
		"batchDiffers":   {kc: fake.New("http://e").On("table * policy retention", batch, nil), cr: retention(common.EntityKindTable, "T", "30d"), exists: true, upToDate: false, cmdsCount: 1},
		"batchInherits":  {kc: fake.New("http://e").On("table * policy retention", batch, nil), cr: retention(common.EntityKindTable, "Inherits", "30d"), exists: false, cmdsCount: 1},
		"batchUnknown":   {kc: fake.New("http://e").On("table * policy retention", batch, nil), cr: retention(common.EntityKindTable, "Nope", "30d"), exists: false, cmdsCount: 1},
		"databaseSingle": {kc: fake.New("http://e").On(".show database ['DB'] policy retention", kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{"RetentionPolicy", "[DB]", `{"SoftDeletePeriod":"3650.00:00:00"}`, nil, "Database"})), nil), cr: retention(common.EntityKindDatabase, "", "3650d"), exists: true, upToDate: true, cmdsCount: 1},
		"mvSingleNull":   {kc: fake.New("http://e").On(".show materialized-view ['MV'] policy retention", kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{"RetentionPolicy", "[DB].[MV]", "null", nil, "MaterializedView"})), nil), cr: retention(common.EntityKindMaterializedView, "MV", "1d"), exists: false, cmdsCount: 1},
		"unresolvedName": {kc: fake.New("http://e"), cr: retention(common.EntityKindTable, "", "1d"), errHas: "entity.name is empty"},
		"invalidSpec":    {kc: fake.New("http://e"), cr: retention(common.EntityKindTable, "T", "soon"), blocked: kerrors.ReasonInvalidSpec},
		"clusterError":   {kc: fake.New("http://e").On("", nil, errors.New("boom")), cr: retention(common.EntityKindTable, "T", "1d"), errHas: errObserve},
		"dbNotFound":     {kc: fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"db not found","@permanent":true}}`)), cr: retention(common.EntityKindTable, "T", "1d"), exists: false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := ext(tc.kc, Retention(), false)
			got, err := e.Observe(context.Background(), tc.cr)
			if tc.blocked != "" {
				var b *kerrors.Blocked
				if !errors.As(err, &b) || b.Reason != tc.blocked {
					t.Fatalf("expected Blocked(%s), got %v", tc.blocked, err)
				}
				return
			}
			if tc.errHas != "" {
				if err == nil || !strings.Contains(err.Error(), tc.errHas) {
					t.Fatalf("expected error containing %q, got %v", tc.errHas, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ResourceExists != tc.exists || got.ResourceUpToDate != tc.upToDate {
				t.Errorf("got %+v, want exists=%v upToDate=%v", got, tc.exists, tc.upToDate)
			}
			if tc.cmdsCount > 0 && len(tc.kc.Commands()) != tc.cmdsCount {
				t.Errorf("commands: %v", tc.kc.Commands())
			}
			if tc.exists && tc.cr.Status.AtProvider.Policy == "" {
				t.Error("status.atProvider.policy must be set")
			}
		})
	}
}

func TestBatchFallback(t *testing.T) {
	kc := fake.New("http://e").
		On("table * policy retention", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_SyntaxError","@message":"Syntax error","@permanent":true}}`)).
		On(".show table ['T'] policy retention", kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{"RetentionPolicy", "[DB].[T]", `{"SoftDeletePeriod":"1.00:00:00"}`, nil, "Table"})), nil)
	e := ext(kc, Retention(), false)
	for i := 0; i < 3; i++ {
		got, err := e.Observe(context.Background(), retention(common.EntityKindTable, "T", "1d"))
		if err != nil || !got.ResourceExists || !got.ResourceUpToDate {
			t.Fatalf("observe %d: %+v %v", i, got, err)
		}
	}
	if kc.Count("table * policy") != 1 {
		t.Errorf("wildcard must be tried once, then remembered as unsupported: %v", kc.Commands())
	}
	if kc.Count(".show table ['T'] policy retention") != 3 || !e.st.batchUnsupported.Load() {
		t.Errorf("single show fallback: %v", kc.Commands())
	}
}

func TestBatchSharedAcrossResources(t *testing.T) {
	rows := make([][]any, 0, 200)
	for i := 0; i < 200; i++ {
		rows = append(rows, []any{"RetentionPolicy", "[DB].[T" + strings.Repeat("x", i%3) + "]", `{"SoftDeletePeriod":"1.00:00:00"}`, nil, "Table"})
	}
	kc := fake.New("http://e").On("table * policy retention", kusto.NewResult(kusto.NewTable("Table_0", showCols, rows...)), nil)
	e := ext(kc, Retention(), false)
	for i := 0; i < 500; i++ {
		if _, err := e.Observe(context.Background(), retention(common.EntityKindTable, "T", "1d")); err != nil {
			t.Fatal(err)
		}
	}
	if n := kc.Count("table * policy"); n != 1 {
		t.Errorf("500 observes must cause exactly one batch load within the TTL, got %d", n)
	}
}

func TestRetentionCreateUpdateDelete(t *testing.T) {
	var stored string
	kc := fake.New("http://e").OnFn("", func(_ string, c cmd.Command) (*kusto.Result, error) {
		text := c.String()
		switch {
		case strings.HasPrefix(text, ".alter table ['T'] policy retention "):
			stored = strings.TrimSuffix(strings.TrimPrefix(text, ".alter table ['T'] policy retention @'"), "'")
			return kusto.NewResult(), nil
		case strings.HasPrefix(text, ".show table ['T'] policy retention"), strings.HasPrefix(text, ".show table * policy retention"):
			var pol any
			if stored != "" {
				pol = stored
			}
			return kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{"RetentionPolicy", "[DB].[T]", pol, nil, "Table"})), nil
		case strings.HasPrefix(text, ".delete table ['T'] policy retention"):
			stored = ""
			return kusto.NewResult(), nil
		}
		return nil, errors.New("unexpected " + text)
	})
	e := ext(kc, Retention(), false)
	cr := retention(common.EntityKindTable, "T", "365d")
	cr.Spec.ForProvider.Recoverability = ptr("Disabled")

	got, _ := e.Observe(context.Background(), cr)
	if got.ResourceExists {
		t.Fatal("must not exist yet")
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if stored != `{"SoftDeletePeriod":"365.00:00:00","Recoverability":"Disabled"}` {
		t.Errorf("stored policy: %s", stored)
	}
	if a, _ := base.Hashes(cr); a != "" {
		t.Error("structural policies must not record text hashes")
	}
	got, _ = e.Observe(context.Background(), cr)
	if !got.ResourceExists || !got.ResourceUpToDate {
		t.Fatalf("after create: %+v", got)
	}
	cr.Spec.ForProvider.SoftDeletePeriod = ptr(common.Timespan("30d"))
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceUpToDate {
		t.Fatal("change must be detected")
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stored, "30.00:00:00") {
		t.Errorf("update not applied: %s", stored)
	}
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if stored != "" {
		t.Error("delete must remove the policy")
	}
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceExists {
		t.Error("after delete the policy is inherited, not owned")
	}
}

func TestUpdatePolicyKQLTolerance(t *testing.T) {
	var stored string
	kc := fake.New("http://e").OnFn("", func(_ string, c cmd.Command) (*kusto.Result, error) {
		text := c.String()
		switch {
		case strings.HasPrefix(text, ".alter table ['Target'] policy update "):
			js := strings.TrimSuffix(strings.TrimPrefix(text, ".alter table ['Target'] policy update @'"), "'")
			// Kusto echoes defaults and reformats the query.
			stored = strings.Replace(js, `"Query":"Src\n| project a"`, `"Query":"Src | project a","IsTransactional":false,"PropagateIngestionProperties":false`, 1)
			return kusto.NewResult(), nil
		case strings.Contains(text, "policy update"):
			var pol any
			if stored != "" {
				pol = stored
			}
			return kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{"UpdatePolicy", "[DB].[Target]", pol, nil, "Table"})), nil
		}
		return nil, errors.New("unexpected " + text)
	})
	e := ext(kc, Update(), false)
	cr := &v1alpha1.UpdatePolicy{ObjectMeta: metav1.ObjectMeta{Name: "u", Namespace: "ns"}}
	cr.Spec.ForProvider = v1alpha1.UpdatePolicyParameters{PolicyTarget: target(common.EntityKindTable, "Target"), Updates: []v1alpha1.UpdatePolicyEntry{{Source: "Src", Query: "Src\n| project a"}}}

	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stored, `"IsEnabled":true`) || !strings.Contains(stored, `"Query":"Src | project a"`) {
		t.Errorf("IsEnabled must default to true and the fake must have reformatted the query: %s", stored)
	}
	if a, o := base.Hashes(cr); a == "" || o == "" {
		t.Fatal("KQL policies must record hashes")
	}
	for i := 0; i < 5; i++ {
		got, err := e.Observe(context.Background(), cr)
		if err != nil || !got.ResourceUpToDate {
			t.Fatalf("observe %d must be up to date despite reformatting: %+v %v", i, got, err)
		}
	}
	// Drift in the cluster query.
	stored = strings.Replace(stored, "Src | project a", "Src | project b", 1)
	e.cache.Invalidate("http://e", "DB")
	got, _ := e.Observe(context.Background(), cr)
	if got.ResourceUpToDate {
		t.Fatal("cluster drift must be detected")
	}
	// Unresolved source.
	cr.Spec.ForProvider.Updates[0].Source = ""
	if _, err := e.Observe(context.Background(), cr); !kerrors.IsBlocked(err) {
		t.Errorf("empty source must be Blocked(InvalidSpec): %v", err)
	}
}

func TestSpecialAlterForms(t *testing.T) {
	ent := target
	cases := map[string]struct {
		def  func() (cmd.Command, error)
		want string
	}{
		"rls": {func() (cmd.Command, error) {
			d := RowLevelSecurity()
			cr := &v1alpha1.RowLevelSecurityPolicy{}
			cr.Spec.ForProvider = v1alpha1.RowLevelSecurityPolicyParameters{PolicyTarget: ent(common.EntityKindTable, "T"), Query: "T | where a == 1"}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter table ['T'] policy row_level_security enable \"T | where a == 1\""},
		"rlsDisabled": {func() (cmd.Command, error) {
			d := RowLevelSecurity()
			cr := &v1alpha1.RowLevelSecurityPolicy{}
			cr.Spec.ForProvider = v1alpha1.RowLevelSecurityPolicyParameters{PolicyTarget: ent(common.EntityKindMaterializedView, "MV"), Enabled: ptr(false), Query: "Q"}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter materialized-view ['MV'] policy row_level_security disable \"Q\""},
		"ingestionTime": {func() (cmd.Command, error) {
			d := IngestionTime()
			cr := &v1alpha1.IngestionTimePolicy{}
			cr.Spec.ForProvider = v1alpha1.IngestionTimePolicyParameters{PolicyTarget: ent(common.EntityKindTable, "T"), Enabled: true}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter table ['T'] policy ingestiontime true"},
		"restrictedView": {func() (cmd.Command, error) {
			d := RestrictedViewAccess()
			cr := &v1alpha1.RestrictedViewAccessPolicy{}
			cr.Spec.ForProvider = v1alpha1.RestrictedViewAccessPolicyParameters{PolicyTarget: ent(common.EntityKindTable, "T"), Enabled: false}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter table ['T'] policy restricted_view_access false"},
		"caching": {func() (cmd.Command, error) {
			d := Caching()
			cr := &v1alpha1.CachingPolicy{}
			cr.Spec.ForProvider = v1alpha1.CachingPolicyParameters{PolicyTarget: ent(common.EntityKindDatabase, ""), Hot: "31d", HotWindows: []v1alpha1.HotWindow{{MinValue: "2026-01-01T00:00:00Z", MaxValue: "2026-02-01T00:00:00Z"}}}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter database ['DB'] policy caching hot = 31d, hot_window = datetime(2026-01-01T00:00:00Z) .. datetime(2026-02-01T00:00:00Z)"},
		"encodingColumn": {func() (cmd.Command, error) {
			d := Encoding()
			cr := &v1alpha1.EncodingPolicy{}
			cr.Spec.ForProvider = v1alpha1.EncodingPolicyParameters{PolicyTarget: ent(common.EntityKindTable, "T"), Column: ptr("C"), Type: "BigObject32"}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter column ['T'].['C'] policy encoding type = \"BigObject32\""},
		"rowOrder": {func() (cmd.Command, error) {
			d := RowOrder()
			cr := &v1alpha1.RowOrderPolicy{}
			cr.Spec.ForProvider = v1alpha1.RowOrderPolicyParameters{PolicyTarget: ent(common.EntityKindTable, "T"), Columns: []v1alpha1.RowOrderColumn{{Name: "a", Direction: "asc"}, {Name: "b", Direction: "desc"}}}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter table ['T'] policy roworder (['a'] asc, ['b'] desc)"},
		"managedIdentity": {func() (cmd.Command, error) {
			d := ManagedIdentity()
			cr := &v1alpha1.ManagedIdentityPolicy{}
			cr.Spec.ForProvider = v1alpha1.ManagedIdentityPolicyParameters{PolicyTarget: ent(common.EntityKindDatabase, ""), Identities: []v1alpha1.ManagedIdentityEntry{{ObjectID: "system", AllowedUsages: []string{"NativeIngestion", "ExternalTable"}}}}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter database ['DB'] policy managed_identity @'[{\"ObjectId\":\"system\",\"AllowedUsages\":\"NativeIngestion, ExternalTable\"}]'"},
		"extentTags": {func() (cmd.Command, error) {
			d := ExtentTagsRetention()
			cr := &v1alpha1.ExtentTagsRetentionPolicy{}
			cr.Spec.ForProvider = v1alpha1.ExtentTagsRetentionPolicyParameters{PolicyTarget: ent(common.EntityKindTable, "T"), Rules: []v1alpha1.ExtentTagsRetentionRule{{TagPrefix: "drop-by:", RetentionPeriod: "3d"}}}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter table ['T'] policy extent_tags_retention @'[{\"TagPrefix\":\"drop-by:\",\"RetentionPeriod\":\"3.00:00:00\"}]'"},
		"queryAcceleration": {func() (cmd.Command, error) {
			d := QueryAcceleration()
			cr := &v1alpha1.QueryAccelerationPolicy{}
			cr.Spec.ForProvider = v1alpha1.QueryAccelerationPolicyParameters{PolicyTarget: ent(common.EntityKindExternalTable, "E"), Enabled: true, Hot: ptr(common.Timespan("1d"))}
			des, _ := d.Desired(cr)
			return d.Policy.AlterCmd(*mustEntity(d, cr), des)
		}, ".alter external table ['E'] policy query_acceleration @'{\"IsEnabled\":true,\"Hot\":\"1.00:00:00\"}'"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c, err := tc.def()
			if err != nil {
				t.Fatal(err)
			}
			if c.String() != tc.want {
				t.Errorf("\n got %s\nwant %s", c.String(), tc.want)
			}
		})
	}
}

func mustEntity[T Policy](d Def[T], cr T) *adxpolicy.Entity {
	e := &external[T]{def: d}
	ent, err := e.entity(cr)
	if err != nil {
		panic(err)
	}
	return &ent
}

func TestDesiredConversions(t *testing.T) {
	// Every Desired function must reject invalid timespans and accept valid ones.
	rp := &v1alpha1.RetentionPolicy{}
	rp.Spec.ForProvider = v1alpha1.RetentionPolicyParameters{PolicyTarget: target(common.EntityKindTable, "T"), SoftDeletePeriod: ptr(common.Timespan("bad"))}
	if _, err := Retention().Desired(rp); err == nil {
		t.Error("retention must reject bad timespan")
	}
	ib := &v1alpha1.IngestionBatchingPolicy{}
	ib.Spec.ForProvider = v1alpha1.IngestionBatchingPolicyParameters{PolicyTarget: target(common.EntityKindDatabase, ""), MaximumBatchingTimeSpan: ptr(common.Timespan("5m")), MaximumNumberOfItems: ptr(int64(500))}
	d, err := IngestionBatching().Desired(ib)
	if err != nil || *d.(ingestionBatchingJSON).MaximumBatchingTimeSpan != "00:05:00" {
		t.Errorf("ingestion batching: %+v %v", d, err)
	}
	si := &v1alpha1.StreamingIngestionPolicy{}
	si.Spec.ForProvider = v1alpha1.StreamingIngestionPolicyParameters{PolicyTarget: target(common.EntityKindTable, "T"), Enabled: true, HintAllocatedRate: ptr("2.5")}
	d, err = StreamingIngestion().Desired(si)
	if err != nil || *d.(streamingIngestionJSON).HintAllocatedRate != 2.5 {
		t.Errorf("streaming: %+v %v", d, err)
	}
	si.Spec.ForProvider.HintAllocatedRate = ptr("x")
	if _, err := StreamingIngestion().Desired(si); err == nil {
		t.Error("bad rate must fail")
	}
	mp := &v1alpha1.MergePolicy{}
	mp.Spec.ForProvider = v1alpha1.MergePolicyParameters{PolicyTarget: target(common.EntityKindTable, "T"), LoopPeriod: ptr(common.Timespan("1h")), Lookback: &v1alpha1.MergeLookback{Kind: "Custom", CustomPeriod: ptr(common.Timespan("2d"))}}
	d, err = Merge().Desired(mp)
	if err != nil || *d.(mergeJSON).Lookback.CustomPeriod != "2.00:00:00" || *d.(mergeJSON).LoopPeriod != "01:00:00" {
		t.Errorf("merge: %+v %v", d, err)
	}
	pp := &v1alpha1.PartitioningPolicy{}
	pp.Spec.ForProvider = v1alpha1.PartitioningPolicyParameters{PolicyTarget: target(common.EntityKindTable, "T"), PartitionKeys: []v1alpha1.PartitionKey{{ColumnName: "ts", Kind: "UniformRange", Properties: v1alpha1.PartitionKeyProperties{RangeSize: ptr(common.Timespan("1d")), Reference: ptr("1970-01-01T00:00:00")}}}}
	d, err = Partitioning().Desired(pp)
	if err != nil || *d.(partitioningJSON).PartitionKeys[0].Properties.RangeSize != "1.00:00:00" {
		t.Errorf("partitioning: %+v %v", d, err)
	}
	ad := &v1alpha1.AutoDeletePolicy{}
	ad.Spec.ForProvider = v1alpha1.AutoDeletePolicyParameters{PolicyTarget: target(common.EntityKindTable, "T"), ExpiryDate: "tomorrow"}
	if _, err := AutoDelete().Desired(ad); err == nil {
		t.Error("auto delete must reject non-ISO dates")
	}
	ad.Spec.ForProvider.ExpiryDate = "2030-01-01"
	if _, err := AutoDelete().Desired(ad); err != nil {
		t.Error(err)
	}
	cp := &v1alpha1.CachingPolicy{}
	cp.Spec.ForProvider = v1alpha1.CachingPolicyParameters{PolicyTarget: target(common.EntityKindTable, "T"), Hot: "1d", HotWindows: []v1alpha1.HotWindow{{MinValue: "x", MaxValue: "y"}}}
	if _, err := Caching().Desired(cp); err == nil {
		t.Error("caching must reject bad window datetimes")
	}
	sh := &v1alpha1.ShardingPolicy{}
	sh.Spec.ForProvider = v1alpha1.ShardingPolicyParameters{PolicyTarget: target(common.EntityKindTable, "T"), MaxRowCount: ptr(int64(1))}
	if d, err := Sharding().Desired(sh); err != nil || *d.(shardingJSON).MaxRowCount != 1 {
		t.Error("sharding")
	}
	if len(SetupAllDefsForTest()) != 17 {
		t.Errorf("expected 17 policy kinds, got %d", len(SetupAllDefsForTest()))
	}
}

// SetupAllDefsForTest lists the Kusto policy names of all kinds (keeps the
// SetupAll list and this list in sync by count).
func SetupAllDefsForTest() []string {
	return []string{
		Retention().Policy.Name, Caching().Policy.Name, Update().Policy.Name, RowLevelSecurity().Policy.Name,
		IngestionBatching().Policy.Name, StreamingIngestion().Policy.Name, Merge().Policy.Name, Sharding().Policy.Name,
		Partitioning().Policy.Name, IngestionTime().Policy.Name, AutoDelete().Policy.Name, RestrictedViewAccess().Policy.Name,
		ExtentTagsRetention().Policy.Name, Encoding().Policy.Name, ManagedIdentity().Policy.Name, RowOrder().Policy.Name, QueryAcceleration().Policy.Name,
	}
}
