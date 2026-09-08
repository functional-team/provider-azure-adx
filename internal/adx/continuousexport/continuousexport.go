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

// Package continuousexport holds the domain logic for the ContinuousExport
// kind: command rendering, parsing of ".show continuous-exports" and the
// comparison of desired against observed state. The query is free text and
// compared through the two stage normalization.
package continuousexport

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/normalize"
	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

// Section is the snapshot cache section for ".show continuous-exports".
const Section = "continuousexports"

// Observed is a continuous export as reported by the cluster. The column set
// (Name, ExternalTableName, Query, ForcedLatency, IntervalBetweenRuns,
// CursorScopedTables, ExportProperties, LastRunTime, IsDisabled,
// LastRunResult, ExportedTo, IsRunning) follows the docs; the key names inside
// ExportProperties are not verified against a cluster (spike S7) and are
// matched case-insensitively.
type Observed struct {
	Name                string
	ExternalTable       string
	Query               string
	ForcedLatency       *timespan.Ticks
	IntervalBetweenRuns *timespan.Ticks
	// CursorScopedTables in Kusto notation ['DB'].['T'].
	CursorScopedTables []string
	Properties         map[string]any
	IsDisabled         bool
	LastRunTime        string
	LastRunResult      string
	ExportedTo         string
	IsRunning          bool
}

// Desired is the normalized desired state.
type Desired struct {
	Name                string
	Database            string
	ExternalTable       string
	OverTables          []string // qualified ['DB'].['T']
	Query               string
	IntervalBetweenRuns timespan.Ticks
	ForcedLatency       *timespan.Ticks
	SizeLimit           *int64
	Distributed         *bool
	Distribution        *string
	DistributionKind    *string
	ParquetRowGroupSize *int64
	ManagedIdentity     *string
	Enabled             bool
}

// FromParams builds the desired state; invalid timespans are an error (the
// caller turns it into a Blocked/InvalidSpec condition).
func FromParams(name string, p v1alpha1.ContinuousExportParameters) (Desired, error) {
	d := Desired{Name: name, Database: p.Database, ExternalTable: strings.TrimSpace(p.ExternalTable), Query: p.Query,
		SizeLimit: p.SizeLimit, Distributed: p.Distributed, Distribution: p.Distribution, DistributionKind: p.DistributionKind,
		ParquetRowGroupSize: p.ParquetRowGroupSize, ManagedIdentity: p.ManagedIdentity, Enabled: p.Enabled == nil || *p.Enabled}
	if d.ExternalTable == "" {
		return d, fmt.Errorf("externalTable is empty (externalTableRef not resolved yet?)")
	}
	ticks, err := timespan.Parse(string(p.IntervalBetweenRuns))
	if err != nil {
		return d, fmt.Errorf("intervalBetweenRuns: %w", err)
	}
	d.IntervalBetweenRuns = ticks
	if p.ForcedLatency != nil {
		fl, err := timespan.Parse(string(*p.ForcedLatency))
		if err != nil {
			return d, fmt.Errorf("forcedLatency: %w", err)
		}
		d.ForcedLatency = &fl
	}
	for _, t := range p.OverTables {
		d.OverTables = append(d.OverTables, Qualify(p.Database, t))
	}
	return d, nil
}

// Qualify renders a table reference in the ['DB'].['T'] notation Kusto uses
// in CursorScopedTables; already qualified names are kept.
func Qualify(db, t string) string {
	t = strings.TrimSpace(t)
	if strings.HasPrefix(t, "[") && strings.Contains(t, "].[") {
		return t
	}
	return normalize.QualifiedTable(db, normalize.UnqualifiedTable(t))
}

// ParseRows parses ".show continuous-exports" / ".show continuous-export CE".
func ParseRows(res *kusto.Result) map[string]Observed { //nolint:gocyclo // one assignment per .show column.
	out := map[string]Observed{}
	if res == nil {
		return out
	}
	for _, r := range res.Rows() {
		name := r.String("Name")
		if name == "" {
			continue
		}
		o := Observed{Name: name, ExternalTable: r.String("ExternalTableName"), Query: r.String("Query"), LastRunTime: r.String("LastRunTime"), LastRunResult: r.String("LastRunResult"), ExportedTo: r.String("ExportedTo")}
		if t, ok := r.Timespan("ForcedLatency"); ok {
			o.ForcedLatency = &t
		}
		if t, ok := r.Timespan("IntervalBetweenRuns"); ok {
			o.IntervalBetweenRuns = &t
		}
		o.IsDisabled, _ = r.Bool("IsDisabled")
		o.IsRunning, _ = r.Bool("IsRunning")
		if raw := r.Dynamic("CursorScopedTables"); len(raw) > 0 {
			var arr []string
			if err := json.Unmarshal(raw, &arr); err == nil {
				for _, t := range arr {
					o.CursorScopedTables = append(o.CursorScopedTables, strings.TrimSpace(t))
				}
			}
		}
		if raw := r.Dynamic("ExportProperties"); len(raw) > 0 {
			var m map[string]any
			if err := json.Unmarshal(raw, &m); err == nil {
				o.Properties = m
			}
		}
		out[name] = o
	}
	return out
}

// ParseOne returns the single continuous export in res.
func ParseOne(res *kusto.Result, name string) (Observed, bool) {
	m := ParseRows(res)
	if o, ok := m[name]; ok {
		return o, true
	}
	if len(m) == 1 {
		for _, o := range m {
			return o, true
		}
	}
	return Observed{}, false
}

// ShowAll is ".show continuous-exports".
func ShowAll() cmd.Command { return cmd.New(".show continuous-exports") }

// Show is ".show continuous-export ['CE']".
func Show(name string) cmd.Command { return cmd.New(".show continuous-export ", cmd.Ident(name)) }

// BuildCreateOrAlter renders the create-or-alter command with the query in
// "<|" form.
func BuildCreateOrAlter(d Desired) cmd.Command {
	parts := []string{".create-or-alter continuous-export ", cmd.Ident(d.Name)}
	if len(d.OverTables) > 0 {
		items := make([]string, 0, len(d.OverTables))
		for _, t := range d.OverTables {
			items = append(items, cmd.Ident(normalize.UnqualifiedTable(t)))
		}
		parts = append(parts, " over ", cmd.List(items))
	}
	props := map[string]string{"intervalBetweenRuns": cmd.Timespan(d.IntervalBetweenRuns)}
	if d.ForcedLatency != nil {
		props["forcedLatency"] = cmd.Timespan(*d.ForcedLatency)
	}
	if d.SizeLimit != nil {
		props["sizeLimit"] = cmd.Int(*d.SizeLimit)
	}
	if d.Distributed != nil {
		props["distributed"] = cmd.Bool(*d.Distributed)
	}
	if d.Distribution != nil {
		props["distribution"] = cmd.Str(*d.Distribution)
	}
	if d.DistributionKind != nil {
		props["distributionKind"] = cmd.Str(*d.DistributionKind)
	}
	if d.ParquetRowGroupSize != nil {
		props["parquetRowGroupSize"] = cmd.Int(*d.ParquetRowGroupSize)
	}
	if d.ManagedIdentity != nil {
		props["managedIdentity"] = cmd.Str(*d.ManagedIdentity)
	}
	parts = append(parts, " to table ", cmd.Ident(d.ExternalTable), cmd.With(props), cmd.Pipe(d.Query))
	return cmd.New(parts...)
}

// BuildToggle renders ".enable" / ".disable continuous-export ['CE']".
func BuildToggle(name string, enabled bool) cmd.Command {
	verb := ".enable continuous-export "
	if !enabled {
		verb = ".disable continuous-export "
	}
	return cmd.New(verb, cmd.Ident(name))
}

// BuildDelete renders ".drop continuous-export ['CE']".
func BuildDelete(name string) cmd.Command {
	return cmd.New(".drop continuous-export ", cmd.Ident(name))
}

// Plan lists what an update has to do.
type Plan struct {
	// Alter is true when create-or-alter must be sent (structural drift).
	Alter bool
	// Toggle is set when the enabled state differs (true = enable).
	Toggle  *bool
	Reasons []string
}

// Empty reports whether nothing needs to be done.
func (p *Plan) Empty() bool { return !p.Alter && p.Toggle == nil }

// String summarizes the plan.
func (p *Plan) String() string { return strings.Join(p.Reasons, "; ") }

func (p *Plan) alter(r string) {
	p.Alter = true
	p.Reasons = append(p.Reasons, r)
}

// Diff compares the structural fields (everything except the query, which is
// compared through the hash mechanism by the caller).
func Diff(d Desired, o Observed) Plan { //nolint:gocyclo // A flat list of field checks.
	var p Plan
	if o.ExternalTable != "" && d.ExternalTable != o.ExternalTable {
		p.alter("externalTable " + o.ExternalTable + " != " + d.ExternalTable)
	}
	if o.IntervalBetweenRuns != nil && *o.IntervalBetweenRuns != d.IntervalBetweenRuns {
		p.alter("intervalBetweenRuns")
	}
	if d.ForcedLatency != nil && (o.ForcedLatency == nil || *o.ForcedLatency != *d.ForcedLatency) {
		p.alter("forcedLatency")
	}
	if len(d.OverTables) > 0 && !sameSet(d.OverTables, o.CursorScopedTables) {
		p.alter("overTables")
	}
	if d.SizeLimit != nil && !propEquals(o.Properties, "SizeLimit", strconv.FormatInt(*d.SizeLimit, 10)) {
		p.alter("sizeLimit")
	}
	if d.Distributed != nil && !propEquals(o.Properties, "Distributed", strconv.FormatBool(*d.Distributed)) {
		p.alter("distributed")
	}
	if d.Distribution != nil && !propEquals(o.Properties, "Distribution", *d.Distribution) {
		p.alter("distribution")
	}
	if d.DistributionKind != nil && !propEquals(o.Properties, "DistributionKind", *d.DistributionKind) {
		p.alter("distributionKind")
	}
	if d.ParquetRowGroupSize != nil && !propEquals(o.Properties, "ParquetRowGroupSize", strconv.FormatInt(*d.ParquetRowGroupSize, 10)) {
		p.alter("parquetRowGroupSize")
	}
	if d.ManagedIdentity != nil && !propEquals(o.Properties, "ManagedIdentity", *d.ManagedIdentity) {
		p.alter("managedIdentity")
	}
	if d.Enabled == o.IsDisabled {
		en := d.Enabled
		p.Toggle = &en
		if en {
			p.Reasons = append(p.Reasons, "enable")
		} else {
			p.Reasons = append(p.Reasons, "disable")
		}
	}
	return p
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	as := append([]string(nil), a...)
	bs := append([]string(nil), b...)
	for i := range as {
		as[i] = strings.ToLower(strings.TrimSpace(as[i]))
	}
	for i := range bs {
		bs[i] = strings.ToLower(strings.TrimSpace(bs[i]))
	}
	sort.Strings(as)
	sort.Strings(bs)
	for i := range as {
		if as[i] != bs[i] {
			return false
		}
	}
	return true
}

// propEquals compares a property case-insensitively (key and value). A
// missing property is treated as equal: the cluster default is unknown and
// must not cause a permanent drift.
func propEquals(m map[string]any, key, want string) bool {
	var got any
	found := false
	for k, v := range m {
		if strings.EqualFold(k, key) {
			got, found = v, v != nil
			break
		}
	}
	if !found {
		return true
	}
	var s string
	switch x := got.(type) {
	case string:
		s = x
	case bool:
		s = strconv.FormatBool(x)
	case float64:
		s = strconv.FormatFloat(x, 'f', -1, 64)
	default:
		b, err := json.Marshal(got)
		if err != nil {
			return false
		}
		s = string(b)
	}
	return strings.EqualFold(strings.TrimSpace(s), strings.TrimSpace(want))
}

// DesiredTexts is the normalized query.
func DesiredTexts(d Desired) []string { return []string{normalize.KQL(d.Query)} }

// ObservedTexts is the normalized query as reported.
func ObservedTexts(o Observed) []string { return []string{normalize.KQL(o.Query)} }

// Observation converts observed state into the status representation.
func Observation(o Observed) v1alpha1.ContinuousExportObservation {
	return v1alpha1.ContinuousExportObservation{ExternalTable: o.ExternalTable, Query: o.Query, IsDisabled: o.IsDisabled, LastRunTime: o.LastRunTime, LastRunResult: o.LastRunResult, ExportedTo: o.ExportedTo, IsRunning: o.IsRunning}
}
