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

// Package emulator starts the Kusto emulator (kustainer) with testcontainers
// for integration tests, or, when ADX_TEST_CLUSTER_URI is set, wires the
// tests to a real development cluster instead. The emulator does not run on
// ARM hosts (deviation A3 in docs/tech-implement.md), so on Apple Silicon the
// dev cluster path is the only local option.
package emulator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
)

// Image is the Kusto emulator image. The tag is bumped by Renovate (see the
// custom manager in renovate.json), keep the "name:tag" form on one line.
const Image = "mcr.microsoft.com/azuredataexplorer/kustainer-linux:latest"

// Environment variables that switch the harness to a real development cluster.
const (
	EnvClusterURI   = "ADX_TEST_CLUSTER_URI"
	EnvClientID     = "ADX_TEST_CLIENT_ID"
	EnvClientSecret = "ADX_TEST_CLIENT_SECRET"
	EnvTenantID     = "ADX_TEST_TENANT_ID"
	EnvDatabase     = "ADX_TEST_DATABASE"
)

const (
	port           = "8080/tcp"
	memoryBytes    = 4 << 30
	startupTimeout = 5 * time.Minute
)

// ErrUnavailable is returned when neither Docker nor a dev cluster is
// reachable. Tests should t.Skip in that case.
var ErrUnavailable = errors.New("kusto emulator unavailable")

// IsUnavailable reports whether err means "no emulator and no cluster".
func IsUnavailable(err error) bool { return errors.Is(err, ErrUnavailable) }

var databaseNameRe = regexp.MustCompile(`^[A-Za-z0-9_]{1,64}$`)

// Emulator is a running emulator container or a connection to a dev cluster.
type Emulator struct {
	container testcontainers.Container
	cfg       kusto.Config
	database  string
}

// Start starts the emulator container, or connects to the dev cluster named
// by ADX_TEST_CLUSTER_URI.
func Start(ctx context.Context) (*Emulator, error) {
	if uri := os.Getenv(EnvClusterURI); uri != "" {
		return startCluster(uri)
	}
	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnavailable, err)
	}
	if err := provider.Health(ctx); err != nil {
		_ = provider.Close()
		return nil, fmt.Errorf("%w: docker daemon not reachable: %w", ErrUnavailable, err)
	}
	_ = provider.Close()

	req := testcontainers.ContainerRequest{
		Image:        Image,
		Env:          map[string]string{"ACCEPT_EULA": "Y"},
		ExposedPorts: []string{port},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.Memory = memoryBytes
		},
		WaitingFor: wait.ForHTTP("/v1/rest/mgmt").
			WithPort(port).
			WithMethod(http.MethodPost).
			WithBody(strings.NewReader(`{"csl":".show cluster"}`)).
			WithHeaders(map[string]string{"Content-Type": "application/json"}).
			WithStatusCodeMatcher(func(status int) bool { return status == http.StatusOK }).
			WithPollInterval(2 * time.Second).
			WithStartupTimeout(startupTimeout),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil {
		return nil, fmt.Errorf("cannot start kusto emulator %s: %w", Image, err)
	}
	host, err := c.Host(ctx)
	if err != nil {
		_ = c.Terminate(ctx)
		return nil, fmt.Errorf("cannot determine emulator host: %w", err)
	}
	mapped, err := c.MappedPort(ctx, port)
	if err != nil {
		_ = c.Terminate(ctx)
		return nil, fmt.Errorf("cannot determine emulator port: %w", err)
	}
	endpoint := fmt.Sprintf("http://%s:%s", host, mapped.Port())
	return &Emulator{
		container: c,
		cfg:       kusto.Config{Endpoint: endpoint, Auth: kusto.Auth{Source: kusto.SourceNone}},
	}, nil
}

func startCluster(uri string) (*Emulator, error) {
	db := os.Getenv(EnvDatabase)
	a := kusto.Auth{
		Source:       kusto.SourceSecret,
		ClientID:     os.Getenv(EnvClientID),
		ClientSecret: os.Getenv(EnvClientSecret),
		TenantID:     os.Getenv(EnvTenantID),
	}
	if db == "" || a.ClientID == "" || a.ClientSecret == "" || a.TenantID == "" {
		return nil, fmt.Errorf("%s is set, so %s, %s, %s and %s must be set too", EnvClusterURI, EnvDatabase, EnvClientID, EnvClientSecret, EnvTenantID)
	}
	return &Emulator{cfg: kusto.Config{Endpoint: strings.TrimRight(uri, "/"), Auth: a}, database: db}, nil
}

// IsEmulator reports whether the tests run against the container (true) or a
// real cluster (false).
func (e *Emulator) IsEmulator() bool { return e.container != nil }

// Endpoint is the cluster URI, e.g. http://localhost:32768.
func (e *Emulator) Endpoint() string { return e.cfg.Endpoint }

// Database returns the database tests should use: name on the emulator, the
// fixed ADX_TEST_DATABASE on a real cluster.
func (e *Emulator) Database(name string) string {
	if e.database != "" {
		return e.database
	}
	return name
}

// Client builds a Kusto client for the emulator or cluster.
func (e *Emulator) Client(opts kusto.Options) (kusto.Client, error) {
	return kusto.New(e.cfg, opts)
}

// CreateDatabase creates a persisted database in the emulator. It is a no-op
// on a real cluster (databases are ARM resources there) and tolerates an
// existing database.
func (e *Emulator) CreateDatabase(ctx context.Context, name string) error {
	if !e.IsEmulator() {
		return nil
	}
	if !databaseNameRe.MatchString(name) {
		return fmt.Errorf("database name %q must match %s (it is used in a file system path)", name, databaseNameRe)
	}
	kc, err := e.Client(kusto.Options{})
	if err != nil {
		return err
	}
	defer closeClient(kc)
	c := cmd.New(".create database ", cmd.Ident(name),
		" persist (@\"/kustodata/dbs/", name, "/md\", @\"/kustodata/dbs/", name, "/data\")")
	if _, err := kc.Mgmt(ctx, "", c); err != nil && !kerrors.IsAlreadyExists(err) {
		return fmt.Errorf("cannot create database %s: %w", name, err)
	}
	return nil
}

// Terminate stops the container. It is a no-op for a real cluster.
func (e *Emulator) Terminate(ctx context.Context) error {
	if e.container == nil {
		return nil
	}
	return e.container.Terminate(ctx)
}

func closeClient(kc kusto.Client) {
	if c, ok := kc.(kusto.Closer); ok {
		_ = c.Close()
	}
}
