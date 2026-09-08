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

// RestrictedViewAccessPolicyParameters are the configurable fields of a RestrictedViewAccessPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table']",message="RestrictedViewAccessPolicy can only target Table"
type RestrictedViewAccessPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Enabled toggles the policy.
	Enabled bool `json:"enabled"`
}

// A RestrictedViewAccessPolicySpec defines the desired state of a RestrictedViewAccessPolicy.
type RestrictedViewAccessPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              RestrictedViewAccessPolicyParameters `json:"forProvider"`
}

// A RestrictedViewAccessPolicyStatus represents the observed state of a RestrictedViewAccessPolicy.
type RestrictedViewAccessPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A RestrictedViewAccessPolicy restricts a table to principals with the UnrestrictedViewer role.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy restricted_view_access"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type RestrictedViewAccessPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RestrictedViewAccessPolicySpec   `json:"spec"`
	Status RestrictedViewAccessPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *RestrictedViewAccessPolicy) GetTarget() *PolicyTarget {
	return &p.Spec.ForProvider.PolicyTarget
}

// GetPolicyObservation returns the observed policy.
func (p *RestrictedViewAccessPolicy) GetPolicyObservation() PolicyObservation {
	return p.Status.AtProvider
}

// SetPolicyObservation stores the observed policy in the status.
func (p *RestrictedViewAccessPolicy) SetPolicyObservation(o PolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// RestrictedViewAccessPolicyList contains a list of RestrictedViewAccessPolicy
type RestrictedViewAccessPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RestrictedViewAccessPolicy `json:"items"`
}

// RestrictedViewAccessPolicy type metadata.
var (
	RestrictedViewAccessPolicyKind             = reflect.TypeOf(RestrictedViewAccessPolicy{}).Name()
	RestrictedViewAccessPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: RestrictedViewAccessPolicyKind}.String()
	RestrictedViewAccessPolicyKindAPIVersion   = RestrictedViewAccessPolicyKind + "." + SchemeGroupVersion.String()
	RestrictedViewAccessPolicyGroupVersionKind = SchemeGroupVersion.WithKind(RestrictedViewAccessPolicyKind)
)

func init() {
	SchemeBuilder.Register(&RestrictedViewAccessPolicy{}, &RestrictedViewAccessPolicyList{})
}
