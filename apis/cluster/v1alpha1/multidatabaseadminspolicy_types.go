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

// MultiDatabaseAdminsPolicyParameters are the configurable fields of a MultiDatabaseAdminsPolicy.
type MultiDatabaseAdminsPolicyParameters struct {
	// Principals is the complete list of allowed principals.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Principals []common.Principal `json:"principals,omitempty"`
}

// A MultiDatabaseAdminsPolicySpec defines the desired state of a MultiDatabaseAdminsPolicy.
type MultiDatabaseAdminsPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              MultiDatabaseAdminsPolicyParameters `json:"forProvider"`
}

// A MultiDatabaseAdminsPolicyStatus represents the observed state of a MultiDatabaseAdminsPolicy.
type MultiDatabaseAdminsPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ClusterPolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A MultiDatabaseAdminsPolicy lists principals that may administer several databases without holding AllDatabasesAdmin. The policy JSON shape is not verified against a cluster (spike), so the comparison relies on the hash annotations.
// There is exactly one multidatabaseadmins policy per cluster, so one managed resource
// per cluster (ProviderConfig) is expected; two managed resources on the same
// cluster policy fight over it and are a user error. The provider principal
// needs the AllDatabasesAdmin cluster role for cluster policies.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxcluster}
type MultiDatabaseAdminsPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultiDatabaseAdminsPolicySpec   `json:"spec"`
	Status MultiDatabaseAdminsPolicyStatus `json:"status,omitempty"`
}

// GetClusterPolicyObservation returns the observed policy.
func (p *MultiDatabaseAdminsPolicy) GetClusterPolicyObservation() ClusterPolicyObservation {
	return p.Status.AtProvider
}

// SetClusterPolicyObservation stores the observed policy in the status.
func (p *MultiDatabaseAdminsPolicy) SetClusterPolicyObservation(o ClusterPolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// MultiDatabaseAdminsPolicyList contains a list of MultiDatabaseAdminsPolicy
type MultiDatabaseAdminsPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultiDatabaseAdminsPolicy `json:"items"`
}

// MultiDatabaseAdminsPolicy type metadata.
var (
	MultiDatabaseAdminsPolicyKind             = reflect.TypeOf(MultiDatabaseAdminsPolicy{}).Name()
	MultiDatabaseAdminsPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: MultiDatabaseAdminsPolicyKind}.String()
	MultiDatabaseAdminsPolicyKindAPIVersion   = MultiDatabaseAdminsPolicyKind + "." + SchemeGroupVersion.String()
	MultiDatabaseAdminsPolicyGroupVersionKind = SchemeGroupVersion.WithKind(MultiDatabaseAdminsPolicyKind)
)

func init() {
	SchemeBuilder.Register(&MultiDatabaseAdminsPolicy{}, &MultiDatabaseAdminsPolicyList{})
}
