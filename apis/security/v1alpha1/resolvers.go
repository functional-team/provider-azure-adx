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

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"

	adxv1alpha1 "github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/common/resolver"
)

// EntityTargets maps entity kinds to the managed resource types a role's
// entity.nameRef / entity.nameSelector may point to.
func EntityTargets() resolver.Targets {
	return resolver.Targets{
		common.EntityKindTable:            {Managed: func() resource.Managed { return &adxv1alpha1.Table{} }, List: func() resource.ManagedList { return &adxv1alpha1.TableList{} }},
		common.EntityKindMaterializedView: {Managed: func() resource.Managed { return &adxv1alpha1.MaterializedView{} }, List: func() resource.ManagedList { return &adxv1alpha1.MaterializedViewList{} }},
		common.EntityKindExternalTable:    {Managed: func() resource.Managed { return &adxv1alpha1.ExternalTable{} }, List: func() resource.ManagedList { return &adxv1alpha1.ExternalTableList{} }},
		common.EntityKindFunction:         {Managed: func() resource.Managed { return &adxv1alpha1.Function{} }, List: func() resource.ManagedList { return &adxv1alpha1.FunctionList{} }},
	}
}

// ResolveReferences of this SecurityRole.
func (mg *SecurityRole) ResolveReferences(ctx context.Context, c client.Reader) error {
	return resolver.ResolveEntity(ctx, c, mg, &mg.Spec.ForProvider.Entity, EntityTargets())
}
