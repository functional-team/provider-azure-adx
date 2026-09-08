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

// Package kerrors classifies errors returned by the Kusto SDK so that the
// reconcilers can react uniformly: NotFound becomes "resource does not
// exist", AlreadyExists is tolerated on create, Throttled and Transient are
// retried by crossplane-runtime's backoff, Unauthorized and Permanent surface
// as-is in the Synced condition.
package kerrors

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	kustoerrors "github.com/Azure/azure-kusto-go/azkustodata/errors"
)

// Class is the provider's view of a Kusto error.
type Class int

// Error classes.
const (
	Unknown Class = iota
	NotFound
	AlreadyExists
	Throttled
	Transient
	Unauthorized
	Permanent
)

var classNames = map[Class]string{
	Unknown: "Unknown", NotFound: "NotFound", AlreadyExists: "AlreadyExists", Throttled: "Throttled",
	Transient: "Transient", Unauthorized: "Unauthorized", Permanent: "Permanent",
}

func (c Class) String() string {
	if s, ok := classNames[c]; ok {
		return s
	}
	return fmt.Sprintf("Class(%d)", int(c))
}

// Detail is what could be extracted from a Kusto REST error body.
type Detail struct {
	StatusCode int
	Code       string
	Type       string
	Message    string
	Permanent  bool
}

// Details extracts the Kusto REST error details from err, if any.
func Details(err error) Detail {
	var d Detail
	var httpErr *kustoerrors.HttpError
	var kErr *kustoerrors.Error
	switch {
	case errors.As(err, &httpErr):
		d.StatusCode = httpErr.StatusCode
		fill(&d, httpErr.UnmarshalREST())
	case errors.As(err, &kErr):
		if m := kErr.UnmarshalREST(); m != nil {
			fill(&d, m)
		} else {
			// No REST body: the v1 endpoint can answer HTTP 200 with an
			// `Exceptions` array (observed on the emulator for `.show table X`
			// on a missing table). The SDK folds that into the message, so
			// the message is all we have to classify on.
			d.Message = messageOf(kErr)
		}
	}
	return d
}

// messageOf returns the wrapped message of a *kustoerrors.Error without the
// Op/Kind prefix the SDK adds in Error().
func messageOf(e *kustoerrors.Error) string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Error()
}

func fill(d *Detail, m map[string]interface{}) {
	if m == nil {
		return
	}
	e, ok := m["error"].(map[string]interface{})
	if !ok {
		return
	}
	d.Code, _ = e["code"].(string)
	d.Type, _ = e["@type"].(string)
	if msg, ok := e["@message"].(string); ok && msg != "" {
		d.Message = msg
	} else {
		d.Message, _ = e["message"].(string)
	}
	d.Permanent, _ = e["@permanent"].(bool)
	// Some responses carry the specific code one level down.
	if inner, ok := e["innererror"].(map[string]interface{}); ok {
		if c, ok := inner["code"].(string); ok && c != "" && (d.Code == "" || strings.HasPrefix(d.Code, "General_")) {
			d.Code = c
		}
	}
}

// Classify maps err to a Class. nil is Unknown.
func Classify(err error) Class { //nolint:gocyclo // A flat decision table is easier to audit than helpers.
	if err == nil {
		return Unknown
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return Transient
	}
	if errors.Is(err, context.Canceled) {
		return Transient
	}

	d := Details(err)
	text := strings.ToLower(d.Code + " " + d.Type + " " + d.Message)

	// Entity-level codes come before the generic HTTP status mapping: they
	// are all HTTP 400 and @permanent, but they carry meaning for us.
	switch {
	case strings.Contains(text, "entitynotfound"),
		strings.Contains(text, "was not found"),
		strings.Contains(text, "does not exist"),
		strings.Contains(text, "doesn't exist"),
		strings.Contains(text, "could not be found"),
		strings.Contains(text, "cannot find"):
		return NotFound
	case strings.Contains(text, "alreadyexists"),
		strings.Contains(text, "already exists"),
		strings.Contains(text, "nameconflict"),
		strings.Contains(text, "entityconflict"):
		return AlreadyExists
	}

	var httpErr *kustoerrors.HttpError
	if errors.As(err, &httpErr) {
		switch {
		case httpErr.IsThrottled(), httpErr.StatusCode == http.StatusTooManyRequests:
			return Throttled
		case httpErr.StatusCode == http.StatusUnauthorized, httpErr.StatusCode == http.StatusForbidden:
			return Unauthorized
		case httpErr.StatusCode == http.StatusNotFound:
			return NotFound
		case httpErr.StatusCode == http.StatusConflict:
			return AlreadyExists
		case httpErr.StatusCode == http.StatusRequestTimeout, httpErr.StatusCode >= 500:
			return Transient
		case d.Permanent, httpErr.StatusCode >= 400:
			return Permanent
		}
	}

	var kErr *kustoerrors.Error
	if errors.As(err, &kErr) {
		switch kErr.Kind { //nolint:exhaustive // Only kinds with a distinct mapping are listed.
		case kustoerrors.KDBNotExist:
			return NotFound
		case kustoerrors.KTimeout, kustoerrors.KIO:
			return Transient
		case kustoerrors.KClientArgs, kustoerrors.KLimitsExceeded:
			return Permanent
		}
		if strings.Contains(strings.ToLower(kErr.Error()), "throttl") {
			return Throttled
		}
		// A v1 `Exceptions` frame is a command-level failure delivered with
		// HTTP 200; the SDK files it as KInternal. Retrying does not help.
		if kErr.Kind == kustoerrors.KInternal && strings.HasPrefix(messageOf(kErr), "exceptions:") {
			return Permanent
		}
		if kustoerrors.Retry(err) {
			return Transient
		}
		return Permanent
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return Transient
	}
	if strings.Contains(strings.ToLower(err.Error()), "throttl") {
		return Throttled
	}
	return Unknown
}

// IsNotFound reports whether err denotes a missing entity or database.
func IsNotFound(err error) bool { return Classify(err) == NotFound }

// IsAlreadyExists reports whether err denotes an entity that already exists.
func IsAlreadyExists(err error) bool { return Classify(err) == AlreadyExists }

// IsThrottled reports whether the cluster throttled the command.
func IsThrottled(err error) bool { return Classify(err) == Throttled }

// IsRetryable reports whether a retry with backoff is worthwhile.
func IsRetryable(err error) bool {
	c := Classify(err)
	return c == Throttled || c == Transient
}

// Message returns the most specific human readable message for err.
func Message(err error) string {
	if err == nil {
		return ""
	}
	if d := Details(err); d.Message != "" {
		if d.Code != "" {
			return d.Code + ": " + d.Message
		}
		return d.Message
	}
	return err.Error()
}

// Guardrail reason codes surfaced in the Synced condition.
const (
	ReasonUnsupportedColumnTypeChange   = "UnsupportedColumnTypeChange"
	ReasonMaterializedViewAlterRejected = "MaterializedViewAlterRejected"
	ReasonImmutableFieldChanged         = "ImmutableFieldChanged"
	ReasonInvalidSpec                   = "InvalidSpec"
	ReasonColumnsMissingInSpec          = "ColumnsMissingInSpec"
	ReasonThrottled                     = "Throttled"
)

// Blocked is a guardrail error: the desired state cannot be applied without
// risking data loss or a Kusto rejection, so the reconciler must not send a
// command. It is returned from Observe, which makes crossplane-runtime set
// Synced=False with the reason and message, and stops Update from being called.
type Blocked struct {
	Reason  string
	Message string
}

func (b *Blocked) Error() string { return b.Reason + ": " + b.Message }

// NewBlocked builds a Blocked error.
func NewBlocked(reason, format string, args ...any) error {
	return &Blocked{Reason: reason, Message: fmt.Sprintf(format, args...)}
}

// IsBlocked reports whether err is a guardrail error.
func IsBlocked(err error) bool {
	var b *Blocked
	return errors.As(err, &b)
}
