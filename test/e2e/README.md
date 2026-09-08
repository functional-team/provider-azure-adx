# End-to-end tests

The e2e suite applies every example under `examples/` (except the provider
configuration) to a kind cluster running Crossplane and this provider, and
asserts with [uptest](https://github.com/crossplane/uptest) that each managed
resource becomes `Ready=True` and `Synced=True`, survives an import round-trip
and is deleted cleanly. It runs against a **real** Azure Data Explorer cluster
because the emulator has no authentication, no Azure Storage and no
materialized-view backfill worth the name.

## Azure resources

| Resource | Purpose |
|---|---|
| Dev cluster (smallest SKU, `Dev(No SLA)_Standard_E2a_v4`) | target of all commands; started before and stopped after the run (`az kusto cluster start/stop`) |
| Database `Telemetry` (or set `ADX_E2E_DATABASE`) | all examples use it |
| Service principal | `AllDatabasesAdmin` on the cluster (Tier 2 kinds, managed-identity external tables), federated credential for GitHub OIDC (`repo:functional-team/provider-azure-adx:ref:refs/heads/main`) |
| Storage account with container `exports` | external tables and continuous exports; the cluster's system managed identity and the SPN need `Storage Blob Data Contributor` |
| `ManagedIdentityPolicy` on the database allowing `NativeIngestion, ExternalTable, AutomatedFlows` for `system` | managed-identity connection strings and exports |

## GitHub configuration

Secrets: `AZURE_CLIENT_ID`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID`,
`ADX_E2E_CLIENT_SECRET` (the SPN secret used by the ProviderConfig inside kind;
the workflow itself authenticates with OIDC).

Variables: `ADX_E2E_CLUSTER_NAME`, `ADX_E2E_RESOURCE_GROUP`,
`ADX_E2E_CLUSTER_URI`, `ADX_E2E_DATABASE`.

## Running locally

```sh
kind create cluster
helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --version 2.0.2
make build            # builds the xpkg; load the image into kind and apply a Provider manifest
export ADX_E2E_CLUSTER_URI=https://mycluster.westeurope.kusto.windows.net
export ADX_E2E_DATABASE=Telemetry
export ADX_E2E_CLIENT_ID=... ADX_E2E_CLIENT_SECRET=... ADX_E2E_TENANT_ID=...
./test/e2e/setup.sh
make e2e
```

Mandatory scenarios covered by the examples (docs/tech-implement.md 11.3):
`SecurityRole` Authoritative and Additive with principal resolution,
`ExternalTable` with managed identity and with a SAS Secret,
`ContinuousExport`, `MaterializedView` with backfill, and 500 policies at once
for throttling behaviour (generate them with a loop over
`examples/policy/retentionpolicy.yaml` when needed).
