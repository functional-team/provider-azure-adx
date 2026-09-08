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

// QueryAccelerationPolicyParameters are the configurable fields of a QueryAccelerationPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['ExternalTable']",message="QueryAccelerationPolicy can only target ExternalTable"
type QueryAccelerationPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Enabled toggles the policy.
	Enabled bool `json:"enabled"`

	// Hot is the period of data kept accelerated.
	// +optional
	Hot *common.Timespan `json:"hot,omitempty"`

	// MaxAge caps how stale accelerated data may be.
	// +optional
	MaxAge *common.Timespan `json:"maxAge,omitempty"`
}

// A QueryAccelerationPolicySpec defines the desired state of a QueryAccelerationPolicy.
type QueryAccelerationPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              QueryAccelerationPolicyParameters `json:"forProvider"`
}

// A QueryAccelerationPolicyStatus represents the observed state of a QueryAccelerationPolicy.
type QueryAccelerationPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A QueryAccelerationPolicy accelerates queries over an external delta table by caching recent data (Tier 3).
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy query_acceleration"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type QueryAccelerationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueryAccelerationPolicySpec   `json:"spec"`
	Status QueryAccelerationPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *QueryAccelerationPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *QueryAccelerationPolicy) GetPolicyObservation() PolicyObservation {
	return p.Status.AtProvider
}

// SetPolicyObservation stores the observed policy in the status.
func (p *QueryAccelerationPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// QueryAccelerationPolicyList contains a list of QueryAccelerationPolicy
type QueryAccelerationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueryAccelerationPolicy `json:"items"`
}

// QueryAccelerationPolicy type metadata.
var (
	QueryAccelerationPolicyKind             = reflect.TypeOf(QueryAccelerationPolicy{}).Name()
	QueryAccelerationPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: QueryAccelerationPolicyKind}.String()
	QueryAccelerationPolicyKindAPIVersion   = QueryAccelerationPolicyKind + "." + SchemeGroupVersion.String()
	QueryAccelerationPolicyGroupVersionKind = SchemeGroupVersion.WithKind(QueryAccelerationPolicyKind)
)

func init() {
	SchemeBuilder.Register(&QueryAccelerationPolicy{}, &QueryAccelerationPolicyList{})
}
