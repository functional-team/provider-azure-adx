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

package table

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema/schematest"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
)

func schemaResult() *kusto.Result {
	return kusto.NewResult(kusto.NewTable("Table_0", []string{"DatabaseSchema"}, []any{schematest.DatabaseJSON}))
}

func newTable(name string, cols ...common.Column) *v1alpha1.Table {
	cr := &v1alpha1.Table{ObjectMeta: metav1.ObjectMeta{Name: "raw-events", Namespace: "ns"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider = v1alpha1.TableParameters{Database: "Telemetry", Columns: cols, SchemaUpdateMode: v1alpha1.SchemaUpdateModeMerge}
	return cr
}

func newExternal(kc kusto.Client, disabled bool) *external {
	return &external{kc: kc, cache: snapshot.New(time.Minute, disabled)}
}

func TestObserve(t *testing.T) {
	ts := common.Column{Name: "Timestamp", Type: "datetime"}
	payload := common.Column{Name: "Payload", Type: "dynamic"}
	cases := map[string]struct {
		kc       *fake.Client
		disabled bool
		cr       *v1alpha1.Table
		want     managed.ExternalObservation
		wantErr  error
		blocked  string
		drift    []string
	}{
		"notFound": {
			kc:   fake.New("http://e").On("schema as json", schemaResult(), nil),
			cr:   newTable("Missing", ts),
			want: managed.ExternalObservation{ResourceExists: false},
		},
		"upToDate": {
			kc:   fake.New("http://e").On("schema as json", schemaResult(), nil),
			cr:   newTable("RawEvents", ts, payload),
			want: managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true},
		},
		"needsUpdate": {
			kc:   fake.New("http://e").On("schema as json", schemaResult(), nil),
			cr:   newTable("RawEvents", ts, payload, common.Column{Name: "Level", Type: "string"}),
			want: managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: false, Diff: "AlterMergeSchema"},
		},
		"driftReported": {
			kc:    fake.New("http://e").On("schema as json", schemaResult(), nil),
			cr:    newTable("RawEvents", ts),
			want:  managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true},
			drift: []string{"Payload"},
		},
		"typeChangeBlocked": {
			kc:      fake.New("http://e").On("schema as json", schemaResult(), nil),
			cr:      newTable("RawEvents", common.Column{Name: "Timestamp", Type: "string"}, payload),
			blocked: kerrors.ReasonUnsupportedColumnTypeChange,
		},
		"databaseNotFound": {
			kc:   fake.New("http://e").On("schema as json", nil, notFound()),
			cr:   newTable("RawEvents", ts),
			want: managed.ExternalObservation{ResourceExists: false},
		},
		"clusterError": {
			kc:      fake.New("http://e").On("schema as json", nil, errors.New("boom")),
			cr:      newTable("RawEvents", ts),
			wantErr: errors.New(errObserve),
		},
		"cacheDisabledSingleShow": {
			kc: fake.New("http://e").On("cslschema", kusto.NewResult(kusto.NewTable("Table_0",
				[]string{"TableName", "Schema", "DatabaseName", "Folder", "DocString"},
				[]any{"RawEvents", "Timestamp:datetime,Payload:dynamic", "Telemetry", "", ""})), nil),
			disabled: true,
			cr:       newTable("RawEvents", ts, payload),
			want:     managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true},
		},
		"cacheDisabledNotFound": {
			kc:       fake.New("http://e").On("cslschema", nil, notFound()),
			disabled: true,
			cr:       newTable("RawEvents", ts),
			want:     managed.ExternalObservation{ResourceExists: false},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := newExternal(tc.kc, tc.disabled)
			got, err := e.Observe(context.Background(), tc.cr)
			if tc.blocked != "" {
				var b *kerrors.Blocked
				if !errors.As(err, &b) || b.Reason != tc.blocked {
					t.Fatalf("expected Blocked(%s), got %v", tc.blocked, err)
				}
				if len(tc.kc.Commands()) != 1 {
					t.Errorf("blocked observe must not send extra commands: %v", tc.kc.Commands())
				}
				return
			}
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error containing %q", tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("-want +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.drift, tc.cr.Status.AtProvider.DriftColumns); diff != "" {
				t.Errorf("drift -want +got:\n%s", diff)
			}
			if got.ResourceExists && tc.cr.GetCondition(xpv2.TypeReady).Status != "True" {
				t.Error("existing table must be Available")
			}
		})
	}
}

func TestObserveUsesCache(t *testing.T) {
	kc := fake.New("http://e").On("schema as json", schemaResult(), nil)
	e := newExternal(kc, false)
	for i := 0; i < 5; i++ {
		if _, err := e.Observe(context.Background(), newTable("RawEvents", common.Column{Name: "Timestamp", Type: "datetime"}, common.Column{Name: "Payload", Type: "dynamic"})); err != nil {
			t.Fatal(err)
		}
		if _, err := e.Observe(context.Background(), newTable("Dim", common.Column{Name: "Id", Type: "string"})); err != nil {
			t.Fatal(err)
		}
	}
	if n := kc.Count("schema as json"); n != 1 {
		t.Errorf("expected one batch load for the database, got %d", n)
	}
}

func TestCreateUpdateDelete(t *testing.T) {
	kc := fake.New("http://e").On("schema as json", schemaResult(), nil).On("", kusto.NewResult(), nil)
	e := newExternal(kc, false)
	cr := newTable("New", common.Column{Name: "A", Type: "string", Docstring: ptr("doc")})
	cr.Spec.ForProvider.Folder = ptr("F")

	// Observe first so the cache is warm, then Create must invalidate it.
	if _, err := e.Observe(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	want := []string{
		".show database ['Telemetry'] schema as json",
		".create table ['New'] (['A']:string) with (folder=\"F\")",
		".alter-merge table ['New'] column-docstrings (['A']:\"doc\")",
	}
	if diff := cmp.Diff(want, kc.Commands()); diff != "" {
		t.Errorf("create commands -want +got:\n%s", diff)
	}
	kc.Reset()
	if _, err := e.Observe(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count("schema as json") != 1 {
		t.Error("create must invalidate the schema cache")
	}

	// Update: add a column to RawEvents.
	kc.Reset()
	upd := newTable("RawEvents", common.Column{Name: "Timestamp", Type: "datetime"}, common.Column{Name: "Payload", Type: "dynamic"}, common.Column{Name: "Level", Type: "string"})
	if _, err := e.Update(context.Background(), upd); err != nil {
		t.Fatal(err)
	}
	// The schema cache is still warm from the Observe above, so only the alter is sent.
	if got := kc.Commands(); len(got) != 1 || got[0] != ".alter-merge table ['RawEvents'] (['Timestamp']:datetime, ['Payload']:dynamic, ['Level']:string)" {
		t.Errorf("update commands: %v", got)
	}

	// Update failure surfaces the step kind.
	failing := fake.New("http://e").On("schema as json", schemaResult(), nil).On(".alter-merge", nil, errors.New("rejected"))
	if _, err := newExternal(failing, false).Update(context.Background(), upd); err == nil {
		t.Error("expected update error")
	}

	// Delete.
	kc.Reset()
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if got := kc.Commands(); len(got) != 1 || got[0] != ".drop table ['New'] ifexists" {
		t.Errorf("delete commands: %v", got)
	}
	if cr.GetCondition(xpv2.TypeReady).Reason != xpv2.ReasonDeleting {
		t.Error("Delete must set Deleting")
	}
	// NotFound on delete is fine.
	gone := fake.New("http://e").On("", nil, notFound())
	if _, err := newExternal(gone, false).Delete(context.Background(), cr); err != nil {
		t.Errorf("NotFound on delete must be ignored: %v", err)
	}
	// AlreadyExists on create is fine.
	exists := fake.New("http://e").On("", nil, alreadyExists())
	if _, err := newExternal(exists, false).Create(context.Background(), cr); err != nil {
		t.Errorf("AlreadyExists on create must be ignored: %v", err)
	}
	if err := e.Disconnect(context.Background()); err != nil {
		t.Error(err)
	}
	_ = test.EquateErrors
}

func ptr(s string) *string { return &s }

func notFound() error {
	return fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"Entity 'X' of kind 'Table' was not found.","@permanent":true}}`)
}

func alreadyExists() error {
	return fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityAlreadyExists","@message":"already exists","@permanent":true}}`)
}
