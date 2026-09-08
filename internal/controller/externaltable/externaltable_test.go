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

package externaltable

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/test"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

var showCols = []string{"TableName", "TableType", "Folder", "DocString", "Properties", "ConnectionStrings", "Partitions", "PathFormat"}

const secretCS = "https://acct2.blob.core.windows.net/c;sig=SUPERSECRET"

func ptr[T any](v T) *T { return &v }

func table(name string) *v1alpha1.ExternalTable {
	cr := &v1alpha1.ExternalTable{ObjectMeta: metav1.ObjectMeta{Name: "et", Namespace: "ns"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider = v1alpha1.ExternalTableParameters{
		Database:   "DB",
		Kind:       "Storage",
		Columns:    []common.Column{{Name: "Timestamp", Type: "datetime"}, {Name: "Payload", Type: "dynamic"}},
		DataFormat: ptr("parquet"),
		ConnectionStrings: []v1alpha1.ExternalTableConnectionString{
			{Value: ptr("https://acct.blob.core.windows.net/exports;managed_identity=system")},
			{SecretKeyRef: &v1alpha1.SecretKeySelector{Name: "storage", Key: "cs"}},
		},
		Properties: &v1alpha1.ExternalTableProperties{Folder: ptr("External")},
	}
	return cr
}

func kube(secret string) client.Client {
	return &test.MockClient{
		MockGet: test.NewMockGetFn(nil, func(obj client.Object) error {
			if s, ok := obj.(*corev1.Secret); ok {
				s.Data = map[string][]byte{"cs": []byte(secret)}
			}
			return nil
		}),
		MockPatch: test.NewMockPatchFn(nil),
	}
}

// cluster keeps the external tables and echoes what .show would report:
// connection strings masked after the first ';', partitions as JSON.
type cluster struct {
	tables map[string][]any
}

func (c *cluster) handle(_ string, command cmd.Command) (*kusto.Result, error) {
	text := command.String()
	switch {
	case strings.HasPrefix(text, ".show external tables"):
		rows := make([][]any, 0, len(c.tables))
		for _, r := range c.tables {
			rows = append(rows, r)
		}
		return kusto.NewResult(kusto.NewTable("Table_0", showCols, rows...)), nil
	case strings.HasPrefix(text, ".show external table ") && strings.HasSuffix(text, "cslschema"):
		return kusto.NewResult(kusto.NewTable("Table_0", []string{"TableName", "Schema", "DatabaseName", "Folder", "DocString"}, []any{"Exports", "Timestamp:datetime,Payload:dynamic", "DB", "External", ""})), nil
	case strings.HasPrefix(text, ".show external table "):
		name := between(text, "['", "']")
		if r, ok := c.tables[name]; ok {
			return kusto.NewResult(kusto.NewTable("Table_0", showCols, r)), nil
		}
		return nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"not found","@permanent":true}}`)
	case strings.HasPrefix(text, ".create-or-alter external table "):
		if strings.Contains(text, "SUPERSECRET") && !strings.Contains(text, "h@'"+secretCS+"'") {
			return nil, fake.HTTPError(400, `{"error":{"code":"General_BadRequest","@message":"secret not obfuscated","@permanent":true}}`)
		}
		name := between(text, "['", "']")
		uris := `["https://acct.blob.core.windows.net/exports;*******","https://acct2.blob.core.windows.net/c;*******"]`
		partitions := "null"
		if strings.Contains(text, "partition by") {
			partitions = `[{"ColumnName":"Date","Kind":"datetime"}]`
		}
		c.tables[name] = []any{name, "Blob", "External", "", `{"Format":"parquet"}`, uris, partitions, ""}
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".drop external table "):
		delete(c.tables, between(text, "['", "']"))
		return kusto.NewResult(), nil
	}
	return nil, fake.HTTPError(400, `{"error":{"code":"General_BadRequest","@message":"unexpected: `+text+`"}}`)
}

func between(s, a, b string) string {
	i := strings.Index(s, a)
	if i < 0 {
		return ""
	}
	s = s[i+len(a):]
	j := strings.Index(s, b)
	if j < 0 {
		return s
	}
	return s[:j]
}

func TestLifecycle(t *testing.T) {
	cl := &cluster{tables: map[string][]any{}}
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := &external{kc: kc, cache: snapshot.New(time.Minute, false), kube: kube(secretCS)}
	cr := table("Exports")
	cr.Spec.ForProvider.PartitionBy = ptr("Date:datetime = bin(Timestamp, 1d)")

	obs, err := e.Observe(context.Background(), cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("initial observe: %+v %v", obs, err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	create := kc.Commands()[len(kc.Commands())-2]
	if !strings.Contains(create, "h@'"+secretCS+"'") || !strings.Contains(create, "kind=storage") || !strings.Contains(create, "partition by") {
		t.Errorf("create command: %s", create)
	}
	if strings.Contains(cmd.New(create).Redacted(), "SUPERSECRET") {
		t.Error("redacted command must not leak the secret")
	}
	if cr.GetAnnotations()[base.AnnotationSecretHash] == "" {
		t.Error("secret hash annotation must be set")
	}
	if a, o := base.Hashes(cr); a == "" || o == "" {
		t.Error("text hashes must be set after create")
	}

	for i := 0; i < 3; i++ {
		obs, err = e.Observe(context.Background(), cr)
		if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
			t.Fatalf("observe %d after create: %+v %v", i, obs, err)
		}
	}
	statusJSON, err := json.Marshal(cr.Status.AtProvider)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(statusJSON), "SUPERSECRET") || strings.Contains(string(statusJSON), "sig=") {
		t.Errorf("status must not contain secrets: %s", statusJSON)
	}
	if len(cr.Status.AtProvider.ConnectionStringURIs) != 2 || len(cr.Status.AtProvider.Columns) != 2 {
		t.Errorf("status: %+v", cr.Status.AtProvider)
	}

	// Secret rotation is detected through the hash.
	e.kube = kube("https://acct2.blob.core.windows.net/c;sig=ROTATED")
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || !strings.Contains(obs.Diff, "connection strings") {
		t.Fatalf("rotated secret must be detected: %+v", obs)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	obs, _ = e.Observe(context.Background(), cr)
	if !obs.ResourceUpToDate {
		t.Fatalf("after update: %+v", obs)
	}

	// Partition drift in the spec.
	cr.Spec.ForProvider.PartitionBy = ptr("Date:datetime = bin(Timestamp, 2d)")
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || !strings.Contains(obs.Diff, "partitionBy") {
		t.Fatalf("partition change must be detected: %+v", obs)
	}

	// URI drift in the cluster.
	cl.tables["Exports"][5] = `["https://other.blob.core.windows.net/x;*******","https://acct2.blob.core.windows.net/c;*******"]`
	e.cache.Invalidate("http://e", "DB")
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || !strings.Contains(obs.Diff, "connection string URIs") {
		t.Fatalf("URI drift must be detected: %+v", obs)
	}

	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if _, ok := cl.tables["Exports"]; ok {
		t.Error("delete must drop the table")
	}
	if err := e.Disconnect(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestSecretErrors(t *testing.T) {
	kc := fake.New("http://e")
	missing := &test.MockClient{MockGet: test.NewMockGetFn(nil, func(obj client.Object) error {
		if s, ok := obj.(*corev1.Secret); ok {
			s.Data = map[string][]byte{"other": []byte("x")}
		}
		return nil
	})}
	e := &external{kc: kc, cache: snapshot.New(time.Minute, false), kube: missing}
	if _, err := e.Observe(context.Background(), table("E")); err == nil || !strings.Contains(err.Error(), "has no key") {
		t.Errorf("missing key must fail: %v", err)
	}
	e.kube = nil
	if _, err := e.Observe(context.Background(), table("E")); err == nil {
		t.Error("no kube client must fail")
	}
	if len(kc.Commands()) != 0 {
		t.Error("no command must be sent when secrets cannot be resolved")
	}
}

func TestCacheDisabledSingleShow(t *testing.T) {
	kc := fake.New("http://e").
		On("cslschema", kusto.NewResult(kusto.NewTable("Table_0", []string{"TableName", "Schema", "DatabaseName", "Folder", "DocString"}, []any{"E", "Timestamp:datetime,Payload:dynamic", "DB", "", ""})), nil).
		On(".show external table ['E']", kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{"E", "Blob", "External", "", `{"Format":"parquet"}`, `["https://acct.blob.core.windows.net/exports;*******","https://acct2.blob.core.windows.net/c;*******"]`, nil, ""})), nil)
	e := &external{kc: kc, cache: snapshot.New(0, true), kube: kube(secretCS)}
	cr := table("E")
	meta.AddAnnotations(cr, map[string]string{base.AnnotationSecretHash: "stale"})
	obs, err := e.Observe(context.Background(), cr)
	if err != nil || !obs.ResourceExists || obs.ResourceUpToDate {
		t.Fatalf("single show with stale secret hash: %+v %v", obs, err)
	}
	if kc.Count(".show external tables") != 0 {
		t.Error("cache disabled must not use the batch command")
	}
}
