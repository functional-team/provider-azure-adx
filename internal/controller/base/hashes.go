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

package base

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"

	"github.com/functional-team/provider-azure-adx/internal/normalize"
)

// Annotations written by the provider.
const (
	// AnnotationAppliedHash is the hash of the normalized desired KQL texts
	// that were last written to the cluster.
	AnnotationAppliedHash = "adx.functional.team/applied-hash"
	// AnnotationObservedHash is the hash of the normalized texts read back
	// from the cluster right after that write.
	AnnotationObservedHash = "adx.functional.team/observed-hash"
	// AnnotationOperationID holds the id of a running async operation
	// (materialized view backfill).
	AnnotationOperationID = "adx.functional.team/operation-id"
	// AnnotationSecretHash is the hash of resolved write-only secrets
	// (external table connection strings) last written.
	AnnotationSecretHash = "adx.functional.team/secret-hash"
)

// Hashes returns the applied and observed hash annotations.
func Hashes(o metav1.Object) (applied, observed string) {
	a := o.GetAnnotations()
	return a[AnnotationAppliedHash], a[AnnotationObservedHash]
}

// SetHashes records the applied/observed hashes for the given normalized
// texts on the object (in memory).
func SetHashes(o metav1.Object, desired, observed []string) {
	meta.AddAnnotations(o, map[string]string{
		AnnotationAppliedHash:  normalize.Hash(desired...),
		AnnotationObservedHash: normalize.Hash(observed...),
	})
}

// TextUpToDate implements the two stage comparison of KQL texts. desired and
// observed are the normalized (stage 1) texts of all free text fields in a
// fixed order. Stage 2 accepts a mismatch when the spec hash still equals
// what was last applied and the cluster text still equals what was read back
// after that write: Kusto reformatted the text, nobody changed it since.
func TextUpToDate(o metav1.Object, desired, observed []string) bool {
	if equal(desired, observed) {
		return true
	}
	applied, obs := Hashes(o)
	return applied != "" && applied == normalize.Hash(desired...) && obs == normalize.Hash(observed...)
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// PersistAnnotations writes the in-memory annotations of o to the API server
// with a merge patch. The patch is applied to a copy so the in-memory status
// (set during Observe) survives; only the resourceVersion is copied back so a
// following status update does not conflict. Create does not need this (the
// reconciler persists annotations after Create), Update does.
func PersistAnnotations(ctx context.Context, c client.Client, o client.Object) error {
	patch := map[string]any{"metadata": map[string]any{"annotations": o.GetAnnotations()}}
	data, err := json.Marshal(patch)
	if err != nil {
		return errors.Wrap(err, "cannot marshal annotation patch")
	}
	cp, ok := o.DeepCopyObject().(client.Object)
	if !ok {
		return errors.New("object is not a client.Object")
	}
	if err := c.Patch(ctx, cp, client.RawPatch(types.MergePatchType, data)); err != nil {
		return errors.Wrap(err, "cannot persist annotations")
	}
	o.SetResourceVersion(cp.GetResourceVersion())
	return nil
}
