#!/usr/bin/env bash
# Changes the KQL body of a function and asserts the change reaches the
# cluster and is then left alone.
#
# The body cannot be compared as a string: Kusto stores it wrapped in braces
# with an added blank line, so what comes back never equals what was sent.
#
#   sent:      RawEvents\n| project Timestamp, DeviceId, ...
#   read back: {\nRawEvents\n| project Timestamp, DeviceId, ...\n\n}
#
# That is why the provider normalizes KQL before comparing, and why an exact
# string assertion would be the wrong instrument. The filter below keeps the
# projected schema, which the update policy on ParsedEvents depends on.
set -euo pipefail
. "$(dirname "${BASH_SOURCE[0]}")/update-check.sh"

NEW_BODY='RawEvents
| where isnotempty(DeviceId)
| project Timestamp, DeviceId, Level = tostring(Payload.level)
'

check_update \
  "function.adx.functional.team/parse-raw-events" \
  "data-platform" \
  "$(jq -n --arg b "$NEW_BODY" '{spec:{forProvider:{body:$b}}}')" \
  '{.status.atProvider.body}' \
  "isnotempty(DeviceId)"
