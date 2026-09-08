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

// FunctionParameters are the configurable fields of a Function.
type FunctionParameters struct {
	// Name of the function in Kusto. Written to the crossplane.io/external-name
	// annotation on the first reconcile if that annotation is empty. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable"
	// +optional
	Name *string `json:"name,omitempty"`

	// Database that holds the function. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// Parameters is the raw parameter list including parentheses, e.g.
	// "(limit:long = 100, T:(x:long))". Whitespace and type aliases are
	// normalized before comparison.
	// +kubebuilder:default="()"
	// +optional
	Parameters string `json:"parameters,omitempty"`

	// Body of the function without the outer curly braces. It is KQL code and
	// is sent verbatim; the provider adds the braces.
	// +kubebuilder:validation:MinLength=1
	Body string `json:"body"`

	// Folder of the function in the database tree.
	// +optional
	Folder *string `json:"folder,omitempty"`

	// Docstring of the function.
	// +optional
	Docstring *string `json:"docstring,omitempty"`

	// View marks the function as a view (usable in wildcard unions). The flag
	// is write-only: Kusto does not report it in .show functions.
	// +optional
	View *bool `json:"view,omitempty"`

	// SkipValidation skips semantic validation of the body on write. Write-only.
	// +optional
	SkipValidation *bool `json:"skipValidation,omitempty"`
}

// FunctionObservation are the observable fields of a Function.
type FunctionObservation struct {
	// Parameters as reported by the cluster.
	// +optional
	Parameters string `json:"parameters,omitempty"`
	// Body as reported by the cluster (Kusto may reformat it).
	// +optional
	Body string `json:"body,omitempty"`
	// Folder as reported by the cluster.
	// +optional
	Folder string `json:"folder,omitempty"`
	// Docstring as reported by the cluster.
	// +optional
	Docstring string `json:"docstring,omitempty"`
}

// A FunctionSpec defines the desired state of a Function.
type FunctionSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              FunctionParameters `json:"forProvider"`
}

// A FunctionStatus represents the observed state of a Function.
type FunctionStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 FunctionObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Function is a stored Kusto function (or view). Whoever may write a
// Function runs arbitrary KQL with the provider principal's permissions.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx}
type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FunctionSpec   `json:"spec"`
	Status FunctionStatus `json:"status,omitempty"`
}

// SpecName returns spec.forProvider.name for the external-name initializer.
func (f *Function) SpecName() string {
	if f.Spec.ForProvider.Name == nil {
		return ""
	}
	return *f.Spec.ForProvider.Name
}

// +kubebuilder:object:root=true

// FunctionList contains a list of Function
type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Function `json:"items"`
}

// Function type metadata.
var (
	FunctionKind             = reflect.TypeOf(Function{}).Name()
	FunctionGroupKind        = schema.GroupKind{Group: Group, Kind: FunctionKind}.String()
	FunctionKindAPIVersion   = FunctionKind + "." + SchemeGroupVersion.String()
	FunctionGroupVersionKind = SchemeGroupVersion.WithKind(FunctionKind)
)

func init() {
	SchemeBuilder.Register(&Function{}, &FunctionList{})
}
