# Permissions of the provider principal

The identity configured in the `ProviderConfig` runs management commands. Grant
the least that covers the kinds you use. Database roles are assigned with
`.add database <DB> <role> (...)` (or the ARM principal assignment in
provider-upbound-azure); cluster roles only exist at ARM level (`Cluster
Principal Assignment`) or in the Azure portal.

| Kinds | Minimum role | Notes |
|---|---|---|
| `Table`, `Function`, `MaterializedView`, `IngestionMapping`, `EntityGroup`, all `policy.*` kinds on tables and materialized views | **Database Admin** on the target database | Database User + Table Admin is enough for tables you already own, but the provider creates new entities, so Database Admin is the practical minimum. |
| Database level policies (`RetentionPolicy` with `entity.kind: Database`, `IngestionBatchingPolicy`, `ManagedIdentityPolicy`, ...) | Database Admin | |
| `ExternalTable` with a managed identity connection string (`;managed_identity=...`) | **AllDatabasesAdmin** | Kusto requires AllDatabasesAdmin for `.create-or-alter external table` with a managed identity, and a `ManagedIdentityPolicy` allowing `ExternalTable` usage. |
| `ContinuousExport` with `managedIdentity` | Database Admin plus a `ManagedIdentityPolicy` with `AutomatedFlows` usage | The identity itself needs write access to the export storage. |
| `SecurityRole` | Database Admin (database and table roles) | Only admins can change principals. |
| `cluster.*` kinds (`WorkloadGroup`, `CapacityPolicy`, ...) | **AllDatabasesAdmin** | Cluster policies are cluster-wide; grant them to a dedicated ProviderConfig. |

Recommendations:

- One `ProviderConfig` per cluster with Database Admin on the databases it
  manages. Use a second `ClusterProviderConfig` with AllDatabasesAdmin only for
  the Tier 2 kinds and managed-identity external tables.
- Workload identity or managed identity over client secrets where possible.
  Rotating a client secret only needs a Secret update; the provider rebuilds
  its pooled client when the credentials change.
- Whoever can create `Function`, `UpdatePolicy`, `RowLevelSecurityPolicy`,
  `MaterializedView` or `ContinuousExport` resources in Kubernetes can run
  arbitrary KQL as the provider principal. Restrict these kinds with Kubernetes
  RBAC to the teams that own the database.

Kubernetes side: the provider needs the standard Crossplane provider RBAC
(managed by Crossplane) plus read access to the Secrets referenced by
ProviderConfigs. Secrets are served from the informer cache by default; disable
that with `--enable-secret-cache=false` if memory matters more than API load.
