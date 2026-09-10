#!/usr/bin/env bash
# The same check for a plain structural field, so the update path is covered
# for more than the KQL case: the docstring is echoed back unchanged, unlike a
# function body.
set -euo pipefail
. "$(dirname "${BASH_SOURCE[0]}")/update-check.sh"

check_update \
  "table.adx.functional.team/raw-events" \
  "data-platform" \
  '{"spec":{"forProvider":{"docstring":"updated by the e2e update check"}}}' \
  '{.status.atProvider.docstring}' \
  "updated by the e2e update check"
