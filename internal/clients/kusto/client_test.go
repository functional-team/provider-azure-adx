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

package kusto

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

func TestRowAccessors(t *testing.T) {
	cols := []string{"S", "B", "L", "T", "D", "Dyn", "Null", "BS", "LS", "TS"}
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	r := NewRow(cols, "str", true, int64(42), now, timespan.Hour, json.RawMessage(`{"a":1}`), nil, "true", "7", "1.00:00:00")

	if r.String("S") != "str" || r.String("Missing") != "" || r.String("Null") != "" {
		t.Error("String")
	}
	if r.String("T") != "2026-09-07T10:00:00Z" || r.String("D") != "01:00:00" || r.String("Dyn") != `{"a":1}` {
		t.Errorf("String conversions: %q %q %q", r.String("T"), r.String("D"), r.String("Dyn"))
	}
	if b, ok := r.Bool("B"); !ok || !b {
		t.Error("Bool")
	}
	if b, ok := r.Bool("BS"); !ok || !b {
		t.Error("Bool from string")
	}
	if _, ok := r.Bool("S"); ok {
		t.Error("Bool from non-bool string must fail")
	}
	if n, ok := r.Int64("L"); !ok || n != 42 {
		t.Error("Int64")
	}
	if n, ok := r.Int64("LS"); !ok || n != 7 {
		t.Error("Int64 from string")
	}
	if d, ok := r.Timespan("D"); !ok || d != timespan.Hour {
		t.Error("Timespan")
	}
	if d, ok := r.Timespan("TS"); !ok || d != timespan.Day {
		t.Error("Timespan from string")
	}
	if tm, ok := r.Time("T"); !ok || !tm.Equal(now) {
		t.Error("Time")
	}
	if string(r.Dynamic("Dyn")) != `{"a":1}` || r.Dynamic("Null") != nil || r.Dynamic("Missing") != nil {
		t.Error("Dynamic")
	}
	if string(r.Dynamic("S")) != "str" {
		t.Error("Dynamic from string passes through")
	}
	if !r.Has("S") || r.Has("Nope") {
		t.Error("Has")
	}
}

func TestResultPrimary(t *testing.T) {
	var nilRes *Result
	if nilRes.Primary() != nil || nilRes.Rows() != nil {
		t.Error("nil result")
	}
	secondary := NewTable("QueryStatus", []string{"X"})
	secondary.Primary = false
	primary := NewTable("PrimaryResult", []string{"Name"}, []any{"T1"}, []any{"T2"})
	res := NewResult(secondary, primary)
	if got := res.Primary().Name; got != "PrimaryResult" {
		t.Errorf("Primary() = %s", got)
	}
	if len(res.Rows()) != 2 || res.Rows()[1].String("Name") != "T2" {
		t.Error("Rows")
	}
	onlySecondary := NewResult(secondary)
	if onlySecondary.Primary().Name != "QueryStatus" {
		t.Error("fallback to first table")
	}
	if NewResult().Primary() != nil {
		t.Error("empty result")
	}
}

func TestCommandError(t *testing.T) {
	inner := errors.New("boom")
	err := &CommandError{Command: ".show x", Err: inner}
	if !errors.Is(err, inner) || err.Error() != ".show x: boom" {
		t.Errorf("CommandError: %v", err)
	}
	if got := short(strings.Repeat("a b ", 100)); len(got) != 120 || !strings.HasSuffix(got, "...") {
		t.Errorf("short() = %d chars", len(got))
	}
}

func TestOpContext(t *testing.T) {
	ctx := WithOp(context.Background(), "Table", "observe")
	if op := OpFromContext(ctx); op.Kind != "Table" || op.Op != "observe" {
		t.Errorf("OpFromContext = %+v", op)
	}
	if op := OpFromContext(context.Background()); op.Kind != "unknown" {
		t.Errorf("default op = %+v", op)
	}
}

func TestKCSB(t *testing.T) {
	cases := map[string]struct {
		cfg Config
		err string
	}{
		"noneHTTP":      {cfg: Config{Endpoint: "http://localhost:8080", Auth: Auth{Source: SourceNone}}},
		"noneHTTPS":     {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceNone}}, err: "only allowed for http://"},
		"secret":        {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceSecret, ClientID: "a", ClientSecret: "b", TenantID: "c"}}},
		"secretMissing": {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceSecret, ClientID: "a"}}, err: "clientId, clientSecret and tenantId"},
		"wiExplicit":    {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceWorkloadIdentity, ClientID: "a", TenantID: "t", TokenFile: "/var/run/token"}}},
		"wiMissing":     {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceWorkloadIdentity}}, err: "workload identity needs"},
		"miSystem":      {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceManagedIdentity, ManagedIdentityType: ManagedIdentitySystemAssigned}}},
		"miUserClient":  {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceManagedIdentity, ManagedIdentityType: ManagedIdentityUserAssigned, ClientID: "id"}}},
		"miUserRes":     {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceManagedIdentity, ManagedIdentityType: ManagedIdentityUserAssigned, ResourceID: "/subscriptions/x"}}},
		"miUserMissing": {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceManagedIdentity, ManagedIdentityType: ManagedIdentityUserAssigned}}, err: "clientId or resourceId"},
		"miBadType":     {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceManagedIdentity, ManagedIdentityType: "Weird"}}, err: "unknown managed identity type"},
		"badSource":     {cfg: Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: "Cert"}}, err: "unknown credentials source"},
		"emptyEndpoint": {cfg: Config{Auth: Auth{Source: SourceNone}}, err: "endpoint is empty"},
		"china":         {cfg: Config{Endpoint: "https://c.kusto.chinacloudapi.cn", Environment: EnvironmentChina, Auth: Auth{Source: SourceManagedIdentity}}},
		"gov":           {cfg: Config{Endpoint: "https://c.kusto.usgovcloudapi.net", Environment: EnvironmentGovernment, Auth: Auth{Source: SourceManagedIdentity}}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("AZURE_CLIENT_ID", "")
			t.Setenv("AZURE_TENANT_ID", "")
			t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "")
			kcsb, err := kcsbFor(tc.cfg)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("expected error containing %q, got %v", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if kcsb.DataSource != tc.cfg.Endpoint {
				t.Errorf("DataSource = %q", kcsb.DataSource)
			}
			if tc.cfg.Environment != "" && kcsb.ClientOptions == nil {
				t.Error("sovereign cloud options not attached")
			}
		})
	}
	t.Run("wiFromEnv", func(t *testing.T) {
		t.Setenv("AZURE_CLIENT_ID", "envclient")
		t.Setenv("AZURE_TENANT_ID", "envtenant")
		t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/tok")
		kcsb, err := kcsbFor(Config{Endpoint: "https://c.kusto.windows.net", Auth: Auth{Source: SourceWorkloadIdentity}})
		if err != nil {
			t.Fatal(err)
		}
		if kcsb.ApplicationClientId != "envclient" || kcsb.AuthorityId != "envtenant" {
			t.Errorf("env defaults not applied: %+v", kcsb)
		}
	})
}

func TestFingerprint(t *testing.T) {
	a := Config{Endpoint: "https://x", Auth: Auth{Source: SourceSecret, ClientID: "a", ClientSecret: "s1", TenantID: "t"}}
	b := a
	b.Auth.ClientSecret = "s2"
	if a.Fingerprint() == b.Fingerprint() {
		t.Error("secret change must change the fingerprint")
	}
	first := a.Fingerprint()
	if first != a.Fingerprint() {
		t.Error("fingerprint must be stable")
	}
}

// New against the emulator endpoint must succeed without any token provider (S1).
func TestNewEmulatorNoAuth(t *testing.T) {
	c, err := New(Config{Endpoint: "http://localhost:8080", Auth: Auth{Source: SourceNone}}, Options{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.Endpoint() != "http://localhost:8080" {
		t.Error("endpoint")
	}
	if cl, ok := c.(Closer); ok {
		_ = cl.Close()
	}
	// http endpoint with credentials must be refused by the SDK.
	if _, err := New(Config{Endpoint: "http://localhost:8080", Auth: Auth{Source: SourceSecret, ClientID: "a", ClientSecret: "b", TenantID: "c"}}, Options{}); err == nil {
		t.Error("expected error for token provider over http")
	}
}
