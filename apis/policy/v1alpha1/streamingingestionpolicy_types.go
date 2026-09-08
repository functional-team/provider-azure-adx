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

// StreamingIngestionPolicyParameters are the configurable fields of a StreamingIngestionPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'Database']",message="StreamingIngestionPolicy can only target Table, Database"
type StreamingIngestionPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Enabled toggles streaming ingestion.
	Enabled bool `json:"enabled"`

	// HintAllocatedRate hints the expected ingestion rate in GB per hour (decimal number as string).
	// +kubebuilder:validation:Pattern=`^[0-9]+(\\.[0-9]+)?$`
	// +optional
	HintAllocatedRate *string `json:"hintAllocatedRate,omitempty"`
}

// A StreamingIngestionPolicySpec defines the desired state of a StreamingIngestionPolicy.
type StreamingIngestionPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              StreamingIngestionPolicyParameters `json:"forProvider"`
}

// A StreamingIngestionPolicyStatus represents the observed state of a StreamingIngestionPolicy.
type StreamingIngestionPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A StreamingIngestionPolicy enables streaming ingestion.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy streamingingestion"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type StreamingIngestionPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StreamingIngestionPolicySpec   `json:"spec"`
	Status StreamingIngestionPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *StreamingIngestionPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *StreamingIngestionPolicy) GetPolicyObservation() PolicyObservation {
	return p.Status.AtProvider
}

// SetPolicyObservation stores the observed policy in the status.
func (p *StreamingIngestionPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// StreamingIngestionPolicyList contains a list of StreamingIngestionPolicy
type StreamingIngestionPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StreamingIngestionPolicy `json:"items"`
}

// StreamingIngestionPolicy type metadata.
var (
	StreamingIngestionPolicyKind             = reflect.TypeOf(StreamingIngestionPolicy{}).Name()
	StreamingIngestionPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: StreamingIngestionPolicyKind}.String()
	StreamingIngestionPolicyKindAPIVersion   = StreamingIngestionPolicyKind + "." + SchemeGroupVersion.String()
	StreamingIngestionPolicyGroupVersionKind = SchemeGroupVersion.WithKind(StreamingIngestionPolicyKind)
)

func init() {
	SchemeBuilder.Register(&StreamingIngestionPolicy{}, &StreamingIngestionPolicyList{})
}
