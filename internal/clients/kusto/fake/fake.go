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

// Package fake provides an in-memory kusto.Client for unit tests. Handlers are
// matched in order against the command text; every call is recorded.
package fake

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	kustoerrors "github.com/Azure/azure-kusto-go/azkustodata/errors"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
)

// Call is one recorded Mgmt invocation.
type Call struct {
	Database string
	Command  string
}

// Handler answers commands whose text contains Match (or all commands if
// Match is empty).
type Handler struct {
	Match  string
	Result *kusto.Result
	Err    error
	// Fn, if set, computes the answer and takes precedence over Result/Err.
	Fn func(database string, c cmd.Command) (*kusto.Result, error)
}

// Client is a scripted kusto.Client.
type Client struct {
	mu       sync.Mutex
	URI      string
	Handlers []Handler
	Calls    []Call
}

var _ kusto.Client = &Client{}

// New returns a fake client with the given endpoint.
func New(endpoint string, handlers ...Handler) *Client {
	return &Client{URI: endpoint, Handlers: handlers}
}

// Endpoint implements kusto.Client.
func (f *Client) Endpoint() string { return f.URI }

// Mgmt implements kusto.Client.
func (f *Client) Mgmt(_ context.Context, database string, c cmd.Command) (*kusto.Result, error) {
	f.mu.Lock()
	f.Calls = append(f.Calls, Call{Database: database, Command: c.String()})
	handlers := append([]Handler(nil), f.Handlers...)
	f.mu.Unlock()
	for _, h := range handlers {
		if h.Match == "" || strings.Contains(c.String(), h.Match) {
			if h.Fn != nil {
				return h.Fn(database, c)
			}
			if h.Err != nil {
				return nil, h.Err
			}
			if h.Result == nil {
				return kusto.NewResult(), nil
			}
			return h.Result, nil
		}
	}
	return nil, fmt.Errorf("fake kusto: unexpected command %q", c.String())
}

// On appends a handler.
func (f *Client) On(match string, res *kusto.Result, err error) *Client {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Handlers = append(f.Handlers, Handler{Match: match, Result: res, Err: err})
	return f
}

// OnFn appends a function handler.
func (f *Client) OnFn(match string, fn func(database string, c cmd.Command) (*kusto.Result, error)) *Client {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Handlers = append(f.Handlers, Handler{Match: match, Fn: fn})
	return f
}

// Commands returns the recorded command texts.
func (f *Client) Commands() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, 0, len(f.Calls))
	for _, c := range f.Calls {
		out = append(out, c.Command)
	}
	return out
}

// Count returns how many recorded commands contain substr.
func (f *Client) Count(substr string) int {
	n := 0
	for _, c := range f.Commands() {
		if strings.Contains(c, substr) {
			n++
		}
	}
	return n
}

// Reset clears recorded calls.
func (f *Client) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Calls = nil
}

// HTTPError builds an SDK HTTP error from a REST body the way the SDK does;
// useful for scripting NotFound/AlreadyExists/Throttled responses.
func HTTPError(status int, body string) error {
	return kustoerrors.HTTP(kustoerrors.OpMgmt, http.StatusText(status), status, io.NopCloser(strings.NewReader(body)), "Mgmt failed")
}
