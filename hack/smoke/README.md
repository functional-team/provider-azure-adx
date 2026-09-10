# Smoke test on a Kubernetes cluster outside Azure

Manual walk-through against a real ADX cluster with a service principal.
Works from any cluster (k3s, kind) because `credentials.source: Secret`
does not depend on Azure identity plumbing.

## Prerequisites

- Crossplane 2.x installed (`helm install crossplane crossplane-stable/crossplane
  -n crossplane-system --create-namespace`).
- An Upbound robot token for `xpkg.upbound.io/functional-team`, as long as the
  repository is private. (A GitHub PAT with `read:packages` works too if you
  point the manifest at the `ghcr.io` copy.)
- A service principal that is Database Admin on the target database:
  `.add database <DB> admins ('aadapp=<clientId>;<tenantId>')`.
- The examples use database `Telemetry`; either create it or replace the
  `database:` fields in the examples you apply.

## Steps

```sh
kubectl -n crossplane-system create secret docker-registry upbound-functional-team \
  --docker-server=xpkg.upbound.io \
  --docker-username=<robot-access-id> --docker-password=<robot-token>

kubectl apply -f hack/smoke/provider.yaml
kubectl get provider.pkg provider-azure-adx -w        # HEALTHY=True

# Fill in the REPLACE_* values first.
kubectl apply -f hack/smoke/providerconfig.yaml

kubectl apply -f examples/adx/table.yaml
kubectl apply -f examples/adx/function.yaml
kubectl -n data-platform get tables.adx.functional.team,functions.adx.functional.team -w
```

`READY=True SYNCED=True` on both means create and observe work. Then check
drift handling: change a column type in `examples/adx/table.yaml` and re-apply;
the resource must show `Synced=False` with the guardrail reason and leave the
table untouched (`docs/guardrails.md`). Delete the resources to test cleanup:

```sh
kubectl delete -f examples/adx/function.yaml -f examples/adx/table.yaml
```

Provider logs: `kubectl -n crossplane-system logs -l pkg.crossplane.io/provider=provider-azure-adx -f`.
