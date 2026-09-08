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

// IngestionBatchingPolicyParameters are the configurable fields of a IngestionBatchingPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'Database']",message="IngestionBatchingPolicy can only target Table, Database"
type IngestionBatchingPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// MaximumBatchingTimeSpan is the maximum delay before a batch is sealed (10s to 15m).
	// +optional
	MaximumBatchingTimeSpan *common.Timespan `json:"maximumBatchingTimeSpan,omitempty"`

	// MaximumNumberOfItems is the maximum number of blobs per batch.
	// +optional
	MaximumNumberOfItems *int64 `json:"maximumNumberOfItems,omitempty"`

	// MaximumRawDataSizeMB is the maximum raw size of a batch in MB.
	// +optional
	MaximumRawDataSizeMB *int64 `json:"maximumRawDataSizeMB,omitempty"`
}

// A IngestionBatchingPolicySpec defines the desired state of a IngestionBatchingPolicy.
type IngestionBatchingPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              IngestionBatchingPolicyParameters `json:"forProvider"`
}

// A IngestionBatchingPolicyStatus represents the observed state of a IngestionBatchingPolicy.
type IngestionBatchingPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A IngestionBatchingPolicy tunes queued ingestion batching.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy ingestionbatching"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type IngestionBatchingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngestionBatchingPolicySpec   `json:"spec"`
	Status IngestionBatchingPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *IngestionBatchingPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *IngestionBatchingPolicy) GetPolicyObservation() PolicyObservation {
	return p.Status.AtProvider
}

// SetPolicyObservation stores the observed policy in the status.
func (p *IngestionBatchingPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// IngestionBatchingPolicyList contains a list of IngestionBatchingPolicy
type IngestionBatchingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngestionBatchingPolicy `json:"items"`
}

// IngestionBatchingPolicy type metadata.
var (
	IngestionBatchingPolicyKind             = reflect.TypeOf(IngestionBatchingPolicy{}).Name()
	IngestionBatchingPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: IngestionBatchingPolicyKind}.String()
	IngestionBatchingPolicyKindAPIVersion   = IngestionBatchingPolicyKind + "." + SchemeGroupVersion.String()
	IngestionBatchingPolicyGroupVersionKind = SchemeGroupVersion.WithKind(IngestionBatchingPolicyKind)
)

func init() {
	SchemeBuilder.Register(&IngestionBatchingPolicy{}, &IngestionBatchingPolicyList{})
}
