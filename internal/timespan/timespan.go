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

// Package timespan parses and formats Kusto timespans. Values are kept as
// .NET ticks (100 ns) because Kusto uses durations such as 1,000,000 days
// ("infinite" retention) that overflow time.Duration.
package timespan

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Ticks is a duration in 100 nanosecond units, the resolution of a .NET TimeSpan.
type Ticks int64

// Tick multiples.
const (
	Microsecond Ticks = 10
	Millisecond Ticks = 1000 * Microsecond
	Second      Ticks = 1000 * Millisecond
	Minute      Ticks = 60 * Second
	Hour        Ticks = 60 * Minute
	Day         Ticks = 24 * Hour
)

var (
	// 30d, 1.5h, 10m, 30s, 100ms, 5microseconds, 3ticks ...
	shorthandRe = regexp.MustCompile(`^(-)?(\d+(?:\.\d+)?)\s*([a-zA-Z]+)$`)
	// [d.]hh:mm[:ss[.fffffff]]
	colonRe = regexp.MustCompile(`^(-)?(?:(\d+)\.)?(\d{1,2}):(\d{1,2})(?::(\d{1,2})(?:\.(\d{1,7}))?)?$`)
	wrapRe  = regexp.MustCompile(`^(?i:time|timespan)\((.*)\)$`)
)

var unitTicks = map[string]Ticks{
	"d": Day, "day": Day, "days": Day,
	"h": Hour, "hr": Hour, "hrs": Hour, "hour": Hour, "hours": Hour,
	"m": Minute, "min": Minute, "mins": Minute, "minute": Minute, "minutes": Minute,
	"s": Second, "sec": Second, "secs": Second, "second": Second, "seconds": Second,
	"ms": Millisecond, "milli": Millisecond, "millis": Millisecond, "millisecond": Millisecond, "milliseconds": Millisecond,
	"microsecond": Microsecond, "microseconds": Microsecond,
	"tick": 1, "ticks": 1,
}

// Parse parses a Kusto timespan literal, the colon form or the .NET form.
func Parse(s string) (Ticks, error) { //nolint:gocyclo // one branch per accepted literal form.
	orig := s
	s = strings.TrimSpace(s)
	if m := wrapRe.FindStringSubmatch(s); m != nil {
		s = strings.TrimSpace(m[1])
	}
	if s == "" {
		return 0, fmt.Errorf("empty timespan")
	}
	if m := shorthandRe.FindStringSubmatch(s); m != nil {
		unit, ok := unitTicks[strings.ToLower(m[3])]
		if !ok {
			return 0, fmt.Errorf("unknown timespan unit %q in %q", m[3], orig)
		}
		f, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid timespan %q: %w", orig, err)
		}
		t := Ticks(f * float64(unit))
		if m[1] == "-" {
			t = -t
		}
		return t, nil
	}
	if m := colonRe.FindStringSubmatch(s); m != nil {
		var t Ticks
		if m[2] != "" {
			d, _ := strconv.ParseInt(m[2], 10, 64)
			t += Ticks(d) * Day
		}
		h, _ := strconv.ParseInt(m[3], 10, 64)
		mi, _ := strconv.ParseInt(m[4], 10, 64)
		t += Ticks(h)*Hour + Ticks(mi)*Minute
		if m[5] != "" {
			sec, _ := strconv.ParseInt(m[5], 10, 64)
			t += Ticks(sec) * Second
		}
		if m[6] != "" {
			frac := m[6] + strings.Repeat("0", 7-len(m[6]))
			f, _ := strconv.ParseInt(frac, 10, 64)
			t += Ticks(f)
		}
		if m[1] == "-" {
			t = -t
		}
		return t, nil
	}
	if d, err := time.ParseDuration(s); err == nil {
		return FromDuration(d), nil
	}
	return 0, fmt.Errorf("invalid timespan %q", orig)
}

// MustParse is Parse for tests and constants; it panics on error.
func MustParse(s string) Ticks {
	t, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}

// FromDuration converts a Go duration to ticks.
func FromDuration(d time.Duration) Ticks { return Ticks(d / 100) }

// Duration converts ticks to a Go duration, saturating at the int64 range.
func (t Ticks) Duration() time.Duration {
	const maxTicks = Ticks(int64(^uint64(0)>>1) / 100)
	switch {
	case t > maxTicks:
		return time.Duration(int64(^uint64(0) >> 1))
	case t < -maxTicks:
		return -time.Duration(int64(^uint64(0) >> 1))
	}
	return time.Duration(t) * 100
}

// DotNet formats ticks as a .NET TimeSpan string: [-][d.]hh:mm:ss[.fffffff].
// This is the form Kusto uses inside policy JSON.
func (t Ticks) DotNet() string {
	neg := t < 0
	if neg {
		t = -t
	}
	days := t / Day
	t -= days * Day
	h := t / Hour
	t -= h * Hour
	m := t / Minute
	t -= m * Minute
	s := t / Second
	frac := t - s*Second
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	if days > 0 {
		fmt.Fprintf(&b, "%d.", days)
	}
	fmt.Fprintf(&b, "%02d:%02d:%02d", h, m, s)
	if frac > 0 {
		fmt.Fprintf(&b, ".%07d", frac)
	}
	return b.String()
}

// String is DotNet.
func (t Ticks) String() string { return t.DotNet() }

// KQL formats ticks as a KQL timespan literal. Values that are a whole
// number of a single unit and smaller than the next unit use the shorthand
// form (30d, 6h, 10m, 30s); everything else uses time(d.hh:mm:ss.fffffff).
func (t Ticks) KQL() string {
	switch {
	case t == 0:
		return "0s"
	case t < 0:
		return "time(" + t.DotNet() + ")"
	case t%Day == 0:
		return fmt.Sprintf("%dd", t/Day)
	case t%Hour == 0 && t < Day:
		return fmt.Sprintf("%dh", t/Hour)
	case t%Minute == 0 && t < Hour:
		return fmt.Sprintf("%dm", t/Minute)
	case t%Second == 0 && t < Minute:
		return fmt.Sprintf("%ds", t/Second)
	}
	return "time(" + t.DotNet() + ")"
}

// Equal reports whether two textual timespans denote the same duration. Unparseable
// inputs are compared as trimmed strings.
func Equal(a, b string) bool {
	ta, errA := Parse(a)
	tb, errB := Parse(b)
	if errA != nil || errB != nil {
		return strings.TrimSpace(a) == strings.TrimSpace(b)
	}
	return ta == tb
}
