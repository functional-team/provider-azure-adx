# provider-azure-adx

`provider-azure-adx` is a [Crossplane](https://crossplane.io/) provider for the
artefacts that live *behind the Kusto endpoint* of an Azure Data Explorer (ADX)
cluster: tables, functions, materialized views, external tables, continuous
exports, ingestion mappings, entity groups, database and table policies,
security roles and cluster-level policies. It manages them declaratively through
Kusto management commands and continuously reconciles drift.

Everything that is an ARM resource (cluster, database, data connections,
principal assignments, database scripts, attached databases, private endpoints)
is out of scope and stays with
[provider-upbound-azure](https://marketplace.upbound.io/providers/upbound/provider-azure-kusto)
(`kusto.azure.upbound.io`). See [`examples/composition`](examples/composition)
for how the two providers compose.

Microsoft's own declarative option, the ARM *database script* (`kusto_script`),
is fire-and-forget: it neither observes nor repairs drift. That gap is what this
provider fills.

Status: pre-release. Packages are published to `ghcr.io/functional-team`
(private until the first public release, see [Roadmap](#roadmap)).

## Requirements

- Crossplane **2.0 or newer**. All managed resources are namespaced (Crossplane v2
  style); there are no cluster-scoped legacy kinds.
- An ADX cluster (or the Kusto emulator for tests). Microsoft Fabric Eventhouse
  uses the same engine and should work, but is untested and unsupported.
- A principal for the provider with the rights listed in
  [docs/permissions.md](docs/permissions.md).

## Install

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-azure-adx
spec:
  package: ghcr.io/functional-team/provider-azure-adx:v0.1.0
```

While the package is private, the `crossplane-system` namespace needs a pull
secret with a GitHub token that has `read:packages`, referenced via
`spec.packagePullSecrets`:

```sh
kubectl -n crossplane-system create secret docker-registry ghcr-functional-team \
  --docker-server=ghcr.io --docker-username=<github-user> --docker-password=<token>
```

```yaml
spec:
  package: ghcr.io/functional-team/provider-azure-adx:v0.1.0
  packagePullSecrets:
    - name: ghcr-functional-team
```

A step-by-step smoke test against a real cluster is in
[`hack/smoke/README.md`](hack/smoke/README.md).

## Connect to a cluster

One `ProviderConfig` (namespaced) or `ClusterProviderConfig` (cluster-scoped) per
ADX cluster. The database is *not* part of the ProviderConfig; every
database-scoped resource names its database in `spec.forProvider.database`.

```yaml
apiVersion: adx.functional.team/v1alpha1
kind: ProviderConfig
metadata:
  name: telemetry-prod
  namespace: data-platform
spec:
  clusterUri: https://telemetry-prod.westeurope.kusto.windows.net
  credentials:
    source: Secret            # Secret | WorkloadIdentity | ManagedIdentity | None
    secretRef:
      namespace: data-platform
      name: adx-sp
      key: credentials        # {"clientId": "...", "clientSecret": "...", "tenantId": "..."}
```

| `credentials.source` | Uses | Notes |
|---|---|---|
| `Secret` | Service principal from a Kubernetes Secret | JSON layout is compatible with the provider-upbound-azure credentials Secret. |
| `WorkloadIdentity` | AKS workload identity (federated token) | `clientId`, `tenantId`, `tokenFile` default to the `AZURE_*` variables injected by the webhook. |
| `ManagedIdentity` | System- or user-assigned managed identity | `type: SystemAssigned` or `UserAssigned` with `clientId` or `resourceId`. |
| `None` | No authentication | Only allowed for `http://` endpoints, i.e. the Kusto emulator. |

`spec.azureEnvironment` selects `AzureChinaCloud` or `AzureUSGovernment`. Both are
wired through to the SDK but untested. See [`examples/provider`](examples/provider).

## Kinds

Every kind has an example under [`examples/`](examples) and a generated reference
in [`docs/api`](docs/api/README.md).

| Group | Kinds |
|---|---|
| `adx.functional.team` | `ProviderConfig`, `ClusterProviderConfig`, `Table`, `Function`, `MaterializedView`, `ExternalTable`, `ContinuousExport`, `IngestionMapping`, `EntityGroup` |
| `policy.adx.functional.team` | `RetentionPolicy`, `CachingPolicy`, `UpdatePolicy`, `RowLevelSecurityPolicy`, `IngestionBatchingPolicy`, `StreamingIngestionPolicy`, `MergePolicy`, `ShardingPolicy`, `PartitioningPolicy`, `IngestionTimePolicy`, `AutoDeletePolicy`, `RestrictedViewAccessPolicy`, `ExtentTagsRetentionPolicy`, `EncodingPolicy`, `ManagedIdentityPolicy`, `RowOrderPolicy`, `QueryAccelerationPolicy` |
| `security.adx.functional.team` | `SecurityRole` |
| `cluster.adx.functional.team` | `WorkloadGroup`, `RequestClassificationPolicy`, `CalloutPolicy`, `CapacityPolicy`, `SandboxPolicy`, `QueryWeakConsistencyPolicy`, `ClusterManagedIdentityPolicy`, `MultiDatabaseAdminsPolicy` |

### Naming

Kusto names are usually PascalCase and would not be valid Kubernetes names. The
`crossplane.io/external-name` annotation is the source of truth for the Kusto
name. You can set it directly, or set `spec.forProvider.name` and the provider
copies it into the annotation on the first reconcile. `metadata.name` is only
the fallback. Both are immutable after the first reconcile: renaming means a
new Kusto entity (the old one is orphaned, not dropped).

### Policies

Policies are separate kinds with an `entity` discriminator, mirroring Kusto:

```yaml
apiVersion: policy.adx.functional.team/v1alpha1
kind: RetentionPolicy
spec:
  forProvider:
    database: Telemetry
    entity:
      kind: Table            # Table | MaterializedView | Database | ExternalTable, depending on the policy
      nameRef:
        name: raw-events     # or name: RawEvents
    softDeletePeriod: 365d
```

- A policy "exists" when it is set *on this entity*. An inherited database policy
  does not count, so a fresh table without its own retention policy is created,
  not adopted.
- Only fields set in the spec are compared. Unset fields keep the Kusto defaults.
- Deleting a policy managed resource runs `.delete ... policy`, which restores
  inheritance. It never drops data.
- Timespans accept Kusto literals (`365d`, `1.5h`, `10m`), the colon form and the
  .NET form (`365.00:00:00`); comparison is by value.

## Guardrails against data loss

Crossplane reconciles without a plan/apply step, so the API is designed to make
destructive changes explicit. Details in [docs/guardrails.md](docs/guardrails.md).

- `Table.spec.forProvider.schemaUpdateMode` defaults to `Merge`: columns are only
  added (`.alter-merge`). Columns that exist in the cluster but not in the spec
  are kept and listed in `status.atProvider.driftColumns`. `Replace` makes the
  spec authoritative and drops columns *and their data*.
- A column type change is never applied. The resource shows `Synced=False` with
  reason `UnsupportedColumnTypeChange` and sends no command.
- Materialized view queries are changed with `.alter materialized-view`. If Kusto
  rejects the change (for example a changed group-by), the resource stays
  `Synced=False` with reason `MaterializedViewAlterRejected`. Nothing is dropped
  or recreated automatically.
- Deleting a `Table` or `MaterializedView` drops it with all data, as Crossplane
  convention dictates. Use `spec.deletionPolicy: Orphan` or
  `managementPolicies` to keep it.
- External table connection strings with secrets are write-only: they are sent
  obfuscated (`h@'...'`) and never appear in status, events or logs. Secret
  rotation is detected through a hash annotation. Prefer managed identity
  (`;managed_identity=system`).
- `SecurityRole` with `mode: Authoritative` (the default) removes principals that
  are not in the spec, including ones added by hand. Use `mode: Additive` when
  several teams share a role.

## KQL is code

Function bodies, update policy queries, row level security queries, materialized
view queries and continuous export queries are KQL and are sent verbatim.
Whoever may create these managed resources runs arbitrary KQL with the
provider principal's permissions. Scope the Kubernetes RBAC for these kinds
accordingly.

Kusto reformats KQL when it stores it. The provider therefore compares text in
two stages: a deterministic normalization (line endings, trailing whitespace,
parameter list spacing, type aliases), and a hash of what was last written and
what was read back right after that write. A resource is up to date when either
stage matches; a change in the spec or in the cluster changes one of the hashes
and is detected. This is what keeps hundreds of functions from being rewritten
every poll.

## Scale and rate limiting

Kusto throttles management commands per cluster, and thousands of managed
resources polling individually would exhaust that budget. The provider:

- polls every **10 minutes** by default (`--poll`),
- observes in batches per database: one `.show database schema as json`, one
  `.show functions`, one `.show table * policy <name>` per database and TTL
  instead of one command per resource (`--observe-cache-ttl`, default 60s;
  `--disable-observe-cache` switches to single `.show` commands),
- limits commands per cluster (`--kusto-commands-per-second`, default 5;
  `--kusto-max-inflight`, default 4) and pools one SDK client per
  ProviderConfig,
- treats HTTP 429 and transient errors as retryable with backoff and exports
  `adx_commands_total`, `adx_command_duration_seconds`,
  `adx_commands_throttled_total` and `adx_observe_cache_requests_total`.

If `adx_commands_throttled_total` grows, raise `--observe-cache-ttl` or `--poll`
and lower `--kusto-commands-per-second`.

## Import existing artefacts

Set the `crossplane.io/external-name` annotation to the Kusto name and start with
`spec.managementPolicies: ["Observe"]`. The provider reads the entity into
`status.atProvider` without writing. Widen the policies when the spec matches.

## Development

```sh
make submodules      # crossplane/build
make reviewable      # generate, lint, unit tests
make build           # provider binary and xpkg
make test-integration  # Kusto emulator via testcontainers (linux/amd64 only)
make run             # run out-of-cluster against the current kubeconfig
```

The Kusto emulator does not run on ARM (Apple Silicon). Integration tests run in
CI on `ubuntu-latest`; locally point them at a development cluster with
`ADX_TEST_CLUSTER_URI`, `ADX_TEST_CLIENT_ID`, `ADX_TEST_CLIENT_SECRET`,
`ADX_TEST_TENANT_ID` and `ADX_TEST_DATABASE`. End-to-end tests against a real
cluster live in [`test/e2e`](test/e2e).

Design documents: [docs/concept.md](docs/concept.md) (decision log),
[docs/tech-implement.md](docs/tech-implement.md) (implementation plan),
[docs/implementation-notes.md](docs/implementation-notes.md) (deviations found
while building), [docs/spikes.md](docs/spikes.md) (what is still unverified
against a real cluster), [docs/migration.md](docs/migration.md) (from the
Terraform provider).

## Roadmap

- M0/M1/M2/M3 code is complete: client layer, Table, Function, all Tier 1
  policies, MaterializedView, ExternalTable, ContinuousExport, IngestionMapping,
  SecurityRole, EntityGroup and the cluster-level Tier 2 kinds.
- Before `v0.1.0`: run the emulator integration suite in CI, run the e2e suite
  against a development cluster and close the spikes in
  [docs/spikes.md](docs/spikes.md). The API groups live under
  `functional.team`; the project is maintained by functional.team and is not
  planned to move to crossplane-contrib (a group rename would break every
  manifest).
- Later: `v1beta1` for Tier 1, SQL/Cosmos external tables, GraphModel.

## Contributing

This project follows the Crossplane
[contributing guide](https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md),
the [code of conduct](CODE_OF_CONDUCT.md) and requires a
[DCO](DCO) sign-off on every commit (`git commit -s`). Maintainers are listed in
[OWNERS.md](OWNERS.md).

## License

Apache-2.0, see [LICENSE](LICENSE).
