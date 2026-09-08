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

// ManagedIdentityEntry allows one identity for usages.
type ManagedIdentityEntry struct {
	// ObjectID of the managed identity, or "system".
	// +kubebuilder:validation:MinLength=1
	ObjectID string `json:"objectId"`
	// AllowedUsages of the identity.
	// +kubebuilder:validation:MinItems=1
	AllowedUsages []string `json:"allowedUsages"`
}

// ManagedIdentityPolicyParameters are the configurable fields of a ManagedIdentityPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Database']",message="ManagedIdentityPolicy can only target Database"
type ManagedIdentityPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Identities is the complete list of allowed identities.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Identities []ManagedIdentityEntry `json:"identities,omitempty"`
}

// A ManagedIdentityPolicySpec defines the desired state of a ManagedIdentityPolicy.
type ManagedIdentityPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              ManagedIdentityPolicyParameters `json:"forProvider"`
}

// A ManagedIdentityPolicyStatus represents the observed state of a ManagedIdentityPolicy.
type ManagedIdentityPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ManagedIdentityPolicy allows managed identities to be used for specific usages in a database.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy managed_identity"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type ManagedIdentityPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedIdentityPolicySpec   `json:"spec"`
	Status ManagedIdentityPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *ManagedIdentityPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *ManagedIdentityPolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *ManagedIdentityPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// ManagedIdentityPolicyList contains a list of ManagedIdentityPolicy
type ManagedIdentityPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedIdentityPolicy `json:"items"`
}

// ManagedIdentityPolicy type metadata.
var (
	ManagedIdentityPolicyKind             = reflect.TypeOf(ManagedIdentityPolicy{}).Name()
	ManagedIdentityPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: ManagedIdentityPolicyKind}.String()
	ManagedIdentityPolicyKindAPIVersion   = ManagedIdentityPolicyKind + "." + SchemeGroupVersion.String()
	ManagedIdentityPolicyGroupVersionKind = SchemeGroupVersion.WithKind(ManagedIdentityPolicyKind)
)

func init() {
	SchemeBuilder.Register(&ManagedIdentityPolicy{}, &ManagedIdentityPolicyList{})
}
