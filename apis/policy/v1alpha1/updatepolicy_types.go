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

// UpdatePolicyEntry is one update policy of a table.
type UpdatePolicyEntry struct {
	// Enabled toggles the policy. Defaults to true.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Source is the table whose ingestion triggers the query.
	// +kubebuilder:validation:MinLength=1
	Source string `json:"source"`
	// SourceRef references a Table managed resource as source.
	// +optional
	SourceRef *xpv2.NamespacedReference `json:"sourceRef,omitempty"`
	// SourceSelector selects a Table managed resource as source.
	// +optional
	SourceSelector *xpv2.NamespacedSelector `json:"sourceSelector,omitempty"`
	// Query is the KQL query (or function call) producing rows for the target table.
	// +kubebuilder:validation:MinLength=1
	Query string `json:"query"`
	// Transactional fails the source ingestion when the update fails.
	// +optional
	Transactional *bool `json:"transactional,omitempty"`
	// PropagateIngestionProperties copies extent tags and creation time.
	// +optional
	PropagateIngestionProperties *bool `json:"propagateIngestionProperties,omitempty"`
	// ManagedIdentity ("system" or an object id) runs the query with a managed identity.
	// +optional
	ManagedIdentity *string `json:"managedIdentity,omitempty"`
}

// UpdatePolicyParameters are the configurable fields of a UpdatePolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table']",message="UpdatePolicy can only target Table"
type UpdatePolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Updates is the complete list of update policies of the table.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Updates []UpdatePolicyEntry `json:"updates,omitempty"`
}

// A UpdatePolicySpec defines the desired state of a UpdatePolicy.
type UpdatePolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              UpdatePolicyParameters `json:"forProvider"`
}

// A UpdatePolicyStatus represents the observed state of a UpdatePolicy.
type UpdatePolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A UpdatePolicy attaches update policies (query-driven ingestion from a source table) to a target table. The list is authoritative for the table.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy update"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type UpdatePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpdatePolicySpec   `json:"spec"`
	Status UpdatePolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *UpdatePolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *UpdatePolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *UpdatePolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// UpdatePolicyList contains a list of UpdatePolicy
type UpdatePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpdatePolicy `json:"items"`
}

// UpdatePolicy type metadata.
var (
	UpdatePolicyKind             = reflect.TypeOf(UpdatePolicy{}).Name()
	UpdatePolicyGroupKind        = schema.GroupKind{Group: Group, Kind: UpdatePolicyKind}.String()
	UpdatePolicyKindAPIVersion   = UpdatePolicyKind + "." + SchemeGroupVersion.String()
	UpdatePolicyGroupVersionKind = SchemeGroupVersion.WithKind(UpdatePolicyKind)
)

func init() {
	SchemeBuilder.Register(&UpdatePolicy{}, &UpdatePolicyList{})
}
