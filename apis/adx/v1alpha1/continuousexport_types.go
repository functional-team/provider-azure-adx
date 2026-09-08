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

// ContinuousExportParameters are the configurable fields of a ContinuousExport.
// +kubebuilder:validation:XValidation:rule="has(self.externalTable) || has(self.externalTableRef) || has(self.externalTableSelector)",message="externalTable, externalTableRef or externalTableSelector is required"
type ContinuousExportParameters struct {
	// Name of the continuous export in Kusto. Written to the
	// crossplane.io/external-name annotation on the first reconcile if that
	// annotation is empty. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable"
	// +optional
	Name *string `json:"name,omitempty"`

	// Database that holds the continuous export. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// ExternalTable is the target external table.
	// +crossplane:generate:reference:type=ExternalTable
	// +optional
	ExternalTable string `json:"externalTable,omitempty"`
	// ExternalTableRef references an ExternalTable managed resource.
	// +optional
	ExternalTableRef *xpv2.NamespacedReference `json:"externalTableRef,omitempty"`
	// ExternalTableSelector selects an ExternalTable managed resource.
	// +optional
	ExternalTableSelector *xpv2.NamespacedSelector `json:"externalTableSelector,omitempty"`

	// OverTables are the fact tables the export cursor is scoped to. When
	// omitted Kusto derives them from the query.
	// +crossplane:generate:reference:type=Table
	// +optional
	OverTables []string `json:"overTables,omitempty"`
	// OverTablesRefs reference Table managed resources.
	// +optional
	OverTablesRefs []xpv2.NamespacedReference `json:"overTablesRefs,omitempty"`
	// OverTablesSelector selects Table managed resources.
	// +optional
	OverTablesSelector *xpv2.NamespacedSelector `json:"overTablesSelector,omitempty"`

	// Query is the KQL query whose results are exported.
	// +kubebuilder:validation:MinLength=1
	Query string `json:"query"`

	// IntervalBetweenRuns between export runs, e.g. "1h".
	IntervalBetweenRuns common.Timespan `json:"intervalBetweenRuns"`

	// ForcedLatency delays the export to include late arriving data.
	// +optional
	ForcedLatency *common.Timespan `json:"forcedLatency,omitempty"`

	// SizeLimit of a single exported file in bytes.
	// +optional
	SizeLimit *int64 `json:"sizeLimit,omitempty"`

	// Distributed exports from multiple nodes concurrently.
	// +optional
	Distributed *bool `json:"distributed,omitempty"`

	// Distribution hint: single, per_node or per_shard.
	// +kubebuilder:validation:Enum=single;per_node;per_shard
	// +optional
	Distribution *string `json:"distribution,omitempty"`

	// DistributionKind hint, e.g. default or uniform.
	// +optional
	DistributionKind *string `json:"distributionKind,omitempty"`

	// ParquetRowGroupSize for parquet exports.
	// +optional
	ParquetRowGroupSize *int64 `json:"parquetRowGroupSize,omitempty"`

	// ManagedIdentity ("system" or an object id) used to write to storage.
	// +optional
	ManagedIdentity *string `json:"managedIdentity,omitempty"`

	// Enabled toggles the export. Defaults to true.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// ContinuousExportObservation are the observable fields of a ContinuousExport.
type ContinuousExportObservation struct {
	// ExternalTable as reported by the cluster.
	// +optional
	ExternalTable string `json:"externalTable,omitempty"`
	// Query as reported by the cluster.
	// +optional
	Query string `json:"query,omitempty"`
	// IsDisabled as reported by the cluster.
	// +optional
	IsDisabled bool `json:"isDisabled,omitempty"`
	// LastRunTime of the export.
	// +optional
	LastRunTime string `json:"lastRunTime,omitempty"`
	// LastRunResult of the export.
	// +optional
	LastRunResult string `json:"lastRunResult,omitempty"`
	// ExportedTo is the ingestion time up to which data was exported.
	// +optional
	ExportedTo string `json:"exportedTo,omitempty"`
	// IsRunning reports whether an export run is in progress.
	// +optional
	IsRunning bool `json:"isRunning,omitempty"`
}

// A ContinuousExportSpec defines the desired state of a ContinuousExport.
type ContinuousExportSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              ContinuousExportParameters `json:"forProvider"`
}

// A ContinuousExportStatus represents the observed state of a ContinuousExport.
type ContinuousExportStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ContinuousExportObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ContinuousExport periodically exports the results of a query to an
// external table.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".spec.forProvider.externalTable"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx}
type ContinuousExport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContinuousExportSpec   `json:"spec"`
	Status ContinuousExportStatus `json:"status,omitempty"`
}

// SpecName returns spec.forProvider.name for the external-name initializer.
func (c *ContinuousExport) SpecName() string {
	if c.Spec.ForProvider.Name == nil {
		return ""
	}
	return *c.Spec.ForProvider.Name
}

// +kubebuilder:object:root=true

// ContinuousExportList contains a list of ContinuousExport
type ContinuousExportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContinuousExport `json:"items"`
}

// ContinuousExport type metadata.
var (
	ContinuousExportKind             = reflect.TypeOf(ContinuousExport{}).Name()
	ContinuousExportGroupKind        = schema.GroupKind{Group: Group, Kind: ContinuousExportKind}.String()
	ContinuousExportKindAPIVersion   = ContinuousExportKind + "." + SchemeGroupVersion.String()
	ContinuousExportGroupVersionKind = SchemeGroupVersion.WithKind(ContinuousExportKind)
)

func init() {
	SchemeBuilder.Register(&ContinuousExport{}, &ContinuousExportList{})
}
