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

package securityrole

import (
	"context"
	"regexp"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/security/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
)

var cols = []string{"Role", "PrincipalType", "PrincipalDisplayName", "PrincipalObjectId", "PrincipalFQN"}

const tenant = "72f988bf-86f1-41af-91ab-2d7cd011db47"

// cluster simulates the principal store of one table role. It renders rows
// the way Kusto does: apps by app id, users by object id with the UPN only in
// the display name.
type cluster struct {
	rows [][]any
}

var literalRe = regexp.MustCompile(`"((?:[^"\\]|\\.)*)"`)

func (c *cluster) row(principal string) []any {
	typ, id, _ := strings.Cut(principal, "=")
	id, _, _ = strings.Cut(id, ";")
	switch typ {
	case "aadapp":
		return []any{"Table RawEvents Ingestor", "AAD Application", "App (app id: " + id + ")", "oid-" + id, "aadapp=" + id + ";" + tenant}
	case "aadgroup":
		return []any{"Table RawEvents Ingestor", "AAD Group", id, "oid-" + id, "aadgroup=oid-" + id + ";" + tenant}
	default:
		return []any{"Table RawEvents Ingestor", "AAD User", "Alice (upn: " + id + ")", "oid-" + id, "aaduser=oid-" + id + ";" + tenant}
	}
}

func (c *cluster) has(principal string) bool {
	r := c.row(principal)
	for _, x := range c.rows {
		if x[3] == r[3] {
			return true
		}
	}
	return false
}

func (c *cluster) handle(_ string, command cmd.Command) (*kusto.Result, error) {
	text := command.String()
	lits := func() []string {
		var out []string
		for _, m := range literalRe.FindAllStringSubmatch(text, -1) {
			if !strings.Contains(m[1], "managed by") {
				out = append(out, m[1])
			}
		}
		return out
	}
	switch {
	case strings.HasPrefix(text, ".show table ['RawEvents'] principals"):
		return kusto.NewResult(kusto.NewTable("Table_0", cols, append([][]any{{"Table RawEvents Admin", "AAD User", "Admin", "oid-admin", "aaduser=oid-admin;" + tenant}}, c.rows...)...)), nil
	case strings.HasPrefix(text, ".set table ['RawEvents'] ingestors none"):
		c.rows = nil
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".set table ['RawEvents'] ingestors "):
		c.rows = nil
		for _, p := range lits() {
			c.rows = append(c.rows, c.row(p))
		}
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".add table ['RawEvents'] ingestors "):
		for _, p := range lits() {
			if !c.has(p) {
				c.rows = append(c.rows, c.row(p))
			}
		}
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".drop table ['RawEvents'] ingestors "):
		for _, p := range lits() {
			target := c.row(p)
			var keep [][]any
			for _, x := range c.rows {
				if x[3] != target[3] && x[4] != p {
					keep = append(keep, x)
				}
			}
			c.rows = keep
		}
		return kusto.NewResult(), nil
	}
	return nil, fake.HTTPError(400, `{"error":{"code":"General_BadRequest","@message":"unexpected: `+text+`"}}`)
}

func role(mode v1alpha1.Mode, principals ...string) *v1alpha1.SecurityRole {
	cr := &v1alpha1.SecurityRole{ObjectMeta: metav1.ObjectMeta{Name: "r", Namespace: "ns"}}
	name := "RawEvents"
	cr.Spec.ForProvider = v1alpha1.SecurityRoleParameters{Database: "DB", Entity: common.EntityReference{Kind: common.EntityKindTable, Name: &name}, Role: "ingestors", Mode: mode}
	for _, p := range principals {
		cr.Spec.ForProvider.Principals = append(cr.Spec.ForProvider.Principals, common.Principal(p))
	}
	d := "managed by crossplane"
	cr.Spec.ForProvider.Description = &d
	return cr
}

const (
	app  = "aadapp=4c7e82bd-6adb-46c3-b413-fdd44834c69b;" + tenant
	user = "aaduser=alice@contoso.com"
)

func TestAuthoritative(t *testing.T) {
	cl := &cluster{}
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := &external{kc: kc}
	cr := role(v1alpha1.ModeAuthoritative, app, user)

	obs, err := e.Observe(context.Background(), cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("empty role must not exist: %+v %v", obs, err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if cmds := kc.Commands(); kc.Count(".set table ['RawEvents'] ingestors (\"aadapp=") != 1 || !strings.Contains(cmds[len(cmds)-1], "skip-results \"managed by crossplane\"") {
		t.Errorf("create commands: %v", kc.Commands())
	}
	obs, err = e.Observe(context.Background(), cr)
	if err != nil || !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("after create: %+v %v (status %+v)", obs, err, cr.Status.AtProvider)
	}
	if len(cr.Status.AtProvider.ResolvedPrincipals) != 2 || len(cr.Status.AtProvider.Unresolved) != 0 {
		t.Fatalf("resolution: %+v", cr.Status.AtProvider)
	}
	// Both principals resolved to the ids Kusto reports, the user via the UPN in the display name.
	for _, r := range cr.Status.AtProvider.ResolvedPrincipals {
		if r.ObjectID == "" || !strings.HasPrefix(r.ObjectID, "oid-") {
			t.Errorf("unresolved object id: %+v", r)
		}
	}
	// Stable across observes.
	for i := 0; i < 3; i++ {
		if obs, _ = e.Observe(context.Background(), cr); !obs.ResourceUpToDate {
			t.Fatalf("observe %d: %+v", i, obs)
		}
	}
	// Someone adds a principal by hand -> drift -> repaired by .set.
	cl.rows = append(cl.rows, cl.row("aaduser=mallory@contoso.com"))
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || !strings.Contains(obs.Diff, "not in the spec") {
		t.Fatalf("manual principal must be detected: %+v", obs)
	}
	kc.Reset()
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count(".set table") != 1 || len(cl.rows) != 2 {
		t.Errorf("update: %v rows=%d", kc.Commands(), len(cl.rows))
	}
	// Spec removes the user -> drift -> repaired.
	cr.Spec.ForProvider.Principals = []common.Principal{app}
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate {
		t.Fatalf("removed principal must be detected: %+v", obs)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if obs, _ = e.Observe(context.Background(), cr); !obs.ResourceUpToDate || len(cl.rows) != 1 {
		t.Fatalf("after removal: %+v rows=%d", obs, len(cl.rows))
	}
	// Delete sets none.
	kc.Reset()
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count("ingestors none") != 1 || len(cl.rows) != 0 {
		t.Errorf("delete: %v", kc.Commands())
	}
}

func TestAuthoritativeEmpty(t *testing.T) {
	cl := &cluster{rows: [][]any{}}
	cl.rows = append(cl.rows, cl.row("aaduser=mallory@contoso.com"))
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := &external{kc: kc}
	cr := role(v1alpha1.ModeAuthoritative)
	// Role has principals, spec wants none: exists, not up to date.
	obs, err := e.Observe(context.Background(), cr)
	if err != nil || !obs.ResourceExists || obs.ResourceUpToDate {
		t.Fatalf("%+v %v", obs, err)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count("ingestors none") != 1 {
		t.Errorf("update must set none: %v", kc.Commands())
	}
	// Owned + empty: still exists and up to date.
	obs, _ = e.Observe(context.Background(), cr)
	if !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("owned empty role: %+v", obs)
	}
}

func TestAdditive(t *testing.T) {
	cl := &cluster{}
	cl.rows = append(cl.rows, cl.row("aaduser=someone.else@contoso.com"))
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := &external{kc: kc}
	cr := role(v1alpha1.ModeAdditive, app, user)

	obs, err := e.Observe(context.Background(), cr)
	if err != nil || obs.ResourceExists {
		t.Fatalf("nothing of ours assigned yet: %+v %v", obs, err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count(".add table ['RawEvents'] ingestors") != 1 || kc.Count(".set ") != 0 {
		t.Errorf("additive create must use .add: %v", kc.Commands())
	}
	obs, _ = e.Observe(context.Background(), cr)
	if !obs.ResourceExists || !obs.ResourceUpToDate {
		t.Fatalf("after add (foreign principal must be tolerated): %+v %+v", obs, cr.Status.AtProvider)
	}
	if len(cl.rows) != 3 {
		t.Fatalf("rows: %d", len(cl.rows))
	}
	// Remove the user from the spec -> only that one is dropped.
	cr.Spec.ForProvider.Principals = []common.Principal{app}
	obs, _ = e.Observe(context.Background(), cr)
	if obs.ResourceUpToDate || !strings.Contains(obs.Diff, "removed from the spec") {
		t.Fatalf("stale principal must be detected: %+v", obs)
	}
	kc.Reset()
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count(".drop table ['RawEvents'] ingestors") != 1 || len(cl.rows) != 2 {
		t.Errorf("update: %v rows=%d", kc.Commands(), len(cl.rows))
	}
	obs, _ = e.Observe(context.Background(), cr)
	if !obs.ResourceUpToDate {
		t.Fatalf("after drop: %+v", obs)
	}
	// Delete drops only our principal, the foreign one stays.
	kc.Reset()
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if kc.Count(".drop table") != 1 || len(cl.rows) != 1 || kc.Count("none") != 0 {
		t.Errorf("additive delete: %v rows=%d", kc.Commands(), len(cl.rows))
	}
}

func TestErrors(t *testing.T) {
	cr := role(v1alpha1.ModeAuthoritative, app)
	cr.Spec.ForProvider.Entity.Name = nil
	if _, err := (&external{kc: fake.New("http://e")}).Observe(context.Background(), cr); err == nil || !strings.Contains(err.Error(), "entity.name is empty") {
		t.Errorf("unresolved entity: %v", err)
	}
	gone := fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"Entity 'RawEvents' of kind 'Table' was not found.","@permanent":true}}`))
	obs, err := (&external{kc: gone}).Observe(context.Background(), role(v1alpha1.ModeAuthoritative, app))
	if err != nil || obs.ResourceExists {
		t.Errorf("missing entity: %+v %v", obs, err)
	}
	if _, err := (&external{kc: gone}).Delete(context.Background(), role(v1alpha1.ModeAuthoritative, app)); err != nil {
		t.Errorf("delete on missing entity must be ok: %v", err)
	}
	if _, err := (&external{kc: gone}).Delete(context.Background(), role(v1alpha1.ModeAdditive, app)); err != nil {
		t.Errorf("additive delete without resolved principals is a no-op: %v", err)
	}
	if err := (&external{}).Disconnect(context.Background()); err != nil {
		t.Error(err)
	}
}
