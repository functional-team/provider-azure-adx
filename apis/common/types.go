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

package common

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// EntityKind is the kind of Kusto entity a policy or security role targets.
type EntityKind string

const (
	// EntityKindDatabase targets the database itself.
	EntityKindDatabase EntityKind = "Database"
	// EntityKindTable targets a table.
	EntityKindTable EntityKind = "Table"
	// EntityKindMaterializedView targets a materialized view.
	EntityKindMaterializedView EntityKind = "MaterializedView"
	// EntityKindExternalTable targets an external table.
	EntityKindExternalTable EntityKind = "ExternalTable"
	// EntityKindFunction targets a stored function.
	EntityKindFunction EntityKind = "Function"
)

// EntityReference addresses the target of a policy or security role inside a
// database. The target may be given by name or resolved from a managed
// resource of this provider (Table, MaterializedView, ExternalTable or
// Function, depending on kind).
// +kubebuilder:validation:XValidation:rule="self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector)",message="name, nameRef or nameSelector is required unless kind is Database"
type EntityReference struct {
	// Kind of the target entity.
	// +kubebuilder:validation:Enum=Database;Table;MaterializedView;ExternalTable;Function
	Kind EntityKind `json:"kind"`

	// Name of the target entity (Kusto name). Ignored for kind Database.
	// +optional
	Name *string `json:"name,omitempty"`

	// NameRef references a managed resource whose external name is used as
	// the target. The referenced kind follows from kind.
	// +optional
	NameRef *xpv2.NamespacedReference `json:"nameRef,omitempty"`

	// NameSelector selects a managed resource whose external name is used as
	// the target. The selected kind follows from kind.
	// +optional
	NameSelector *xpv2.NamespacedSelector `json:"nameSelector,omitempty"`
}

// Timespan is a Kusto timespan. Accepted forms are Kusto literals such as
// "30d", "1.5h", "10m", "30s", "time(1.02:03:04)", the colon form
// "hh:mm:ss" and the .NET form "d.hh:mm:ss[.fffffff]". Comparison against the
// cluster is done on the parsed value, so "1d" and "1.00:00:00" are equal.
// +kubebuilder:validation:MinLength=1
type Timespan string

// Principal is a Kusto principal in FQN notation, e.g.
// "aaduser=alice@contoso.com", "aadapp=<appId>;<tenantId>" or
// "aadgroup=<name-or-objectId>;<tenantId>".
// +kubebuilder:validation:MinLength=1
type Principal string

// ColumnType is a Kusto scalar data type. Aliases (boolean, date, uuid,
// uniqueid, double, time) are accepted and normalized to the canonical type
// before comparison, because the cluster always reports canonical types.
// +kubebuilder:validation:Enum=bool;boolean;datetime;date;dynamic;guid;uuid;uniqueid;int;long;real;double;decimal;string;timespan;time
type ColumnType string

// Column is one column of a table schema.
type Column struct {
	// Name of the column.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Name string `json:"name"`

	// Type of the column.
	Type ColumnType `json:"type"`

	// Docstring of the column.
	// +optional
	Docstring *string `json:"docstring,omitempty"`
}

// ObservedColumn is a column as reported by the cluster.
type ObservedColumn struct {
	// Name of the column.
	Name string `json:"name"`
	// Type of the column (canonical Kusto type).
	Type string `json:"type"`
	// Docstring of the column.
	// +optional
	Docstring string `json:"docstring,omitempty"`
}
