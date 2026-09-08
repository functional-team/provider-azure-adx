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

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// CapacityPolicyParameters are the configurable fields of a CapacityPolicy.
type CapacityPolicyParameters struct {
	// Policy is the capacity policy JSON as documented by Kusto (IngestionCapacity, ExtentsMergeCapacity, ExportCapacity, MaterializedViewsCapacity, ...). Only the keys present here are compared with the cluster; unknown keys are passed through.
	// +kubebuilder:pruning:PreserveUnknownFields
	Policy apiextensionsv1.JSON `json:"policy"`
}

// A CapacityPolicySpec defines the desired state of a CapacityPolicy.
type CapacityPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              CapacityPolicyParameters `json:"forProvider"`
}

// A CapacityPolicyStatus represents the observed state of a CapacityPolicy.
type CapacityPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ClusterPolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A CapacityPolicy tunes the concurrency limits of cluster operations (ingestion, merge, export, materialized views, ...). The policy always exists on a cluster, so the managed resource only ever alters it; deleting the managed resource leaves the last applied values in place (Kusto has no delete for this policy).
// There is exactly one capacity policy per cluster, so one managed resource
// per cluster (ProviderConfig) is expected; two managed resources on the same
// cluster policy fight over it and are a user error. The provider principal
// needs the AllDatabasesAdmin cluster role for cluster policies.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxcluster}
type CapacityPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CapacityPolicySpec   `json:"spec"`
	Status CapacityPolicyStatus `json:"status,omitempty"`
}

// GetClusterPolicyObservation returns the observed policy.
func (p *CapacityPolicy) GetClusterPolicyObservation() ClusterPolicyObservation {
	return p.Status.AtProvider
}

// SetClusterPolicyObservation stores the observed policy in the status.
func (p *CapacityPolicy) SetClusterPolicyObservation(o ClusterPolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// CapacityPolicyList contains a list of CapacityPolicy
type CapacityPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CapacityPolicy `json:"items"`
}

// CapacityPolicy type metadata.
var (
	CapacityPolicyKind             = reflect.TypeOf(CapacityPolicy{}).Name()
	CapacityPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: CapacityPolicyKind}.String()
	CapacityPolicyKindAPIVersion   = CapacityPolicyKind + "." + SchemeGroupVersion.String()
	CapacityPolicyGroupVersionKind = SchemeGroupVersion.WithKind(CapacityPolicyKind)
)

func init() {
	SchemeBuilder.Register(&CapacityPolicy{}, &CapacityPolicyList{})
}
