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
)

// ShardingPolicyParameters are the configurable fields of a ShardingPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'MaterializedView', 'Database']",message="ShardingPolicy can only target Table, MaterializedView, Database"
type ShardingPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// MaxRowCount caps the rows per extent.
	// +optional
	MaxRowCount *int64 `json:"maxRowCount,omitempty"`

	// MaxExtentSizeInMb caps the compressed size per extent.
	// +optional
	MaxExtentSizeInMb *int64 `json:"maxExtentSizeInMb,omitempty"`

	// MaxOriginalSizeInMb caps the original size per extent.
	// +optional
	MaxOriginalSizeInMb *int64 `json:"maxOriginalSizeInMb,omitempty"`

	// ShardEngineMaxRowCount caps rows per extent created by the shard engine.
	// +optional
	ShardEngineMaxRowCount *int64 `json:"shardEngineMaxRowCount,omitempty"`

	// ShardEngineMaxExtentSizeInMb caps the compressed size per shard engine extent.
	// +optional
	ShardEngineMaxExtentSizeInMb *int64 `json:"shardEngineMaxExtentSizeInMb,omitempty"`

	// ShardEngineMaxOriginalSizeInMb caps the original size per shard engine extent.
	// +optional
	ShardEngineMaxOriginalSizeInMb *int64 `json:"shardEngineMaxOriginalSizeInMb,omitempty"`
}

// A ShardingPolicySpec defines the desired state of a ShardingPolicy.
type ShardingPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              ShardingPolicyParameters `json:"forProvider"`
}

// A ShardingPolicyStatus represents the observed state of a ShardingPolicy.
type ShardingPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ShardingPolicy controls extent (shard) sizing.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy sharding"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type ShardingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ShardingPolicySpec   `json:"spec"`
	Status ShardingPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *ShardingPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *ShardingPolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *ShardingPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// ShardingPolicyList contains a list of ShardingPolicy
type ShardingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ShardingPolicy `json:"items"`
}

// ShardingPolicy type metadata.
var (
	ShardingPolicyKind             = reflect.TypeOf(ShardingPolicy{}).Name()
	ShardingPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: ShardingPolicyKind}.String()
	ShardingPolicyKindAPIVersion   = ShardingPolicyKind + "." + SchemeGroupVersion.String()
	ShardingPolicyGroupVersionKind = SchemeGroupVersion.WithKind(ShardingPolicyKind)
)

func init() {
	SchemeBuilder.Register(&ShardingPolicy{}, &ShardingPolicyList{})
}
