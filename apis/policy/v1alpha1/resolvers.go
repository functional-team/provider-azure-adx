//go:build !angryjet

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

// Hand-written reference resolvers call into crossplane-runtime with the
// managed resource, which needs the angryjet-generated methods. The build tag
// hides this file while angryjet type-checks the package (see apis/generate.go).

package v1alpha1

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"

	adxv1alpha1 "github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/common/resolver"
)

// EntityTargets maps entity kinds to the managed resource types a policy's
// entity.nameRef / entity.nameSelector may point to.
func EntityTargets() resolver.Targets {
	return resolver.Targets{
		common.EntityKindTable:            {Managed: func() resource.Managed { return &adxv1alpha1.Table{} }, List: func() resource.ManagedList { return &adxv1alpha1.TableList{} }},
		common.EntityKindMaterializedView: {Managed: func() resource.Managed { return &adxv1alpha1.MaterializedView{} }, List: func() resource.ManagedList { return &adxv1alpha1.MaterializedViewList{} }},
		common.EntityKindExternalTable:    {Managed: func() resource.Managed { return &adxv1alpha1.ExternalTable{} }, List: func() resource.ManagedList { return &adxv1alpha1.ExternalTableList{} }},
		common.EntityKindFunction:         {Managed: func() resource.Managed { return &adxv1alpha1.Function{} }, List: func() resource.ManagedList { return &adxv1alpha1.FunctionList{} }},
	}
}

// ResolveReferences of this RetentionPolicy.
func (mg *RetentionPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this CachingPolicy.
func (mg *CachingPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this UpdatePolicy.
func (mg *UpdatePolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	if err := resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets()); err != nil {
		return err
	}
	r := reference.NewAPINamespacedResolver(c, mg)
	for i := range mg.Spec.ForProvider.Updates {
		u := &mg.Spec.ForProvider.Updates[i]
		rsp, err := r.Resolve(ctx, reference.NamespacedResolutionRequest{
			CurrentValue: u.Source,
			Reference:    u.SourceRef,
			Selector:     u.SourceSelector,
			To:           reference.To{Managed: &adxv1alpha1.Table{}, List: &adxv1alpha1.TableList{}},
			Extract:      reference.ExternalName(),
			Namespace:    mg.GetNamespace(),
		})
		if err != nil {
			return errors.Wrapf(err, "spec.forProvider.updates[%d].source", i)
		}
		u.Source = rsp.ResolvedValue
		u.SourceRef = rsp.ResolvedReference
	}
	return nil
}

// ResolveReferences of this RowLevelSecurityPolicy.
func (mg *RowLevelSecurityPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this IngestionBatchingPolicy.
func (mg *IngestionBatchingPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this StreamingIngestionPolicy.
func (mg *StreamingIngestionPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this MergePolicy.
func (mg *MergePolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this ShardingPolicy.
func (mg *ShardingPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this PartitioningPolicy.
func (mg *PartitioningPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this IngestionTimePolicy.
func (mg *IngestionTimePolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this AutoDeletePolicy.
func (mg *AutoDeletePolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this RestrictedViewAccessPolicy.
func (mg *RestrictedViewAccessPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this ExtentTagsRetentionPolicy.
func (mg *ExtentTagsRetentionPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this EncodingPolicy.
func (mg *EncodingPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this ManagedIdentityPolicy.
func (mg *ManagedIdentityPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this RowOrderPolicy.
func (mg *RowOrderPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}

// ResolveReferences of this QueryAccelerationPolicy.
func (mg *QueryAccelerationPolicy) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}
