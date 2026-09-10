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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/common"
)

// MergeLookback limits the extents considered for merging.
type MergeLookback struct {
	// Kind of lookback.
	// +kubebuilder:validation:Enum=Default;All;HotCache;Custom
	Kind string `json:"kind"`
	// CustomPeriod is the lookback period for kind Custom.
	// +optional
	CustomPeriod *common.Timespan `json:"customPeriod,omitempty"`
}

// MergePolicyParameters are the configurable fields of a MergePolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'MaterializedView', 'Database']",message="MergePolicy can only target Table, MaterializedView, Database"
type MergePolicyParameters struct {
	PolicyTarget `json:",inline"`

	// RowCountUpperBoundForMerge caps the row count of a merged extent.
	// +optional
	RowCountUpperBoundForMerge *int64 `json:"rowCountUpperBoundForMerge,omitempty"`

	// OriginalSizeMBUpperBoundForMerge caps the original size of a merged extent.
	// +optional
	OriginalSizeMBUpperBoundForMerge *int64 `json:"originalSizeMBUpperBoundForMerge,omitempty"`

	// MaxExtentsToMerge caps the number of extents per merge.
	// +optional
	MaxExtentsToMerge *int64 `json:"maxExtentsToMerge,omitempty"`

	// MaxRangeInHours caps the creation time span of merged extents.
	// +optional
	MaxRangeInHours *int64 `json:"maxRangeInHours,omitempty"`

	// AllowRebuild enables rebuild operations.
	// +optional
	AllowRebuild *bool `json:"allowRebuild,omitempty"`

	// AllowMerge enables merge operations.
	// +optional
	AllowMerge *bool `json:"allowMerge,omitempty"`

	// Lookback limits which extents are considered.
	// +optional
	Lookback *MergeLookback `json:"lookback,omitempty"`
}

// A MergePolicySpec defines the desired state of a MergePolicy.
type MergePolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              MergePolicyParameters `json:"forProvider"`
}

// A MergePolicyStatus represents the observed state of a MergePolicy.
type MergePolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A MergePolicy controls extent merging.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy merge"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type MergePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MergePolicySpec   `json:"spec"`
	Status MergePolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *MergePolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *MergePolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *MergePolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// MergePolicyList contains a list of MergePolicy
type MergePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MergePolicy `json:"items"`
}

// MergePolicy type metadata.
var (
	MergePolicyKind             = reflect.TypeOf(MergePolicy{}).Name()
	MergePolicyGroupKind        = schema.GroupKind{Group: Group, Kind: MergePolicyKind}.String()
	MergePolicyKindAPIVersion   = MergePolicyKind + "." + SchemeGroupVersion.String()
	MergePolicyGroupVersionKind = SchemeGroupVersion.WithKind(MergePolicyKind)
)

func init() {
	SchemeBuilder.Register(&MergePolicy{}, &MergePolicyList{})
}
