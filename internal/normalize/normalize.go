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

// Package normalize implements the deterministic (stage 1) drift
// normalization for KQL text and Kusto identifiers, see
// docs/tech-implement.md section 8. It never collapses whitespace inside a
// line, because whitespace inside string literals is semantic.
package normalize

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
)

// KQL normalizes free text KQL: CRLF to LF, trailing whitespace per line
// removed, leading and trailing blank lines removed.
func KQL(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	start, end := 0, len(lines)
	for start < end && lines[start] == "" {
		start++
	}
	for end > start && lines[end-1] == "" {
		end--
	}
	return strings.Join(lines[start:end], "\n")
}

// Body normalizes a function or view body: one pair of outer curly braces is
// removed (Kusto echoes bodies with braces), then KQL normalization applies.
func Body(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		s = s[1 : len(s)-1]
	}
	return KQL(strings.TrimSpace(s))
}

// EqualKQL reports whether two KQL texts are equal after normalization.
func EqualKQL(a, b string) bool { return KQL(a) == KQL(b) }

// EqualBody reports whether two bodies are equal after normalization.
func EqualBody(a, b string) bool { return Body(a) == Body(b) }

// Params canonicalizes a function parameter list: whitespace outside string
// literals is removed and type aliases are mapped to their canonical Kusto
// type, so "(a:string, b:int = 5)" equals "(a:string,b:int=5)". Tabular
// parameters "T:(x:long)" and "*" are passed through token by token.
func Params(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "()"
	}
	if !strings.HasPrefix(s, "(") {
		s = "(" + s + ")"
	}
	var b strings.Builder
	b.Grow(len(s))
	scan(s, func(c byte, inString bool) {
		if inString || !isSpace(c) {
			b.WriteByte(c)
		}
	})
	return canonicalizeTypes(b.String())
}

// EqualParams reports whether two parameter lists are equal after canonicalization.
func EqualParams(a, b string) bool { return Params(a) == Params(b) }

// canonicalizeTypes rewrites type names that follow a ':' outside string
// literals. Types are identifiers up to the next ',', ')', '=' or '('.
func canonicalizeTypes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inType := false
	var typ strings.Builder
	flush := func() {
		if typ.Len() > 0 {
			b.WriteString(ColumnType(typ.String()))
			typ.Reset()
		}
		inType = false
	}
	scan(s, func(c byte, inString bool) {
		switch {
		case inString:
			flush()
			b.WriteByte(c)
		case inType && (isIdent(c)):
			typ.WriteByte(c)
		case inType:
			flush()
			if c == ':' {
				inType = true
			}
			b.WriteByte(c)
		case c == ':':
			inType = true
			b.WriteByte(c)
		default:
			b.WriteByte(c)
		}
	})
	flush()
	return b.String()
}

// scan walks s and calls fn for every byte with a flag telling whether the
// byte is part of a string literal (delimiters included). Handles "...",
// '...', verbatim @'...'/@"..." (doubled quotes) and backslash escapes.
func scan(s string, fn func(c byte, inString bool)) { //nolint:gocyclo // a small hand written scanner; splitting it would obscure the state machine.
	var quote byte
	verbatim, escaped := false, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			fn(c, true)
			switch {
			case escaped:
				escaped = false
			case !verbatim && c == '\\':
				escaped = true
			case c == quote:
				if verbatim && i+1 < len(s) && s[i+1] == quote {
					i++
					fn(s[i], true)
					continue
				}
				quote = 0
			}
			continue
		}
		if c == '"' || c == '\'' {
			quote = c
			verbatim = i > 0 && s[i-1] == '@'
			fn(c, true)
			continue
		}
		fn(c, false)
	}
}

func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' }

func isIdent(c byte) bool {
	return c == '_' || c == '.' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// canonicalColumnTypes maps each canonical Kusto type to its aliases and the
// .NET type names that appear in ".show ... schema as json" Type fields.
var canonicalColumnTypes = map[string][]string{
	"bool":     {"boolean", "system.boolean", "system.sbyte"},
	"datetime": {"date", "system.datetime"},
	"dynamic":  {"system.object"},
	"guid":     {"uuid", "uniqueid", "system.guid"},
	"int":      {"system.int32"},
	"long":     {"system.int64"},
	"real":     {"double", "system.double"},
	"decimal":  {"system.data.sqltypes.sqldecimal"},
	"string":   {"system.string"},
	"timespan": {"time", "system.timespan"},
}

var columnTypeAliases = func() map[string]string {
	m := map[string]string{}
	for canonical, aliases := range canonicalColumnTypes {
		m[canonical] = canonical
		for _, a := range aliases {
			m[a] = canonical
		}
	}
	return m
}()

// ColumnType maps a Kusto type alias or .NET type name to the canonical Kusto
// type. Unknown names are returned lower-cased and trimmed.
func ColumnType(t string) string {
	k := strings.ToLower(strings.TrimSpace(t))
	if c, ok := columnTypeAliases[k]; ok {
		return c
	}
	return k
}

// QualifiedTable renders a table in the form Kusto uses in .show outputs:
// ['DB'].['T'].
func QualifiedTable(database, table string) string {
	return cmd.Qualified(database, table)
}

// UnqualifiedTable strips a ['DB'].['T'], [DB].[T] or ['T'] wrapper and returns
// the bare table name; plain names are returned unchanged.
func UnqualifiedTable(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "].["); i >= 0 && strings.HasPrefix(s, "[") {
		s = s[i+2:]
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") && !strings.HasPrefix(s, "['") && !strings.HasPrefix(s, "[\"") {
		s = s[1 : len(s)-1]
	}
	s = strings.TrimPrefix(s, "['")
	s = strings.TrimSuffix(s, "']")
	s = strings.TrimPrefix(s, "[\"")
	s = strings.TrimSuffix(s, "\"]")
	return s
}

// Hash returns a stable 128 bit hex digest over parts. Used for the
// applied-hash / observed-hash annotations (stage 2 drift tolerance).
func Hash(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:32]
}

// Trim is strings.TrimSpace for string pointers; nil becomes "".
func Trim(p *string) string {
	if p == nil {
		return ""
	}
	return strings.TrimSpace(*p)
}
