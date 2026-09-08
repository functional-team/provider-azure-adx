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

// ExtentTagsRetentionRule drops tags with a prefix after a period.
type ExtentTagsRetentionRule struct {
	// TagPrefix of the tags to drop, e.g. "drop-by:".
	// +kubebuilder:validation:MinLength=1
	TagPrefix string `json:"tagPrefix"`
	// RetentionPeriod after which matching tags are removed.
	RetentionPeriod common.Timespan `json:"retentionPeriod"`
}

// ExtentTagsRetentionPolicyParameters are the configurable fields of a ExtentTagsRetentionPolicy.
// +kubebuilder:validation:XValidation:rule="self.entity.kind in ['Table', 'Database']",message="ExtentTagsRetentionPolicy can only target Table, Database"
type ExtentTagsRetentionPolicyParameters struct {
	PolicyTarget `json:",inline"`

	// Rules is the complete list of tag retention rules.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Rules []ExtentTagsRetentionRule `json:"rules,omitempty"`
}

// A ExtentTagsRetentionPolicySpec defines the desired state of a ExtentTagsRetentionPolicy.
type ExtentTagsRetentionPolicySpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              ExtentTagsRetentionPolicyParameters `json:"forProvider"`
}

// A ExtentTagsRetentionPolicyStatus represents the observed state of a ExtentTagsRetentionPolicy.
type ExtentTagsRetentionPolicyStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 PolicyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ExtentTagsRetentionPolicy removes extent tags with a prefix after a retention period.
// Only fields set in the spec are compared with the cluster; unset fields keep
// the Kusto defaults. Deleting the managed resource deletes the policy on the
// entity (".delete ... policy extent_tags_retention"), which restores inheritance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DATABASE",type="string",JSONPath=".spec.forProvider.database"
// +kubebuilder:printcolumn:name="ENTITY-KIND",type="string",JSONPath=".spec.forProvider.entity.kind"
// +kubebuilder:printcolumn:name="ENTITY",type="string",JSONPath=".spec.forProvider.entity.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,adx,adxpolicy}
type ExtentTagsRetentionPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExtentTagsRetentionPolicySpec   `json:"spec"`
	Status ExtentTagsRetentionPolicyStatus `json:"status,omitempty"`
}

// GetTarget returns the database and entity the policy applies to.
func (p *ExtentTagsRetentionPolicy) GetTarget() *PolicyTarget {
	return &p.Spec.ForProvider.PolicyTarget
}

// GetPolicyObservation returns the observed policy.
func (p *ExtentTagsRetentionPolicy) GetPolicyObservation() PolicyObservation {
	return p.Status.AtProvider
}

// SetPolicyObservation stores the observed policy in the status.
func (p *ExtentTagsRetentionPolicy) SetPolicyObservation(o PolicyObservation) {
	p.Status.AtProvider = o
}

// +kubebuilder:object:root=true

// ExtentTagsRetentionPolicyList contains a list of ExtentTagsRetentionPolicy
type ExtentTagsRetentionPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExtentTagsRetentionPolicy `json:"items"`
}

// ExtentTagsRetentionPolicy type metadata.
var (
	ExtentTagsRetentionPolicyKind             = reflect.TypeOf(ExtentTagsRetentionPolicy{}).Name()
	ExtentTagsRetentionPolicyGroupKind        = schema.GroupKind{Group: Group, Kind: ExtentTagsRetentionPolicyKind}.String()
	ExtentTagsRetentionPolicyKindAPIVersion   = ExtentTagsRetentionPolicyKind + "." + SchemeGroupVersion.String()
	ExtentTagsRetentionPolicyGroupVersionKind = SchemeGroupVersion.WithKind(ExtentTagsRetentionPolicyKind)
)

func init() {
	SchemeBuilder.Register(&ExtentTagsRetentionPolicy{}, &ExtentTagsRetentionPolicyList{})
}
