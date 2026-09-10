#!/usr/bin/env bash
# Two updates in one patch: a plain structural field (the table docstring,
# echoed back unchanged, unlike a function body) and a new column that carries
# a docstring of its own.
#
# The second one is here because of issue #1. Adding a column to a table whose
# other columns already have docstrings used to strip theirs: the update sends
# ".alter table T column-docstrings (...)" with only the columns it considers
# changed, and that verb removes the docstring of every column it does not
# name. The next Observe saw the loss as drift, wrote the old set back, dropped
# the new one again, and the table flipped between two states on every poll --
# reporting Synced=True the whole time, because each write did succeed.
#
# The example table therefore carries docstrings on two columns before this
# patch adds a third. What is asserted is not that the new docstring arrives
# (it did even with the bug) but that all three are present at the same time.
set -euo pipefail
. "$(dirname "${BASH_SOURCE[0]}")/update-check.sh"

KUBECTL="${KUBECTL:-kubectl}"
RESOURCE="table.adx.functional.team/raw-events"
NS="data-platform"

read -r -d '' PATCH <<'JSON' || true
{"spec":{"forProvider":{
  "docstring":"updated by the e2e update check",
  "columns":[
    {"name":"Timestamp","type":"datetime"},
    {"name":"DeviceId","type":"string","docstring":"Device that produced the event"},
    {"name":"Payload","type":"dynamic","docstring":"Raw JSON body"},
    {"name":"IsProcessed","type":"bool","docstring":"Added by the e2e update check"}
  ]}}}
JSON

check_update "$RESOURCE" "$NS" "$PATCH" \
  '{.status.atProvider.docstring}' \
  "updated by the e2e update check"

echo "waiting for all three column docstrings to be set at the same time"
observed=""
for _ in $(seq 1 24); do
  observed=$("$KUBECTL" -n "$NS" get "$RESOURCE" -o 'jsonpath={.status.atProvider.columns[*].docstring}')
  if [[ "$observed" == *"Device that produced the event"* &&
        "$observed" == *"Raw JSON body"* &&
        "$observed" == *"Added by the e2e update check"* ]]; then
    echo "all column docstrings survived the update"
    exit 0
  fi
  sleep 5
done

echo "FAIL: the column docstrings never held all three values at once." >&2
echo "That is issue #1: the update replaces the whole docstring set instead of" >&2
echo "merging into it, so the table oscillates between two states." >&2
echo "last read: $observed" >&2
"$KUBECTL" -n "$NS" get "$RESOURCE" -o yaml >&2 || true
exit 1
