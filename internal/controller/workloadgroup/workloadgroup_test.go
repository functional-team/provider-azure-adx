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

package workloadgroup

import (
	"context"
	"errors"
	"strings"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/cluster/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
)

var cols = []string{"WorkloadGroupName", "WorkloadGroup"}

type cluster struct{ groups map[string]string }

func (c *cluster) handle(db string, command cmd.Command) (*kusto.Result, error) {
	if db != "" {
		return nil, errors.New("cluster commands must run with an empty database")
	}
	text := command.String()
	name := func() string {
		s := text[strings.Index(text, "['")+2:] //nolint:gocritic // test fake, commands are well-formed
		return s[:strings.Index(s, "']")]       //nolint:gocritic // test fake, commands are well-formed
	}
	switch {
	case strings.HasPrefix(text, ".show workload_group "):
		n := name()
		js, ok := c.groups[n]
		if !ok {
			return nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"Workload group '`+n+`' was not found.","@permanent":true}}`)
		}
		return kusto.NewResult(kusto.NewTable("Table_0", cols, []any{n, js})), nil
	case strings.HasPrefix(text, ".create-or-alter workload_group "):
		n := name()
		js := text[strings.Index(text, "@'")+2:]
		c.groups[n] = strings.TrimSuffix(js, "'")
		return kusto.NewResult(kusto.NewTable("Table_0", cols, []any{n, c.groups[n]})), nil
	case strings.HasPrefix(text, ".drop workload_group "):
		delete(c.groups, name())
		return kusto.NewResult(), nil
	}
	return nil, errors.New("unexpected " + text)
}

func wg(name, js string) *v1alpha1.WorkloadGroup {
	cr := &v1alpha1.WorkloadGroup{ObjectMeta: metav1.ObjectMeta{Name: "wg", Namespace: "ns"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider.WorkloadGroup = apiextensionsv1.JSON{Raw: []byte(js)}
	return cr
}

func TestLifecycle(t *testing.T) {
	cl := &cluster{groups: map[string]string{"default": `{"RequestLimitsPolicy":{"MaxMemoryPerQueryPerNode":{"IsRelaxable":true,"Value":6442450944}}}`}}
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := &external{kc: kc}
	cr := wg("Ingestion", `{"RequestRateLimitPolicies":[{"IsEnabled":true,"Scope":"WorkloadGroup","LimitKind":"ConcurrentRequests","Properties":{"MaxConcurrentRequests":20}}]}`)

	got, err := e.Observe(context.Background(), cr)
	if err != nil || got.ResourceExists {
		t.Fatalf("initial: %+v %v", got, err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(kc.Commands()[1], ".create-or-alter workload_group ['Ingestion'] @'{") {
		t.Errorf("create: %s", kc.Commands()[1])
	}
	got, err = e.Observe(context.Background(), cr)
	if err != nil || !got.ResourceExists || !got.ResourceUpToDate || cr.Status.AtProvider.WorkloadGroup == "" {
		t.Fatalf("after create: %+v %v %q", got, err, cr.Status.AtProvider.WorkloadGroup)
	}
	if cr.GetCondition(xpv2.TypeReady).Status != "True" {
		t.Error("must be Available")
	}
	// Kusto adds defaults; subset comparison stays up to date.
	cl.groups["Ingestion"] = `{"RequestRateLimitPolicies":[{"IsEnabled":true,"Scope":"WorkloadGroup","LimitKind":"ConcurrentRequests","Properties":{"MaxConcurrentRequests":20}}],"RequestQueuingPolicy":{"IsEnabled":false}}`
	got, _ = e.Observe(context.Background(), cr)
	if !got.ResourceUpToDate {
		t.Fatalf("server defaults must not be drift: %+v", got)
	}
	// Drift.
	cl.groups["Ingestion"] = strings.Replace(cl.groups["Ingestion"], `"MaxConcurrentRequests":20`, `"MaxConcurrentRequests":5`, 1)
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceUpToDate || !strings.Contains(got.Diff, "MaxConcurrentRequests") {
		t.Fatalf("drift must be detected: %+v", got)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	got, _ = e.Observe(context.Background(), cr)
	if !got.ResourceUpToDate {
		t.Fatal("update must repair drift")
	}
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if _, ok := cl.groups["Ingestion"]; ok {
		t.Error("delete must drop the group")
	}
	if cr.GetCondition(xpv2.TypeReady).Reason != xpv2.ReasonDeleting {
		t.Error("Delete must set Deleting")
	}
	// Built-in groups are altered but never dropped.
	def := wg("default", `{"RequestLimitsPolicy":{"MaxMemoryPerQueryPerNode":{"IsRelaxable":true,"Value":6442450944}}}`)
	got, _ = e.Observe(context.Background(), def)
	if !got.ResourceExists || !got.ResourceUpToDate {
		t.Fatalf("default group: %+v", got)
	}
	before := len(kc.Commands())
	if _, err := e.Delete(context.Background(), def); err != nil {
		t.Fatal(err)
	}
	if len(kc.Commands()) != before || cl.groups["default"] == "" {
		t.Error("built-in group must not be dropped")
	}
	// NotFound on delete of a custom group is fine.
	if _, err := e.Delete(context.Background(), wg("Gone", `{"a":1}`)); err != nil {
		t.Errorf("drop of missing group: %v", err)
	}
	if err := e.Disconnect(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestSpecAndErrors(t *testing.T) {
	boom := &external{kc: fake.New("http://e").On("", nil, errors.New("boom"))}
	if _, err := boom.Observe(context.Background(), wg("X", `{"a":1}`)); err == nil || !strings.Contains(err.Error(), errObserve) {
		t.Errorf("observe error: %v", err)
	}
	if _, err := boom.Create(context.Background(), wg("X", `{"a":1}`)); err == nil {
		t.Error("create error expected")
	}
	if _, err := boom.Update(context.Background(), wg("X", `{"a":1}`)); err == nil {
		t.Error("update error expected")
	}
	if _, err := boom.Delete(context.Background(), wg("X", `{"a":1}`)); err == nil {
		t.Error("delete error expected")
	}
	for _, bad := range []string{"", "null", "[1,2]", "{not json"} {
		if _, err := boom.Observe(context.Background(), wg("X", bad)); !kerrors.IsBlocked(err) {
			t.Errorf("%q must be Blocked(InvalidSpec): %v", bad, err)
		}
	}
	cr := &v1alpha1.WorkloadGroup{}
	if cr.SpecName() != "" {
		t.Error("empty spec name")
	}
	n := "WG"
	cr.Spec.ForProvider.Name = &n
	if cr.SpecName() != "WG" {
		t.Error("spec name")
	}
}
