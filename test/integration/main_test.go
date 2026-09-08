//go:build integration

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

// Package integration runs the provider's external clients against the Kusto
// emulator (or a development cluster, see test/emulator). Run with
// go test -tags integration ./test/integration/...
package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/test/emulator"
)

var (
	emu *emulator.Emulator
	kc  *counting
	db  string
)

// counting records every command so tests can assert how many commands a
// scenario needs (batch observe, no update loops).
type counting struct {
	kusto.Client
	mu   sync.Mutex
	cmds []string
}

func (c *counting) Mgmt(ctx context.Context, database string, command cmd.Command) (*kusto.Result, error) {
	c.mu.Lock()
	c.cmds = append(c.cmds, command.String())
	c.mu.Unlock()
	return c.Client.Mgmt(ctx, database, command)
}

func (c *counting) count(sub string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for _, s := range c.cmds {
		if strings.Contains(s, sub) {
			n++
		}
	}
	return n
}

func (c *counting) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cmds = nil
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	e, err := emulator.Start(ctx)
	if err != nil {
		if emulator.IsUnavailable(err) {
			fmt.Println("SKIP integration tests:", err)
			os.Exit(0)
		}
		fmt.Println("cannot start emulator:", err)
		os.Exit(1)
	}
	emu = e
	raw, err := emu.Client(kusto.Options{CommandsPerSecond: 20, MaxInflight: 8})
	if err != nil {
		fmt.Println("cannot create client:", err)
		_ = emu.Terminate(ctx)
		os.Exit(1)
	}
	kc = &counting{Client: raw}
	db = emu.Database("ProviderIT")
	if err := emu.CreateDatabase(ctx, db); err != nil {
		fmt.Println("cannot create database:", err)
		_ = emu.Terminate(ctx)
		os.Exit(1)
	}
	code := m.Run()
	_ = emu.Terminate(context.Background())
	os.Exit(code)
}

// run executes a raw command in the test database.
func run(t *testing.T, text string) *kusto.Result {
	t.Helper()
	res, err := kc.Mgmt(context.Background(), db, cmd.New(text))
	if err != nil {
		t.Fatalf("%s: %v", text, err)
	}
	return res
}
