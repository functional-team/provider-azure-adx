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
	"strings"
	"testing"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/security/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
)

func TestEntity(t *testing.T) {
	for kind, want := range map[common.EntityKind]string{
		common.EntityKindDatabase: "database ['DB']", common.EntityKindTable: "table ['T']", common.EntityKindMaterializedView: "materialized-view ['T']",
		common.EntityKindExternalTable: "external table ['T']", common.EntityKindFunction: "function ['T']",
	} {
		got, err := (Entity{Kind: kind, Database: "DB", Name: "T"}).Render()
		if err != nil || got != want {
			t.Errorf("%s: %q %v", kind, got, err)
		}
	}
	if _, err := (Entity{Kind: "Weird"}).Render(); err == nil {
		t.Error("unknown kind")
	}
	if (Entity{Kind: common.EntityKindDatabase, Database: "DB"}).Display() != "['DB']" || (Entity{Kind: common.EntityKindTable, Database: "DB", Name: "T"}).Display() != "['DB'].['T']" {
		t.Error("Display")
	}
}

func TestParseAndFilter(t *testing.T) {
	res := kusto.NewResult(kusto.NewTable("Table_0", []string{"Role", "PrincipalType", "PrincipalDisplayName", "PrincipalObjectId", "PrincipalFQN"},
		[]any{"Table RawEvents Admin", "AAD User", "Admin", "oid-a", "aaduser=oid-a;t"},
		[]any{"Table RawEvents Ingestor", "AAD Application", "App (app id: 1111)", "oid-1", "aadapp=1111;t"},
		[]any{"Database UnrestrictedViewer", "AAD User", "Bob", "oid-b", "aaduser=oid-b;t"},
		[]any{"", "", "", "", ""}))
	rows := ParsePrincipals(res)
	if len(rows) != 3 {
		t.Fatalf("rows: %d", len(rows))
	}
	if got := FilterRole(rows, "ingestors"); len(got) != 1 || got[0].ObjectID != "oid-1" {
		t.Errorf("ingestors: %+v", got)
	}
	if got := FilterRole(rows, "unrestrictedviewers"); len(got) != 1 || got[0].DisplayName != "Bob" {
		t.Errorf("unrestrictedviewers: %+v", got)
	}
	if got := FilterRole(rows, "admins"); len(got) != 1 {
		t.Errorf("admins: %+v", got)
	}
	if ParsePrincipals(nil) != nil {
		t.Error("nil result")
	}
	if (Row{FQN: "aadapp=X"}).ID() != "aadapp=x" {
		t.Error("ID falls back to FQN")
	}
}

func TestMatch(t *testing.T) {
	rows := []Row{
		{Type: "AAD Application", DisplayName: "MyApp (app id: 1111-2222)", ObjectID: "oid-app", FQN: "aadapp=1111-2222;tenant"},
		{Type: "AAD User", DisplayName: "Alice (upn: alice@contoso.com)", ObjectID: "oid-alice", FQN: "aaduser=oid-alice;tenant"},
		{Type: "AAD Group", DisplayName: "data-ingest", ObjectID: "oid-group", FQN: "aadgroup=oid-group;tenant"},
	}
	cases := map[string]struct {
		spec string
		want string
	}{
		"fqnExact":        {"aadapp=1111-2222;tenant", "oid-app"},
		"fqnNoTenant":     {"aadapp=1111-2222", "oid-app"},
		"fqnCase":         {"AADAPP=1111-2222;TENANT", "oid-app"},
		"userByUPN":       {"aaduser=alice@contoso.com", "oid-alice"},
		"userByObjectID":  {"aaduser=oid-alice;tenant", "oid-alice"},
		"groupByName":     {"aadgroup=data-ingest;contoso.com", "oid-group"},
		"groupByObjectID": {"aadgroup=oid-group", "oid-group"},
		"wrongTypeUPN":    {"aadapp=alice@contoso.com", ""},
		"unknown":         {"aaduser=nobody@contoso.com", ""},
		"malformed":       {"nonsense", ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r, ok := Match(tc.spec, rows)
			if (tc.want == "") == ok || (ok && r.ObjectID != tc.want) {
				t.Errorf("Match(%q) = %+v, %v; want %q", tc.spec, r, ok, tc.want)
			}
		})
	}
	if _, err := ParsePrincipal("aaduser="); err == nil {
		t.Error("empty id must fail")
	}
	p, _ := ParsePrincipal(" AADApp = 1 ; t ")
	if p.Type != "aadapp" || p.ID != "1" || p.Tenant != "t" {
		t.Errorf("ParsePrincipal: %+v", p)
	}
}

func TestEvaluateAuthoritative(t *testing.T) {
	rows := []Row{{Type: "AAD Application", ObjectID: "oid-app", FQN: "aadapp=1111;t"}, {Type: "AAD User", DisplayName: "Alice (upn: alice@x)", ObjectID: "oid-alice", FQN: "aaduser=oid-alice;t"}}
	ev := Evaluate(Input{Spec: []string{"aadapp=1111", "aaduser=alice@x"}, Rows: rows})
	if !ev.Exists || !ev.UpToDate || len(ev.Resolved) != 2 || len(ev.Unresolved) != 0 {
		t.Fatalf("all resolved: %+v", ev)
	}
	// Extra principal in the cluster.
	ev = Evaluate(Input{Spec: []string{"aadapp=1111"}, Rows: rows})
	if !ev.Exists || ev.UpToDate || len(ev.Extra) != 1 || !strings.Contains(ev.Diff, "not in the spec") {
		t.Fatalf("extra: %+v", ev)
	}
	// Principal missing in the cluster, previously resolved.
	ev = Evaluate(Input{Spec: []string{"aadapp=1111", "aaduser=bob@x"}, Rows: rows[:1], Resolved: []Resolved{{Spec: "aaduser=bob@x", ObjectID: "oid-bob"}}})
	if ev.UpToDate || len(ev.Absent) != 1 || !strings.Contains(ev.Diff, "missing") {
		t.Fatalf("absent: %+v", ev)
	}
	// Unresolvable spec entry: after the write (SpecUnchanged) one unmatched row accounts for it.
	ev = Evaluate(Input{Spec: []string{"aadapp=1111", "aaduser=unknown@x"}, Rows: append(rows[:1], Row{Type: "AAD User", ObjectID: "oid-u", FQN: "aaduser=oid-u;t"}), SpecUnchanged: true, BeforeIDs: []string{"oid-app", "oid-u"}})
	if !ev.UpToDate || len(ev.Unresolved) != 1 {
		t.Fatalf("settled unresolved: %+v", ev)
	}
	ev = Evaluate(Input{Spec: []string{"aadapp=1111", "aaduser=unknown@x"}, Rows: rows[:1], SpecUnchanged: true})
	if ev.UpToDate {
		t.Fatalf("unresolved without a matching row must write: %+v", ev)
	}
	// Single new row attributed via BeforeIDs.
	ev = Evaluate(Input{Spec: []string{"aaduser=unknown@x"}, Rows: []Row{{Type: "AAD User", ObjectID: "oid-new", FQN: "aaduser=oid-new;t"}}, BeforeIDs: []string{}})
	if !ev.UpToDate || len(ev.Resolved) != 1 || ev.Resolved[0].ObjectID != "oid-new" {
		t.Fatalf("attribution: %+v", ev)
	}
	// Empty spec.
	ev = Evaluate(Input{Spec: nil, Rows: rows})
	if !ev.Exists || ev.UpToDate {
		t.Fatalf("empty spec with rows: %+v", ev)
	}
	ev = Evaluate(Input{Spec: nil, Rows: nil, Owned: true})
	if !ev.Exists || !ev.UpToDate {
		t.Fatalf("owned empty: %+v", ev)
	}
	ev = Evaluate(Input{Spec: nil, Rows: nil})
	if ev.Exists {
		t.Fatalf("not owned empty: %+v", ev)
	}
}

func TestEvaluateAdditive(t *testing.T) {
	foreign := Row{Type: "AAD User", ObjectID: "oid-f", FQN: "aaduser=oid-f;t"}
	app := Row{Type: "AAD Application", ObjectID: "oid-app", FQN: "aadapp=1111;t"}
	ev := Evaluate(Input{Additive: true, Spec: []string{"aadapp=1111"}, Rows: []Row{foreign, app}})
	if !ev.Exists || !ev.UpToDate || len(ev.Extra) != 1 {
		t.Fatalf("foreign tolerated: %+v", ev)
	}
	ev = Evaluate(Input{Additive: true, Spec: []string{"aadapp=1111", "aaduser=new@x"}, Rows: []Row{foreign}})
	if ev.Exists || ev.UpToDate || len(ev.ToAdd) != 2 {
		t.Fatalf("nothing of ours: %+v", ev)
	}
	ev = Evaluate(Input{Additive: true, Spec: []string{"aadapp=1111"}, Rows: []Row{foreign, app}, Resolved: []Resolved{{Spec: "aaduser=old@x", ObjectID: "oid-f", FQN: "aaduser=oid-f;t"}}})
	if ev.UpToDate || len(ev.Stale) != 1 || len(ev.ToDrop) != 1 || ev.ToDrop[0] != "aaduser=oid-f;t" {
		t.Fatalf("stale: %+v", ev)
	}
}

func TestBuildAndHash(t *testing.T) {
	e := Entity{Kind: common.EntityKindTable, Database: "DB", Name: "T"}
	d := "desc"
	set, _ := BuildSet(e, "ingestors", []string{"aadapp=1;t", "aaduser=a@x"}, &d)
	if set.String() != `.set table ['T'] ingestors ("aadapp=1;t", "aaduser=a@x") skip-results "desc"` {
		t.Error(set)
	}
	none, _ := BuildSet(e, "ingestors", nil, &d)
	if none.String() != ".set table ['T'] ingestors none skip-results" {
		t.Error(none)
	}
	add, _ := BuildAdd(e, "admins", []string{"aadgroup=g;t"}, nil)
	if add.String() != `.add table ['T'] admins ("aadgroup=g;t") skip-results` {
		t.Error(add)
	}
	drop, _ := BuildDrop(e, "admins", []string{"aadgroup=g;t"})
	if drop.String() != `.drop table ['T'] admins ("aadgroup=g;t") skip-results` {
		t.Error(drop)
	}
	show, _ := ShowPrincipals(Entity{Kind: common.EntityKindDatabase, Database: "DB"})
	if show.String() != ".show database ['DB'] principals" {
		t.Error(show)
	}
	if _, err := BuildSet(Entity{Kind: "Weird"}, "admins", nil, nil); err == nil {
		t.Error("bad entity")
	}
	if SpecHash([]string{"B", " a"}) != SpecHash([]string{"a", "b"}) || SpecHash([]string{"a"}) == SpecHash([]string{"b"}) {
		t.Error("SpecHash must be order and case insensitive")
	}
	if got := IDs([]Row{{ObjectID: "B"}, {ObjectID: "a"}}); got[0] != "a" || got[1] != "b" {
		t.Error("IDs sorted lower-case")
	}
	obs := Observation(Evaluation{Resolved: []Resolved{{Spec: "x", ObjectID: "o"}}, Unresolved: []string{"y"}}, []Row{{FQN: "f"}, {ObjectID: "o2"}})
	if len(obs.ResolvedPrincipals) != 1 || len(obs.Principals) != 2 || obs.Unresolved[0] != "y" {
		t.Errorf("Observation: %+v", obs)
	}
	if got := FromObservation(v1alpha1.SecurityRoleObservation{ResolvedPrincipals: []v1alpha1.ResolvedPrincipal{{Spec: "x", ObjectID: "o"}}}); len(got) != 1 || got[0].ObjectID != "o" {
		t.Error("FromObservation")
	}
}
