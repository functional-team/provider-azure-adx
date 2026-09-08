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

package timespan

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		in   string
		want Ticks
		err  bool
	}{
		"days":            {in: "30d", want: 30 * Day},
		"fractionalHours": {in: "1.5h", want: 90 * Minute},
		"minutes":         {in: "10m", want: 10 * Minute},
		"seconds":         {in: "30s", want: 30 * Second},
		"millis":          {in: "100ms", want: 100 * Millisecond},
		"ticks":           {in: "3ticks", want: 3},
		"colon":           {in: "00:10:00", want: 10 * Minute},
		"colonNoSeconds":  {in: "01:30", want: 90 * Minute},
		"dotnetDays":      {in: "365.00:00:00", want: 365 * Day},
		"dotnetFraction":  {in: "1.02:03:04.5000000", want: Day + 2*Hour + 3*Minute + 4*Second + 500*Millisecond},
		"shortFraction":   {in: "00:00:00.5", want: 500 * Millisecond},
		"timeWrapper":     {in: "time(1.02:03:04)", want: Day + 2*Hour + 3*Minute + 4*Second},
		"timespanWrapper": {in: "timespan(6h)", want: 6 * Hour},
		"infiniteKusto":   {in: "1000000d", want: 1000000 * Day},
		"infiniteDotNet":  {in: "1000000.00:00:00", want: 1000000 * Day},
		"goDuration":      {in: "1h30m", want: 90 * Minute},
		"negative":        {in: "-1h", want: -Hour},
		"whitespace":      {in: "  30d ", want: 30 * Day},
		"empty":           {in: "", err: true},
		"garbage":         {in: "soon", err: true},
		"badUnit":         {in: "3 fortnights", err: true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := Parse(tc.in)
			if tc.err {
				if err == nil {
					t.Fatalf("Parse(%q): expected error, got %v", tc.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q): unexpected error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("Parse(%q) = %d (%s), want %d (%s)", tc.in, got, got, tc.want, tc.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	cases := map[string]struct {
		in     Ticks
		dotnet string
		kql    string
	}{
		"zero":      {in: 0, dotnet: "00:00:00", kql: "0s"},
		"days":      {in: 30 * Day, dotnet: "30.00:00:00", kql: "30d"},
		"hours":     {in: 6 * Hour, dotnet: "06:00:00", kql: "6h"},
		"minutes":   {in: 10 * Minute, dotnet: "00:10:00", kql: "10m"},
		"seconds":   {in: 90 * Second, dotnet: "00:01:30", kql: "time(00:01:30)"},
		"seconds2":  {in: 30 * Second, dotnet: "00:00:30", kql: "30s"},
		"fraction":  {in: 500 * Millisecond, dotnet: "00:00:00.5000000", kql: "time(00:00:00.5000000)"},
		"mixed":     {in: Day + 2*Hour + 3*Minute + 4*Second, dotnet: "1.02:03:04", kql: "time(1.02:03:04)"},
		"infinite":  {in: 1000000 * Day, dotnet: "1000000.00:00:00", kql: "1000000d"},
		"negative":  {in: -Hour, dotnet: "-01:00:00", kql: "time(-01:00:00)"},
		"roundtrip": {in: MustParse("1.5h"), dotnet: "01:30:00", kql: "time(01:30:00)"},
		"minutes90": {in: 90 * Minute, dotnet: "01:30:00", kql: "time(01:30:00)"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.in.DotNet(); got != tc.dotnet {
				t.Errorf("DotNet() = %q, want %q", got, tc.dotnet)
			}
			if got := tc.in.KQL(); got != tc.kql {
				t.Errorf("KQL() = %q, want %q", got, tc.kql)
			}
			// Every formatted value must parse back to the same ticks.
			if back, err := Parse(tc.in.DotNet()); err != nil || back != tc.in {
				t.Errorf("Parse(DotNet()) = %v, %v; want %v", back, err, tc.in)
			}
			if back, err := Parse(tc.in.KQL()); err != nil || back != tc.in {
				t.Errorf("Parse(KQL()) = %v, %v; want %v", back, err, tc.in)
			}
		})
	}
}

func TestDurationSaturates(t *testing.T) {
	if got := (1000000 * Day).Duration(); got != time.Duration(int64(^uint64(0)>>1)) {
		t.Errorf("expected saturation, got %v", got)
	}
	if got := (2 * Hour).Duration(); got != 2*time.Hour {
		t.Errorf("got %v, want 2h", got)
	}
}

func TestEqual(t *testing.T) {
	if !Equal("1d", "1.00:00:00") {
		t.Error("1d should equal 1.00:00:00")
	}
	if !Equal("365d", "365.00:00:00") {
		t.Error("365d should equal 365.00:00:00")
	}
	if Equal("1d", "2d") {
		t.Error("1d should not equal 2d")
	}
	if !Equal("foo", " foo ") {
		t.Error("unparseable values compare as trimmed strings")
	}
}
