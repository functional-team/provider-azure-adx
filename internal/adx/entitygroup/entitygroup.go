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

// Package entitygroup holds the domain logic for the EntityGroup kind (Tier
// 3). Entities are raw KQL entity expressions and are passed through
// verbatim; comparison is order-insensitive on trimmed strings with the two
// stage hash tolerance on top, because the cluster may echo the expressions
// in a canonical form.
//
// NOTE: the column layout of ".show entity_group" is not verified against a
// cluster; the parser accepts a dynamic array or a comma separated string in
// the "Entities" column and the schema JSON form.
package entitygroup

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
)

// Observed is an entity group as reported by the cluster.
type Observed struct {
	Name     string
	Entities []string
}

// FromSchema converts a schema JSON entity group.
func FromSchema(g schema.EntityGroup) Observed {
	return Observed{Name: g.Name, Entities: Normalize(g.Entities)}
}

// ParseShow parses the result of ".show entity_group ['EG']" (columns Name,
// Entities). Entities may be a dynamic array or a comma separated string.
func ParseShow(res *kusto.Result, name string) (Observed, bool) {
	rows := res.Rows()
	if len(rows) == 0 {
		return Observed{}, false
	}
	r := rows[0]
	o := Observed{Name: r.String("Name")}
	if o.Name == "" {
		o.Name = name
	}
	if raw := r.Dynamic("Entities"); len(raw) > 0 {
		var list []string
		if err := json.Unmarshal(raw, &list); err == nil {
			o.Entities = Normalize(list)
			return o, true
		}
	}
	if s := strings.TrimSpace(r.String("Entities")); s != "" {
		o.Entities = Normalize(splitEntities(s))
	}
	return o, true
}

// splitEntities splits "a, b" at commas outside parentheses and quotes.
func splitEntities(s string) []string {
	var out []string
	depth := 0
	var quote byte
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"':
			quote = c
		case c == '(':
			depth++
		case c == ')':
			depth--
		case c == ',' && depth == 0:
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

// Normalize trims, drops empty entries and sorts the entities.
func Normalize(entities []string) []string {
	out := make([]string, 0, len(entities))
	for _, e := range entities {
		if e = strings.TrimSpace(e); e != "" {
			out = append(out, e)
		}
	}
	sort.Strings(out)
	return out
}

// Equal compares two entity lists as sets of trimmed strings.
func Equal(desired, observed []string) bool {
	d, o := Normalize(desired), Normalize(observed)
	if len(d) != len(o) {
		return false
	}
	for i := range d {
		if d[i] != o[i] {
			return false
		}
	}
	return true
}

// Texts returns the normalized entity list for the hash based tolerance.
func Texts(entities []string) []string { return Normalize(entities) }

// BuildCreateOrAlter renders ".create-or-alter entity_group ['EG'] (e1, e2)".
func BuildCreateOrAlter(name string, entities []string) cmd.Command {
	return cmd.New(".create-or-alter entity_group ", cmd.Ident(name), " ", cmd.List(Normalize(entities)))
}

// BuildDelete renders ".drop entity_group ['EG']".
func BuildDelete(name string) cmd.Command {
	return cmd.New(".drop entity_group ", cmd.Ident(name))
}

// ShowOne renders ".show entity_group ['EG']".
func ShowOne(name string) cmd.Command {
	return cmd.New(".show entity_group ", cmd.Ident(name))
}

// Observation converts observed state into the status representation.
func Observation(o Observed) v1alpha1.EntityGroupObservation {
	return v1alpha1.EntityGroupObservation{Entities: o.Entities}
}
