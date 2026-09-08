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

package kusto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-kusto-go/azkustodata"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

// Credential sources, mirrored from the ProviderConfig API without importing it.
const (
	SourceNone             = "None"
	SourceSecret           = "Secret"
	SourceWorkloadIdentity = "WorkloadIdentity"
	SourceManagedIdentity  = "ManagedIdentity"

	ManagedIdentitySystemAssigned = "SystemAssigned"
	ManagedIdentityUserAssigned   = "UserAssigned"

	EnvironmentPublic     = "AzurePublicCloud"
	EnvironmentChina      = "AzureChinaCloud"
	EnvironmentGovernment = "AzureUSGovernment"
)

// Auth is the resolved credential material of a ProviderConfig.
type Auth struct {
	Source string

	// Secret (service principal)
	ClientID     string
	ClientSecret string
	TenantID     string

	// WorkloadIdentity
	TokenFile string

	// ManagedIdentity
	ManagedIdentityType string
	ResourceID          string
}

// Config describes one cluster connection.
type Config struct {
	Endpoint    string
	Environment string
	Auth        Auth
}

// Fingerprint hashes the whole config including secrets so the pool can
// detect changed credentials.
func (c Config) Fingerprint() string {
	h := sha256.New()
	for _, s := range []string{c.Endpoint, c.Environment, c.Auth.Source, c.Auth.ClientID, c.Auth.ClientSecret, c.Auth.TenantID, c.Auth.TokenFile, c.Auth.ManagedIdentityType, c.Auth.ResourceID} {
		h.Write([]byte(s))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

// kcsbFor maps a Config onto the SDK connection string builder.
func kcsbFor(cfg Config) (*azkustodata.ConnectionStringBuilder, error) { //nolint:gocyclo // one switch per credential source is easier to audit than helpers.
	if strings.TrimSpace(cfg.Endpoint) == "" {
		return nil, fmt.Errorf("cluster endpoint is empty")
	}
	kcsb := azkustodata.NewConnectionStringBuilder(cfg.Endpoint)
	a := cfg.Auth
	switch a.Source {
	case SourceNone, "":
		if !strings.HasPrefix(strings.ToLower(cfg.Endpoint), "http://") {
			return nil, fmt.Errorf("credentials source None is only allowed for http:// endpoints (emulator)")
		}
	case SourceSecret:
		if a.ClientID == "" || a.ClientSecret == "" || a.TenantID == "" {
			return nil, fmt.Errorf("secret credentials need clientId, clientSecret and tenantId")
		}
		kcsb = kcsb.WithAadAppKey(a.ClientID, a.ClientSecret, a.TenantID)
	case SourceWorkloadIdentity:
		clientID := firstNonEmpty(a.ClientID, os.Getenv("AZURE_CLIENT_ID"))
		tenantID := firstNonEmpty(a.TenantID, os.Getenv("AZURE_TENANT_ID"))
		tokenFile := firstNonEmpty(a.TokenFile, os.Getenv("AZURE_FEDERATED_TOKEN_FILE"))
		if clientID == "" || tenantID == "" || tokenFile == "" {
			return nil, fmt.Errorf("workload identity needs clientId, tenantId and tokenFile (or AZURE_CLIENT_ID, AZURE_TENANT_ID, AZURE_FEDERATED_TOKEN_FILE)")
		}
		kcsb = kcsb.WithKubernetesWorkloadIdentity(clientID, tokenFile, tenantID)
	case SourceManagedIdentity:
		switch a.ManagedIdentityType {
		case ManagedIdentitySystemAssigned, "":
			kcsb = kcsb.WithSystemManagedIdentity()
		case ManagedIdentityUserAssigned:
			switch {
			case a.ClientID != "":
				kcsb = kcsb.WithUserAssignedIdentityClientId(a.ClientID)
			case a.ResourceID != "":
				kcsb = kcsb.WithUserAssignedIdentityResourceId(a.ResourceID)
			default:
				return nil, fmt.Errorf("user assigned managed identity needs clientId or resourceId")
			}
		default:
			return nil, fmt.Errorf("unknown managed identity type %q", a.ManagedIdentityType)
		}
	default:
		return nil, fmt.Errorf("unknown credentials source %q", a.Source)
	}
	if cc, ok := cloudFor(cfg.Environment); ok {
		kcsb = kcsb.AttachPolicyClientOptions(&azcore.ClientOptions{Cloud: cc})
	}
	return kcsb, nil
}

func cloudFor(env string) (cloud.Configuration, bool) {
	switch env {
	case EnvironmentChina:
		return cloud.AzureChina, true
	case EnvironmentGovernment:
		return cloud.AzureGovernment, true
	}
	return cloud.Configuration{}, false
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
