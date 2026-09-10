# End-to-end tests

The e2e suite applies every example under `examples/` (except the provider
configuration) to a kind cluster running Crossplane and this provider, and
asserts with [uptest](https://github.com/crossplane/uptest) that each managed
resource becomes `Ready=True` and `Synced=True`, survives an import round-trip
and is deleted cleanly. It runs against a **real** Azure Data Explorer cluster
because the emulator has no authentication, no Azure Storage and no
materialized-view backfill worth the name.

Three examples are excluded in `Makefile` (`UPTEST_INPUT_MANIFESTS`):
`clustermanagedidentitypolicy.yaml` names a second identity for
`DataConnection` that is not attached to the dev cluster;
`queryaccelerationpolicy.yaml` applies to an external table `ExportsDelta` that
no example creates, and query acceleration only works on Delta Lake external
tables, so it needs delta data in storage rather than just a container; and
`function.yaml`'s `parameters` value (`"(limit:long = 100)"`) collides with
chainsaw/uptest's own templating, which parses any string starting with `(` and
ending with `)` as a JMESPath expression. The Function *kind* is still
exercised, by the parameterless functions the update-policy and
row-level-security examples ship.

## Azure resources

| Resource | Purpose |
|---|---|
| Dev cluster (smallest SKU, `Dev(No SLA)_Standard_E2a_v4`) | target of all commands; started before and stopped after the run (`az kusto cluster start/stop`) |
| Database `Telemetry` (or set `ADX_E2E_DATABASE`) | all examples use it |
| Service principal | `AllDatabasesAdmin` on the cluster (Tier 2 kinds, managed-identity external tables), federated credential for GitHub OIDC (`repo:functional-team/provider-azure-adx:ref:refs/heads/main`) |
| Storage account with container `exports` | external tables and continuous exports; the cluster's system managed identity and the SPN need `Storage Blob Data Contributor` |
| `ManagedIdentityPolicy` on the database allowing `NativeIngestion, ExternalTable, AutomatedFlows` for `system` | managed-identity connection strings and exports |

## GitHub configuration

Secrets: `ADX_E2E_STORAGE_ACCOUNT` and `ADX_E2E_STORAGE_SAS` (the SAS query
string only, no leading `?`) point the external-table example at the real
storage account -- one connection string authenticates as the cluster's managed
identity, the other with the SAS. Plus `AZURE_CLIENT_ID`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID`,
`ADX_E2E_CLIENT_SECRET` (the SPN secret used by the ProviderConfig inside kind;
the workflow itself authenticates with OIDC), and `ADX_E2E_CLUSTER_NAME`,
`ADX_E2E_RESOURCE_GROUP`, `ADX_E2E_CLUSTER_URI`, `ADX_E2E_DATABASE` (not
sensitive by nature, but kept as secrets rather than variables so the dev
cluster's name/URI aren't visible in plaintext to anyone with access to the
Actions settings).

Variable: `ADX_E2E_ENABLED` (`true` to run the job at all). A job-level `if:`
cannot reference the `secrets` context, so this plain on/off switch is a
variable while the actual cluster config stays in secrets.

## What the run asserts, and why not with uptest's update step

Three post-assert hooks run inside the apply phase:

- `assert-table-update.sh` and `assert-kql-update.sh` change a field and
  require it to reach the cluster and then be left alone -- a structural field
  on a Table, and a function's KQL body, which Kusto stores reformatted so it
  can only be checked by behaviour.
- `assert-no-drift.sh` generalizes the second half to every kind: once
  everything has converged it watches past two poll intervals and fails on any
  `UpdatedExternalResource` event. That check is why e2e runs the provider with
  `--poll=1m` (`runtimeconfig.yaml`); at the 10m default a run finishes before
  slow drift could ever show.

uptest's own update step is not used. Its template patches the resource
without `--namespace`, while every managed resource here is namespaced, so the
patch never finds the object; its retry loop then spins forever because it
increments with the non-POSIX `((attempt++))` under `/usr/bin/sh`, and the step
dies at the timeout (run 34452246023). Its assert and delete templates do pass
the namespace, so this looks like an oversight worth reporting upstream.

## Running locally

```sh
kind create cluster
helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --version 2.0.2
make build            # builds the xpkg; load the image into kind and apply a Provider manifest
export ADX_E2E_CLUSTER_URI=https://mycluster.westeurope.kusto.windows.net
export ADX_E2E_DATABASE=Telemetry
export ADX_E2E_CLIENT_ID=... ADX_E2E_CLIENT_SECRET=... ADX_E2E_TENANT_ID=...
./test/e2e/rewrite-examples.sh   # only if ADX_E2E_DATABASE is not "Telemetry"
./test/e2e/setup.sh
make e2e
```

Mandatory scenarios covered by the examples (docs/tech-implement.md 11.3):
`SecurityRole` Authoritative and Additive with principal resolution,
`ExternalTable` with managed identity and with a SAS Secret,
`ContinuousExport`, `MaterializedView` with backfill, and 500 policies at once
for throttling behaviour (generate them with a loop over
`examples/policy/retentionpolicy.yaml` when needed).
