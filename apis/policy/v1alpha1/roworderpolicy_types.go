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

// RowOrderColumn is one sort column.
type RowOrderColumn struct {
	// Name of the column.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Direction of the sort.
	// +kubebuilder:validation:Enum=asc;desc
	Direction string `json:"direction"`
}

// RowOrderPolicyParameters are the configurable fields of a RowOrderPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'MaterializedView']",message="RowOrderPolicy can only target Table, MaterializedView"
type RowOrderPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Columns and their sort direction.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Columns []RowOrderColumn `json:"columns,omitempty"`
}

// A RowOrderPolicySpec defines the desired state of a RowOrderPolicy.
type RowOrderPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              RowOrderPolicyParameters `json:"forProvider"`
}

// A RowOrderPolicyStatus represents the observed state of a RowOrderPolicy.
type RowOrderPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A RowOrderPolicy orders rows inside extents by columns (Tier 3, unverified against a cluster).
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy roworder"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type RowOrderPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RowOrderPolicySpec   `json:"spec"`
	Status RowOrderPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *RowOrderPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *RowOrderPolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *RowOrderPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// RowOrderPolicyList contains a list of RowOrderPolicy
type RowOrderPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RowOrderPolicy `json:"items"`
}

// RowOrderPolicy type metadata.
var (
	RowOrderPolicyKind             = reflect.TypeOf(RowOrderPolicy{}).Name()
	RowOrderPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: RowOrderPolicyKind}.String()
	RowOrderPolicyKindAPIVersion   = RowOrderPolicyKind + "." + SchemeGroupVersion.String()
	RowOrderPolicyGroupVersionKind = SchemeGroupVersion.WithKind(RowOrderPolicyKind)
)

func init() {
	SchemeBuilder.Register(&RowOrderPolicy{}, &RowOrderPolicyList{})
}
