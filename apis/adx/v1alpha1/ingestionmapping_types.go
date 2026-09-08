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

// IngestionMappingKind is the data format an ingestion mapping applies to.
// +kubebuilder:validation:Enum=csv;json;avro;apacheavro;parquet;orc;w3clogfile
type IngestionMappingKind string

// IngestionMappingColumn maps one source field to a table column.
type IngestionMappingColumn struct {
	// Column is the target column of the table.
	// +kubebuilder:validation:MinLength=1
	Column string `json:"column"`

	// DataType of the target column. Optional; Kusto derives it from the
	// table when omitted.
	// +optional
	DataType *common.ColumnType `json:"dataType,omitempty"`

	// Properties are the format specific mapping properties, e.g. path
	// ("$.ts"), transform, ordinal, constValue or field. Keys are case
	// insensitive and written to Kusto with an upper-case first letter.
	// +optional
	Properties map[string]string `json:"properties,omitempty"`
}

// IngestionMappingParameters are the configurable fields of an IngestionMapping.
// +kubebuilder:validation:XValidation:rule="has(self.table) || has(self.tableRef) || has(self.tableSelector)",message="table, tableRef or tableSelector is required"
type IngestionMappingParameters struct {
	// Name of the mapping in Kusto. Written to the crossplane.io/external-name
	// annotation on the first reconcile if that annotation is empty. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable"
	// +optional
	Name *string `json:"name,omitempty"`

	// Database that holds the table. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// Table the mapping belongs to. Immutable once set.
	// +crossplane:generate:reference:type=Table
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="table is immutable"
	// +optional
	Table string `json:"table,omitempty"`

	// TableRef references a Table managed resource.
	// +optional
	TableRef *xpv2.NamespacedReference `json:"tableRef,omitempty"`

	// TableSelector selects a Table managed resource.
	// +optional
	TableSelector *xpv2.NamespacedSelector `json:"tableSelector,omitempty"`

	// Kind is the data format of the mapping. Immutable.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="kind is immutable"
	Kind IngestionMappingKind `json:"kind"`

	// Mapping is the ordered list of column mappings.
	// +kubebuilder:validation:MinItems=1
	Mapping []IngestionMappingColumn `json:"mapping"`
}

// IngestionMappingObservation are the observable fields of an IngestionMapping.
type IngestionMappingObservation struct {
	// Kind as reported by the cluster.
	// +optional
	Kind string `json:"kind,omitempty"`
	// Table as reported by the cluster.
	// +optional
	Table string `json:"table,omitempty"`
	// LastUpdatedOn as reported by the cluster.
	// +optional
	LastUpdatedOn string `json:"lastUpdatedOn,omitempty"`
	// Mapping is the mapping JSON as reported by the cluster.
	// +optional
	Mapping string `json:"mapping,omitempty"`
}

// An IngestionMappingSpec defines the desired state of an IngestionMapping.
type IngestionMappingSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              IngestionMappingParameters `json:"forProvider"`
}

// An IngestionMappingStatus represents the observed state of an IngestionMapping.
type IngestionMappingStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 IngestionMappingObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// An IngestionMapping is a pre-created ingestion mapping of a table. Its
// identity is (database, table, kind, name).
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="TABLE",type="string",JSONPath=".spec.forProvider.table"
// +kubebuilder:printcolumn:name="KIND",type="string",JSONPath=".spec.forProvider.kind"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx}
type IngestionMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngestionMappingSpec   `json:"spec"`
	Status IngestionMappingStatus `json:"status,omitempty"`
}

// SpecName returns spec.forProvider.name for the external-name initializer.
func (m *IngestionMapping) SpecName() string {
	if m.Spec.ForProvider.Name == nil {
		return ""
	}
	return *m.Spec.ForProvider.Name
}

// +kubebuilder:object:root=true

// IngestionMappingList contains a list of IngestionMapping
type IngestionMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngestionMapping `json:"items"`
}

// IngestionMapping type metadata.
var (
	IngestionMappingKindName         = reflect.TypeOf(IngestionMapping{}).Name()
	IngestionMappingGroupKind        = schema.GroupKind{Group: Group, Kind: IngestionMappingKindName}.String()
	IngestionMappingKindAPIVersion   = IngestionMappingKindName + "." + SchemeGroupVersion.String()
	IngestionMappingGroupVersionKind = SchemeGroupVersion.WithKind(IngestionMappingKindName)
)

func init() {
	SchemeBuilder.Register(&IngestionMapping{}, &IngestionMappingList{})
}
