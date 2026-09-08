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

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	changelogsv1alpha1 "github.com/crossplane/crossplane-runtime/v2/apis/changelogs/proto/v1alpha1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/feature"
	"github.com/crossplane/crossplane-runtime/v2/pkg/gate"
	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/customresourcesgate"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/statemetrics"

	"github.com/functional-team/provider-azure-adx/apis"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	adx "github.com/functional-team/provider-azure-adx/internal/controller"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
	"github.com/functional-team/provider-azure-adx/internal/version"
)

func main() {
	var (
		app            = kingpin.New(filepath.Base(os.Args[0]), "Azure Data Explorer (Kusto) support for Crossplane.").DefaultEnvars()
		debug          = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		leaderElection = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").Envar("LEADER_ELECTION").Bool()

		syncInterval            = app.Flag("sync", "How often all resources will be double-checked for drift from the desired state.").Short('s').Default("1h").Duration()
		pollInterval            = app.Flag("poll", "How often individual resources will be checked for drift from the desired state. Kusto management commands are rate limited per cluster, so keep this high for large fleets.").Default("10m").Duration()
		pollStateMetricInterval = app.Flag("poll-state-metric", "State metric recording interval").Default("5s").Duration()

		maxReconcileRate = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may checked for drift from the desired state.").Default("10").Int()

		observeCacheTTL     = app.Flag("observe-cache-ttl", "How long batch-observe results (e.g. all tables of a database) are cached. 0 disables the cache.").Default("60s").Duration()
		disableObserveCache = app.Flag("disable-observe-cache", "Disable the batch-observe cache and run one .show per resource.").Default("false").Bool()
		kustoCommandsPerSec = app.Flag("kusto-commands-per-second", "Sustained rate of management commands per cluster.").Default("5").Float64()
		kustoMaxInflight    = app.Flag("kusto-max-inflight", "Maximum concurrent management commands per cluster.").Default("4").Int()
		kustoClientIdle     = app.Flag("kusto-client-idle-timeout", "Idle time after which a pooled Kusto client is closed.").Default("30m").Duration()

		enableManagementPolicies = app.Flag("enable-management-policies", "Enable support for Management Policies.").Default("true").Envar("ENABLE_MANAGEMENT_POLICIES").Bool()
		enableChangeLogs         = app.Flag("enable-changelogs", "Enable support for capturing change logs during reconciliation.").Default("false").Envar("ENABLE_CHANGE_LOGS").Bool()
		changelogsSocketPath     = app.Flag("changelogs-socket-path", "Path for changelogs socket (if enabled)").Default("/var/run/changelogs/changelogs.sock").Envar("CHANGELOGS_SOCKET_PATH").String()

		enableSecretCache = app.Flag("enable-secret-cache", "Enable caching of Secret objects. When true, Secrets are served from the informer cache instead of direct API calls. This reduces API server load but increases memory usage.").Default("true").Envar("ENABLE_SECRET_CACHE").Bool()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-azure-adx"))
	if *debug {
		// The controller-runtime is *very* verbose even at info level, so we only
		// provide it a real logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	} else {
		ctrl.SetLogger(zap.New(zap.WriteTo(io.Discard)))
	}

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	var clientOpts client.Options
	if !*enableSecretCache {
		clientOpts = client.Options{
			Cache: &client.CacheOptions{
				DisableFor: []client.Object{&corev1.Secret{}},
			},
		}
	}

	scheme := runtime.NewScheme()
	kingpin.FatalIfError(clientgoscheme.AddToScheme(scheme), "Cannot add clientgo scheme")
	kingpin.FatalIfError(apis.AddToScheme(scheme), "Cannot add Azure ADX APIs to scheme")
	kingpin.FatalIfError(apiextensionsv1.AddToScheme(scheme), "Cannot add CustomResourceDefinition to scheme")

	mgr, err := ctrl.NewManager(ratelimiter.LimitRESTConfig(cfg, *maxReconcileRate), ctrl.Options{
		Scheme: scheme,
		Cache: cache.Options{
			SyncPeriod: syncInterval,
			ByObject: map[client.Object]cache.ByObject{
				&apiextensionsv1.CustomResourceDefinition{}: {
					Transform: customresourcesgate.TransformStripCRDSchema,
				},
			},
		},
		Client: clientOpts,

		LeaderElection:             *leaderElection,
		LeaderElectionID:           "crossplane-leader-election-provider-azure-adx",
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")

	metricRecorder := managed.NewMRMetricRecorder()
	stateMetrics := statemetrics.NewMRStateMetrics()

	metrics.Registry.MustRegister(metricRecorder)
	metrics.Registry.MustRegister(stateMetrics)
	kingpin.FatalIfError(kusto.RegisterMetrics(metrics.Registry), "Cannot register Kusto metrics")
	kingpin.FatalIfError(snapshot.RegisterMetrics(metrics.Registry), "Cannot register cache metrics")

	o := controller.Options{
		Logger:                  log,
		MaxConcurrentReconciles: *maxReconcileRate,
		PollInterval:            *pollInterval,
		GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
		Features:                &feature.Flags{},
		Gate:                    new(gate.Gate[schema.GroupVersionKind]),
		MetricOptions: &controller.MetricOptions{
			PollStateMetricInterval: *pollStateMetricInterval,
			MRMetrics:               metricRecorder,
			MRStateMetrics:          stateMetrics,
		},
	}

	if *enableManagementPolicies {
		o.Features.Enable(feature.EnableBetaManagementPolicies)
		log.Info("Beta feature enabled", "flag", feature.EnableBetaManagementPolicies)
	}

	if *enableChangeLogs {
		o.Features.Enable(feature.EnableAlphaChangeLogs)
		log.Info("Alpha feature enabled", "flag", feature.EnableAlphaChangeLogs)

		conn, err := grpc.NewClient("unix://"+*changelogsSocketPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
		kingpin.FatalIfError(err, "failed to create change logs client connection at %s", *changelogsSocketPath)

		clo := controller.ChangeLogOptions{
			ChangeLogger: managed.NewGRPCChangeLogger(
				changelogsv1alpha1.NewChangeLogServiceClient(conn),
				managed.WithProviderVersion(fmt.Sprintf("provider-azure-adx:%s", version.Version))),
		}
		o.ChangeLogOptions = &clo
	}

	kustoOpts := kusto.Options{CommandsPerSecond: *kustoCommandsPerSec, MaxInflight: *kustoMaxInflight, Logger: log.WithValues("component", "kusto")}
	deps := base.Deps{
		Pool:  kusto.NewPool(*kustoClientIdle, func(c kusto.Config) (kusto.Client, error) { return kusto.New(c, kustoOpts) }),
		Cache: snapshot.New(*observeCacheTTL, *disableObserveCache),
		Log:   log,
	}
	log.Info("Kusto client settings", "poll", pollInterval.String(), "observe-cache-ttl", observeCacheTTL.String(), "cache-disabled", deps.Cache.Disabled(), "commands-per-second", *kustoCommandsPerSec, "max-inflight", *kustoMaxInflight)

	kingpin.FatalIfError(customresourcesgate.Setup(mgr, o), "Cannot setup CRD gate controller")
	kingpin.FatalIfError(adx.SetupGated(mgr, o, deps), "Cannot setup Azure ADX controllers")
	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}
