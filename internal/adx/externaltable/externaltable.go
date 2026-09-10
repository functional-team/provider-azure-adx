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

// Package externaltable holds the domain logic for the ExternalTable kind:
// command rendering (connection strings always obfuscated), parsing of
// ".show external tables" and the comparison of desired against observed
// state. Secrets in connection strings are write-only and tracked by hash.
package externaltable

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/adx/table"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/normalize"
)

// Section is the snapshot cache section for ".show external tables".
const Section = "externaltables"

// Normalized external table kinds.
const (
	kindStorage = "storage"
	kindDelta   = "delta"
)

// Column is a normalized column (canonical type).
type Column struct {
	Name string
	Type string
}

// Observed is an external table as reported by the cluster.
//
// Column shapes of ".show external tables" (TableName, TableType, Folder,
// DocString, Properties, ConnectionStrings, Partitions, PathFormat) are taken
// from the docs; the exact JSON layout of Properties, ConnectionStrings and
// Partitions is not verified against a cluster (spike S7), which is why the
// parser is tolerant and partitions are compared through hashes only.
type Observed struct {
	Name string
	// Kind normalized to "storage", "delta" or the lower-cased TableType.
	Kind string
	// Columns are only set when loaded via cslschema (ColumnsLoaded).
	Columns       []Column
	ColumnsLoaded bool
	Folder        string
	DocString     string
	// Properties is the parsed Properties JSON (keys as reported).
	Properties map[string]any
	// ConnectionStringURIs are the parts before the first ';' of each
	// connection string; Kusto masks the rest.
	ConnectionStringURIs []string
	// Partitions is the compacted Partitions JSON ("" when none).
	Partitions string
	PathFormat string
}

// Desired is the normalized desired state. ConnectionStrings hold the
// resolved secrets and must never be logged or put into status.
type Desired struct {
	Name              string
	Kind              string
	Columns           []Column
	PartitionBy       string
	PathFormat        string
	DataFormat        string
	ConnectionStrings []string
	Properties        *v1alpha1.ExternalTableProperties
}

// FromParams builds the desired state from the spec and the resolved
// connection strings.
func FromParams(name string, p v1alpha1.ExternalTableParameters, connectionStrings []string) Desired {
	d := Desired{Name: name, Kind: normalizeKind(string(p.Kind)), ConnectionStrings: connectionStrings, Properties: p.Properties}
	if d.Kind == "" {
		d.Kind = kindStorage
	}
	for _, c := range p.Columns {
		d.Columns = append(d.Columns, Column{Name: c.Name, Type: normalize.ColumnType(string(c.Type))})
	}
	d.PartitionBy = normalize.Trim(p.PartitionBy)
	d.PathFormat = normalize.Trim(p.PathFormat)
	d.DataFormat = strings.ToLower(normalize.Trim(p.DataFormat))
	return d
}

// normalizeKind maps spec kinds and .show TableType values onto one vocabulary.
func normalizeKind(k string) string {
	switch strings.ToLower(strings.TrimSpace(k)) {
	case kindStorage, "blob", "adl", "adls":
		return kindStorage
	case kindDelta:
		return kindDelta
	case "":
		return ""
	default:
		return strings.ToLower(strings.TrimSpace(k))
	}
}

// URIOf returns the storage location of a connection string: everything
// before the credential. A credential follows either a ';' (an account key or
// managed_identity=...) or a '?' (a SAS token, which is a query string).
//
// Both have to go, because neither is readable back: the service masks a SAS
// as "******" (observed 2026-09-10, e2e run 34439601775):
//
//	sent:      https://acct.blob.core.windows.net/exports?se=...&sig=...
//	read back: https://acct.blob.core.windows.net/exports?******
//
// Comparing those made the external table differ on every observe, so it was
// rewritten once per poll interval, forever. The location alone identifies the
// storage; whether the credential still matches is tracked by SecretHash.
func URIOf(cs string) string {
	cs = strings.TrimSpace(cs)
	if i := strings.IndexAny(cs, ";?"); i >= 0 {
		cs = cs[:i]
	}
	return strings.TrimRight(cs, "/")
}

// SecretHash hashes the resolved connection strings (stored in an annotation
// so a rotated secret is detected without ever comparing secrets to .show).
func SecretHash(connectionStrings []string) string {
	return normalize.Hash(connectionStrings...)
}

// ParseRows parses ".show external tables" / ".show external table E".
func ParseRows(res *kusto.Result) map[string]Observed {
	out := map[string]Observed{}
	if res == nil {
		return out
	}
	for _, r := range res.Rows() {
		name := r.String("TableName")
		if name == "" {
			continue
		}
		o := Observed{Name: name, Kind: normalizeKind(r.String("TableType")), Folder: r.String("Folder"), DocString: r.String("DocString"), PathFormat: r.String("PathFormat")}
		if raw := r.Dynamic("Properties"); len(raw) > 0 {
			var m map[string]any
			if err := json.Unmarshal(raw, &m); err == nil {
				o.Properties = m
			}
		}
		o.ConnectionStringURIs = parseConnectionStrings(r.Dynamic("ConnectionStrings"))
		o.Partitions = compact(r.Dynamic("Partitions"))
		out[name] = o
	}
	return out
}

// parseConnectionStrings accepts a JSON array of strings, a JSON string or a
// plain string and returns the URI parts.
func parseConnectionStrings(raw json.RawMessage) []string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return nil
	}
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		out := make([]string, 0, len(arr))
		for _, cs := range arr {
			out = append(out, URIOf(cs))
		}
		return out
	}
	var one string
	if err := json.Unmarshal(raw, &one); err == nil {
		return []string{URIOf(one)}
	}
	return []string{URIOf(s)}
}

func compact(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" || s == "[]" || s == "{}" {
		return ""
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return s
	}
	return string(b)
}

// ParseOne returns the single external table in res.
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

// ParseCslSchema parses ".show external table E cslschema" into columns.
func ParseCslSchema(res *kusto.Result) ([]Column, bool) {
	o, ok := table.ParseCslSchema(res)
	if !ok {
		return nil, false
	}
	cols := make([]Column, 0, len(o.Columns))
	for _, c := range o.Columns {
		cols = append(cols, Column{Name: c.Name, Type: c.Type})
	}
	return cols, true
}

// ShowAll is ".show external tables".
func ShowAll() cmd.Command { return cmd.New(".show external tables") }

// Show is ".show external table ['E']".
func Show(name string) cmd.Command { return cmd.New(".show external table ", cmd.Ident(name)) }

// ShowCslSchema is ".show external table ['E'] cslschema".
func ShowCslSchema(name string) cmd.Command {
	return cmd.New(".show external table ", cmd.Ident(name), " cslschema")
}

// BuildCreateOrAlter renders the idempotent create-or-alter command. All
// connection strings are obfuscated literals.
func BuildCreateOrAlter(d Desired) cmd.Command {
	parts := []string{".create-or-alter external table ", cmd.Ident(d.Name)}
	if len(d.Columns) > 0 {
		defs := make([]cmd.ColumnDef, 0, len(d.Columns))
		for _, c := range d.Columns {
			defs = append(defs, cmd.ColumnDef{Name: c.Name, Type: c.Type})
		}
		parts = append(parts, " ", cmd.Schema(defs))
	}
	parts = append(parts, " kind=", d.Kind)
	if d.Kind != kindDelta {
		if d.PartitionBy != "" {
			parts = append(parts, " partition by (", d.PartitionBy, ")")
		}
		if d.PathFormat != "" {
			parts = append(parts, " pathformat=(", d.PathFormat, ")")
		}
		if d.DataFormat != "" {
			parts = append(parts, " dataformat=", d.DataFormat)
		}
	}
	css := make([]string, 0, len(d.ConnectionStrings))
	for _, cs := range d.ConnectionStrings {
		css = append(css, cmd.Obfuscated(cs))
	}
	parts = append(parts, " ", cmd.List(css), cmd.With(props(d.Properties)))
	return cmd.New(parts...)
}

func props(p *v1alpha1.ExternalTableProperties) map[string]string {
	m := map[string]string{}
	if p == nil {
		return m
	}
	if p.Folder != nil {
		m["folder"] = cmd.Str(*p.Folder)
	}
	if p.DocString != nil {
		m["docString"] = cmd.Str(*p.DocString)
	}
	if p.Compressed != nil {
		m["compressed"] = cmd.Bool(*p.Compressed)
	}
	if p.CompressionType != nil {
		m["compressionType"] = cmd.Str(*p.CompressionType)
	}
	if p.IncludeHeaders != nil {
		m["includeHeaders"] = *p.IncludeHeaders
	}
	if p.NamePrefix != nil {
		m["namePrefix"] = cmd.Str(*p.NamePrefix)
	}
	if p.FileExtension != nil {
		m["fileExtension"] = cmd.Str(*p.FileExtension)
	}
	if p.Encoding != nil {
		m["encoding"] = cmd.Str(*p.Encoding)
	}
	return m
}

// BuildDelete renders ".drop external table ['E']".
func BuildDelete(name string) cmd.Command {
	return cmd.New(".drop external table ", cmd.Ident(name))
}

// Diff is the result of a structural comparison.
type Diff struct {
	Equal   bool
	Reasons []string
}

func (d *Diff) add(r string) {
	d.Equal = false
	d.Reasons = append(d.Reasons, r)
}

// String joins the reasons.
func (d *Diff) String() string { return strings.Join(d.Reasons, "; ") }

// Compare checks kind, columns (when loaded), data format, connection string
// URIs and the set properties. Partitions and path format are not compared
// here (hash mechanism, see DesiredTexts).
func Compare(d Desired, o Observed) Diff { //nolint:gocyclo // A flat list of field checks.
	diff := Diff{Equal: true}
	if o.Kind != "" && d.Kind != o.Kind {
		diff.add("kind " + o.Kind + " != " + d.Kind)
	}
	if len(d.Columns) > 0 && o.ColumnsLoaded && !sameColumns(d.Columns, o.Columns) {
		diff.add("columns differ")
	}
	if d.DataFormat != "" {
		if f := strings.ToLower(propString(o.Properties, "Format")); f != "" && f != d.DataFormat {
			diff.add("dataFormat " + f + " != " + d.DataFormat)
		}
	}
	if !sameURIs(d.ConnectionStrings, o.ConnectionStringURIs) {
		diff.add("connection string URIs differ")
	}
	if p := d.Properties; p != nil {
		if p.Folder != nil && strings.TrimSpace(*p.Folder) != strings.TrimSpace(o.Folder) {
			diff.add("folder")
		}
		if p.DocString != nil && strings.TrimSpace(*p.DocString) != strings.TrimSpace(o.DocString) {
			diff.add("docString")
		}
		check := func(name string, want *string) {
			if want == nil {
				return
			}
			if got, ok := prop(o.Properties, name); ok && !strings.EqualFold(strings.TrimSpace(asString(got)), strings.TrimSpace(*want)) {
				diff.add(name)
			}
		}
		if p.Compressed != nil {
			if got, ok := prop(o.Properties, "Compressed"); ok && !strings.EqualFold(asString(got), strconv.FormatBool(*p.Compressed)) {
				diff.add("compressed")
			}
		}
		check("CompressionType", p.CompressionType)
		check("IncludeHeaders", p.IncludeHeaders)
		check("NamePrefix", p.NamePrefix)
		check("FileExtension", p.FileExtension)
		check("Encoding", p.Encoding)
	}
	return diff
}

func sameColumns(d, o []Column) bool {
	if len(d) != len(o) {
		return false
	}
	for i := range d {
		if d[i].Name != o[i].Name || d[i].Type != normalize.ColumnType(o[i].Type) {
			return false
		}
	}
	return true
}

// sameURIs compares connection strings by count and URI part, order-insensitive.
func sameURIs(desired, observedURIs []string) bool {
	if len(desired) != len(observedURIs) {
		return false
	}
	seen := map[string]int{}
	for _, u := range observedURIs {
		seen[strings.ToLower(URIOf(u))]++
	}
	for _, cs := range desired {
		k := strings.ToLower(URIOf(cs))
		if seen[k] == 0 {
			return false
		}
		seen[k]--
	}
	return true
}

// prop looks a key up case-insensitively.
func prop(m map[string]any, key string) (any, bool) {
	if m == nil {
		return nil, false
	}
	if v, ok := m[key]; ok {
		return v, v != nil
	}
	for k, v := range m {
		if strings.EqualFold(k, key) {
			return v, v != nil
		}
	}
	return nil, false
}

func propString(m map[string]any, key string) string {
	v, ok := prop(m, key)
	if !ok {
		return ""
	}
	return asString(v)
}

func asString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		return strconv.FormatBool(x)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}

// DesiredTexts are the free text fields compared through the hash mechanism.
func DesiredTexts(d Desired) []string {
	return []string{normalize.KQL(d.PartitionBy), normalize.KQL(d.PathFormat)}
}

// ObservedTexts mirror DesiredTexts with what the cluster reports (partitions
// come back as JSON, so stage 1 equality is not expected).
func ObservedTexts(o Observed) []string {
	return []string{o.Partitions, normalize.KQL(o.PathFormat)}
}

// Observation converts observed state into the status representation.
func Observation(o Observed) v1alpha1.ExternalTableObservation {
	obs := v1alpha1.ExternalTableObservation{Kind: o.Kind, Folder: o.Folder, DocString: o.DocString, DataFormat: propString(o.Properties, "Format"), ConnectionStringURIs: o.ConnectionStringURIs, Partitions: o.Partitions, PathFormat: o.PathFormat}
	for _, c := range o.Columns {
		obs.Columns = append(obs.Columns, common.ObservedColumn{Name: c.Name, Type: c.Type})
	}
	return obs
}
