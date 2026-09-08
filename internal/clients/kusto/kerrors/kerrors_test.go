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

package kerrors

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	kustoerrors "github.com/Azure/azure-kusto-go/azkustodata/errors"
)

// fixture builds an *HttpError the way the SDK does from a recorded REST body.
func fixture(t *testing.T, status int, file string) error {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", file))
	if err != nil {
		t.Fatal(err)
	}
	return kustoerrors.HTTP(kustoerrors.OpMgmt, http.StatusText(status), status, io.NopCloser(strings.NewReader(string(body))), "Mgmt failed")
}

func TestClassifyFixtures(t *testing.T) {
	cases := map[string]struct {
		status int
		file   string
		want   Class
		code   string
	}{
		"entityNotFound":       {400, "entity_not_found.json", NotFound, "BadRequest_EntityNotFound"},
		"tableNotFoundMessage": {400, "table_not_found_message.json", NotFound, "General_BadRequest"},
		"databaseNotFound":     {400, "database_not_found.json", NotFound, "BadRequest_EntityNotFound"},
		"alreadyExists":        {400, "already_exists.json", AlreadyExists, "BadRequest_EntityAlreadyExists"},
		"alreadyExistsMessage": {400, "already_exists_message.json", AlreadyExists, "General_BadRequest"},
		"throttled":            {429, "throttled.json", Throttled, "TooManyRequests"},
		"forbidden":            {403, "forbidden.json", Unauthorized, "Forbidden"},
		"unauthorized":         {401, "unauthorized.json", Unauthorized, ""},
		"badRequestPermanent":  {400, "bad_request_permanent.json", Permanent, "General_BadRequest"},
		"syntaxError":          {400, "syntax_error.json", Permanent, "BadRequest_SyntaxError"},
		"internalServerError":  {500, "internal_server_error.json", Transient, "General_InternalServerError"},
		"serviceUnavailable":   {503, "service_unavailable.json", Transient, ""},
		"nonJSONBody":          {502, "non_json.txt", Transient, ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := fixture(t, tc.status, tc.file)
			if got := Classify(err); got != tc.want {
				t.Errorf("Classify() = %s, want %s (details: %+v)", got, tc.want, Details(err))
			}
			if d := Details(err); d.Code != tc.code {
				t.Errorf("Details().Code = %q, want %q", d.Code, tc.code)
			}
			// Wrapping must not change the classification.
			if got := Classify(fmt.Errorf("mgmt: %w", err)); got != tc.want {
				t.Errorf("Classify(wrapped) = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestClassifyNonHTTP(t *testing.T) {
	cases := map[string]struct {
		err  error
		want Class
	}{
		"nil":            {nil, Unknown},
		"dbNotExist":     {kustoerrors.ES(kustoerrors.OpMgmt, kustoerrors.KDBNotExist, "db missing"), NotFound},
		"timeout":        {kustoerrors.ES(kustoerrors.OpMgmt, kustoerrors.KTimeout, "timeout"), Transient},
		"io":             {kustoerrors.ES(kustoerrors.OpServConn, kustoerrors.KIO, "conn reset"), Transient},
		"clientArgs":     {kustoerrors.ES(kustoerrors.OpMgmt, kustoerrors.KClientArgs, "bad arg"), Permanent},
		"ctxDeadline":    {context.DeadlineExceeded, Transient},
		"ctxCanceled":    {context.Canceled, Transient},
		"plainThrottled": {errors.New("request was throttled by the cluster"), Throttled},
		"plainUnknown":   {errors.New("something odd"), Unknown},
		"otherRetryable": {kustoerrors.ES(kustoerrors.OpMgmt, kustoerrors.KHTTPError, "gateway hiccup"), Transient},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Classify(tc.err); got != tc.want {
				t.Errorf("Classify() = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestHelpers(t *testing.T) {
	nf := fixture(t, 400, "entity_not_found.json")
	if !IsNotFound(nf) || IsAlreadyExists(nf) || IsRetryable(nf) {
		t.Error("IsNotFound helpers wrong")
	}
	if !IsRetryable(fixture(t, 429, "throttled.json")) || !IsThrottled(fixture(t, 429, "throttled.json")) {
		t.Error("throttled must be retryable")
	}
	if msg := Message(nf); !strings.Contains(msg, "BadRequest_EntityNotFound") || !strings.Contains(msg, "'RawEvents'") {
		t.Errorf("Message() = %q", msg)
	}
	if Message(nil) != "" || Message(errors.New("x")) != "x" {
		t.Error("Message fallback wrong")
	}
	if Unknown.String() != "Unknown" || Class(99).String() != "Class(99)" {
		t.Error("String() wrong")
	}
}

func TestBlocked(t *testing.T) {
	err := NewBlocked(ReasonUnsupportedColumnTypeChange, "column %s type change %s->%s is not supported", "A", "int", "long")
	if !IsBlocked(err) || !IsBlocked(fmt.Errorf("wrap: %w", err)) || IsBlocked(errors.New("x")) {
		t.Error("IsBlocked wrong")
	}
	want := "UnsupportedColumnTypeChange: column A type change int->long is not supported"
	if err.Error() != want {
		t.Errorf("got %q want %q", err.Error(), want)
	}
}
