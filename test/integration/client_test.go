//go:build integration

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

package integration

import (
	"context"
	"testing"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
)

func TestClientSmoke(t *testing.T) {
	res := run(t, ".show cluster")
	if len(res.Rows()) == 0 {
		t.Fatal(".show cluster returned no rows")
	}
	res = run(t, ".show database "+cmd.Ident(db)+" schema as json")
	if len(res.Rows()) != 1 || res.Rows()[0].String("DatabaseSchema") == "" {
		t.Fatal("schema as json returned nothing")
	}
}

// TestErrorClassification records what the emulator returns for the error
// classes the reconcilers depend on (spikes S2, S10).
func TestErrorClassification(t *testing.T) {
	ctx := context.Background()
	// S11: the emulator answers `.show table X cslschema` for a missing table
	// with success and no rows; a real cluster answers with an EntityNotFound
	// error. The reconcilers treat both as "does not exist".
	res, err := kc.Mgmt(ctx, db, cmd.New(".show table ", cmd.Ident("DoesNotExist"), " cslschema"))
	switch {
	case err == nil:
		t.Logf("S11: missing-table .show succeeded with %d rows", len(res.Rows()))
		if len(res.Rows()) != 0 {
			t.Errorf("missing table: expected no rows, got %d", len(res.Rows()))
		}
	case kerrors.Classify(err) == kerrors.NotFound:
		t.Logf("S11: missing-table .show failed with NotFound: %v", err)
	default:
		t.Errorf("missing table: Classify = %s (%+v): %v", kerrors.Classify(err), kerrors.Details(err), err)
	}

	_, err = kc.Mgmt(ctx, db, cmd.New(".show table ", cmd.Ident("DoesNotExist"), " cslschemaa"))
	if got := kerrors.Classify(err); got != kerrors.Permanent {
		t.Errorf("syntax error: Classify = %s (%+v)", got, kerrors.Details(err))
	}

	// .create table on an existing table succeeds (S10).
	run(t, ".create table "+cmd.Ident("Idempotent")+" (['A']:string)")
	if _, err := kc.Mgmt(ctx, db, cmd.New(".create table ", cmd.Ident("Idempotent"), " (['A']:string)")); err != nil {
		t.Errorf(".create table on an existing table must succeed: %v", err)
	}

	// .create function twice: record the already-exists code (S2).
	run(t, ".create function "+cmd.Ident("Dup")+"() { print 1 }")
	_, err = kc.Mgmt(ctx, db, cmd.New(".create function ", cmd.Ident("Dup"), "() { print 1 }"))
	if err == nil {
		t.Log("S2: .create function on an existing function succeeded (idempotent)")
	} else {
		d := kerrors.Details(err)
		t.Logf("S2: already-exists response: code=%q type=%q message=%q -> %s", d.Code, d.Type, d.Message, kerrors.Classify(err))
		if kerrors.Classify(err) != kerrors.AlreadyExists {
			t.Errorf("expected AlreadyExists classification, got %s", kerrors.Classify(err))
		}
	}
}
