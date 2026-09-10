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

// SandboxRule configures one sandbox kind.
type SandboxRule struct {
	// SandboxKind of the rule.
	// +kubebuilder:validation:Enum=PythonExecution;RExecution
	SandboxKind string `json:"sandboxKind"`
	// IsEnabled allows sandboxes of this kind to run on the cluster's nodes.
	// Defaults to false in Kusto, so without it the sandboxes stay off.
	// +optional
	IsEnabled *bool `json:"isEnabled,omitempty"`
	// InitializeOnStartup creates sandboxes eagerly on node startup.
	// +optional
	InitializeOnStartup *bool `json:"initializeOnStartup,omitempty"`
	// TargetCountPerNode is how many sandboxes of this kind may run per node.
	// Between one and twice the processors per node; Kusto defaults to 16.
	// +optional
	TargetCountPerNode *int64 `json:"targetCountPerNode,omitempty"`
	// MaxCPURatePerSandbox is the maximum CPU rate one sandbox may use, as a
	// percentage of all available cores (1-100). Kusto defaults to 50.
	// +optional
	MaxCPURatePerSandbox *int64 `json:"maxCpuRatePerSandbox,omitempty"`
	// MaxMemoryMbPerSandbox in MB.
	// +optional
	MaxMemoryMbPerSandbox *int64 `json:"maxMemoryMbPerSandbox,omitempty"`
}

// SandboxPolicyParameters are the configurable fields of a SandboxPolicy.
type SandboxPolicyParameters struct {
	// Sandboxes is the complete list of sandbox settings.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Sandboxes []SandboxRule `json:"sandboxes,omitempty"`
}

// A SandboxPolicySpec defines the desired state of a SandboxPolicy.
type SandboxPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              SandboxPolicyParameters `json:"forProvider"`
}

// A SandboxPolicyStatus represents the observed state of a SandboxPolicy.
type SandboxPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ClusterPolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A SandboxPolicy configures the sandboxes used by the python() and r() plugins. The list is authoritative for the cluster.
// There is exactly one sandbox policy per cluster, so one managed resource
// per cluster (ProviderConfig) is expected; two managed resources on the same
// cluster policy fight over it and are a user error. The provider principal
// needs the AllDatabasesAdmin cluster role for cluster policies.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxcluster}
type SandboxPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SandboxPolicySpec   `json:"spec"`
	Status SandboxPolicyStatus `json:"status,omitempty"`
}

// GetClusterPolicyObservation returns the observed policy.
func (p *SandboxPolicy) GetClusterPolicyObservation() ClusterPolicyObservation {
	return p.Status.AtProvider
}

// SetClusterPolicyObservation stores the observed policy in the status.
func (p *SandboxPolicy) SetClusterPolicyObservation(o ClusterPolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// SandboxPolicyList contains a list of SandboxPolicy
type SandboxPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SandboxPolicy `json:"items"`
}

// SandboxPolicy type metadata.
var (
	SandboxPolicyKind             = reflect.TypeOf(SandboxPolicy{}).Name()
	SandboxPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: SandboxPolicyKind}.String()
	SandboxPolicyKindAPIVersion   = SandboxPolicyKind + "." + SchemeGroupVersion.String()
	SandboxPolicyGroupVersionKind = SchemeGroupVersion.WithKind(SandboxPolicyKind)
)

func init() {
	SchemeBuilder.Register(&SandboxPolicy{}, &SandboxPolicyList{})
}
