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

// IngestionTimePolicyParameters are the configurable fields of a IngestionTimePolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table']",message="IngestionTimePolicy can only target Table"
type IngestionTimePolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Enabled toggles the policy.
	Enabled bool `json:"enabled"`
}

// A IngestionTimePolicySpec defines the desired state of a IngestionTimePolicy.
type IngestionTimePolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              IngestionTimePolicyParameters `json:"forProvider"`
}

// A IngestionTimePolicyStatus represents the observed state of a IngestionTimePolicy.
type IngestionTimePolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A IngestionTimePolicy adds the hidden ingestion_time() column.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy ingestiontime"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type IngestionTimePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngestionTimePolicySpec   `json:"spec"`
	Status IngestionTimePolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *IngestionTimePolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *IngestionTimePolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *IngestionTimePolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// IngestionTimePolicyList contains a list of IngestionTimePolicy
type IngestionTimePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngestionTimePolicy `json:"items"`
}

// IngestionTimePolicy type metadata.
var (
	IngestionTimePolicyKind             = reflect.TypeOf(IngestionTimePolicy{}).Name()
	IngestionTimePolicyGroupKind        = schema.GroupKind{Group: Group, Kind: IngestionTimePolicyKind}.String()
	IngestionTimePolicyKindAPIVersion   = IngestionTimePolicyKind + "." + SchemeGroupVersion.String()
	IngestionTimePolicyGroupVersionKind = SchemeGroupVersion.WithKind(IngestionTimePolicyKind)
)

func init() {
	SchemeBuilder.Register(&IngestionTimePolicy{}, &IngestionTimePolicyList{})
}
