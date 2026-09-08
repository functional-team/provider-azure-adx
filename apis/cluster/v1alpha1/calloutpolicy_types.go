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

// CalloutRule allows or denies callouts of one type to URIs matching a regex.
type CalloutRule struct {
	// CalloutType of the rule.
	// +kubebuilder:validation:Enum=kusto;sql;cosmosdb;external_data;azure_digital_twins;sandbox_artifacts;webapi;mysql;postgresql;genevametrics;azure_openai
	CalloutType string `json:"calloutType"`
	// CalloutURIRegex matches the target URIs.
	// +kubebuilder:validation:MinLength=1
	CalloutURIRegex string `json:"calloutUriRegex"`
	// CanCall allows (true) or denies (false) matching callouts.
	CanCall bool `json:"canCall"`
}

// CalloutPolicyParameters are the configurable fields of a CalloutPolicy.
type CalloutPolicyParameters struct {
	// Callouts is the complete list of callout rules.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Callouts []CalloutRule `json:"callouts,omitempty"`
}

// A CalloutPolicySpec defines the desired state of a CalloutPolicy.
type CalloutPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              CalloutPolicyParameters `json:"forProvider"`
}

// A CalloutPolicyStatus represents the observed state of a CalloutPolicy.
type CalloutPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ClusterPolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A CalloutPolicy controls which external endpoints (SQL, Cosmos DB, external data, sandboxes, ...) the cluster may call. The list is authoritative for the cluster.
// There is exactly one callout policy per cluster, so one managed resource
// per cluster (ProviderConfig) is expected; two managed resources on the same
// cluster policy fight over it and are a user error. The provider principal
// needs the AllDatabasesAdmin cluster role for cluster policies.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxcluster}
type CalloutPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CalloutPolicySpec   `json:"spec"`
	Status CalloutPolicyStatus `json:"status,omitempty"`
}

// GetClusterPolicyObservation returns the observed policy.
func (p *CalloutPolicy) GetClusterPolicyObservation() ClusterPolicyObservation {
	return p.Status.AtProvider
}

// SetClusterPolicyObservation stores the observed policy in the status.
func (p *CalloutPolicy) SetClusterPolicyObservation(o ClusterPolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// CalloutPolicyList contains a list of CalloutPolicy
type CalloutPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CalloutPolicy `json:"items"`
}

// CalloutPolicy type metadata.
var (
	CalloutPolicyKind             = reflect.TypeOf(CalloutPolicy{}).Name()
	CalloutPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: CalloutPolicyKind}.String()
	CalloutPolicyKindAPIVersion   = CalloutPolicyKind + "." + SchemeGroupVersion.String()
	CalloutPolicyGroupVersionKind = SchemeGroupVersion.WithKind(CalloutPolicyKind)
)

func init() {
	SchemeBuilder.Register(&CalloutPolicy{}, &CalloutPolicyList{})
}
