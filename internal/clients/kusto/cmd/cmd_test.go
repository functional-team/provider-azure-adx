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

package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

var update = flag.Bool("update", false, "update golden files")

func TestQuoting(t *testing.T) {
	cases := map[string]struct {
		got  string
		want string
	}{
		"identSimple":        {Ident("RawEvents"), "['RawEvents']"},
		"identSpace":         {Ident("My Table"), "['My Table']"},
		"identQuote":         {Ident("it's"), `['it\'s']`},
		"identBackslash":     {Ident(`a\b`), `['a\\b']`},
		"identNewline":       {Ident("a\nb"), `['a\nb']`},
		"identBracketClose":  {Ident("a']; .drop table x; --"), `['a\']; .drop table x; --']`},
		"qualified":          {Qualified("DB", "T"), "['DB'].['T']"},
		"str":                {Str(`say "hi"`), `"say \"hi\""`},
		"strNewline":         {Str("a\nb"), `"a\nb"`},
		"strControl":         {Str("a\x01b"), `"a\u0001b"`},
		"verbatim":           {Verbatim("it's"), "@'it''s'"},
		"verbatimBackslash":  {Verbatim(`\d+`), `@'\d+'`},
		"obfuscated":         {Obfuscated("https://x;sig=abc'def"), "h@'https://x;sig=abc''def'"},
		"obfuscatedNewline":  {Obfuscated("a\nb"), `h"a\nb"`},
		"bool":               {Bool(true), "true"},
		"int":                {Int(-5), "-5"},
		"timespanDays":       {Timespan(30 * timespan.Day), "30d"},
		"timespanOdd":        {Timespan(timespan.Hour + timespan.Second), "time(01:00:01)"},
		"body":               {Body("T | take 1"), "{\nT | take 1\n}"},
		"pipe":               {Pipe("T | take 1"), "\n<| T | take 1"},
		"withEmpty":          {With(nil), ""},
		"withSorted":         {With(map[string]string{"folder": Str("F"), "docstring": Str("D")}), ` with (docstring="D", folder="F")`},
		"schema":             {Schema([]ColumnDef{{"A", "string"}, {"B", "long"}}), "(['A']:string, ['B']:long)"},
		"list":               {List([]string{"'a'", "'b'"}), "('a', 'b')"},
		"rawJSONVerbatim":    {RawJSON(`{"a":"it's"}`), `@'{"a":"it''s"}'`},
		"rawJSONWithNewline": {RawJSON("{\n}"), `"{\n}"`},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("got  %s\nwant %s", tc.got, tc.want)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	got, err := JSON(map[string]any{"SoftDeletePeriod": "365.00:00:00", "Recoverability": "Enabled"})
	if err != nil {
		t.Fatal(err)
	}
	want := `@'{"Recoverability":"Enabled","SoftDeletePeriod":"365.00:00:00"}'`
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
	if _, err := JSON(make(chan int)); err == nil {
		t.Error("expected marshal error")
	}
}

func TestRedacted(t *testing.T) {
	c := New(".create external table ", Ident("E"), " (", Obfuscated("https://acct;sig=SECRET"), ", ", Obfuscated("x\ny"), ") ", Str("keep"))
	r := c.Redacted()
	if strings.Contains(r, "SECRET") || strings.Contains(r, "x\ny") {
		t.Errorf("secret leaked: %s", r)
	}
	if !strings.Contains(r, "h@'***'") || !strings.Contains(r, `h"***"`) {
		t.Errorf("placeholders missing: %s", r)
	}
	if !strings.Contains(r, `"keep"`) {
		t.Errorf("non-secret literal must be kept: %s", r)
	}
	if !strings.Contains(c.String(), "SECRET") {
		t.Error("String() must return the raw command")
	}
}

func TestStatement(t *testing.T) {
	c := New(".show table ", Ident("T"), " cslschema")
	if got := c.Statement().String(); got != c.String() {
		t.Errorf("statement text %q != command %q", got, c.String())
	}
	if !New().IsZero() || New("x").IsZero() {
		t.Error("IsZero mismatch")
	}
}

// TestGolden renders representative commands and compares them with golden
// files under testdata. Run with -update to regenerate.
func TestGolden(t *testing.T) {
	cases := map[string]Command{
		"create_table": New(".create table ", Ident("Raw Events"), " ",
			Schema([]ColumnDef{{"Timestamp", "datetime"}, {"Pay'load", "dynamic"}}),
			With(map[string]string{"folder": Str("Raw"), "docstring": Str(`Landing "table"`)})),
		"create_or_alter_function": New(".create-or-alter function",
			With(map[string]string{"docstring": Str("d"), "folder": Str("f"), "skipvalidation": Bool(false)}),
			" ", Ident("F"), "(limit:long = 100) ", Body("RawEvents\n| take limit")),
		"alter_policy_json": New(".alter table ", Ident("T"), " policy retention ", RawJSON(`{"SoftDeletePeriod":"365.00:00:00","Recoverability":"Enabled"}`)),
		"external_table": New(".create-or-alter external table ", Ident("E"), " ", Schema([]ColumnDef{{"A", "string"}}),
			" kind=storage dataformat=parquet ", List([]string{Obfuscated("https://acct.blob.core.windows.net/c;managed_identity=system")}),
			With(map[string]string{"folder": Str("External")})),
		"continuous_export": New(".create-or-alter continuous-export ", Ident("CE"), " over (", Ident("T"), ") to table ", Ident("E"),
			With(map[string]string{"intervalBetweenRuns": Timespan(timespan.Hour), "forcedLatency": Timespan(10 * timespan.Minute)}),
			Pipe("T | where Level == \"Error\"")),
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join("testdata", name+".kql")
			if *update {
				if err := os.WriteFile(path, []byte(c.String()), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read golden: %v (run with -update to create)", err)
			}
			if string(want) != c.String() {
				t.Errorf("golden mismatch\n--- want\n%s\n--- got\n%s", want, c.String())
			}
		})
	}
}
