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

// SchemaUpdateMode controls how schema changes are applied.
type SchemaUpdateMode string

const (
	// SchemaUpdateModeMerge only adds columns (.alter-merge table). Columns that
	// exist in the cluster but not in the spec are kept and reported in
	// status.atProvider.driftColumns.
	SchemaUpdateModeMerge SchemaUpdateMode = "Merge"
	// SchemaUpdateModeReplace makes the spec authoritative (.alter table):
	// columns missing from the spec are dropped together with their data.
	SchemaUpdateModeReplace SchemaUpdateMode = "Replace"
)

// TableParameters are the configurable fields of a Table.
type TableParameters struct {
	// Name of the table in Kusto. Written to the crossplane.io/external-name
	// annotation on the first reconcile if that annotation is empty; the
	// annotation remains the source of truth. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable"
	// +optional
	Name *string `json:"name,omitempty"`

	// Database that holds the table. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// Columns of the table in order.
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=name
	Columns []common.Column `json:"columns"`

	// Folder of the table in the database tree.
	// +optional
	Folder *string `json:"folder,omitempty"`

	// Docstring of the table.
	// +optional
	Docstring *string `json:"docstring,omitempty"`

	// SchemaUpdateMode selects Merge (default, add-only, safe) or Replace
	// (authoritative, drops columns and their data).
	// +kubebuilder:validation:Enum=Merge;Replace
	// +kubebuilder:default=Merge
	// +optional
	SchemaUpdateMode SchemaUpdateMode `json:"schemaUpdateMode,omitempty"`
}

// TableObservation are the observable fields of a Table.
type TableObservation struct {
	// Columns as reported by the cluster.
	// +optional
	Columns []common.ObservedColumn `json:"columns,omitempty"`
	// Folder as reported by the cluster.
	// +optional
	Folder string `json:"folder,omitempty"`
	// Docstring as reported by the cluster.
	// +optional
	Docstring string `json:"docstring,omitempty"`
	// DriftColumns exist in the cluster but not in the spec. They are kept in
	// Merge mode and would be dropped in Replace mode.
	// +optional
	DriftColumns []string `json:"driftColumns,omitempty"`
}

// A TableSpec defines the desired state of a Table.
type TableSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              TableParameters `json:"forProvider"`
}

// A TableStatus represents the observed state of a Table.
type TableStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 TableObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Table is a Kusto table with a managed schema. Deleting the managed
// resource drops the table and all of its data (Crossplane default deletion
// semantics); use deletionPolicy Orphan or managementPolicies to keep it.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx}
type Table struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TableSpec   `json:"spec"`
	Status TableStatus `json:"status,omitempty"`
}

// SpecName returns spec.forProvider.name for the external-name initializer.
func (t *Table) SpecName() string {
	if t.Spec.ForProvider.Name == nil {
		return ""
	}
	return *t.Spec.ForProvider.Name
}

// +kubebuilder:object:root=true

// TableList contains a list of Table
type TableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Table `json:"items"`
}

// Table type metadata.
var (
	TableKind             = reflect.TypeOf(Table{}).Name()
	TableGroupKind        = schema.GroupKind{Group: Group, Kind: TableKind}.String()
	TableKindAPIVersion   = TableKind + "." + SchemeGroupVersion.String()
	TableGroupVersionKind = SchemeGroupVersion.WithKind(TableKind)
)

func init() {
	SchemeBuilder.Register(&Table{}, &TableList{})
}
