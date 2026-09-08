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

// AutoDeletePolicyParameters are the configurable fields of a AutoDeletePolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table']",message="AutoDeletePolicy can only target Table"
type AutoDeletePolicyParameters struct {
	PolicyTarget `json:",inline"`

	// ExpiryDate (ISO 8601 date or datetime) after which the table is dropped.
	// +kubebuilder:validation:MinLength=1
	ExpiryDate string `json:"expiryDate"`

	// DeleteIfNotEmpty also drops the table when it still holds data.
	// +optional
	DeleteIfNotEmpty *bool `json:"deleteIfNotEmpty,omitempty"`
}

// A AutoDeletePolicySpec defines the desired state of a AutoDeletePolicy.
type AutoDeletePolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              AutoDeletePolicyParameters `json:"forProvider"`
}

// A AutoDeletePolicyStatus represents the observed state of a AutoDeletePolicy.
type AutoDeletePolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A AutoDeletePolicy drops the table at an expiry date.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy auto_delete"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type AutoDeletePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoDeletePolicySpec   `json:"spec"`
	Status AutoDeletePolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *AutoDeletePolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *AutoDeletePolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *AutoDeletePolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// AutoDeletePolicyList contains a list of AutoDeletePolicy
type AutoDeletePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutoDeletePolicy `json:"items"`
}

// AutoDeletePolicy type metadata.
var (
	AutoDeletePolicyKind             = reflect.TypeOf(AutoDeletePolicy{}).Name()
	AutoDeletePolicyGroupKind        = schema.GroupKind{Group: Group, Kind: AutoDeletePolicyKind}.String()
	AutoDeletePolicyKindAPIVersion   = AutoDeletePolicyKind + "." + SchemeGroupVersion.String()
	AutoDeletePolicyGroupVersionKind = SchemeGroupVersion.WithKind(AutoDeletePolicyKind)
)

func init() {
	SchemeBuilder.Register(&AutoDeletePolicy{}, &AutoDeletePolicyList{})
}
