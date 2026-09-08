#!/usr/bin/env bash
# uptest setup script: creates the namespace, the credentials Secret and the
# ProviderConfig the examples reference. Expects a kubeconfig with Crossplane
# and provider-azure-adx installed and these variables:
#   ADX_E2E_CLUSTER_URI   https://<cluster>.<region>.kusto.windows.net
#   ADX_E2E_DATABASE      database the examples use (default Telemetry)
#   ADX_E2E_CLIENT_ID / ADX_E2E_CLIENT_SECRET / ADX_E2E_TENANT_ID  service principal
set -euo pipefail

: "${ADX_E2E_CLUSTER_URI:?ADX_E2E_CLUSTER_URI is required}"
: "${ADX_E2E_CLIENT_ID:?ADX_E2E_CLIENT_ID is required}"
: "${ADX_E2E_CLIENT_SECRET:?ADX_E2E_CLIENT_SECRET is required}"
: "${ADX_E2E_TENANT_ID:?ADX_E2E_TENANT_ID is required}"
ADX_E2E_DATABASE="${ADX_E2E_DATABASE:-Telemetry}"

kubectl create namespace data-platform --dry-run=client -o yaml | kubectl apply -f -

kubectl -n data-platform create secret generic adx-sp \
  --from-literal=credentials="{\"clientId\":\"${ADX_E2E_CLIENT_ID}\",\"clientSecret\":\"${ADX_E2E_CLIENT_SECRET}\",\"tenantId\":\"${ADX_E2E_TENANT_ID}\"}" \
  --dry-run=client -o yaml | kubectl apply -f -

cat <<YAML | kubectl apply -f -
apiVersion: adx.functional.team/v1alpha1
kind: ProviderConfig
metadata:
  name: telemetry-prod
  namespace: data-platform
spec:
  clusterUri: ${ADX_E2E_CLUSTER_URI}
  credentials:
    source: Secret
    secretRef:
      namespace: data-platform
      name: adx-sp
      key: credentials
YAML

echo "e2e setup done: ProviderConfig telemetry-prod -> ${ADX_E2E_CLUSTER_URI} (database ${ADX_E2E_DATABASE})"
