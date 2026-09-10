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

package policy

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/functional-team/provider-azure-adx/internal/normalize"
	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

// Result of a comparison. Equal covers the structural fields; KQL text
// fields are returned as normalized pairs for the two stage check.
type Result struct {
	Equal         bool
	Diff          string
	DesiredTexts  []string
	ObservedTexts []string
}

type comparer struct {
	kql, sets map[string]bool
	// aliases maps a JSON key to desired values the service resolves to
	// something of its own, e.g. ObjectId "system" comes back as the GUID of
	// the cluster's system-assigned identity. Any observed value is accepted
	// for those, because the desired side names an identity rather than a
	// value.
	aliases map[string]map[string]bool
	// listSubset requires every desired array element to be present in the
	// observed array instead of demanding equal length: the callout policy is
	// returned with 19 immutable built-in rules alongside the two we manage.
	listSubset bool
	res        *Result
}

// Options tune Compare for a policy whose observed shape is not a plain echo
// of what was sent.
type Options struct {
	KQLFields   []string
	SetFields   []string
	AliasFields map[string][]string
	ListSubset  bool
}

// Compare checks that every field set in desired equals the observed policy
// JSON. Observed fields missing from desired are ignored (Kusto defaults).
// Timespans compare by value ("365d" == "365.00:00:00"), datetimes by instant,
// numbers numerically, {"Value": x} wrappers are unwrapped, and fields listed
// in kqlFields are collected as normalized texts instead of compared.
func Compare(desired any, observed json.RawMessage, kqlFields, setFields []string) (Result, error) {
	return CompareWithOptions(desired, observed, Options{KQLFields: kqlFields, SetFields: setFields})
}

// CompareWithOptions is Compare with the extra allowances described by Options.
func CompareWithOptions(desired any, observed json.RawMessage, o Options) (Result, error) {
	db, err := json.Marshal(desired)
	if err != nil {
		return Result{}, fmt.Errorf("cannot marshal desired policy: %w", err)
	}
	var dv, ov any
	if err := json.Unmarshal(db, &dv); err != nil {
		return Result{}, err
	}
	if !IsNull(observed) {
		if err := json.Unmarshal(observed, &ov); err != nil {
			return Result{}, fmt.Errorf("cannot parse observed policy JSON: %w", err)
		}
	}
	aliases := make(map[string]map[string]bool, len(o.AliasFields))
	for k, vs := range o.AliasFields {
		aliases[k] = toSet(vs)
	}
	c := &comparer{kql: toSet(o.KQLFields), sets: toSet(o.SetFields), aliases: aliases, listSubset: o.ListSubset, res: &Result{Equal: true}}
	c.subset("", dv, ov)
	return *c.res, nil
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

func (c *comparer) fail(path, format string, args ...any) {
	if c.res.Equal {
		c.res.Equal = false
		c.res.Diff = strings.TrimPrefix(path, ".") + ": " + fmt.Sprintf(format, args...)
	}
}

func lastKey(path string) string {
	if i := strings.LastIndex(path, "."); i >= 0 {
		return path[i+1:]
	}
	return path
}

func (c *comparer) subset(path string, d, o any) { //nolint:gocyclo // A type switch over JSON kinds.
	switch dv := d.(type) {
	case nil:
		if o != nil {
			c.fail(path, "want null, got %v", o)
		}
	case map[string]any:
		om, ok := o.(map[string]any)
		if !ok {
			c.fail(path, "want object, got %v", o)
			return
		}
		keys := make([]string, 0, len(dv))
		for k := range dv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			ov, present := om[k]
			if !present {
				if isEmpty(dv[k]) {
					continue
				}
				c.fail(path+"."+k, "missing in cluster (want %v)", dv[k])
				continue
			}
			c.subset(path+"."+k, dv[k], ov)
		}
	case []any:
		oa, ok := o.([]any)
		if !ok {
			c.fail(path, "want array, got %v", o)
			return
		}
		if c.listSubset {
			c.subsetOfList(path, dv, oa)
			return
		}
		if len(oa) != len(dv) {
			c.fail(path, "want %d elements, got %d", len(dv), len(oa))
			return
		}
		for i := range dv {
			c.subset(fmt.Sprintf("%s[%d]", path, i), dv[i], oa[i])
		}
	case string:
		key := lastKey(strings.TrimRight(path, "0123456789[]"))
		if c.aliases[key][dv] {
			// The desired value names something the service resolves itself,
			// so whatever it reports back is the right answer by definition.
			return
		}
		if c.kql[key] {
			c.res.DesiredTexts = append(c.res.DesiredTexts, normalize.KQL(dv))
			c.res.ObservedTexts = append(c.res.ObservedTexts, normalize.KQL(asString(o)))
			return
		}
		os := asString(o)
		if c.sets[key] {
			if csvSet(dv) != csvSet(os) {
				c.fail(path, "want %q, got %q", dv, os)
			}
			return
		}
		if !stringsEqual(dv, os) {
			c.fail(path, "want %q, got %q", dv, os)
		}
	case bool:
		switch ov := o.(type) {
		case bool:
			if ov != dv {
				c.fail(path, "want %v, got %v", dv, ov)
			}
		case string:
			if strings.ToLower(ov) != strconv.FormatBool(dv) {
				c.fail(path, "want %v, got %q", dv, ov)
			}
		default:
			c.fail(path, "want %v, got %v", dv, o)
		}
	case float64:
		switch ov := o.(type) {
		case float64:
			if ov != dv {
				c.fail(path, "want %v, got %v", dv, ov)
			}
		case string:
			f, err := strconv.ParseFloat(ov, 64)
			if err != nil || f != dv {
				c.fail(path, "want %v, got %q", dv, ov)
			}
		default:
			c.fail(path, "want %v, got %v", dv, o)
		}
	default:
		c.fail(path, "unsupported desired type %T", d)
	}
}

func isEmpty(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case []any:
		return len(x) == 0
	case map[string]any:
		return len(x) == 0
	}
	return false
}

// asString flattens observed scalars and {"Value": x} wrappers to text.
func asString(o any) string {
	switch ov := o.(type) {
	case nil:
		return ""
	case string:
		return ov
	case bool:
		return strconv.FormatBool(ov)
	case float64:
		return strconv.FormatFloat(ov, 'f', -1, 64)
	case map[string]any:
		if v, ok := ov["Value"]; ok {
			return asString(v)
		}
	}
	b, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprint(o)
	}
	return string(b)
}

func stringsEqual(a, b string) bool {
	a, b = strings.TrimSpace(a), strings.TrimSpace(b)
	if a == b {
		return true
	}
	if ta, err := timespan.Parse(a); err == nil {
		if tb, err := timespan.Parse(b); err == nil {
			return ta == tb
		}
	}
	if da, ok := parseTime(a); ok {
		if dbt, ok := parseTime(b); ok {
			return da.Equal(dbt)
		}
	}
	return false
}

var timeLayouts = []string{time.RFC3339Nano, "2006-01-02T15:04:05.9999999", "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"}

func parseTime(s string) (time.Time, bool) {
	for _, l := range timeLayouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// csvSet canonicalizes "A, B" into a sorted, lower-cased set string.
func csvSet(s string) string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.ToLower(strings.TrimSpace(p)); p != "" {
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return strings.Join(out, ",")
}

// subsetOfList requires every desired element to appear somewhere in observed,
// ignoring extras the service adds and ignoring order. Each element is matched
// with a throwaway comparer so a mismatch on one candidate does not poison the
// result.
func (c *comparer) subsetOfList(path string, desired, observed []any) {
	for i, d := range desired {
		found := false
		for _, o := range observed {
			probe := &comparer{kql: c.kql, sets: c.sets, aliases: c.aliases, listSubset: c.listSubset, res: &Result{Equal: true}}
			probe.subset("", d, o)
			if probe.res.Equal {
				found = true
				break
			}
		}
		if !found {
			c.fail(fmt.Sprintf("%s[%d]", path, i), "not found in cluster (want %v)", d)
			return
		}
	}
}
