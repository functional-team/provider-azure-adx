# Migrating from the Terraform provider

The Terraform provider [`favoretti/adx`](https://registry.terraform.io/providers/favoretti/adx)
(source: github.com/favoretti/terraform-provider-adx, v0.0.43 at the time of
writing) manages a subset of the same artefacts. This chapter maps its
resources onto the kinds of provider-azure-adx and lists the behavioural
differences that matter when moving.

## Model differences

| Topic | Terraform provider | provider-azure-adx |
|---|---|---|
| Connection | provider block with cluster URI and credentials; one provider alias per cluster | `ProviderConfig` / `ClusterProviderConfig` per cluster, referenced by every resource |
| Database | attribute `database_name` on every resource | `spec.forProvider.database` on every database-scoped kind |
| Entity name | attribute `name` | `crossplane.io/external-name` annotation or `spec.forProvider.name` (copied into the annotation) |
| Plan/apply | plan shows destructive changes before apply | no plan step; guardrails in the API instead (`schemaUpdateMode`, type changes blocked, no auto-recreate), see [guardrails.md](guardrails.md) |
| Drift detection | on `terraform plan`/`refresh` | continuous (default every 10 minutes), batched per database |
| Policies on tables vs. materialized views | separate resource types (`adx_table_retention_policy`, `adx_materialized_view_retention_policy`) | one kind per policy type with `entity.kind: Table \| MaterializedView \| Database` |
| Import | `terraform import` | set the external-name annotation and start with `managementPolicies: ["Observe"]` |
| Secrets | in state (plaintext) | Kubernetes Secrets, sent obfuscated, never in status |

## Resource mapping

| Terraform resource | Kind | Notes |
|---|---|---|
| `adx_table` | `Table` | `column {}` blocks -> `columns[]`; `table_schema` (string) and `from_query` have no equivalent (F7: structured schema only). `merge_on_update` -> `schemaUpdateMode: Merge` (default) / `Replace`. |
| `adx_function` | `Function` | `parameters`, `body`, `folder`, `docstring` map 1:1; `view`/`skipvalidation` are write-only. |
| `adx_materialized_view` | `MaterializedView` | `source_table_name` -> `sourceTable`/`sourceTableRef`; `backfill`, `effective_date_time` create-only; `lookback` is a Kusto timespan; `enabled` toggles `.enable/.disable`. |
| `adx_table_mapping` | `IngestionMapping` | `kind` (csv/json/avro/parquet/orc/w3clogfile), `mapping[]` with `column`, `dataType`, `properties`; identity is (database, table, kind, name). |
| `adx_external_table` | `ExternalTable` | `kind: Storage \| Delta`; connection strings via `value` (managed identity) or `secretKeyRef` (SAS, write-only); `partitionBy`/`pathFormat` are raw Kusto expressions. |
| `adx_table_continuous_export` | `ContinuousExport` | `externalTable`/`externalTableRef`, `overTables`, `intervalBetweenRuns`, `forcedLatency`, `sizeLimit`, `distributed`, `managedIdentity`, `enabled`. |
| `adx_table_retention_policy` / `adx_materialized_view_retention_policy` | `RetentionPolicy` | `entity.kind: Table` or `MaterializedView`; database level via `entity.kind: Database`. |
| `adx_table_caching_policy` / `adx_materialized_view_caching_policy` | `CachingPolicy` | `data_hot_span`/`index_hot_span` -> `hot` (Kusto sets both); `hotWindows[]` optional. |
| `adx_table_update_policy` | `UpdatePolicy` | one resource holds the complete list `updates[]` of a table (authoritative); `source` may be a `sourceRef`. |
| `adx_table_row_level_security_policy` / `adx_materialized_view_row_level_security_policy` | `RowLevelSecurityPolicy` | `enabled`, `query`. |
| `adx_table_ingestion_batching_policy` | `IngestionBatchingPolicy` | table or database level. |
| `adx_table_ingestion_time_policy` | `IngestionTimePolicy` | `enabled`. |
| `adx_table_streaming_ingestion_policy` | `StreamingIngestionPolicy` | `enabled`, `hintAllocatedRate` (string decimal). |
| `adx_table_partitioning_policy` | `PartitioningPolicy` | `partitionKeys[]` with `kind: Hash \| UniformRange`, `effectiveDateTime`. |
| `adx_table_restricted_view_access_policy` | `RestrictedViewAccessPolicy` | `enabled`. |
| `adx_merge_policy` | `MergePolicy` | table, materialized view or database level. |
| `adx_column_encoding_policy` | `EncodingPolicy` | `entity.kind: Table` plus `column`; database level without `column`. |
| `adx_table_security_role` | `SecurityRole` | `entity.kind: Table` (also Database, ExternalTable, MaterializedView, Function); `mode: Authoritative` matches the Terraform behaviour of owning the whole role, `Additive` only manages the listed principals. Principals are matched to object ids after the first write (see `status.atProvider.resolvedPrincipals`). |
| `adx_workload_group` | `WorkloadGroup` (group `cluster.adx.functional.team`) | the workload group JSON goes into `workloadGroup` as an object. Needs AllDatabasesAdmin. |
| `adx_cluster_request_classification_policy` | `RequestClassificationPolicy` | `enabled`, `query`. Needs AllDatabasesAdmin. |

Kinds without a Terraform counterpart: `EntityGroup`, `ShardingPolicy`,
`AutoDeletePolicy`, `ExtentTagsRetentionPolicy`, `ManagedIdentityPolicy`,
`RowOrderPolicy`, `QueryAccelerationPolicy`, and the cluster policies
`CalloutPolicy`, `CapacityPolicy`, `SandboxPolicy`, `QueryWeakConsistencyPolicy`,
`ClusterManagedIdentityPolicy`, `MultiDatabaseAdminsPolicy`.

## Migration steps

1. Create a `ProviderConfig` per cluster with the same principal Terraform used
   (or better, workload identity). Grant Database Admin on the databases.
2. For every Terraform resource create the managed resource with
   `crossplane.io/external-name` set to the Kusto name and
   `spec.managementPolicies: ["Observe"]`. Wait for `Ready=True` and compare
   `status.atProvider` with the Terraform state.
3. Fill the spec until an observe with `managementPolicies: ["*"]` would be a
   no-op: for tables that means listing every column (or accepting
   `driftColumns` in Merge mode), for functions the exact body and parameters,
   for policies only the fields you care about.
4. Switch `managementPolicies` to `["*"]`. The provider now repairs drift.
5. Remove the resources from Terraform state (`terraform state rm`) instead of
   destroying them; a `terraform destroy` would drop tables with their data.

Timespans: Terraform used strings like `"365d"` and `"P365D"` depending on the
resource; the kinds here accept Kusto literals, colon and .NET forms and
compare by value, so `365d` and `365.00:00:00` are the same.
