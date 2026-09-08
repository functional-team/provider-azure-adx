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

package normalize

import "testing"

func TestKQL(t *testing.T) {
	cases := map[string]struct{ in, want string }{
		"crlf":            {"a\r\nb\r\n", "a\nb"},
		"trailingSpaces":  {"a   \n\tb\t\n", "a\n\tb"},
		"blankLines":      {"\n\n  \na\n\n", "a"},
		"innerSpacesKept": {"a  |  b", "a  |  b"},
		"stringLiteral":   {"T | where x == \"a  b\"", "T | where x == \"a  b\""},
		"empty":           {"  \n ", ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := KQL(tc.in); got != tc.want {
				t.Errorf("KQL(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestBody(t *testing.T) {
	cases := map[string]struct{ in, want string }{
		"braces":          {"{ StormEvents | take myLimit}", "StormEvents | take myLimit"},
		"bracesNewlines":  {"{\nRawEvents\n| take limit\n}", "RawEvents\n| take limit"},
		"noBraces":        {"RawEvents | take 1", "RawEvents | take 1"},
		"nestedBraceKept": {"{ let x = dynamic({\"a\":1}); x }", "let x = dynamic({\"a\":1}); x"},
		"indentedFirst":   {"{\n  T\n  | take 1\n}", "T\n  | take 1"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Body(tc.in); got != tc.want {
				t.Errorf("Body(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
	// A6: Kusto reformats "{ StormEvents | take myLimit}" to "{ StormEvents | take myLimit }".
	// The difference is trailing whitespace before the brace, which stage 1 removes.
	if !EqualBody("{ StormEvents | take myLimit}", "{ StormEvents | take myLimit }") {
		t.Error("A6 example should normalize equal")
	}
	// Interior whitespace is never collapsed (it is semantic inside string literals),
	// which is why stage 2 (hash tolerance) stays mandatory.
	if EqualBody("T  | take 1", "T | take 1") {
		t.Error("stage 1 must not collapse interior whitespace")
	}
	if !EqualBody("{\nT | take 1\n}", "T | take 1\r\n") {
		t.Error("braces + CRLF should normalize equal")
	}
}

func TestParams(t *testing.T) {
	cases := map[string]struct{ in, want string }{
		"empty":           {"", "()"},
		"spacesOnly":      {"( )", "()"},
		"simple":          {"(a:string, b:int = 5)", "(a:string,b:int=5)"},
		"aliases":         {"(a: boolean, b:double, c:date, d:uuid, e:time)", "(a:bool,b:real,c:datetime,d:guid,e:timespan)"},
		"defaultString":   {"(a:string = \"x  y\", b:long)", "(a:string=\"x  y\",b:long)"},
		"defaultVerbatim": {"(a:string = @'x  y', b:long)", "(a:string=@'x  y',b:long)"},
		"defaultColon":    {"(a:string = \"a:boolean\")", "(a:string=\"a:boolean\")"},
		"tabular":         {"(T:(x:long, y: double), limit:long)", "(T:(x:long,y:real),limit:long)"},
		"star":            {"(T:(*), n:int)", "(T:(*),n:int)"},
		"noParens":        {"a:string", "(a:string)"},
		"dynamicDefault":  {"(d:dynamic = dynamic({\"a\": 1}))", "(d:dynamic=dynamic({\"a\":1}))"},
		"escapedQuote":    {"(a:string = \"q\\\"x  y\")", "(a:string=\"q\\\"x  y\")"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Params(tc.in); got != tc.want {
				t.Errorf("Params(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
	if !EqualParams("(a:string, b:int = 5)", "(a:string,b:int=5)") {
		t.Error("EqualParams should be true")
	}
}

func TestColumnType(t *testing.T) {
	cases := map[string]string{
		"boolean": "bool", "Bool": "bool", "date": "datetime", "uuid": "guid", "uniqueid": "guid",
		"double": "real", "time": "timespan", "string": "string", "System.String": "string",
		"System.Int64": "long", "System.Object": "dynamic", "System.Data.SqlTypes.SqlDecimal": "decimal",
		" long ": "long", "weird": "weird", //nolint:gocritic // the surrounding whitespace is what is being tested
	}
	for in, want := range cases {
		if got := ColumnType(in); got != want {
			t.Errorf("ColumnType(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestQualifiedTable(t *testing.T) {
	if got := QualifiedTable("DB", "T"); got != "['DB'].['T']" {
		t.Errorf("got %q", got)
	}
	for in, want := range map[string]string{
		"['DB'].['T']": "T", "T": "T", "['T']": "T", "[\"DB\"].[\"T\"]": "T", " ['DB'].['Raw Events'] ": "Raw Events",
		"[DB].[T]": "T", "[T]": "T", "[DB].[Raw Events]": "Raw Events",
	} {
		if got := UnqualifiedTable(in); got != want {
			t.Errorf("UnqualifiedTable(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHash(t *testing.T) {
	a := Hash("x", "y")
	if len(a) != 32 || a != Hash("x", "y") {
		t.Errorf("hash not stable or wrong length: %q", a)
	}
	if Hash("xy", "") == Hash("x", "y") {
		t.Error("parts must be separated")
	}
	if Trim(nil) != "" {
		t.Error("Trim(nil) must be empty")
	}
	s := " a "
	if Trim(&s) != "a" {
		t.Error("Trim should trim")
	}
}
