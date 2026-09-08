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

package entitygroup

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema/schematest"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

func group(name string, entities ...string) *v1alpha1.EntityGroup {
	cr := &v1alpha1.EntityGroup{ObjectMeta: metav1.ObjectMeta{Name: "eg", Namespace: "ns"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider = v1alpha1.EntityGroupParameters{Database: "Telemetry", Entities: entities}
	return cr
}

func schemaResult() *kusto.Result {
	return kusto.NewResult(kusto.NewTable("Table_0", []string{"DatabaseSchema"}, []any{schematest.DatabaseJSON}))
}

func TestLifecycle(t *testing.T) {
	kc := fake.New("http://e").
		On("schema as json", schemaResult(), nil).
		On(".show entity_group ['New']", kusto.NewResult(kusto.NewTable("Table_0", []string{"Name", "Entities"}, []any{"New", json.RawMessage(`["cluster('c').database('x')"]`)})), nil).
		On("", kusto.NewResult(), nil)
	e := &external{kc: kc, cache: snapshot.New(time.Minute, false)}

	// Existing group from the schema fixture: up to date, order-insensitive.
	obs, err := e.Observe(context.Background(), group("EG", "cluster('c').database('d')"))
	if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("existing: %+v %v", obs, err)
	}
	obs, _ = e.Observe(context.Background(), group("EG", "cluster('c').database('other')"))
	if !obs.ResourceExists || obs.ResourceUpToDate {
		t.Fatalf("drift must be detected: %+v", obs)
	}
	// Missing group.
	cr := group("New", "cluster('c').database('x')")
	if obs, _ = e.Observe(context.Background(), cr); obs.ResourceExists {
		t.Fatal("New must not exist")
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if cmds := kc.Commands(); cmds[len(cmds)-2] != ".create-or-alter entity_group ['New'] (cluster('c').database('x'))" {
		t.Errorf("create commands: %v", cmds)
	}
	if a, o := base.Hashes(cr); a == "" || a != o {
		t.Errorf("hashes after create: %q %q", a, o)
	}
	kc.Reset()
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if got := kc.Commands(); len(got) != 1 || got[0] != ".drop entity_group ['New']" {
		t.Errorf("delete commands: %v", got)
	}
	if err := e.Disconnect(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestCacheDisabled(t *testing.T) {
	kc := fake.New("http://e").On(".show entity_group ['EG']", kusto.NewResult(kusto.NewTable("Table_0", []string{"Name", "Entities"}, []any{"EG", "cluster('c').database('d')"})), nil)
	e := &external{kc: kc, cache: snapshot.New(0, true)}
	obs, err := e.Observe(context.Background(), group("EG", "cluster('c').database('d')"))
	if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("single show: %+v %v", obs, err)
	}
	gone := fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"not found","@permanent":true}}`))
	obs, err = (&external{kc: gone, cache: snapshot.New(0, true)}).Observe(context.Background(), group("EG", "x"))
	if err != nil || obs.ResourceExists {
		t.Fatalf("not found: %+v %v", obs, err)
	}
}
