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

// WorkloadGroupParameters are the configurable fields of a WorkloadGroup.
type WorkloadGroupParameters struct {
	// Name of the workload group in Kusto. Written to the
	// crossplane.io/external-name annotation on the first reconcile if that
	// annotation is empty. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable"
	// +optional
	Name *string `json:"name,omitempty"`

	// WorkloadGroup is the workload group JSON as documented by Kusto:
	// RequestLimitsPolicy, RequestRateLimitPolicies,
	// RequestRateLimitsEnforcementPolicy and RequestQueuingPolicy. Only the
	// keys present here are compared with the cluster; unknown keys are
	// passed through unchanged.
	// +kubebuilder:pruning:PreserveUnknownFields
	WorkloadGroup apiextensionsv1.JSON `json:"workloadGroup"`
}

// WorkloadGroupObservation are the observable fields of a WorkloadGroup.
type WorkloadGroupObservation struct {
	// WorkloadGroup is the workload group JSON as reported by the cluster.
	// +optional
	WorkloadGroup string `json:"workloadGroup,omitempty"`
}

// A WorkloadGroupSpec defines the desired state of a WorkloadGroup.
type WorkloadGroupSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              WorkloadGroupParameters `json:"forProvider"`
}

// A WorkloadGroupStatus represents the observed state of a WorkloadGroup.
type WorkloadGroupStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 WorkloadGroupObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A WorkloadGroup is a Kusto workload group (request limits, rate limits,
// queuing). The built-in groups "default" and "internal" can be altered but
// not dropped: deleting their managed resource leaves them in place. The
// provider principal needs the AllDatabasesAdmin cluster role.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxcluster}
type WorkloadGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkloadGroupSpec   `json:"spec"`
	Status WorkloadGroupStatus `json:"status,omitempty"`
}

// SpecName returns spec.forProvider.name for the external-name initializer.
func (w *WorkloadGroup) SpecName() string {
	if w.Spec.ForProvider.Name == nil {
		return ""
	}
	return *w.Spec.ForProvider.Name
}

// +kubebuilder:object:root=true

// WorkloadGroupList contains a list of WorkloadGroup
type WorkloadGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkloadGroup `json:"items"`
}

// WorkloadGroup type metadata.
var (
	WorkloadGroupKind             = reflect.TypeOf(WorkloadGroup{}).Name()
	WorkloadGroupGroupKind        = schema.GroupKind{Group: Group, Kind: WorkloadGroupKind}.String()
	WorkloadGroupKindAPIVersion   = WorkloadGroupKind + "." + SchemeGroupVersion.String()
	WorkloadGroupGroupVersionKind = SchemeGroupVersion.WithKind(WorkloadGroupKind)
)

func init() {
	SchemeBuilder.Register(&WorkloadGroup{}, &WorkloadGroupList{})
}
