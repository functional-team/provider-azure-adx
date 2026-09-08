#!/usr/bin/env bash
# Points the examples at the dev cluster's database. The examples are written
# for a database named "Telemetry"; a dev cluster usually has another one.
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

if [ "${ADX_E2E_DATABASE}" = "Telemetry" ]; then
  echo "examples already use database Telemetry, nothing to rewrite"
  exit 0
fi

# -i.bak plus a delete pass keeps this working on both GNU and BSD sed.
find "${REPO_ROOT}/examples" -name '*.yaml' -not -path "${REPO_ROOT}/examples/provider/*" \
  -exec sed -i.bak "s/database: Telemetry/database: ${ADX_E2E_DATABASE}/g" {} +
find "${REPO_ROOT}/examples" -name '*.yaml.bak' -delete

echo "examples rewritten to database ${ADX_E2E_DATABASE}"
grep -rl "database: ${ADX_E2E_DATABASE}" "${REPO_ROOT}/examples" | sort
