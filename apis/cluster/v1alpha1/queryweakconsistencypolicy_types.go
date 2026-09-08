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

// QueryWeakConsistencyPolicyParameters are the configurable fields of a QueryWeakConsistencyPolicy.
type QueryWeakConsistencyPolicyParameters struct {
	// PercentageOfNodes that serve weakly consistent queries (-1 for default).
	// +optional
	PercentageOfNodes *int64 `json:"percentageOfNodes,omitempty"`

	// MinimumNumberOfNodes for weak consistency (-1 for default).
	// +optional
	MinimumNumberOfNodes *int64 `json:"minimumNumberOfNodes,omitempty"`

	// MaximumNumberOfNodes for weak consistency (-1 for default).
	// +optional
	MaximumNumberOfNodes *int64 `json:"maximumNumberOfNodes,omitempty"`

	// SuperSlackerNumberOfNodesThreshold (-1 for default).
	// +optional
	SuperSlackerNumberOfNodesThreshold *int64 `json:"superSlackerNumberOfNodesThreshold,omitempty"`

	// EnableMetadataPrefetch prefetches metadata on weak consistency nodes.
	// +optional
	EnableMetadataPrefetch *bool `json:"enableMetadataPrefetch,omitempty"`

	// MaximumLagAllowedInMinutes of metadata (-1 for default).
	// +optional
	MaximumLagAllowedInMinutes *int64 `json:"maximumLagAllowedInMinutes,omitempty"`

	// RefreshPeriodInSeconds of metadata (-1 for default).
	// +optional
	RefreshPeriodInSeconds *int64 `json:"refreshPeriodInSeconds,omitempty"`
}

// A QueryWeakConsistencyPolicySpec defines the desired state of a QueryWeakConsistencyPolicy.
type QueryWeakConsistencyPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              QueryWeakConsistencyPolicyParameters `json:"forProvider"`
}

// A QueryWeakConsistencyPolicyStatus represents the observed state of a QueryWeakConsistencyPolicy.
type QueryWeakConsistencyPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ClusterPolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A QueryWeakConsistencyPolicy configures weak consistency query nodes. The policy always exists on a cluster; deleting the managed resource leaves the last applied values in place (Kusto has no delete for this policy).
// There is exactly one query_weak_consistency policy per cluster, so one managed resource
// per cluster (ProviderConfig) is expected; two managed resources on the same
// cluster policy fight over it and are a user error. The provider principal
// needs the AllDatabasesAdmin cluster role for cluster policies.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxcluster}
type QueryWeakConsistencyPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueryWeakConsistencyPolicySpec   `json:"spec"`
	Status QueryWeakConsistencyPolicyStatus `json:"status,omitempty"`
}

// GetClusterPolicyObservation returns the observed policy.
func (p *QueryWeakConsistencyPolicy) GetClusterPolicyObservation() ClusterPolicyObservation {
	return p.Status.AtProvider
}

// SetClusterPolicyObservation stores the observed policy in the status.
func (p *QueryWeakConsistencyPolicy) SetClusterPolicyObservation(o ClusterPolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// QueryWeakConsistencyPolicyList contains a list of QueryWeakConsistencyPolicy
type QueryWeakConsistencyPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueryWeakConsistencyPolicy `json:"items"`
}

// QueryWeakConsistencyPolicy type metadata.
var (
	QueryWeakConsistencyPolicyKind             = reflect.TypeOf(QueryWeakConsistencyPolicy{}).Name()
	QueryWeakConsistencyPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: QueryWeakConsistencyPolicyKind}.String()
	QueryWeakConsistencyPolicyKindAPIVersion   = QueryWeakConsistencyPolicyKind + "." + SchemeGroupVersion.String()
	QueryWeakConsistencyPolicyGroupVersionKind = SchemeGroupVersion.WithKind(QueryWeakConsistencyPolicyKind)
)

func init() {
	SchemeBuilder.Register(&QueryWeakConsistencyPolicy{}, &QueryWeakConsistencyPolicyList{})
}
