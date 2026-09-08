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

// RetentionPolicyParameters are the configurable fields of a RetentionPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'MaterializedView', 'Database']",message="RetentionPolicy can only target Table, MaterializedView, Database"
type RetentionPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// SoftDeletePeriod is how long data is kept before soft deletion, e.g. "365d" or "1000000d" for unlimited.
	// +optional
	SoftDeletePeriod *common.Timespan `json:"softDeletePeriod,omitempty"`

	// Recoverability enables recovery of deleted data for 14 days.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	Recoverability *string `json:"recoverability,omitempty"`
}

// A RetentionPolicySpec defines the desired state of a RetentionPolicy.
type RetentionPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              RetentionPolicyParameters `json:"forProvider"`
}

// A RetentionPolicyStatus represents the observed state of a RetentionPolicy.
type RetentionPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A RetentionPolicy controls how long data is kept (soft delete) and whether it is recoverable. Deleting the managed resource removes the policy from the entity, which falls back to the inherited policy.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy retention"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type RetentionPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RetentionPolicySpec   `json:"spec"`
	Status RetentionPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *RetentionPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *RetentionPolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *RetentionPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// RetentionPolicyList contains a list of RetentionPolicy
type RetentionPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RetentionPolicy `json:"items"`
}

// RetentionPolicy type metadata.
var (
	RetentionPolicyKind             = reflect.TypeOf(RetentionPolicy{}).Name()
	RetentionPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: RetentionPolicyKind}.String()
	RetentionPolicyKindAPIVersion   = RetentionPolicyKind + "." + SchemeGroupVersion.String()
	RetentionPolicyGroupVersionKind = SchemeGroupVersion.WithKind(RetentionPolicyKind)
)

func init() {
	SchemeBuilder.Register(&RetentionPolicy{}, &RetentionPolicyList{})
}
