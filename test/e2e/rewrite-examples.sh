#!/usr/bin/env bash
# Points the examples at the dev environment: the database name, and the AAD
# principals in the SecurityRole examples. The examples are written for a
# database named "Telemetry" and for placeholder principals (alice@contoso.com
# and an app id in Microsoft's own tenant) -- correct as documentation, but no
# cluster can resolve them, so e2e substitutes the e2e service principal.
#
# This deliberately does not live in setup.sh: uptest snapshots the example
# manifests into its scratch directory before it runs the setup script, so a
# rewrite from there never reaches the manifests that actually get applied. In
# e2e run 34257424955 every database-scoped resource came up with database
# "Telemetry" and failed with BadRequest_EntityNotFound, while the
# cluster-scoped policies (which name no database) were Ready and Synced. Run
# this before `make e2e`.
set -euo pipefail

ADX_E2E_DATABASE="${ADX_E2E_DATABASE:-Telemetry}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# -i.bak plus a delete pass keeps this working on both GNU and BSD sed.
edit() {
  find "${REPO_ROOT}/examples" -name '*.yaml' -not -path "${REPO_ROOT}/examples/provider/*" \
    -exec sed -i.bak "$1" {} +
  find "${REPO_ROOT}/examples" -name '*.yaml.bak' -delete
}

if [ "${ADX_E2E_DATABASE}" != "Telemetry" ]; then
  edit "s/database: Telemetry/database: ${ADX_E2E_DATABASE}/g"
  echo "database -> ${ADX_E2E_DATABASE}"
fi

# The SecurityRole examples name principals that exist in no real tenant. Point
# them at the e2e service principal instead; the group has no counterpart, so
# its list entry goes away entirely.
if [ -n "${ADX_E2E_CLIENT_ID:-}" ] && [ -n "${ADX_E2E_TENANT_ID:-}" ]; then
  sp="aadapp=${ADX_E2E_CLIENT_ID};${ADX_E2E_TENANT_ID}"
  edit "s|aadapp=4c7e82bd-6adb-46c3-b413-fdd44834c69b;72f988bf-86f1-41af-91ab-2d7cd011db47|${sp}|g"
  edit "s|aaduser=alice@contoso.com|${sp}|g"
  edit "/aadgroup=data-ingest;contoso.com/d"
  echo "SecurityRole principals -> the e2e service principal"
else
  echo "ADX_E2E_CLIENT_ID/ADX_E2E_TENANT_ID unset, leaving principals as they are" >&2
fi
