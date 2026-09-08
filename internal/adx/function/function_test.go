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

package function

import (
	"testing"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
)

func ptr[T any](v T) *T { return &v }

func rows(vals ...[]any) *kusto.Result {
	return kusto.NewResult(kusto.NewTable("Table_0", []string{"Name", "Parameters", "Body", "Folder", "DocString"}, vals...))
}

func TestBuild(t *testing.T) {
	d := FromParams("TakeSome", v1alpha1.FunctionParameters{Parameters: "(limit:long = 100)", Body: "RawEvents\n| take limit", Folder: ptr("Parsing"), Docstring: ptr("doc"), View: ptr(false), SkipValidation: ptr(true)})
	want := ".create-or-alter function with (docstring=\"doc\", folder=\"Parsing\", skipvalidation=true, view=false) ['TakeSome'](limit:long = 100) {\nRawEvents\n| take limit\n}"
	if got := BuildCreateOrAlter(d).String(); got != want {
		t.Errorf("\n got %s\nwant %s", got, want)
	}
	minimal := FromParams("F", v1alpha1.FunctionParameters{Body: "T"})
	if got := BuildCreateOrAlter(minimal).String(); got != ".create-or-alter function ['F']() {\nT\n}" {
		t.Errorf("minimal: %s", got)
	}
	if BuildDelete("F").String() != ".drop function ['F'] ifexists" || ShowFunction("F").String() != ".show function ['F']" || ShowFunctions().String() != ".show functions" {
		t.Error("simple commands")
	}
}

func TestParse(t *testing.T) {
	res := rows([]any{"A", "(x:long)", "{ T | take x }", "F", "d"}, []any{"B", "()", "{ T }", "", ""}, []any{"", "", "", "", ""})
	m := ParseRows(res)
	if len(m) != 2 || m["A"].Body != "{ T | take x }" || m["A"].Folder != "F" || m["B"].Parameters != "()" {
		t.Errorf("ParseRows: %+v", m)
	}
	if o, ok := ParseOne(rows([]any{"A", "()", "{ T }", "", ""}), "A"); !ok || o.Name != "A" {
		t.Error("ParseOne by name")
	}
	if o, ok := ParseOne(rows([]any{"Renamed", "()", "{ T }", "", ""}), "A"); !ok || o.Name != "Renamed" {
		t.Error("ParseOne single row fallback")
	}
	if _, ok := ParseOne(kusto.NewResult(), "A"); ok {
		t.Error("ParseOne empty")
	}
	if ParseRows(nil) == nil {
		t.Error("nil result must give empty map")
	}
	obs := Observation(m["A"])
	if obs.Body != "{ T | take x }" || obs.Docstring != "d" {
		t.Error("Observation")
	}
}

func TestTextsAndMetadata(t *testing.T) {
	d := FromParams("F", v1alpha1.FunctionParameters{Parameters: "(a:string, b:int = 5)", Body: "T\r\n| take b  \r\n"})
	o := Observed{Name: "F", Parameters: "(a:string,b:int=5)", Body: "{\nT\n| take b\n}"}
	dt, ot := DesiredTexts(d), ObservedTexts(o)
	if dt[0] != ot[0] || dt[1] != ot[1] {
		t.Errorf("stage 1 should make these equal: %q vs %q", dt, ot)
	}
	if !MetadataUpToDate(d, Observed{Folder: "whatever"}) {
		t.Error("unset folder must be ignored")
	}
	d.Folder = ptr("X")
	if MetadataUpToDate(d, Observed{Folder: "Y"}) || !MetadataUpToDate(d, Observed{Folder: " X "}) {
		t.Error("folder comparison")
	}
	d.Docstring = ptr("doc")
	if MetadataUpToDate(d, Observed{Folder: "X"}) {
		t.Error("docstring mismatch must be detected")
	}
}
