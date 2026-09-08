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

// ExternalTableStorageKind is the storage kind of an external table.
type ExternalTableStorageKind string

const (
	// ExternalTableKindStorage reads files (csv, json, parquet, ...) from blob or ADLS storage.
	ExternalTableKindStorage ExternalTableStorageKind = "Storage"
	// ExternalTableKindDelta reads a Delta Lake table from storage.
	ExternalTableKindDelta ExternalTableStorageKind = "Delta"
)

// SecretKeySelector selects a key of a Secret. The namespace defaults to the
// namespace of the managed resource.
type SecretKeySelector struct {
	// Name of the secret.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Namespace of the secret; defaults to the namespace of the managed resource.
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// Key to select.
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`
}

// ExternalTableConnectionString is one storage connection string. Exactly one
// of value and secretKeyRef must be set. Connection strings are always sent
// obfuscated (h@'...') so Kusto hides everything after the first ';' in .show
// output and command logs; only the URI part is compared and reported.
// +kubebuilder:validation:XValidation:rule="has(self.value) != has(self.secretKeyRef)",message="exactly one of value and secretKeyRef must be set"
type ExternalTableConnectionString struct {
	// Value is the connection string, e.g.
	// "https://acct.blob.core.windows.net/exports;managed_identity=system".
	// Use it for identity based access without secrets.
	// +optional
	Value *string `json:"value,omitempty"`
	// SecretKeyRef reads the connection string (with SAS token or account
	// key) from a Secret. Write-only: a changed secret is detected through a
	// hash, the value never appears in status, events or logs.
	// +optional
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// ExternalTableProperties are the optional "with (...)" properties.
type ExternalTableProperties struct {
	// Folder of the external table in the database tree.
	// +optional
	Folder *string `json:"folder,omitempty"`
	// DocString of the external table.
	// +optional
	DocString *string `json:"docString,omitempty"`
	// Compressed marks the files as compressed.
	// +optional
	Compressed *bool `json:"compressed,omitempty"`
	// CompressionType, e.g. gzip or snappy.
	// +optional
	CompressionType *string `json:"compressionType,omitempty"`
	// IncludeHeaders for delimited formats: None, All or FirstFile.
	// +kubebuilder:validation:Enum=None;All;FirstFile
	// +optional
	IncludeHeaders *string `json:"includeHeaders,omitempty"`
	// NamePrefix of exported files.
	// +optional
	NamePrefix *string `json:"namePrefix,omitempty"`
	// FileExtension of the files, e.g. ".parquet".
	// +optional
	FileExtension *string `json:"fileExtension,omitempty"`
	// Encoding of text files, e.g. UTF8NoBOM.
	// +optional
	Encoding *string `json:"encoding,omitempty"`
}

// ExternalTableParameters are the configurable fields of an ExternalTable.
// +kubebuilder:validation:XValidation:rule="self.kind != 'Storage' || (has(self.columns) && size(self.columns) > 0)",message="columns are required for kind Storage"
// +kubebuilder:validation:XValidation:rule="self.kind != 'Storage' || has(self.dataFormat)",message="dataFormat is required for kind Storage"
type ExternalTableParameters struct {
	// Name of the external table in Kusto. Written to the
	// crossplane.io/external-name annotation on the first reconcile if that
	// annotation is empty. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable"
	// +optional
	Name *string `json:"name,omitempty"`

	// Database that holds the external table. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// Kind of the external table. Immutable.
	// +kubebuilder:validation:Enum=Storage;Delta
	// +kubebuilder:default=Storage
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="kind is immutable"
	// +optional
	Kind ExternalTableStorageKind `json:"kind,omitempty"`

	// Columns of the external table. Required for Storage; optional for Delta
	// (schema is inferred from the delta log when omitted).
	// +listType=map
	// +listMapKey=name
	// +optional
	Columns []common.Column `json:"columns,omitempty"`

	// PartitionBy is the raw Kusto partition expression, e.g.
	// "Date:datetime = bin(Timestamp, 1d)". Compared through the hash
	// mechanism because Kusto reports partitions as JSON.
	// +optional
	PartitionBy *string `json:"partitionBy,omitempty"`

	// PathFormat is the raw Kusto path format expression, e.g.
	// "datetime_pattern(\"yyyy/MM/dd\", Date)".
	// +optional
	PathFormat *string `json:"pathFormat,omitempty"`

	// DataFormat of the files (Storage kind).
	// +kubebuilder:validation:Enum=csv;tsv;tsve;psv;scsv;sohsv;json;multijson;avro;apacheavro;parquet;orc;w3clogfile;txt;raw;sstream
	// +optional
	DataFormat *string `json:"dataFormat,omitempty"`

	// ConnectionStrings of the storage containers, at least one.
	// +kubebuilder:validation:MinItems=1
	ConnectionStrings []ExternalTableConnectionString `json:"connectionStrings"`

	// Properties are the optional "with (...)" settings.
	// +optional
	Properties *ExternalTableProperties `json:"properties,omitempty"`
}

// ExternalTableObservation are the observable fields of an ExternalTable.
// Connection strings are reported as URI only (everything after the first
// ';' is a secret and is masked by Kusto as well).
type ExternalTableObservation struct {
	// Kind as reported by the cluster (Blob, Adl, Delta, ...).
	// +optional
	Kind string `json:"kind,omitempty"`
	// Columns as reported by the cluster (only loaded when the spec has columns).
	// +optional
	Columns []common.ObservedColumn `json:"columns,omitempty"`
	// Folder as reported by the cluster.
	// +optional
	Folder string `json:"folder,omitempty"`
	// DocString as reported by the cluster.
	// +optional
	DocString string `json:"docString,omitempty"`
	// DataFormat as reported by the cluster.
	// +optional
	DataFormat string `json:"dataFormat,omitempty"`
	// ConnectionStringURIs are the URI parts of the connection strings.
	// +optional
	ConnectionStringURIs []string `json:"connectionStringUris,omitempty"`
	// Partitions is the partition definition JSON as reported by the cluster.
	// +optional
	Partitions string `json:"partitions,omitempty"`
	// PathFormat as reported by the cluster.
	// +optional
	PathFormat string `json:"pathFormat,omitempty"`
}

// An ExternalTableSpec defines the desired state of an ExternalTable.
type ExternalTableSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              ExternalTableParameters `json:"forProvider"`
}

// An ExternalTableStatus represents the observed state of an ExternalTable.
type ExternalTableStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ExternalTableObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// An ExternalTable is a Kusto external table over blob/ADLS storage or a
// Delta Lake table. Secrets in connection strings are write-only; prefer
// managed identity (";managed_identity=system"). Creating an external table
// with a managed identity requires the AllDatabasesAdmin role.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="KIND",type="string",JSONPath=".spec.forProvider.kind"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx}
type ExternalTable struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalTableSpec   `json:"spec"`
	Status ExternalTableStatus `json:"status,omitempty"`
}

// SpecName returns spec.forProvider.name for the external-name initializer.
func (t *ExternalTable) SpecName() string {
	if t.Spec.ForProvider.Name == nil {
		return ""
	}
	return *t.Spec.ForProvider.Name
}

// +kubebuilder:object:root=true

// ExternalTableList contains a list of ExternalTable
type ExternalTableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalTable `json:"items"`
}

// ExternalTable type metadata.
var (
	ExternalTableKind             = reflect.TypeOf(ExternalTable{}).Name()
	ExternalTableGroupKind        = schema.GroupKind{Group: Group, Kind: ExternalTableKind}.String()
	ExternalTableKindAPIVersion   = ExternalTableKind + "." + SchemeGroupVersion.String()
	ExternalTableGroupVersionKind = SchemeGroupVersion.WithKind(ExternalTableKind)
)

func init() {
	SchemeBuilder.Register(&ExternalTable{}, &ExternalTableList{})
}
