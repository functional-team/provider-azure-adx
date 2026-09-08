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

// Package resolver resolves EntityReference name references. It lives in its
// own package because the target constructors are functions, which the
// deepcopy generator must not see.
package resolver

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"

	"github.com/functional-team/provider-azure-adx/apis/common"
)

// Target describes how to look up managed resources of one entity kind.
type Target struct {
	Managed func() resource.Managed
	List    func() resource.ManagedList
}

// Targets maps entity kinds to their managed resource types.
type Targets map[common.EntityKind]Target

// ResolveEntity resolves ref.nameRef / ref.nameSelector into ref.name using
// the target registered for ref.kind. Database targets and plain names are
// no-ops.
func ResolveEntity(ctx context.Context, c client.Reader, from resource.Managed, ref *common.EntityReference, targets Targets) error {
	if ref == nil || ref.Kind == common.EntityKindDatabase {
		return nil
	}
	if ref.NameRef == nil && ref.NameSelector == nil {
		return nil
	}
	t, ok := targets[ref.Kind]
	if !ok {
		return errors.Errorf("entity kind %s cannot be resolved by reference; set entity.name", ref.Kind)
	}
	r := reference.NewAPINamespacedResolver(c, from)
	rsp, err := r.Resolve(ctx, reference.NamespacedResolutionRequest{
		CurrentValue: reference.FromPtrValue(ref.Name),
		Reference:    ref.NameRef,
		Selector:     ref.NameSelector,
		To:           reference.To{Managed: t.Managed(), List: t.List()},
		Extract:      reference.ExternalName(),
		Namespace:    from.GetNamespace(),
	})
	if err != nil {
		return errors.Wrap(err, "entity.name")
	}
	ref.Name = reference.ToPtrValue(rsp.ResolvedValue)
	ref.NameRef = rsp.ResolvedReference
	return nil
}
