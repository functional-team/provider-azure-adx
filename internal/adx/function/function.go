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

// Package function holds the domain logic for the Function kind. Bodies and
// parameter lists are free text that Kusto reformats, so comparison uses the
// two stage normalization from docs/tech-implement.md section 8.
package function

import (
	"strings"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/normalize"
)

// Section is the snapshot cache section for ".show functions".
const Section = "functions"

// Observed is a function as reported by the cluster.
type Observed struct {
	Name       string
	Parameters string
	Body       string
	Folder     string
	Docstring  string
}

// Desired is the normalized desired state.
type Desired struct {
	Name           string
	Parameters     string
	Body           string
	Folder         *string
	Docstring      *string
	View           *bool
	SkipValidation *bool
}

// FromParams builds the desired state from the spec.
func FromParams(name string, p v1alpha1.FunctionParameters) Desired {
	params := strings.TrimSpace(p.Parameters)
	if params == "" {
		params = "()"
	}
	return Desired{Name: name, Parameters: params, Body: p.Body, Folder: p.Folder, Docstring: p.Docstring, View: p.View, SkipValidation: p.SkipValidation}
}

// ParseRows parses the rows of ".show functions" / ".show function" /
// ".create-or-alter function" (columns Name, Parameters, Body, Folder, DocString).
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
		out[name] = Observed{Name: name, Parameters: r.String("Parameters"), Body: r.String("Body"), Folder: r.String("Folder"), Docstring: r.String("DocString")}
	}
	return out
}

// ParseOne returns the single function in res.
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

// ShowFunctions is ".show functions".
func ShowFunctions() cmd.Command { return cmd.New(".show functions") }

// ShowFunction is ".show function ['F']".
func ShowFunction(name string) cmd.Command {
	return cmd.New(".show function ", cmd.Ident(name))
}

// BuildCreateOrAlter renders ".create-or-alter function with (...) ['F'](params) { body }".
func BuildCreateOrAlter(d Desired) cmd.Command {
	props := map[string]string{}
	if d.Folder != nil {
		props["folder"] = cmd.Str(*d.Folder)
	}
	if d.Docstring != nil {
		props["docstring"] = cmd.Str(*d.Docstring)
	}
	if d.View != nil {
		props["view"] = cmd.Bool(*d.View)
	}
	if d.SkipValidation != nil {
		props["skipvalidation"] = cmd.Bool(*d.SkipValidation)
	}
	return cmd.New(".create-or-alter function", cmd.With(props), " ", cmd.Ident(d.Name), d.Parameters, " ", cmd.Body(d.Body))
}

// BuildDelete renders ".drop function ['F'] ifexists".
func BuildDelete(name string) cmd.Command {
	return cmd.New(".drop function ", cmd.Ident(name), " ifexists")
}

// DesiredTexts returns the stage 1 normalized free text fields in a fixed order.
func DesiredTexts(d Desired) []string {
	return []string{normalize.Params(d.Parameters), normalize.Body(d.Body)}
}

// ObservedTexts returns the stage 1 normalized free text fields in a fixed order.
func ObservedTexts(o Observed) []string {
	return []string{normalize.Params(o.Parameters), normalize.Body(o.Body)}
}

// MetadataUpToDate compares folder and docstring; unset spec fields are ignored.
func MetadataUpToDate(d Desired, o Observed) bool {
	if d.Folder != nil && strings.TrimSpace(*d.Folder) != strings.TrimSpace(o.Folder) {
		return false
	}
	if d.Docstring != nil && strings.TrimSpace(*d.Docstring) != strings.TrimSpace(o.Docstring) {
		return false
	}
	return true
}

// Observation converts observed state into the status representation.
func Observation(o Observed) v1alpha1.FunctionObservation {
	return v1alpha1.FunctionObservation{Parameters: o.Parameters, Body: o.Body, Folder: o.Folder, Docstring: o.Docstring}
}
