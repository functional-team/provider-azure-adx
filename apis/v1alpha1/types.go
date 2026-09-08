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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// CredentialsSource selects how the provider authenticates against the cluster.
type CredentialsSource string

const (
	// CredentialsSourceNone sends no credentials. Only allowed for http://
	// endpoints, i.e. the Kusto emulator.
	CredentialsSourceNone CredentialsSource = "None"
	// CredentialsSourceSecret reads a service principal (clientId, clientSecret,
	// tenantId) from a Kubernetes Secret.
	CredentialsSourceSecret CredentialsSource = "Secret"
	// CredentialsSourceWorkloadIdentity uses Azure Workload Identity federation
	// (AKS) with the projected service account token.
	CredentialsSourceWorkloadIdentity CredentialsSource = "WorkloadIdentity"
	// CredentialsSourceManagedIdentity uses a system- or user-assigned managed
	// identity of the node/pod.
	CredentialsSourceManagedIdentity CredentialsSource = "ManagedIdentity"
)

// AzureEnvironment names a sovereign cloud.
type AzureEnvironment string

const (
	// AzurePublicCloud is the default global Azure cloud.
	AzurePublicCloud AzureEnvironment = "AzurePublicCloud"
	// AzureChinaCloud is Azure China (21Vianet).
	AzureChinaCloud AzureEnvironment = "AzureChinaCloud"
	// AzureUSGovernment is Azure US Government.
	AzureUSGovernment AzureEnvironment = "AzureUSGovernment"
)

// WorkloadIdentityCredentials configure Azure Workload Identity. All fields are
// optional; empty values fall back to the AZURE_CLIENT_ID, AZURE_TENANT_ID and
// AZURE_FEDERATED_TOKEN_FILE environment variables injected by the webhook.
type WorkloadIdentityCredentials struct {
	// ClientID of the federated application.
	// +optional
	ClientID string `json:"clientId,omitempty"`
	// TenantID of the federated application.
	// +optional
	TenantID string `json:"tenantId,omitempty"`
	// TokenFile is the path of the projected service account token.
	// +optional
	TokenFile string `json:"tokenFile,omitempty"`
}

// ManagedIdentityCredentials configure a managed identity.
type ManagedIdentityCredentials struct {
	// Type of the managed identity.
	// +kubebuilder:validation:Enum=SystemAssigned;UserAssigned
	// +kubebuilder:default=SystemAssigned
	Type string `json:"type"`
	// ClientID of a user-assigned identity. Mutually exclusive with resourceId.
	// +optional
	ClientID string `json:"clientId,omitempty"`
	// ResourceID of a user-assigned identity. Mutually exclusive with clientId.
	// +optional
	ResourceID string `json:"resourceId,omitempty"`
}

// ProviderCredentials required to authenticate.
// +kubebuilder:validation:XValidation:rule="self.source != 'Secret' || has(self.secretRef)",message="secretRef is required when source is Secret"
// +kubebuilder:validation:XValidation:rule="self.source != 'ManagedIdentity' || has(self.managedIdentity)",message="managedIdentity is required when source is ManagedIdentity"
type ProviderCredentials struct {
	// Source of the provider credentials.
	// +kubebuilder:validation:Enum=None;Secret;WorkloadIdentity;ManagedIdentity
	Source CredentialsSource `json:"source"`

	// SecretRef points to a Secret key holding a JSON document with the keys
	// clientId, clientSecret and tenantId. The format is compatible with the
	// credentials Secret of provider-upbound-azure (additional keys are ignored).
	// +optional
	SecretRef *xpv2.SecretKeySelector `json:"secretRef,omitempty"`

	// WorkloadIdentity settings, used when source is WorkloadIdentity.
	// +optional
	WorkloadIdentity *WorkloadIdentityCredentials `json:"workloadIdentity,omitempty"`

	// ManagedIdentity settings, used when source is ManagedIdentity.
	// +optional
	ManagedIdentity *ManagedIdentityCredentials `json:"managedIdentity,omitempty"`
}

// A ProviderConfigSpec defines the desired state of a ProviderConfig.
// +kubebuilder:validation:XValidation:rule="self.credentials.source != 'None' || self.clusterUri.startsWith('http://')",message="credentials.source None is only allowed for http:// endpoints (Kusto emulator)"
type ProviderConfigSpec struct {
	// ClusterURI is the Kusto query endpoint of the cluster, e.g.
	// https://mycluster.westeurope.kusto.windows.net. Immutable.
	// +kubebuilder:validation:Pattern=`^https?://[^/\s]+/?$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="clusterUri is immutable"
	ClusterURI string `json:"clusterUri"`

	// AzureEnvironment selects the sovereign cloud. Only AzurePublicCloud is
	// tested; the others are wired through to the SDK but unverified.
	// +kubebuilder:validation:Enum=AzurePublicCloud;AzureChinaCloud;AzureUSGovernment
	// +kubebuilder:default=AzurePublicCloud
	// +optional
	AzureEnvironment AzureEnvironment `json:"azureEnvironment,omitempty"`

	// Credentials required to authenticate to this provider.
	Credentials ProviderCredentials `json:"credentials"`
}

// A ProviderConfigStatus defines the status of a Provider.
type ProviderConfigStatus struct {
	xpv2.ProviderConfigStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="CLUSTER",type="string",JSONPath=".spec.clusterUri"
// +kubebuilder:printcolumn:name="SOURCE",type="string",JSONPath=".spec.credentials.source"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,provider,adx}
// A ProviderConfig configures how the provider connects to one Azure Data
// Explorer cluster: the endpoint plus the identity used for management commands.
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="CONFIG-NAME",type="string",JSONPath=".providerConfigRef.name"
// +kubebuilder:printcolumn:name="RESOURCE-KIND",type="string",JSONPath=".resourceRef.kind"
// +kubebuilder:printcolumn:name="RESOURCE-NAME",type="string",JSONPath=".resourceRef.name"
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,provider,adx}
// A ProviderConfigUsage indicates that a resource is using a ProviderConfig or a
// ClusterProviderConfig. There is deliberately no cluster scoped usage type, Usages always live in
// the namespace of the MR that created them, and record which kind of config they refer to.
type ProviderConfigUsage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	xpv2.TypedProviderConfigUsage `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfigUsageList contains a list of ProviderConfigUsage
type ProviderConfigUsageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfigUsage `json:"items"`
}

// +kubebuilder:object:root=true

// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="CLUSTER",type="string",JSONPath=".spec.clusterUri"
// +kubebuilder:printcolumn:name="SOURCE",type="string",JSONPath=".spec.credentials.source"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Cluster,categories={crossplane,provider,adx}
// A ClusterProviderConfig is the cluster-scoped variant of ProviderConfig. It
// can be referenced by managed resources in any namespace.
type ClusterProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterProviderConfigList contains a list of ClusterProviderConfig.
type ClusterProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterProviderConfig `json:"items"`
}
