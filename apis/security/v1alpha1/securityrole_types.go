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

// Role is a Kusto security role.
// +kubebuilder:validation:Enum=admins;users;viewers;unrestrictedviewers;ingestors;monitors
type Role string

// Mode controls whether the managed resource owns the whole principal list of
// the role or only its own entries.
// +kubebuilder:validation:Enum=Authoritative;Additive
type Mode string

const (
	// ModeAuthoritative makes spec.principals the complete list of the role
	// (.set): principals added by others are removed.
	ModeAuthoritative Mode = "Authoritative"
	// ModeAdditive only adds and removes the principals listed in the spec
	// (.add/.drop) and leaves other principals alone.
	ModeAdditive Mode = "Additive"
)

// SecurityRoleParameters are the configurable fields of a SecurityRole.
// +kubebuilder:validation:XValidation:rule="self.entity.kind != 'Table' || self.role in ['admins','ingestors']",message="tables only have the admins and ingestors roles"
// +kubebuilder:validation:XValidation:rule="!(self.entity.kind in ['ExternalTable','MaterializedView','Function']) || self.role == 'admins'",message="external tables, materialized views and functions only have the admins role"
// +kubebuilder:validation:XValidation:rule="self.entity.kind == oldSelf.entity.kind",message="entity.kind is immutable"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name",message="entity.name is immutable once set"
type SecurityRoleParameters struct {
	// Database that holds the entity. Immutable.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="database is immutable"
	Database string `json:"database"`

	// Entity the role belongs to: Database, Table, ExternalTable,
	// MaterializedView or Function.
	Entity common.EntityReference `json:"entity"`

	// Role to manage. Immutable.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="role is immutable"
	Role Role `json:"role"`

	// Mode is Authoritative (default; the spec is the complete principal
	// list of the role) or Additive (only the listed principals are managed).
	// +kubebuilder:default=Authoritative
	// +optional
	Mode Mode `json:"mode,omitempty"`

	// Principals in Kusto notation, e.g. "aaduser=alice@contoso.com",
	// "aadapp=<appId>;<tenantId>", "aadgroup=<objectId>;<tenantId>". An empty
	// list in Authoritative mode removes every principal from the role.
	// +listType=set
	// +optional
	Principals []common.Principal `json:"principals,omitempty"`

	// Description recorded with the role assignment.
	// +optional
	Description *string `json:"description,omitempty"`
}

// ResolvedPrincipal maps one spec entry to the principal Kusto reports.
type ResolvedPrincipal struct {
	// Spec is the principal string from the spec.
	Spec string `json:"spec"`
	// ObjectID is the Entra object id Kusto reports for it.
	ObjectID string `json:"objectId"`
	// FQN is the fully qualified principal name Kusto reports.
	// +optional
	FQN string `json:"fqn,omitempty"`
	// DisplayName Kusto reports.
	// +optional
	DisplayName string `json:"displayName,omitempty"`
	// Type of the principal as reported (e.g. AAD User).
	// +optional
	Type string `json:"type,omitempty"`
}

// SecurityRoleObservation are the observable fields of a SecurityRole.
type SecurityRoleObservation struct {
	// ResolvedPrincipals maps spec entries to the object ids Kusto reports.
	// The provider compares sets of object ids, so the notation used in the
	// spec (UPN, app id, object id) does not matter.
	// +optional
	ResolvedPrincipals []ResolvedPrincipal `json:"resolvedPrincipals,omitempty"`
	// Unresolved lists spec entries that could not be mapped to a principal
	// reported by the cluster (see the provider documentation).
	// +optional
	Unresolved []string `json:"unresolved,omitempty"`
	// Principals are the FQNs currently assigned to the role in the cluster.
	// +optional
	Principals []string `json:"principals,omitempty"`
}

// A SecurityRoleSpec defines the desired state of a SecurityRole.
type SecurityRoleSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              SecurityRoleParameters `json:"forProvider"`
}

// A SecurityRoleStatus represents the observed state of a SecurityRole.
type SecurityRoleStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 SecurityRoleObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A SecurityRole assigns principals to a role of a database, table, external
// table, materialized view or function. Deleting an Authoritative resource
// removes all principals of that role; deleting an Additive one removes only
// the principals it added.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="ROLE",type="string",JSONPath=".spec.forProvider.role"
// +kubebuilder:printcolumn:name="MODE",type="string",JSONPath=".spec.forProvider.mode"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx}
type SecurityRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecurityRoleSpec   `json:"spec"`
	Status SecurityRoleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecurityRoleList contains a list of SecurityRole
type SecurityRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityRole `json:"items"`
}

// SecurityRole type metadata.
var (
	SecurityRoleKind             = reflect.TypeOf(SecurityRole{}).Name()
	SecurityRoleGroupKind        = schema.GroupKind{Group: Group, Kind: SecurityRoleKind}.String()
	SecurityRoleKindAPIVersion   = SecurityRoleKind + "." + SchemeGroupVersion.String()
	SecurityRoleGroupVersionKind = SchemeGroupVersion.WithKind(SecurityRoleKind)
)

func init() {
	SchemeBuilder.Register(&SecurityRole{}, &SecurityRoleList{})
}
