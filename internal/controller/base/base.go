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

// Package base holds the glue shared by all managed resource controllers:
// resolving the ProviderConfig into a pooled Kusto client, the standard
// reconciler options, the spec.forProvider.name initializer and the hash
// annotations used for KQL drift tolerance.
package base

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/feature"
	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/crossplane-runtime/v2/pkg/statemetrics"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	apisv1alpha1 "github.com/functional-team/provider-azure-adx/apis/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
)

const (
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCPC       = "cannot get ClusterProviderConfig"
	errGetCreds     = "cannot get credentials"
	errNewClient    = "cannot create Kusto client"
)

// Deps are the provider-wide dependencies handed to every controller.
type Deps struct {
	Pool  *kusto.Pool
	Cache *snapshot.Cache
	Log   logging.Logger
}

// Connector resolves a managed resource's ProviderConfig into a Kusto client.
type Connector struct {
	Kube  client.Client
	Usage *resource.ProviderConfigUsageTracker
	Deps  Deps
}

// NewConnector builds a connector for mgr.
func NewConnector(mgr ctrl.Manager, d Deps) *Connector {
	return &Connector{
		Kube:  mgr.GetClient(),
		Usage: resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
		Deps:  d,
	}
}

// Connect tracks the ProviderConfig usage and returns the pooled client.
func (c *Connector) Connect(ctx context.Context, mg resource.ModernManaged) (kusto.Client, error) {
	if err := c.Usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}
	key, spec, err := c.providerConfig(ctx, mg)
	if err != nil {
		return nil, err
	}
	cfg, err := c.config(ctx, spec)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}
	kc, err := c.Deps.Pool.Get(key, cfg)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}
	return kc, nil
}

func (c *Connector) providerConfig(ctx context.Context, mg resource.ModernManaged) (string, apisv1alpha1.ProviderConfigSpec, error) {
	ref := mg.GetProviderConfigReference()
	if ref == nil {
		return "", apisv1alpha1.ProviderConfigSpec{}, errors.New("providerConfigRef is not set")
	}
	switch ref.Kind {
	case apisv1alpha1.ProviderConfigKind, "":
		pc := &apisv1alpha1.ProviderConfig{}
		if err := c.Kube.Get(ctx, types.NamespacedName{Name: ref.Name, Namespace: mg.GetNamespace()}, pc); err != nil {
			return "", apisv1alpha1.ProviderConfigSpec{}, errors.Wrap(err, errGetPC)
		}
		return "ProviderConfig/" + mg.GetNamespace() + "/" + ref.Name, pc.Spec, nil
	case apisv1alpha1.ClusterProviderConfigKind:
		cpc := &apisv1alpha1.ClusterProviderConfig{}
		if err := c.Kube.Get(ctx, types.NamespacedName{Name: ref.Name}, cpc); err != nil {
			return "", apisv1alpha1.ProviderConfigSpec{}, errors.Wrap(err, errGetCPC)
		}
		return "ClusterProviderConfig/" + ref.Name, cpc.Spec, nil
	}
	return "", apisv1alpha1.ProviderConfigSpec{}, errors.Errorf("unsupported provider config kind: %s", ref.Kind)
}

// servicePrincipal is the JSON layout of the credentials Secret.
type servicePrincipal struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	TenantID     string `json:"tenantId"`
}

func (c *Connector) config(ctx context.Context, spec apisv1alpha1.ProviderConfigSpec) (kusto.Config, error) {
	cfg := kusto.Config{
		Endpoint:    strings.TrimRight(spec.ClusterURI, "/"),
		Environment: string(spec.AzureEnvironment),
		Auth:        kusto.Auth{Source: string(spec.Credentials.Source)},
	}
	switch spec.Credentials.Source {
	case apisv1alpha1.CredentialsSourceSecret:
		if spec.Credentials.SecretRef == nil {
			return cfg, errors.New("credentials.secretRef is required for source Secret")
		}
		data, err := resource.ExtractSecret(ctx, c.Kube, xpv2.CommonCredentialSelectors{SecretRef: spec.Credentials.SecretRef})
		if err != nil {
			return cfg, err
		}
		var sp servicePrincipal
		if err := json.Unmarshal(data, &sp); err != nil {
			return cfg, fmt.Errorf("credentials secret must be a JSON document with clientId, clientSecret and tenantId: %w", err)
		}
		cfg.Auth.ClientID, cfg.Auth.ClientSecret, cfg.Auth.TenantID = sp.ClientID, sp.ClientSecret, sp.TenantID
	case apisv1alpha1.CredentialsSourceWorkloadIdentity:
		if wi := spec.Credentials.WorkloadIdentity; wi != nil {
			cfg.Auth.ClientID, cfg.Auth.TenantID, cfg.Auth.TokenFile = wi.ClientID, wi.TenantID, wi.TokenFile
		}
	case apisv1alpha1.CredentialsSourceManagedIdentity:
		if mi := spec.Credentials.ManagedIdentity; mi != nil {
			cfg.Auth.ManagedIdentityType, cfg.Auth.ClientID, cfg.Auth.ResourceID = mi.Type, mi.ClientID, mi.ResourceID
		}
	case apisv1alpha1.CredentialsSourceNone:
	default:
		return cfg, errors.Errorf("unsupported credentials source %q", spec.Credentials.Source)
	}
	return cfg, nil
}

// Kind bundles what the generic registration needs to know about a managed
// resource type.
type Kind[T resource.ModernManaged] struct {
	GVK    schema.GroupVersionKind
	Object T
	List   resource.ManagedList
}

// SetupGated registers setup to run once the CRD for k is present (safe-start).
func SetupGated[T resource.ModernManaged](o controller.Options, k Kind[T], setup func() error) error {
	o.Gate.Register(func() {
		if err := setup(); err != nil {
			panic(errors.Wrapf(err, "cannot setup %s controller", k.GVK.Kind))
		}
	}, k.GVK)
	return nil
}

// Register creates the managed reconciler for k with the standard options
// (initializers, poll interval, recorder, metrics, management policies, change
// logs) plus extra options, and adds it to mgr.
func Register[T resource.ModernManaged](mgr ctrl.Manager, o controller.Options, k Kind[T], conn managed.TypedExternalConnector[T], extra ...managed.ReconcilerOption) error {
	name := managed.ControllerName(k.GVK.GroupKind().String())

	opts := []managed.ReconcilerOption{
		managed.WithTypedExternalConnector[T](conn),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))), //nolint:staticcheck // TODO(jbw976) Crossplane needs to update to the new events API, see https://github.com/crossplane/crossplane/issues/7152
		managed.WithInitializers(NewSpecNameAsExternalName(mgr.GetClient()), managed.NewNameAsExternalName(mgr.GetClient())),
	}
	if o.Features.Enabled(feature.EnableBetaManagementPolicies) {
		opts = append(opts, managed.WithManagementPolicies())
	}
	if o.Features.Enabled(feature.EnableAlphaChangeLogs) && o.ChangeLogOptions != nil {
		opts = append(opts, managed.WithChangeLogger(o.ChangeLogOptions.ChangeLogger))
	}
	if o.MetricOptions != nil {
		opts = append(opts, managed.WithMetricRecorder(o.MetricOptions.MRMetrics))
		if o.MetricOptions.MRStateMetrics != nil {
			rec := statemetrics.NewMRStateRecorder(mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, k.List, o.MetricOptions.PollStateMetricInterval)
			if err := mgr.Add(rec); err != nil {
				return errors.Wrapf(err, "cannot register MR state metrics recorder for kind %s", k.GVK.Kind)
			}
		}
	}
	opts = append(opts, extra...)

	r := managed.NewReconciler(mgr, resource.ManagedKind(k.GVK), opts...)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(k.Object).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}
