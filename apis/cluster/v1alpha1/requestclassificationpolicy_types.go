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

// RequestClassificationPolicyParameters are the configurable fields of a RequestClassificationPolicy.
type RequestClassificationPolicyParameters struct {
	// Enabled toggles the policy. Defaults to true.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Query is the classification function body, e.g. iff(request_properties.current_principal == "aadapp=...", "Ingestion", "default").
	// +kubebuilder:validation:MinLength=1
	Query string `json:"query"`
}

// A RequestClassificationPolicySpec defines the desired state of a RequestClassificationPolicy.
type RequestClassificationPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              RequestClassificationPolicyParameters `json:"forProvider"`
}

// A RequestClassificationPolicyStatus represents the observed state of a RequestClassificationPolicy.
type RequestClassificationPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ClusterPolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A RequestClassificationPolicy routes requests to workload groups through a classification function (KQL). The query is compared with the two stage KQL normalization.
// There is exactly one request_classification policy per cluster, so one managed resource
// per cluster (ProviderConfig) is expected; two managed resources on the same
// cluster policy fight over it and are a user error. The provider principal
// needs the AllDatabasesAdmin cluster role for cluster policies.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxcluster}
type RequestClassificationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RequestClassificationPolicySpec   `json:"spec"`
	Status RequestClassificationPolicyStatus `json:"status,omitempty"`
}

// GetClusterPolicyObservation returns the observed policy.
func (p *RequestClassificationPolicy) GetClusterPolicyObservation() ClusterPolicyObservation {
	return p.Status.AtProvider
}

// SetClusterPolicyObservation stores the observed policy in the status.
func (p *RequestClassificationPolicy) SetClusterPolicyObservation(o ClusterPolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// RequestClassificationPolicyList contains a list of RequestClassificationPolicy
type RequestClassificationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RequestClassificationPolicy `json:"items"`
}

// RequestClassificationPolicy type metadata.
var (
	RequestClassificationPolicyKind             = reflect.TypeOf(RequestClassificationPolicy{}).Name()
	RequestClassificationPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: RequestClassificationPolicyKind}.String()
	RequestClassificationPolicyKindAPIVersion   = RequestClassificationPolicyKind + "." + SchemeGroupVersion.String()
	RequestClassificationPolicyGroupVersionKind = SchemeGroupVersion.WithKind(RequestClassificationPolicyKind)
)

func init() {
	SchemeBuilder.Register(&RequestClassificationPolicy{}, &RequestClassificationPolicyList{})
}
