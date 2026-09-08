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
	"errors"
	"testing"
	"time"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
)

type stubClient struct {
	endpoint string
	closed   bool
}

func (s *stubClient) Mgmt(context.Context, string, cmd.Command) (*Result, error) {
	return NewResult(), nil
}
func (s *stubClient) Endpoint() string { return s.endpoint }
func (s *stubClient) Close() error     { s.closed = true; return nil }

func TestPool(t *testing.T) {
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	created := 0
	p := NewPool(30*time.Minute, func(cfg Config) (Client, error) {
		created++
		if cfg.Endpoint == "bad" {
			return nil, errors.New("bad endpoint")
		}
		return &stubClient{endpoint: cfg.Endpoint}, nil
	})
	p.now = func() time.Time { return now }

	cfg := Config{Endpoint: "https://a", Auth: Auth{Source: SourceSecret, ClientID: "a", ClientSecret: "s", TenantID: "t"}}
	c1, err := p.Get("ProviderConfig/ns/pc", cfg)
	if err != nil {
		t.Fatal(err)
	}
	c2, _ := p.Get("ProviderConfig/ns/pc", cfg)
	if c1 != c2 || created != 1 {
		t.Errorf("expected reuse, created=%d", created)
	}

	// Changed secret -> rebuilt, old closed.
	cfg2 := cfg
	cfg2.Auth.ClientSecret = "rotated"
	c3, _ := p.Get("ProviderConfig/ns/pc", cfg2)
	if c3 == c1 || created != 2 || !c1.(*stubClient).closed {
		t.Errorf("expected rebuild on fingerprint change, created=%d closed=%v", created, c1.(*stubClient).closed)
	}

	// Error path leaves the pool untouched.
	if _, err := p.Get("other", Config{Endpoint: "bad", Auth: Auth{Source: SourceNone}}); err == nil {
		t.Error("expected error")
	}
	if p.Len() != 1 {
		t.Errorf("Len = %d", p.Len())
	}

	// Idle eviction.
	now = now.Add(31 * time.Minute)
	p.Get("fresh", Config{Endpoint: "https://b", Auth: Auth{Source: SourceNone}}) //nolint:errcheck // exercised above
	if p.Len() != 1 || !c3.(*stubClient).closed {
		t.Errorf("expected idle eviction, len=%d", p.Len())
	}
	p.Remove("fresh")
	if p.Len() != 0 {
		t.Error("Remove")
	}
}
