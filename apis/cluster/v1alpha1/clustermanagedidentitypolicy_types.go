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

// ClusterManagedIdentityEntry allows one identity for usages.
type ClusterManagedIdentityEntry struct {
	// ObjectID of the managed identity, or "system".
	// +kubebuilder:validation:MinLength=1
	ObjectID string `json:"objectId"`
	// AllowedUsages of the identity (e.g. NativeIngestion, ExternalTable, DataConnection, AutomatedFlows, SandboxArtifacts, AzureAI, All).
	// +kubebuilder:validation:MinItems=1
	AllowedUsages []string `json:"allowedUsages"`
}

// ClusterManagedIdentityPolicyParameters are the configurable fields of a ClusterManagedIdentityPolicy.
type ClusterManagedIdentityPolicyParameters struct {
	// Identities is the complete list of allowed identities.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Identities []ClusterManagedIdentityEntry `json:"identities,omitempty"`
}

// A ClusterManagedIdentityPolicySpec defines the desired state of a ClusterManagedIdentityPolicy.
type ClusterManagedIdentityPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              ClusterManagedIdentityPolicyParameters `json:"forProvider"`
}

// A ClusterManagedIdentityPolicyStatus represents the observed state of a ClusterManagedIdentityPolicy.
type ClusterManagedIdentityPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ClusterPolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ClusterManagedIdentityPolicy allows managed identities to be used for specific usages at the cluster level (all databases). The list is authoritative for the cluster.
// There is exactly one managed_identity policy per cluster, so one managed resource
// per cluster (ProviderConfig) is expected; two managed resources on the same
// cluster policy fight over it and are a user error. The provider principal
// needs the AllDatabasesAdmin cluster role for cluster policies.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxcluster}
type ClusterManagedIdentityPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterManagedIdentityPolicySpec   `json:"spec"`
	Status ClusterManagedIdentityPolicyStatus `json:"status,omitempty"`
}

// GetClusterPolicyObservation returns the observed policy.
func (p *ClusterManagedIdentityPolicy) GetClusterPolicyObservation() ClusterPolicyObservation {
	return p.Status.AtProvider
}

// SetClusterPolicyObservation stores the observed policy in the status.
func (p *ClusterManagedIdentityPolicy) SetClusterPolicyObservation(o ClusterPolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// ClusterManagedIdentityPolicyList contains a list of ClusterManagedIdentityPolicy
type ClusterManagedIdentityPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterManagedIdentityPolicy `json:"items"`
}

// ClusterManagedIdentityPolicy type metadata.
var (
	ClusterManagedIdentityPolicyKind             = reflect.TypeOf(ClusterManagedIdentityPolicy{}).Name()
	ClusterManagedIdentityPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: ClusterManagedIdentityPolicyKind}.String()
	ClusterManagedIdentityPolicyKindAPIVersion   = ClusterManagedIdentityPolicyKind + "." + SchemeGroupVersion.String()
	ClusterManagedIdentityPolicyGroupVersionKind = SchemeGroupVersion.WithKind(ClusterManagedIdentityPolicyKind)
)

func init() {
	SchemeBuilder.Register(&ClusterManagedIdentityPolicy{}, &ClusterManagedIdentityPolicyList{})
}
