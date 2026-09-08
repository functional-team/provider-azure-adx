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

// HotWindow is a datetime range kept in the hot cache.
type HotWindow struct {
	// MinValue is the inclusive start, ISO 8601 (e.g. 2026-01-01T00:00:00Z).
	MinValue string `json:"minValue"`
	// MaxValue is the exclusive end, ISO 8601.
	MaxValue string `json:"maxValue"`
}

// CachingPolicyParameters are the configurable fields of a CachingPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'MaterializedView', 'Database']",message="CachingPolicy can only target Table, MaterializedView, Database"
type CachingPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Hot is the hot cache span applied to data and index, e.g. "31d".
	Hot common.Timespan `json:"hot"`

	// HotWindows are additional datetime ranges kept in hot cache.
	// +optional
	HotWindows []HotWindow `json:"hotWindows,omitempty"`
}

// A CachingPolicySpec defines the desired state of a CachingPolicy.
type CachingPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              CachingPolicyParameters `json:"forProvider"`
}

// A CachingPolicyStatus represents the observed state of a CachingPolicy.
type CachingPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A CachingPolicy sets the hot cache period and optional hot windows.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy caching"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type CachingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CachingPolicySpec   `json:"spec"`
	Status CachingPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *CachingPolicy) GetTarget() *PolicyTarget { return &p.Spec.ForProvider.PolicyTarget }

// GetPolicyObservation returns the observed policy.
func (p *CachingPolicy) GetPolicyObservation() PolicyObservation { return p.Status.AtProvider }

// SetPolicyObservation stores the observed policy in the status.
func (p *CachingPolicy) SetPolicyObservation(o PolicyObservation) { p.Status.AtProvider = o }

// +kubebuilder:object:root=true

// CachingPolicyList contains a list of CachingPolicy
type CachingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CachingPolicy `json:"items"`
}

// CachingPolicy type metadata.
var (
	CachingPolicyKind             = reflect.TypeOf(CachingPolicy{}).Name()
	CachingPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: CachingPolicyKind}.String()
	CachingPolicyKindAPIVersion   = CachingPolicyKind + "." + SchemeGroupVersion.String()
	CachingPolicyGroupVersionKind = SchemeGroupVersion.WithKind(CachingPolicyKind)
)

func init() {
	SchemeBuilder.Register(&CachingPolicy{}, &CachingPolicyList{})
}
