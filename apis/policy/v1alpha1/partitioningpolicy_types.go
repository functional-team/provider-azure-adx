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

// PartitionKey is one partition key.
type PartitionKey struct {
	// ColumnName of the key.
	// +kubebuilder:validation:MinLength=1
	ColumnName string `json:"columnName"`
	// Kind of the key.
	// +kubebuilder:validation:Enum=Hash;UniformRange
	Kind string `json:"kind"`
	// Properties of the key.
	Properties PartitionKeyProperties `json:"properties"`
}

// PartitionKeyProperties are the kind specific settings of a partition key.
type PartitionKeyProperties struct {
	// Function is the hash function (Hash keys).
	// +kubebuilder:validation:Enum=XxHash64;StringHash
	// +optional
	Function *string `json:"function,omitempty"`
	// MaxPartitionCount is the number of hash partitions (Hash keys).
	// +optional
	MaxPartitionCount *int64 `json:"maxPartitionCount,omitempty"`
	// Seed of the hash function (Hash keys).
	// +optional
	Seed *int64 `json:"seed,omitempty"`
	// PartitionAssignmentMode for Hash keys.
	// +kubebuilder:validation:Enum=Default;Uniform
	// +optional
	PartitionAssignmentMode *string `json:"partitionAssignmentMode,omitempty"`
	// Reference datetime (ISO 8601) that anchors the ranges (UniformRange keys).
	// +optional
	Reference *string `json:"reference,omitempty"`
	// RangeSize of a partition (UniformRange keys).
	// +optional
	RangeSize *common.Timespan `json:"rangeSize,omitempty"`
	// OverrideCreationTime sets the extent creation time to the range start.
	// +optional
	OverrideCreationTime *bool `json:"overrideCreationTime,omitempty"`
}

// PartitioningPolicyParameters are the configurable fields of a PartitioningPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'MaterializedView']",message="PartitioningPolicy can only target Table, MaterializedView"
type PartitioningPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// PartitionKeys are the partition keys (at most one hash and one uniform range key).
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=2
	// +optional
	PartitionKeys []PartitionKey `json:"partitionKeys,omitempty"`

	// EffectiveDateTime applies the policy to extents created after this ISO 8601 datetime.
	// +optional
	EffectiveDateTime *string `json:"effectiveDateTime,omitempty"`
}

// A PartitioningPolicySpec defines the desired state of a PartitioningPolicy.
type PartitioningPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              PartitioningPolicyParameters `json:"forProvider"`
}

// A PartitioningPolicyStatus represents the observed state of a PartitioningPolicy.
type PartitioningPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A PartitioningPolicy partitions extents by hash or uniform range keys.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy partitioning"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type PartitioningPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PartitioningPolicySpec   `json:"spec"`
	Status PartitioningPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *PartitioningPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *PartitioningPolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *PartitioningPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// PartitioningPolicyList contains a list of PartitioningPolicy
type PartitioningPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PartitioningPolicy `json:"items"`
}

// PartitioningPolicy type metadata.
var (
	PartitioningPolicyKind             = reflect.TypeOf(PartitioningPolicy{}).Name()
	PartitioningPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: PartitioningPolicyKind}.String()
	PartitioningPolicyKindAPIVersion   = PartitioningPolicyKind + "." + SchemeGroupVersion.String()
	PartitioningPolicyGroupVersionKind = SchemeGroupVersion.WithKind(PartitioningPolicyKind)
)

func init() {
	SchemeBuilder.Register(&PartitioningPolicy{}, &PartitioningPolicyList{})
}
