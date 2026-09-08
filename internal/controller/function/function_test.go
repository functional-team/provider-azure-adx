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

package function

import (
	"context"
	"strings"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

var cols = []string{"Name", "Parameters", "Body", "Folder", "DocString"}

func fn(name, params, body string) *v1alpha1.Function {
	cr := &v1alpha1.Function{ObjectMeta: metav1.ObjectMeta{Name: "f", Namespace: "ns"}}
	meta.SetExternalName(cr, name)
	cr.Spec.ForProvider = v1alpha1.FunctionParameters{Database: "DB", Parameters: params, Body: body}
	return cr
}

// cluster simulates Kusto: it stores functions and echoes bodies reformatted
// ("{ body }" with the interior whitespace collapsed, like the A6 example).
type cluster struct {
	funcs map[string][]any
}

func newCluster() *cluster { return &cluster{funcs: map[string][]any{}} }

func (c *cluster) reformat(body string) string {
	return "{ " + strings.Join(strings.Fields(body), " ") + " }"
}

func (c *cluster) handle(_ string, command cmd.Command) (*kusto.Result, error) {
	text := command.String()
	switch {
	case strings.HasPrefix(text, ".show functions"):
		rows := make([][]any, 0, len(c.funcs))
		for _, r := range c.funcs {
			rows = append(rows, r)
		}
		return kusto.NewResult(kusto.NewTable("Table_0", cols, rows...)), nil
	case strings.HasPrefix(text, ".create-or-alter function"):
		// ['Name'](params) {\nbody\n}
		start := strings.Index(text, "['") + 2
		end := strings.Index(text[start:], "']") + start
		name := text[start:end]
		rest := text[end+2:]
		brace := strings.Index(rest, "{")
		params := strings.TrimSpace(rest[:brace])
		body := strings.TrimSuffix(strings.TrimPrefix(rest[brace:], "{"), "}")
		params = strings.ReplaceAll(strings.ReplaceAll(params, " ", ""), "\t", "")
		row := []any{name, params, c.reformat(body), "", ""}
		c.funcs[name] = row
		return kusto.NewResult(kusto.NewTable("Table_0", cols, row)), nil
	case strings.HasPrefix(text, ".drop function"):
		start := strings.Index(text, "['") + 2
		end := strings.Index(text[start:], "']") + start
		delete(c.funcs, text[start:end])
		return kusto.NewResult(), nil
	}
	return nil, fake.HTTPError(400, `{"error":{"code":"General_BadRequest","@message":"unexpected: `+text+`"}}`)
}

func newExternal(kc kusto.Client) *external {
	return &external{kc: kc, cache: snapshot.New(time.Minute, false)}
}

// TestNoUpdateLoop is the M1 acceptance criterion in miniature: a function
// with "ugly" formatting is created once and then observed as up to date
// although the cluster echoes a reformatted body, and a real change on
// either side is still detected.
func TestNoUpdateLoop(t *testing.T) {
	cl := newCluster()
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := newExternal(kc)
	cr := fn("Ugly", "(a:string,   b:int = 5)", "RawEvents\n|  where   Level == \"Error\"   \n| take b")

	// Not there yet.
	obs, err := e.Observe(context.Background(), cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("initial observe: %+v %v", obs, err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	applied, observed := base.Hashes(cr)
	if applied == "" || observed == "" || applied == observed {
		t.Fatalf("hashes after create: %q %q", applied, observed)
	}

	// Ten observes: up to date every time, no writes.
	writes := kc.Count(".create-or-alter")
	for i := 0; i < 10; i++ {
		obs, err := e.Observe(context.Background(), cr)
		if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
			t.Fatalf("observe %d: %+v %v", i, obs, err)
		}
	}
	if kc.Count(".create-or-alter") != writes {
		t.Error("observe must not write")
	}
	if cr.Status.AtProvider.Body == "" {
		t.Error("atProvider should carry the observed body")
	}

	// Spec change -> update needed, one write, then stable again.
	cr.Spec.ForProvider.Body = "RawEvents | take 1"
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate {
		t.Fatal("spec change must be detected")
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	obs, _ = e.Observe(context.Background(), cr)
	if !obs.ResourceUpToDate {
		t.Fatal("after update the function must be up to date")
	}

	// Drift in the cluster -> detected.
	cl.funcs["Ugly"][2] = "{ RawEvents | take 2 }"
	e.cache.Invalidate("http://e", "DB")
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate {
		t.Fatal("cluster drift must be detected")
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	obs, _ = e.Observe(context.Background(), cr)
	if !obs.ResourceUpToDate {
		t.Fatal("drift repaired")
	}

	// Metadata drift.
	cr.Spec.ForProvider.Folder = ptr("Parsing")
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate {
		t.Fatal("folder change must be detected")
	}

	// Delete.
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceExists {
		t.Fatal("deleted function must not exist")
	}
}

func TestObserveCacheDisabledAndErrors(t *testing.T) {
	kc := fake.New("http://e").On(".show function ['F']", kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"F", "()", "{ T }", "", ""})), nil)
	e := &external{kc: kc, cache: snapshot.New(0, true)}
	obs, err := e.Observe(context.Background(), fn("F", "()", "T"))
	if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("single show: %+v %v", obs, err)
	}
	gone := fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"not found","@permanent":true}}`))
	obs, err = newExternal(gone).Observe(context.Background(), fn("F", "()", "T"))
	if err != nil || obs.ResourceExists {
		t.Fatalf("not found: %+v %v", obs, err)
	}
	if _, err := newExternal(gone).Delete(context.Background(), fn("F", "()", "T")); err != nil {
		t.Errorf("delete not found must be ok: %v", err)
	}
	broken := fake.New("http://e").On("", nil, fake.HTTPError(500, `{"error":{"@message":"boom"}}`))
	if _, err := newExternal(broken).Observe(context.Background(), fn("F", "()", "T")); err == nil {
		t.Error("expected observe error")
	}
	if _, err := newExternal(broken).Create(context.Background(), fn("F", "()", "T")); err == nil {
		t.Error("expected create error")
	}
	// Write without echo falls back to .show function.
	quiet := fake.New("http://e").On(".create-or-alter", kusto.NewResult(), nil).On(".show function", kusto.NewResult(kusto.NewTable("Table_0", cols, []any{"F", "()", "{ T }", "", ""})), nil)
	cr := fn("F", "()", "T")
	if _, err := newExternal(quiet).Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if a, _ := base.Hashes(cr); a == "" {
		t.Error("hashes must be set after fallback read")
	}
}

func ptr[T any](v T) *T { return &v }
