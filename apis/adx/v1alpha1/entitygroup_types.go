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

// EntityGroupParameters are the configurable fields of an EntityGroup.
type EntityGroupParameters struct {
	// Name of the entity group in Kusto. Written to the
	// crossplane.io/external-name annotation on the first reconcile if that
	// annotation is empty. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable"
	// +optional
	Name *string `json:"name,omitempty"`

	// Database that holds the entity group. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// Entities are raw KQL entity expressions such as
	// cluster('c').database('d') or database('d').table('t'). They are code
	// and are sent verbatim; order does not matter.
	// +kubebuilder:validation:MinItems=1
	Entities []string `json:"entities"`
}

// EntityGroupObservation are the observable fields of an EntityGroup.
type EntityGroupObservation struct {
	// Entities as reported by the cluster.
	// +optional
	Entities []string `json:"entities,omitempty"`
}

// An EntityGroupSpec defines the desired state of an EntityGroup.
type EntityGroupSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              EntityGroupParameters `json:"forProvider"`
}

// An EntityGroupStatus represents the observed state of an EntityGroup.
type EntityGroupStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 EntityGroupObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// An EntityGroup is a named set of Kusto entities usable with the macro-expand
// operator (Tier 3).
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx}
type EntityGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntityGroupSpec   `json:"spec"`
	Status EntityGroupStatus `json:"status,omitempty"`
}

// SpecName returns spec.forProvider.name for the external-name initializer.
func (e *EntityGroup) SpecName() string {
	if e.Spec.ForProvider.Name == nil {
		return ""
	}
	return *e.Spec.ForProvider.Name
}

// +kubebuilder:object:root=true

// EntityGroupList contains a list of EntityGroup
type EntityGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntityGroup `json:"items"`
}

// EntityGroup type metadata.
var (
	EntityGroupKind             = reflect.TypeOf(EntityGroup{}).Name()
	EntityGroupGroupKind        = schema.GroupKind{Group: Group, Kind: EntityGroupKind}.String()
	EntityGroupKindAPIVersion   = EntityGroupKind + "." + SchemeGroupVersion.String()
	EntityGroupGroupVersionKind = SchemeGroupVersion.WithKind(EntityGroupKind)
)

func init() {
	SchemeBuilder.Register(&EntityGroup{}, &EntityGroupList{})
}
