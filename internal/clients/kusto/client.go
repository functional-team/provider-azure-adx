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

// Package kusto wraps the Azure Kusto Go SDK: authentication from
// ProviderConfig credentials, a client pool per ProviderConfig, rate limiting
// per cluster, metrics and a small result model that the domain packages
// parse without depending on the SDK.
package kusto

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-kusto-go/azkustodata"
	"github.com/Azure/azure-kusto-go/azkustodata/query"
	"golang.org/x/time/rate"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

// Client executes management commands against one cluster.
type Client interface {
	// Mgmt runs a management command in the context of database ("" for the
	// cluster level).
	Mgmt(ctx context.Context, database string, c cmd.Command) (*Result, error)
	// Endpoint is the cluster URI.
	Endpoint() string
}

// Closer is implemented by clients that hold connections.
type Closer interface {
	Close() error
}

// CommandError wraps an SDK error with the (redacted) command that failed. It
// unwraps to the original error so kerrors.Classify keeps working.
type CommandError struct {
	Command string
	Err     error
}

func (e *CommandError) Error() string {
	return e.Command + ": " + kerrors.Message(e.Err)
}

// Unwrap returns the SDK error.
func (e *CommandError) Unwrap() error { return e.Err }

// Op labels a command for metrics and logs. Kind is the managed resource kind,
// Op is observe|create|update|delete.
type Op struct {
	Kind string
	Op   string
}

type opKey struct{}

// WithOp attaches metric labels to ctx.
func WithOp(ctx context.Context, kind, op string) context.Context {
	return context.WithValue(ctx, opKey{}, Op{Kind: kind, Op: op})
}

// OpFromContext returns the labels attached with WithOp.
func OpFromContext(ctx context.Context) Op {
	if o, ok := ctx.Value(opKey{}).(Op); ok {
		return o
	}
	return Op{Kind: "unknown", Op: "unknown"}
}

// Options tune a client.
type Options struct {
	// CommandsPerSecond limits the sustained command rate per cluster.
	CommandsPerSecond float64
	// MaxInflight limits concurrent commands per cluster.
	MaxInflight int
	// Logger receives debug logs for every command (redacted).
	Logger logging.Logger
}

func (o Options) withDefaults() Options {
	if o.CommandsPerSecond <= 0 {
		o.CommandsPerSecond = 5
	}
	if o.MaxInflight <= 0 {
		o.MaxInflight = 4
	}
	if o.Logger == nil {
		o.Logger = logging.NewNopLogger()
	}
	return o
}

type sdkClient struct {
	kc       *azkustodata.Client
	endpoint string
	lim      *rate.Limiter
	sem      chan struct{}
	log      logging.Logger
}

// New builds a client for cfg.
func New(cfg Config, o Options) (Client, error) {
	kcsb, err := kcsbFor(cfg)
	if err != nil {
		return nil, err
	}
	kc, err := azkustodata.New(kcsb)
	if err != nil {
		return nil, fmt.Errorf("cannot create kusto client: %w", err)
	}
	return wrap(kc, cfg.Endpoint, o), nil
}

func wrap(kc *azkustodata.Client, endpoint string, o Options) *sdkClient {
	o = o.withDefaults()
	return &sdkClient{
		kc:       kc,
		endpoint: endpoint,
		lim:      rate.NewLimiter(rate.Limit(o.CommandsPerSecond), int(o.CommandsPerSecond)+1),
		sem:      make(chan struct{}, o.MaxInflight),
		log:      o.Logger.WithValues("endpoint", endpoint),
	}
}

func (c *sdkClient) Endpoint() string { return c.endpoint }

func (c *sdkClient) Close() error { return c.kc.Close() }

func (c *sdkClient) Mgmt(ctx context.Context, database string, command cmd.Command) (*Result, error) {
	if command.IsZero() {
		return nil, fmt.Errorf("empty command")
	}
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if err := c.lim.Wait(ctx); err != nil {
		return nil, err
	}
	op := OpFromContext(ctx)
	start := time.Now()
	ds, err := c.kc.Mgmt(ctx, database, command.Statement())
	dur := time.Since(start)
	if err != nil {
		class := kerrors.Classify(err)
		observe(op, class.String(), dur)
		c.log.Debug("kusto command failed", "database", database, "command", short(command.Redacted()), "class", class.String(), "duration", dur.String(), "error", kerrors.Message(err))
		return nil, &CommandError{Command: short(command.Redacted()), Err: err}
	}
	observe(op, "ok", dur)
	c.log.Debug("kusto command ok", "database", database, "command", short(command.Redacted()), "duration", dur.String())
	return fromDataset(ds), nil
}

func short(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 120 {
		return s[:117] + "..."
	}
	return s
}

// Result is the tables of one management command response.
type Result struct {
	Tables []Table
}

// Table is one result table.
type Table struct {
	Name    string
	Kind    string
	Primary bool
	Columns []string
	Rows    []Row
}

// Row is one result row with access by column name.
type Row struct {
	cols map[string]int
	vals []any
}

// NewTable builds a table from column names and row values; used by tests and
// fakes. Values keep their Go type (string, bool, int64, float64, time.Time,
// timespan.Ticks, json.RawMessage).
func NewTable(name string, cols []string, rows ...[]any) Table {
	idx := index(cols)
	t := Table{Name: name, Primary: true, Columns: cols}
	for _, r := range rows {
		t.Rows = append(t.Rows, Row{cols: idx, vals: r})
	}
	return t
}

// NewResult builds a result from tables.
func NewResult(tables ...Table) *Result { return &Result{Tables: tables} }

// NewRow builds a row; used by tests.
func NewRow(cols []string, vals ...any) Row { return Row{cols: index(cols), vals: vals} }

func index(cols []string) map[string]int {
	idx := make(map[string]int, len(cols))
	for i, c := range cols {
		idx[c] = i
	}
	return idx
}

// Primary returns the primary result table (or the first table).
func (r *Result) Primary() *Table {
	if r == nil {
		return nil
	}
	for i := range r.Tables {
		if r.Tables[i].Primary {
			return &r.Tables[i]
		}
	}
	if len(r.Tables) > 0 {
		return &r.Tables[0]
	}
	return nil
}

// Rows returns the rows of the primary table.
func (r *Result) Rows() []Row {
	if t := r.Primary(); t != nil {
		return t.Rows
	}
	return nil
}

// Has reports whether the row has a column.
func (r Row) Has(col string) bool {
	_, ok := r.cols[col]
	return ok
}

func (r Row) raw(col string) (any, bool) {
	i, ok := r.cols[col]
	if !ok || i >= len(r.vals) {
		return nil, false
	}
	return r.vals[i], true
}

// String returns the column as text ("" if missing or null). Non-string
// values are rendered with fmt or as JSON for dynamic columns.
func (r Row) String(col string) string { //nolint:gocyclo // type switch over the value kinds the SDK can return.
	v, ok := r.raw(col)
	if !ok || v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case *string:
		if x == nil {
			return ""
		}
		return *x
	case json.RawMessage:
		return string(x)
	case []byte:
		return string(x)
	case timespan.Ticks:
		return x.DotNet()
	case time.Time:
		return x.UTC().Format(time.RFC3339Nano)
	case fmt.Stringer:
		return x.String()
	}
	return fmt.Sprint(v)
}

// Bool returns the column as bool; ok is false if missing, null or not a bool.
func (r Row) Bool(col string) (val, ok bool) {
	v, found := r.raw(col)
	if !found || v == nil {
		return false, false
	}
	switch x := v.(type) {
	case bool:
		return x, true
	case *bool:
		if x == nil {
			return false, false
		}
		return *x, true
	case string:
		switch strings.ToLower(strings.TrimSpace(x)) {
		case "true":
			return true, true
		case "false":
			return false, true
		}
	}
	return false, false
}

// Int64 returns the column as int64.
func (r Row) Int64(col string) (int64, bool) { //nolint:gocyclo // type switch over the value kinds the SDK can return.
	v, found := r.raw(col)
	if !found || v == nil {
		return 0, false
	}
	switch x := v.(type) {
	case int64:
		return x, true
	case int32:
		return int64(x), true
	case int:
		return int64(x), true
	case float64:
		return int64(x), true
	case *int64:
		if x != nil {
			return *x, true
		}
	case *int32:
		if x != nil {
			return int64(*x), true
		}
	case string:
		var n int64
		if _, err := fmt.Sscan(x, &n); err == nil {
			return n, true
		}
	}
	return 0, false
}

// Timespan returns the column as ticks.
func (r Row) Timespan(col string) (timespan.Ticks, bool) {
	v, found := r.raw(col)
	if !found || v == nil {
		return 0, false
	}
	switch x := v.(type) {
	case timespan.Ticks:
		return x, true
	case time.Duration:
		return timespan.FromDuration(x), true
	case *time.Duration:
		if x != nil {
			return timespan.FromDuration(*x), true
		}
	case string:
		if t, err := timespan.Parse(x); err == nil {
			return t, true
		}
	}
	return 0, false
}

// Time returns the column as time.
func (r Row) Time(col string) (time.Time, bool) {
	v, found := r.raw(col)
	if !found || v == nil {
		return time.Time{}, false
	}
	switch x := v.(type) {
	case time.Time:
		return x, true
	case *time.Time:
		if x != nil {
			return *x, true
		}
	case string:
		if t, err := time.Parse(time.RFC3339Nano, x); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// Dynamic returns the column as raw JSON (nil if missing or null).
func (r Row) Dynamic(col string) json.RawMessage {
	v, found := r.raw(col)
	if !found || v == nil {
		return nil
	}
	switch x := v.(type) {
	case json.RawMessage:
		return x
	case []byte:
		return json.RawMessage(x)
	case string:
		s := strings.TrimSpace(x)
		if s == "" || s == "null" {
			return nil
		}
		return json.RawMessage(s)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return b
}

// fromDataset converts an SDK dataset into the result model.
func fromDataset(ds query.Dataset) *Result {
	res := &Result{}
	if ds == nil {
		return res
	}
	for _, t := range ds.Tables() {
		cols := make([]string, 0, len(t.Columns()))
		for _, c := range t.Columns() {
			cols = append(cols, c.Name())
		}
		tbl := Table{Name: t.Name(), Kind: t.Kind(), Primary: t.IsPrimaryResult(), Columns: cols}
		idx := index(cols)
		for _, row := range t.Rows() {
			vals := make([]any, len(cols))
			for i, v := range row.Values() {
				if i >= len(vals) {
					break
				}
				vals[i] = convertValue(v)
			}
			tbl.Rows = append(tbl.Rows, Row{cols: idx, vals: vals})
		}
		res.Tables = append(res.Tables, tbl)
	}
	return res
}

// convertValue flattens SDK values (mostly pointers to Go types) into plain
// Go values understood by the Row accessors.
func convertValue(v interface{ GetValue() interface{} }) any { //nolint:gocyclo // type switch over the value kinds the SDK can return.
	if v == nil {
		return nil
	}
	raw := v.GetValue()
	switch x := raw.(type) {
	case nil:
		return nil
	case *bool:
		if x == nil {
			return nil
		}
		return *x
	case *int32:
		if x == nil {
			return nil
		}
		return *x
	case *int64:
		if x == nil {
			return nil
		}
		return *x
	case *float64:
		if x == nil {
			return nil
		}
		return *x
	case *time.Time:
		if x == nil {
			return nil
		}
		return *x
	case *time.Duration:
		if x == nil {
			return nil
		}
		return timespan.FromDuration(*x)
	case time.Duration:
		return timespan.FromDuration(x)
	case *string:
		if x == nil {
			return nil
		}
		return *x
	case []byte:
		if len(x) == 0 {
			return nil
		}
		return json.RawMessage(x)
	}
	return raw
}
