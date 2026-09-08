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

// Package materializedview holds the domain logic for the MaterializedView
// kind: command builders (including the async backfill create and the
// operation polling commands), tolerant parsing of ".show materialized-views"
// and the diff between desired and observed state.
package materializedview

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/normalize"
	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

// Section is the snapshot cache section for ".show materialized-views".
const Section = "materializedviews"

// Operation states reported by ".show operations".
const (
	StateInProgress         = "InProgress"
	StateScheduled          = "Scheduled"
	StateCompleted          = "Completed"
	StateFailed             = "Failed"
	StateThrottled          = "Throttled"
	StateAbandoned          = "Abandoned"
	StateCanceled           = "Canceled"
	StateBadInput           = "BadInput"
	StatePartiallySucceeded = "PartiallySucceeded"
)

// Observed is a materialized view as reported by the cluster.
type Observed struct {
	Name              string
	SourceTable       string
	Query             string
	IsEnabled         bool
	IsHealthy         bool
	Folder            string
	Docstring         string
	AutoUpdateSchema  *bool
	EffectiveDateTime string
	Lookback          *timespan.Ticks
	LookbackColumn    string
	LastRunResult     string
}

// Desired is the normalized desired state.
type Desired struct {
	Name                      string
	SourceTable               string
	SourceMaterializedView    string
	Query                     string
	Backfill                  bool
	EffectiveDateTime         *string
	UpdateExtentsCreationTime *bool
	Lookback                  *timespan.Ticks
	LookbackColumn            *string
	AutoUpdateSchema          *bool
	DimensionTables           []string
	Folder                    *string
	Docstring                 *string
	AllowWithoutRLS           *bool
	MaxSourceRecords          *int64
	Concurrency               *int64
	Enabled                   bool
}

var datetimeRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}([T ]\d{2}:\d{2}(:\d{2}(\.\d+)?)?(Z|[+-]\d{2}:\d{2})?)?$`)

// FromParams builds the desired state from the spec. It fails on unparseable
// timespans/datetimes and when no source is set (reference not resolved yet).
func FromParams(name string, p v1alpha1.MaterializedViewParameters) (Desired, error) { //nolint:gocyclo // one validation per optional field.
	d := Desired{
		Name:                      name,
		Query:                     p.Query,
		Backfill:                  p.Backfill != nil && *p.Backfill,
		EffectiveDateTime:         p.EffectiveDateTime,
		UpdateExtentsCreationTime: p.UpdateExtentsCreationTime,
		LookbackColumn:            p.LookbackColumn,
		AutoUpdateSchema:          p.AutoUpdateSchema,
		DimensionTables:           p.DimensionTables,
		Folder:                    p.Folder,
		Docstring:                 p.Docstring,
		AllowWithoutRLS:           p.AllowMaterializedViewsWithoutRowLevelSecurity,
		MaxSourceRecords:          p.MaxSourceRecordsForSingleIngest,
		Concurrency:               p.Concurrency,
		Enabled:                   p.Enabled == nil || *p.Enabled,
	}
	if p.SourceTable != nil {
		d.SourceTable = strings.TrimSpace(*p.SourceTable)
	}
	if p.SourceMaterializedView != nil {
		d.SourceMaterializedView = strings.TrimSpace(*p.SourceMaterializedView)
	}
	if d.SourceTable == "" && d.SourceMaterializedView == "" {
		return d, fmt.Errorf("sourceTable is empty; set it or wait for sourceTableRef to resolve")
	}
	if p.Lookback != nil {
		t, err := timespan.Parse(string(*p.Lookback))
		if err != nil {
			return d, fmt.Errorf("lookback: %w", err)
		}
		d.Lookback = &t
	}
	if p.EffectiveDateTime != nil && !datetimeRe.MatchString(strings.TrimSpace(*p.EffectiveDateTime)) {
		return d, fmt.Errorf("effectiveDateTime %q is not an ISO 8601 datetime", *p.EffectiveDateTime)
	}
	return d, nil
}

// ParseRows parses the rows of ".show materialized-views" / ".show
// materialized-view" by column name; unknown columns are ignored.
func ParseRows(res *kusto.Result) map[string]Observed {
	out := map[string]Observed{}
	if res == nil {
		return out
	}
	for _, r := range res.Rows() {
		name := r.String("Name")
		if name == "" {
			continue
		}
		o := Observed{
			Name:              name,
			SourceTable:       r.String("SourceTable"),
			Query:             r.String("Query"),
			IsEnabled:         true,
			Folder:            r.String("Folder"),
			Docstring:         r.String("DocString"),
			EffectiveDateTime: r.String("EffectiveDateTime"),
			LookbackColumn:    r.String("LookbackColumn"),
			LastRunResult:     r.String("LastRunResult"),
		}
		if b, ok := r.Bool("IsEnabled"); ok {
			o.IsEnabled = b
		}
		if b, ok := r.Bool("IsHealthy"); ok {
			o.IsHealthy = b
		}
		if b, ok := r.Bool("AutoUpdateSchema"); ok {
			o.AutoUpdateSchema = &b
		}
		if t, ok := r.Timespan("Lookback"); ok {
			o.Lookback = &t
		}
		out[name] = o
	}
	return out
}

// ParseOne returns the view named name from res (or the single row).
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

// ShowAll is ".show materialized-views".
func ShowAll() cmd.Command { return cmd.New(".show materialized-views") }

// ShowOne is ".show materialized-view ['MV']".
func ShowOne(name string) cmd.Command {
	return cmd.New(".show materialized-view ", cmd.Ident(name))
}

func source(d Desired) string {
	if d.SourceMaterializedView != "" {
		return " on materialized-view " + cmd.Ident(d.SourceMaterializedView)
	}
	return " on table " + cmd.Ident(d.SourceTable)
}

// alterProps are the properties accepted by both create and alter.
func alterProps(d Desired) map[string]string {
	props := map[string]string{}
	if d.Lookback != nil {
		props["lookback"] = cmd.Timespan(*d.Lookback)
	}
	if d.LookbackColumn != nil {
		props["lookback_column"] = cmd.Str(*d.LookbackColumn)
	}
	if d.AutoUpdateSchema != nil {
		props["autoUpdateSchema"] = cmd.Bool(*d.AutoUpdateSchema)
	}
	if len(d.DimensionTables) > 0 {
		items := make([]string, 0, len(d.DimensionTables))
		for _, t := range d.DimensionTables {
			items = append(items, cmd.Str(t))
		}
		props["dimensionTables"] = "dynamic([" + strings.Join(items, ", ") + "])"
	}
	if d.Folder != nil {
		props["folder"] = cmd.Str(*d.Folder)
	}
	if d.Docstring != nil {
		props["docString"] = cmd.Str(*d.Docstring)
	}
	return props
}

// BuildCreate renders ".create [async] ifnotexists materialized-view with (...)
// ['MV'] on table ['Src'] { query }". async is used only with backfill.
func BuildCreate(d Desired) cmd.Command {
	props := alterProps(d)
	if d.Backfill {
		props["backfill"] = "true"
	}
	if d.EffectiveDateTime != nil {
		props["effectiveDateTime"] = "datetime(" + strings.TrimSpace(*d.EffectiveDateTime) + ")"
	}
	if d.UpdateExtentsCreationTime != nil {
		props["updateExtentsCreationTime"] = cmd.Bool(*d.UpdateExtentsCreationTime)
	}
	if d.AllowWithoutRLS != nil {
		props["allowMaterializedViewsWithoutRowLevelSecurity"] = cmd.Bool(*d.AllowWithoutRLS)
	}
	if d.MaxSourceRecords != nil {
		props["MaxSourceRecordsForSingleIngest"] = cmd.Int(*d.MaxSourceRecords)
	}
	if d.Concurrency != nil {
		props["Concurrency"] = cmd.Int(*d.Concurrency)
	}
	verb := ".create ifnotexists materialized-view"
	if d.Backfill {
		verb = ".create async ifnotexists materialized-view"
	}
	return cmd.New(verb, cmd.With(props), " ", cmd.Ident(d.Name), source(d), " ", cmd.Body(d.Query))
}

// BuildAlter renders ".alter materialized-view with (...) ['MV'] on table ['Src'] { query }".
func BuildAlter(d Desired) cmd.Command {
	return cmd.New(".alter materialized-view", cmd.With(alterProps(d)), " ", cmd.Ident(d.Name), source(d), " ", cmd.Body(d.Query))
}

// BuildEnable renders ".enable materialized-view ['MV']" or the disable form.
func BuildEnable(name string, enabled bool) cmd.Command {
	if enabled {
		return cmd.New(".enable materialized-view ", cmd.Ident(name))
	}
	return cmd.New(".disable materialized-view ", cmd.Ident(name))
}

// BuildDelete renders ".drop materialized-view ['MV'] ifexists".
func BuildDelete(name string) cmd.Command {
	return cmd.New(".drop materialized-view ", cmd.Ident(name), " ifexists")
}

var operationIDRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ValidOperationID reports whether id looks like a Kusto operation GUID. Only
// such ids are ever interpolated into commands.
func ValidOperationID(id string) bool { return operationIDRe.MatchString(id) }

// ShowOperation is ".show operations <id>" (recent operations, < 6 hours).
func ShowOperation(id string) (cmd.Command, error) {
	if !ValidOperationID(id) {
		return cmd.Command{}, fmt.Errorf("invalid operation id %q", id)
	}
	return cmd.New(".show operations ", id), nil
}

// ShowOperationHistoric is the log form ".show operations | where
// tostring(OperationId) == '<id>'" used when the direct lookup returns no row.
func ShowOperationHistoric(id string) (cmd.Command, error) {
	if !ValidOperationID(id) {
		return cmd.Command{}, fmt.Errorf("invalid operation id %q", id)
	}
	return cmd.New(".show operations | where tostring(OperationId) == ", cmd.Str(id)), nil
}

// Operation is the state of an async operation.
type Operation struct {
	ID     string
	State  string
	Status string
}

// Running reports whether the operation is still going.
func (o Operation) Running() bool {
	return o.State == StateInProgress || o.State == StateScheduled
}

// ParseOperation parses the first row of a ".show operations" result.
func ParseOperation(res *kusto.Result) (Operation, bool) {
	rows := res.Rows()
	if len(rows) == 0 {
		return Operation{}, false
	}
	r := rows[0]
	return Operation{ID: r.String("OperationId"), State: r.String("State"), Status: r.String("Status")}, true
}

// ParseOperationID extracts the OperationId from the result of an async command.
func ParseOperationID(res *kusto.Result) string {
	rows := res.Rows()
	if len(rows) == 0 {
		return ""
	}
	return strings.TrimSpace(rows[0].String("OperationId"))
}

// StepKind names an update step.
type StepKind string

// Update steps.
const (
	StepAlter   StepKind = "Alter"
	StepEnable  StepKind = "Enable"
	StepDisable StepKind = "Disable"
)

// Step is one command of an update plan.
type Step struct {
	Kind    StepKind
	Command cmd.Command
}

// Plan is the ordered list of commands needed to reach the desired state.
type Plan struct {
	Steps []Step
}

// Empty reports whether nothing needs to be sent.
func (p Plan) Empty() bool { return len(p.Steps) == 0 }

// String summarizes the plan.
func (p Plan) String() string {
	kinds := make([]string, 0, len(p.Steps))
	for _, s := range p.Steps {
		kinds = append(kinds, string(s.Kind))
	}
	return strings.Join(kinds, ",")
}

// Diff compares desired and observed state. queryEqual is the result of the
// two stage text comparison of the query (the caller has the hashes). Only
// fields set in the spec are compared; dimensionTables and the create-only
// properties are write-only.
func Diff(d Desired, o Observed, queryEqual bool) Plan {
	var p Plan
	if !queryEqual || !propsUpToDate(d, o) {
		p.Steps = append(p.Steps, Step{Kind: StepAlter, Command: BuildAlter(d)})
	}
	if d.Enabled != o.IsEnabled {
		kind := StepEnable
		if !d.Enabled {
			kind = StepDisable
		}
		p.Steps = append(p.Steps, Step{Kind: kind, Command: BuildEnable(d.Name, d.Enabled)})
	}
	return p
}

func propsUpToDate(d Desired, o Observed) bool { //nolint:gocyclo // a flat list of field checks.
	if d.Lookback != nil && (o.Lookback == nil || *o.Lookback != *d.Lookback) {
		return false
	}
	if d.LookbackColumn != nil && strings.TrimSpace(*d.LookbackColumn) != strings.TrimSpace(o.LookbackColumn) {
		return false
	}
	if d.AutoUpdateSchema != nil && (o.AutoUpdateSchema == nil || *o.AutoUpdateSchema != *d.AutoUpdateSchema) {
		return false
	}
	if d.Folder != nil && strings.TrimSpace(*d.Folder) != strings.TrimSpace(o.Folder) {
		return false
	}
	if d.Docstring != nil && strings.TrimSpace(*d.Docstring) != strings.TrimSpace(o.Docstring) {
		return false
	}
	return true
}

// DesiredTexts returns the stage 1 normalized free text fields.
func DesiredTexts(d Desired) []string { return []string{normalize.KQL(d.Query)} }

// ObservedTexts returns the stage 1 normalized free text fields.
func ObservedTexts(o Observed) []string { return []string{normalize.KQL(o.Query)} }

// Observation converts observed state into the status representation.
func Observation(o Observed) v1alpha1.MaterializedViewObservation {
	enabled, healthy := o.IsEnabled, o.IsHealthy
	obs := v1alpha1.MaterializedViewObservation{
		SourceTable:       o.SourceTable,
		Query:             o.Query,
		IsEnabled:         &enabled,
		IsHealthy:         &healthy,
		Folder:            o.Folder,
		Docstring:         o.Docstring,
		EffectiveDateTime: o.EffectiveDateTime,
		LastRunResult:     o.LastRunResult,
	}
	if o.Lookback != nil {
		obs.Lookback = o.Lookback.DotNet()
	}
	return obs
}
