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

// Package ingestionmapping holds the domain logic for the IngestionMapping
// kind: command rendering, parsing of ".show table T ingestion mappings" and
// the comparison of mapping JSON documents.
package ingestionmapping

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/normalize"
)

// SectionPrefix is the snapshot cache section prefix; the section is
// "mappings:<table>" because there is no database-wide mapping command.
const SectionPrefix = "mappings:"

// Section returns the cache section for a table.
func Section(table string) string { return SectionPrefix + table }

// Column is one entry of the Kusto mapping JSON.
type Column struct {
	Column     string            `json:"Column"`
	DataType   string            `json:"DataType,omitempty"`
	Properties map[string]string `json:"Properties,omitempty"`
}

// Observed is a mapping as reported by the cluster.
type Observed struct {
	Name          string
	Kind          string
	Table         string
	LastUpdatedOn string
	Mapping       []Column
	Raw           string
}

// Desired is the normalized desired state.
type Desired struct {
	Name    string
	Table   string
	Kind    string
	Mapping []Column
}

// Key identifies a mapping within a table: kind and name (Kusto allows the
// same name for different kinds).
func Key(kind, name string) string {
	return strings.ToLower(strings.TrimSpace(kind)) + "/" + name
}

// CanonicalKey writes a property key the way Kusto does: first letter upper
// case, rest as given ("path" -> "Path", "constValue" -> "ConstValue").
func CanonicalKey(k string) string {
	k = strings.TrimSpace(k)
	if k == "" {
		return k
	}
	r := []rune(k)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// FromParams normalizes the spec.
func FromParams(name string, p v1alpha1.IngestionMappingParameters) Desired {
	d := Desired{Name: name, Table: p.Table, Kind: strings.ToLower(string(p.Kind))}
	for _, c := range p.Mapping {
		col := Column{Column: c.Column}
		if c.DataType != nil {
			col.DataType = normalize.ColumnType(string(*c.DataType))
		}
		if len(c.Properties) > 0 {
			col.Properties = make(map[string]string, len(c.Properties))
			for k, v := range c.Properties {
				col.Properties[CanonicalKey(k)] = v
			}
		}
		d.Mapping = append(d.Mapping, col)
	}
	return d
}

// MappingJSON renders the mapping array as JSON.
func MappingJSON(cols []Column) (string, error) {
	b, err := json.Marshal(cols)
	if err != nil {
		return "", fmt.Errorf("cannot marshal ingestion mapping: %w", err)
	}
	return string(b), nil
}

// ParseMapping parses the Kusto mapping JSON (case insensitive keys; Kusto
// adds keys such as CsvDataType that are ignored).
func ParseMapping(raw string) ([]Column, error) { //nolint:gocyclo // tolerant parser over several JSON spellings.
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return nil, nil
	}
	var generic []map[string]any
	if err := json.Unmarshal([]byte(raw), &generic); err != nil {
		return nil, fmt.Errorf("cannot parse ingestion mapping JSON: %w", err)
	}
	out := make([]Column, 0, len(generic))
	for _, g := range generic {
		var c Column
		for k, v := range g {
			switch strings.ToLower(k) {
			case "column":
				c.Column, _ = v.(string)
			case "datatype":
				if s, ok := v.(string); ok {
					c.DataType = normalize.ColumnType(s)
				}
			case "properties":
				if m, ok := v.(map[string]any); ok && len(m) > 0 {
					c.Properties = make(map[string]string, len(m))
					for pk, pv := range m {
						if pv == nil {
							continue
						}
						if s, ok := pv.(string); ok {
							c.Properties[CanonicalKey(pk)] = s
						} else {
							c.Properties[CanonicalKey(pk)] = fmt.Sprint(pv)
						}
					}
				}
			}
		}
		out = append(out, c)
	}
	return out, nil
}

// ParseRows parses ".show table T ingestion mappings" (columns Name, Kind,
// Mapping, LastUpdatedOn, Database, Table) into Key -> Observed.
func ParseRows(res *kusto.Result) (map[string]Observed, error) {
	out := map[string]Observed{}
	if res == nil {
		return out, nil
	}
	for _, r := range res.Rows() {
		name := r.String("Name")
		if name == "" {
			continue
		}
		raw := r.String("Mapping")
		cols, err := ParseMapping(raw)
		if err != nil {
			return nil, fmt.Errorf("mapping %q: %w", name, err)
		}
		o := Observed{Name: name, Kind: strings.ToLower(r.String("Kind")), Table: r.String("Table"), LastUpdatedOn: r.String("LastUpdatedOn"), Mapping: cols, Raw: compact(raw)}
		out[Key(o.Kind, name)] = o
	}
	return out, nil
}

// ParseOne returns the mapping of kind/name in res (single .show), falling
// back to the only row when the kind column is missing.
func ParseOne(res *kusto.Result, kind, name string) (Observed, bool, error) {
	m, err := ParseRows(res)
	if err != nil {
		return Observed{}, false, err
	}
	if o, ok := m[Key(kind, name)]; ok {
		return o, true, nil
	}
	if len(m) == 1 {
		for _, o := range m {
			if o.Name == name {
				return o, true, nil
			}
		}
	}
	return Observed{}, false, nil
}

func compact(raw string) string {
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return strings.TrimSpace(raw)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return strings.TrimSpace(raw)
	}
	return string(b)
}

// ShowAll is ".show table ['T'] ingestion mappings".
func ShowAll(table string) cmd.Command {
	return cmd.New(".show table ", cmd.Ident(table), " ingestion mappings")
}

// ShowOne is ".show table ['T'] ingestion <kind> mapping "Name"".
func ShowOne(table, kind, name string) cmd.Command {
	return cmd.New(".show table ", cmd.Ident(table), " ingestion ", strings.ToLower(kind), " mapping ", cmd.Str(name))
}

// BuildCreateOrAlter is ".create-or-alter table ['T'] ingestion <kind> mapping "Name" '[...]'".
func BuildCreateOrAlter(d Desired) (cmd.Command, error) {
	js, err := MappingJSON(d.Mapping)
	if err != nil {
		return cmd.Command{}, err
	}
	return cmd.New(".create-or-alter table ", cmd.Ident(d.Table), " ingestion ", d.Kind, " mapping ", cmd.Str(d.Name), " ", cmd.RawJSON(js)), nil
}

// BuildDelete is ".drop table ['T'] ingestion <kind> mapping "Name"".
func BuildDelete(table, kind, name string) cmd.Command {
	return cmd.New(".drop table ", cmd.Ident(table), " ingestion ", strings.ToLower(kind), " mapping ", cmd.Str(name))
}

// Equal compares the desired mapping with the observed one element by
// element. Only fields set in the spec are compared; property keys are
// matched case insensitively. The second return value describes the first
// difference.
func Equal(d Desired, o Observed) (bool, string) { //nolint:gocyclo // element-wise comparison with only-set-fields semantics.
	if len(d.Mapping) != len(o.Mapping) {
		return false, fmt.Sprintf("mapping has %d columns in the cluster, %d in the spec", len(o.Mapping), len(d.Mapping))
	}
	for i, dc := range d.Mapping {
		oc := o.Mapping[i]
		if dc.Column != oc.Column {
			return false, fmt.Sprintf("mapping[%d].column: want %q, got %q", i, dc.Column, oc.Column)
		}
		if dc.DataType != "" && dc.DataType != oc.DataType {
			return false, fmt.Sprintf("mapping[%d].dataType: want %q, got %q", i, dc.DataType, oc.DataType)
		}
		lower := make(map[string]string, len(oc.Properties))
		for k, v := range oc.Properties {
			lower[strings.ToLower(k)] = v
		}
		keys := make([]string, 0, len(dc.Properties))
		for k := range dc.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			ov, ok := lower[strings.ToLower(k)]
			if !ok {
				return false, fmt.Sprintf("mapping[%d].properties.%s: missing in cluster", i, k)
			}
			if strings.TrimSpace(ov) != strings.TrimSpace(dc.Properties[k]) {
				return false, fmt.Sprintf("mapping[%d].properties.%s: want %q, got %q", i, k, dc.Properties[k], ov)
			}
		}
	}
	return true, ""
}

// Observation converts observed state into the status representation.
func Observation(o Observed) v1alpha1.IngestionMappingObservation {
	return v1alpha1.IngestionMappingObservation{Kind: o.Kind, Table: o.Table, LastUpdatedOn: o.LastUpdatedOn, Mapping: o.Raw}
}
