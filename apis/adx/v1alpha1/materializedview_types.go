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

// MaterializedViewParameters are the configurable fields of a MaterializedView.
// +kubebuilder:validation:XValidation:rule="has(self.sourceTable) || has(self.sourceTableRef) || has(self.sourceTableSelector) || has(self.sourceMaterializedView)",message="one of sourceTable, sourceTableRef, sourceTableSelector or sourceMaterializedView is required"
// +kubebuilder:validation:XValidation:rule="!(has(self.sourceMaterializedView) && (has(self.sourceTable) || has(self.sourceTableRef) || has(self.sourceTableSelector)))",message="sourceMaterializedView and sourceTable are mutually exclusive"
type MaterializedViewParameters struct {
	// Name of the materialized view in Kusto. Written to the
	// crossplane.io/external-name annotation on the first reconcile if that
	// annotation is empty. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable"
	// +optional
	Name *string `json:"name,omitempty"`

	// Database that holds the view. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// SourceTable the view is defined over. Immutable once set; changing the
	// source means a new view.
	// +crossplane:generate:reference:type=Table
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="sourceTable is immutable"
	// +optional
	SourceTable *string `json:"sourceTable,omitempty"`

	// SourceTableRef references a Table managed resource as source.
	// +optional
	SourceTableRef *xpv2.NamespacedReference `json:"sourceTableRef,omitempty"`

	// SourceTableSelector selects a Table managed resource as source.
	// +optional
	SourceTableSelector *xpv2.NamespacedSelector `json:"sourceTableSelector,omitempty"`

	// SourceMaterializedView is another materialized view used as source
	// (materialized view over materialized view). Immutable.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="sourceMaterializedView is immutable"
	// +optional
	SourceMaterializedView *string `json:"sourceMaterializedView,omitempty"`

	// Query is the KQL aggregation query over the source. It is sent verbatim.
	// Kusto only accepts limited changes to an existing view's query (no
	// change of the group-by expressions, column names or types); rejected
	// changes surface as Synced=False with the Kusto message.
	// +kubebuilder:validation:MinLength=1
	Query string `json:"query"`

	// Backfill materializes existing source data. Create-only: the view is
	// created asynchronously and the backfill operation is polled until it
	// completes.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="backfill is create-only"
	// +optional
	Backfill *bool `json:"backfill,omitempty"`

	// EffectiveDateTime (ISO 8601) limits backfill to records ingested after
	// this time. Create-only.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="effectiveDateTime is create-only"
	// +optional
	EffectiveDateTime *string `json:"effectiveDateTime,omitempty"`

	// UpdateExtentsCreationTime sets the creation time of backfilled extents
	// to the ingestion time of the source records. Create-only.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="updateExtentsCreationTime is create-only"
	// +optional
	UpdateExtentsCreationTime *bool `json:"updateExtentsCreationTime,omitempty"`

	// Lookback limits the period of source data considered during
	// materialization, e.g. "6h".
	// +optional
	Lookback *common.Timespan `json:"lookback,omitempty"`

	// LookbackColumn is the datetime column the lookback applies to. Kusto
	// does not allow changing it after it has been set.
	// +optional
	LookbackColumn *string `json:"lookbackColumn,omitempty"`

	// AutoUpdateSchema propagates source schema changes to the view.
	// +optional
	AutoUpdateSchema *bool `json:"autoUpdateSchema,omitempty"`

	// DimensionTables are tables joined in the query that must not be treated
	// as fact tables. Write-only: Kusto does not report them.
	// +optional
	DimensionTables []string `json:"dimensionTables,omitempty"`

	// Folder of the view in the database tree.
	// +optional
	Folder *string `json:"folder,omitempty"`

	// Docstring of the view.
	// +optional
	Docstring *string `json:"docstring,omitempty"`

	// AllowMaterializedViewsWithoutRowLevelSecurity allows creating the view
	// over a source table with a row level security policy. Write-only.
	// +optional
	AllowMaterializedViewsWithoutRowLevelSecurity *bool `json:"allowMaterializedViewsWithoutRowLevelSecurity,omitempty"`

	// MaxSourceRecordsForSingleIngest caps the records per backfill ingest
	// operation. Create-only.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="maxSourceRecordsForSingleIngest is create-only"
	// +optional
	MaxSourceRecordsForSingleIngest *int64 `json:"maxSourceRecordsForSingleIngest,omitempty"`

	// Concurrency of the backfill. Create-only.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="concurrency is create-only"
	// +optional
	Concurrency *int64 `json:"concurrency,omitempty"`

	// Enabled toggles materialization (.enable/.disable materialized-view).
	// Defaults to true.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// MaterializedViewObservation are the observable fields of a MaterializedView.
type MaterializedViewObservation struct {
	// SourceTable as reported by the cluster.
	// +optional
	SourceTable string `json:"sourceTable,omitempty"`
	// Query as reported by the cluster (Kusto may reformat it).
	// +optional
	Query string `json:"query,omitempty"`
	// IsEnabled as reported by the cluster.
	// +optional
	IsEnabled *bool `json:"isEnabled,omitempty"`
	// IsHealthy as reported by the cluster.
	// +optional
	IsHealthy *bool `json:"isHealthy,omitempty"`
	// Folder as reported by the cluster.
	// +optional
	Folder string `json:"folder,omitempty"`
	// Docstring as reported by the cluster.
	// +optional
	Docstring string `json:"docstring,omitempty"`
	// Lookback as reported by the cluster (.NET timespan format).
	// +optional
	Lookback string `json:"lookback,omitempty"`
	// EffectiveDateTime as reported by the cluster.
	// +optional
	EffectiveDateTime string `json:"effectiveDateTime,omitempty"`
	// LastRunResult of the materialization as reported by the cluster.
	// +optional
	LastRunResult string `json:"lastRunResult,omitempty"`
	// OperationID of a running backfill operation.
	// +optional
	OperationID string `json:"operationId,omitempty"`
	// OperationState of the backfill operation (InProgress, Scheduled, ...).
	// +optional
	OperationState string `json:"operationState,omitempty"`
}

// A MaterializedViewSpec defines the desired state of a MaterializedView.
type MaterializedViewSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              MaterializedViewParameters `json:"forProvider"`
}

// A MaterializedViewStatus represents the observed state of a MaterializedView.
type MaterializedViewStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 MaterializedViewObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A MaterializedView is a Kusto materialized view. With backfill the view is
// created asynchronously; the resource stays in Creating until the backfill
// operation completes. Deleting the managed resource drops the view.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="SOURCE",type="string",JSONPath=".status.atProvider.sourceTable"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx}
type MaterializedView struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaterializedViewSpec   `json:"spec"`
	Status MaterializedViewStatus `json:"status,omitempty"`
}

// SpecName returns spec.forProvider.name for the external-name initializer.
func (m *MaterializedView) SpecName() string {
	if m.Spec.ForProvider.Name == nil {
		return ""
	}
	return *m.Spec.ForProvider.Name
}

// +kubebuilder:object:root=true

// MaterializedViewList contains a list of MaterializedView
type MaterializedViewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaterializedView `json:"items"`
}

// MaterializedView type metadata.
var (
	MaterializedViewKind             = reflect.TypeOf(MaterializedView{}).Name()
	MaterializedViewGroupKind        = schema.GroupKind{Group: Group, Kind: MaterializedViewKind}.String()
	MaterializedViewKindAPIVersion   = MaterializedViewKind + "." + SchemeGroupVersion.String()
	MaterializedViewGroupVersionKind = SchemeGroupVersion.WithKind(MaterializedViewKind)
)

func init() {
	SchemeBuilder.Register(&MaterializedView{}, &MaterializedViewList{})
}
