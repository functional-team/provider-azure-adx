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

package ingestionmapping

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
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
)

var cols = []string{"Name", "Kind", "Mapping", "LastUpdatedOn", "Database", "Table"}

func mapping(name, table string, props map[string]string) *v1alpha1.IngestionMapping {
	cr := &v1alpha1.IngestionMapping{ObjectMeta: metav1.ObjectMeta{Name: "m", Namespace: "ns"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider = v1alpha1.IngestionMappingParameters{Database: "DB", Table: table, Kind: "json", Mapping: []v1alpha1.IngestionMappingColumn{{Column: "Timestamp", Properties: props}}}
	return cr
}

// cluster keeps mappings per table and answers show/create/drop.
type cluster struct{ maps map[string]map[string]string } // table -> kind/name -> json

func newCluster() *cluster { return &cluster{maps: map[string]map[string]string{}} }

func (c *cluster) handle(_ string, command cmd.Command) (*kusto.Result, error) {
	text := command.String()
	table := func() string {
		s := text[strings.Index(text, "['")+2:] //nolint:gocritic // test fake, commands are well-formed
		return s[:strings.Index(s, "']")]       //nolint:gocritic // test fake, commands are well-formed
	}
	switch {
	case strings.HasPrefix(text, ".show table") && strings.HasSuffix(text, "ingestion mappings"):
		rows := make([][]any, 0, len(c.maps[table()]))
		for key, js := range c.maps[table()] {
			kind, name, _ := strings.Cut(key, "/")
			rows = append(rows, []any{name, strings.ToUpper(kind[:1]) + kind[1:], js, "2026-09-07T10:00:00Z", "DB", table()})
		}
		return kusto.NewResult(kusto.NewTable("Table_0", cols, rows...)), nil
	case strings.HasPrefix(text, ".show table"):
		// .show table ['T'] ingestion json mapping "Name"
		parts := strings.Split(text, " ")
		kind := parts[4]
		name := strings.Trim(parts[6], "\"")
		js, ok := c.maps[table()][kind+"/"+name]
		if !ok {
			return nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"Mapping was not found","@permanent":true}}`)
		}
		return kusto.NewResult(kusto.NewTable("Table_0", cols, []any{name, kind, js, "", "DB", table()})), nil
	case strings.HasPrefix(text, ".create-or-alter table"):
		parts := strings.SplitN(text, " ", 8)
		kind := parts[4]
		name := strings.Trim(parts[6], "\"")
		js := strings.TrimSuffix(strings.TrimPrefix(parts[7], "@'"), "'")
		// Kusto adds CsvDataType: null to every column.
		js = strings.ReplaceAll(js, "\"Column\":", "\"CsvDataType\":null,\"Column\":")
		if c.maps[table()] == nil {
			c.maps[table()] = map[string]string{}
		}
		c.maps[table()][kind+"/"+name] = js
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".drop table"):
		parts := strings.Split(text, " ")
		delete(c.maps[table()], parts[4]+"/"+strings.Trim(parts[6], "\""))
		return kusto.NewResult(), nil
	}
	return nil, errors.New("unexpected " + text)
}

func newExternal(kc kusto.Client, disabled bool) *external {
	return &external{kc: kc, cache: snapshot.New(time.Minute, disabled)}
}

func TestLifecycle(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		cl := newCluster()
		kc := fake.New("http://e").OnFn("", cl.handle)
		e := newExternal(kc, disabled)
		cr := mapping("RawJson", "RawEvents", map[string]string{"path": "$.ts"})

		got, err := e.Observe(context.Background(), cr)
		if err != nil || got.ResourceExists {
			t.Fatalf("initial observe (disabled=%v): %+v %v", disabled, got, err)
		}
		if _, err := e.Create(context.Background(), cr); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(cl.maps["RawEvents"]["json/RawJson"], `"Path":"$.ts"`) {
			t.Errorf("stored mapping: %s", cl.maps["RawEvents"]["json/RawJson"])
		}
		got, err = e.Observe(context.Background(), cr)
		if err != nil || !got.ResourceExists || !got.ResourceUpToDate {
			t.Fatalf("after create: %+v %v", got, err)
		}
		if cr.Status.AtProvider.Mapping == "" || cr.Status.AtProvider.Kind != "json" || cr.GetCondition(xpv2.TypeReady).Status != "True" {
			t.Errorf("status: %+v", cr.Status.AtProvider)
		}
		// Spec change -> update.
		cr.Spec.ForProvider.Mapping[0].Properties["path"] = "$.time"
		got, _ = e.Observe(context.Background(), cr)
		if got.ResourceUpToDate || got.Diff == "" {
			t.Fatalf("change must be detected: %+v", got)
		}
		if _, err := e.Update(context.Background(), cr); err != nil {
			t.Fatal(err)
		}
		got, _ = e.Observe(context.Background(), cr)
		if !got.ResourceUpToDate {
			t.Fatal("after update must be up to date")
		}
		// Drift in the cluster.
		cl.maps["RawEvents"]["json/RawJson"] = strings.Replace(cl.maps["RawEvents"]["json/RawJson"], "$.time", "$.other", 1)
		e.cache.Invalidate("http://e", "DB")
		if got, _ = e.Observe(context.Background(), cr); got.ResourceUpToDate {
			t.Fatal("cluster drift must be detected")
		}
		// Delete.
		if _, err := e.Delete(context.Background(), cr); err != nil {
			t.Fatal(err)
		}
		if got, _ = e.Observe(context.Background(), cr); got.ResourceExists {
			t.Fatal("deleted mapping must not exist")
		}
		if cr.GetCondition(xpv2.TypeReady).Reason != xpv2.ReasonDeleting {
			t.Error("Delete must set Deleting")
		}
	}
}

func TestCacheSharedPerTable(t *testing.T) {
	cl := newCluster()
	cl.maps["RawEvents"] = map[string]string{"json/A": `[{"Column":"Timestamp","Properties":{"Path":"$.ts"}}]`, "csv/B": `[{"Column":"Timestamp","Properties":{"Ordinal":"0"}}]`}
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := newExternal(kc, false)
	a := mapping("A", "RawEvents", map[string]string{"path": "$.ts"})
	b := mapping("B", "RawEvents", map[string]string{"ordinal": "0"})
	b.Spec.ForProvider.Kind = "csv"
	for i := 0; i < 5; i++ {
		if got, err := e.Observe(context.Background(), a); err != nil || !got.ResourceUpToDate {
			t.Fatalf("a: %+v %v", got, err)
		}
		if got, err := e.Observe(context.Background(), b); err != nil || !got.ResourceUpToDate {
			t.Fatalf("b: %+v %v", got, err)
		}
	}
	if n := kc.Count("ingestion mappings"); n != 1 {
		t.Errorf("expected one batch load per table, got %d", n)
	}
	// Same name, other kind, must not match.
	c := mapping("A", "RawEvents", map[string]string{"path": "$.ts"})
	c.Spec.ForProvider.Kind = "csv"
	if got, _ := e.Observe(context.Background(), c); got.ResourceExists {
		t.Error("identity includes the kind")
	}
}

func TestErrors(t *testing.T) {
	cr := mapping("M", "", nil)
	if _, err := newExternal(fake.New("http://e"), false).Observe(context.Background(), cr); err == nil || !strings.Contains(err.Error(), "table is empty") {
		t.Errorf("unresolved table: %v", err)
	}
	broken := fake.New("http://e").On("", nil, errors.New("boom"))
	if _, err := newExternal(broken, false).Observe(context.Background(), mapping("M", "T", nil)); err == nil {
		t.Error("expected observe error")
	}
	if _, err := newExternal(broken, false).Create(context.Background(), mapping("M", "T", nil)); err == nil {
		t.Error("expected create error")
	}
	if _, err := newExternal(broken, false).Update(context.Background(), mapping("M", "T", nil)); err == nil {
		t.Error("expected update error")
	}
	if _, err := newExternal(broken, false).Delete(context.Background(), mapping("M", "T", nil)); err == nil {
		t.Error("expected delete error")
	}
	notFound := fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"Table 'T' was not found","@permanent":true}}`))
	if got, err := newExternal(notFound, false).Observe(context.Background(), mapping("M", "T", nil)); err != nil || got.ResourceExists {
		t.Errorf("missing table means missing mapping: %+v %v", got, err)
	}
	if _, err := newExternal(notFound, false).Delete(context.Background(), mapping("M", "T", nil)); err != nil {
		t.Errorf("delete not found is fine: %v", err)
	}
	exists := fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityAlreadyExists","@message":"already exists","@permanent":true}}`))
	if _, err := newExternal(exists, false).Create(context.Background(), mapping("M", "T", nil)); err != nil {
		t.Errorf("already exists on create is fine: %v", err)
	}
	if err := newExternal(exists, false).Disconnect(context.Background()); err != nil {
		t.Error(err)
	}
}
