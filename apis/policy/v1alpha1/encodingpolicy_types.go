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

// EncodingPolicyParameters are the configurable fields of a EncodingPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'Database']",message="EncodingPolicy can only target Table, Database"
type EncodingPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Column applies the policy to one column of the table instead of the whole entity.
	// +optional
	Column *string `json:"column,omitempty"`

	// Type is the encoding policy type, e.g. Identifier, BigObject32 or Vector16.
	// +kubebuilder:validation:MinLength=1
	Type string `json:"type"`
}

// A EncodingPolicySpec defines the desired state of a EncodingPolicy.
type EncodingPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              EncodingPolicyParameters `json:"forProvider"`
}

// A EncodingPolicyStatus represents the observed state of a EncodingPolicy.
type EncodingPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A EncodingPolicy sets the encoding policy type of a database, table or single column.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy encoding"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type EncodingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EncodingPolicySpec   `json:"spec"`
	Status EncodingPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *EncodingPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *EncodingPolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *EncodingPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// EncodingPolicyList contains a list of EncodingPolicy
type EncodingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EncodingPolicy `json:"items"`
}

// EncodingPolicy type metadata.
var (
	EncodingPolicyKind             = reflect.TypeOf(EncodingPolicy{}).Name()
	EncodingPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: EncodingPolicyKind}.String()
	EncodingPolicyKindAPIVersion   = EncodingPolicyKind + "." + SchemeGroupVersion.String()
	EncodingPolicyGroupVersionKind = SchemeGroupVersion.WithKind(EncodingPolicyKind)
)

func init() {
	SchemeBuilder.Register(&EncodingPolicy{}, &EncodingPolicyList{})
}
