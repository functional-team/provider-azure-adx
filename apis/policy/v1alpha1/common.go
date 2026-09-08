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

package v1alpha1

import (
	"github.com/functional-team/provider-azure-adx/apis/common"
)

// PolicyTarget addresses the database entity a policy applies to. It is
// embedded by every policy kind.
// +kubebuilder:validation:XValidation:rule="self.entity.kind == oldSelf.entity.kind",message="entity.kind is immutable"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name",message="entity.name is immutable once set"
type PolicyTarget struct {
	// Database that holds the entity. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// Entity the policy applies to. Which kinds are allowed depends on the
	// policy type and is validated per kind.
	Entity common.EntityReference `json:"entity"`
}

// PolicyObservation is the observed state shared by all policy kinds.
type PolicyObservation struct {
	// Entity the policy was read from, in Kusto notation.
	// +optional
	Entity string `json:"entity,omitempty"`
	// Policy is the policy JSON as reported by the cluster for this entity
	// (not the effective/inherited policy).
	// +optional
	Policy string `json:"policy,omitempty"`
}
