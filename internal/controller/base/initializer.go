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

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

// SpecNamer is implemented by managed resources with an optional
// spec.forProvider.name that names the Kusto entity.
type SpecNamer interface {
	SpecName() string
}

// SpecNameAsExternalName copies spec.forProvider.name into the external-name
// annotation when the annotation is empty. It runs before
// managed.NameAsExternalName, so metadata.name is only the fallback. Kusto
// names are usually PascalCase and would not be valid Kubernetes names.
type SpecNameAsExternalName struct{ client client.Client }

// NewSpecNameAsExternalName returns the initializer.
func NewSpecNameAsExternalName(c client.Client) *SpecNameAsExternalName {
	return &SpecNameAsExternalName{client: c}
}

// Initialize implements managed.Initializer.
func (i *SpecNameAsExternalName) Initialize(ctx context.Context, mg resource.Managed) error {
	if meta.GetExternalName(mg) != "" {
		return nil
	}
	n, ok := mg.(SpecNamer)
	if !ok || n.SpecName() == "" {
		return nil
	}
	meta.SetExternalName(mg, n.SpecName())
	return errors.Wrap(i.client.Update(ctx, mg), "cannot update managed resource")
}
