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

// Package cmd builds Kusto management commands. It is the only place in the
// provider that is allowed to concatenate user supplied values into command
// text; every identifier, string, secret and JSON document goes through one
// of the quoting helpers below.
package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/Azure/azure-kusto-go/azkustodata"
	"github.com/Azure/azure-kusto-go/azkustodata/kql"

	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

// Command is a fully rendered management command.
type Command struct {
	text string
}

// New concatenates already quoted parts into a command.
func New(parts ...string) Command {
	return Command{text: strings.Join(parts, "")}
}

// String returns the raw command text, including secrets. Never log it; use
// Redacted.
func (c Command) String() string { return c.text }

// IsZero reports whether the command is empty.
func (c Command) IsZero() bool { return c.text == "" }

var (
	obfuscatedVerbatimRe = regexp.MustCompile(`h@'(?:[^']|'')*'`)
	obfuscatedQuotedRe   = regexp.MustCompile(`h"(?:[^"\\]|\\.)*"`)
)

// Redacted returns the command text with obfuscated literals (h@'...' and
// h"...") replaced by a placeholder. Safe for logs and events.
func (c Command) Redacted() string {
	s := obfuscatedVerbatimRe.ReplaceAllString(c.text, "h@'***'")
	return obfuscatedQuotedRe.ReplaceAllString(s, `h"***"`)
}

// Statement converts the command into the SDK statement type. The SDK's
// kql.Builder only accepts compile time constants through its safe API, so
// the (already quoted) command text is attached via AddUnsafe.
func (c Command) Statement() azkustodata.Statement {
	return kql.New("").AddUnsafe(c.text)
}

// Ident renders an entity name as a bracket quoted identifier: ['name'].
// Backslashes, single quotes and control characters are escaped.
func Ident(name string) string {
	return "['" + escape(name, '\'') + "']"
}

// Qualified renders ['DB'].['Name'].
func Qualified(database, name string) string {
	return Ident(database) + "." + Ident(name)
}

// Str renders a double quoted Kusto string literal.
func Str(s string) string {
	return `"` + escape(s, '"') + `"`
}

// Verbatim renders a verbatim string literal @'...' (single quotes doubled).
// Verbatim literals cannot contain newlines; callers must use Str for those.
func Verbatim(s string) string {
	return "@'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// Obfuscated renders a string literal with the h prefix so Kusto hides it in
// .show output and command logs. Used for connection strings and secrets.
func Obfuscated(s string) string {
	if strings.ContainsAny(s, "\n\r") {
		return "h" + Str(s)
	}
	return "h" + Verbatim(s)
}

// JSON serializes v with encoding/json and renders it as a verbatim literal.
// If the JSON contains characters that cannot appear in a verbatim literal
// the quoted form is used instead.
func JSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("cannot marshal policy JSON: %w", err)
	}
	return RawJSON(string(b)), nil
}

// RawJSON renders an already serialized JSON document as a string literal.
func RawJSON(js string) string {
	if strings.ContainsAny(js, "\n\r") {
		return Str(js)
	}
	return Verbatim(js)
}

// Bool renders a boolean literal.
func Bool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// Int renders an integer literal.
func Int(i int64) string { return fmt.Sprintf("%d", i) }

// Timespan renders a timespan literal.
func Timespan(t timespan.Ticks) string { return t.KQL() }

// Body wraps a KQL body in curly braces on their own lines. The body itself is
// never modified: it is code, and whitespace inside string literals matters.
func Body(kqlText string) string {
	return "{\n" + kqlText + "\n}"
}

// Pipe renders the "<| query" suffix used by commands that take a query.
func Pipe(query string) string {
	return "\n<| " + query
}

// With renders a "with (k=v, ...)" property list from pre-rendered values,
// sorted by key for deterministic output. Returns "" for an empty map.
func With(props map[string]string) string {
	if len(props) == 0 {
		return ""
	}
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+props[k])
	}
	return " with (" + strings.Join(parts, ", ") + ")"
}

// ColumnDef is a column for schema rendering.
type ColumnDef struct {
	Name string
	Type string
}

// Schema renders a column list: (['A']:string, ['B']:long).
func Schema(cols []ColumnDef) string {
	parts := make([]string, 0, len(cols))
	for _, c := range cols {
		parts = append(parts, Ident(c.Name)+":"+c.Type)
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

// List renders a parenthesized, comma separated list of pre-rendered items.
func List(items []string) string {
	return "(" + strings.Join(items, ", ") + ")"
}

// escape applies Kusto string literal escaping. quote is the delimiter that
// must additionally be escaped.
func escape(s string, quote rune) string {
	var b strings.Builder
	b.Grow(len(s) + 8)
	for _, r := range s {
		switch {
		case r == '\\':
			b.WriteString(`\\`)
		case r == quote:
			b.WriteRune('\\')
			b.WriteRune(r)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20 || r == 0x7f || r == utf8.RuneError:
			fmt.Fprintf(&b, `\u%04x`, r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
