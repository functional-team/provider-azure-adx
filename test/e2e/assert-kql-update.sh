#!/usr/bin/env bash
# Asserts that a changed KQL body reaches the cluster, and that the provider
# then stops writing.
#
# The body cannot be compared as a string: Kusto stores it wrapped in braces
# with an added blank line, so what comes back never equals what was sent.
#
#   sent:      RawEvents\n| project Timestamp, DeviceId, ...
#   read back: {\nRawEvents\n| project Timestamp, DeviceId, ...\n\n}
#
# That is why the provider normalizes KQL before comparing, and why uptest's
# own update step (which greps for an exact value in status.atProvider) is the
# wrong instrument here. Assert the two things that actually matter instead:
# the new code is in the cluster, and the resource does not keep drifting
# afterwards.
set -euo pipefail

KUBECTL="${KUBECTL:-kubectl}"
NS=data-platform
FN=parse-raw-events
MARKER="isnotempty(DeviceId)"
# Longer than two poll intervals (--poll=1m in e2e).
QUIET_SECONDS="${KQL_QUIET_SECONDS:-150}"

# Keeps the projected schema, which the update policy on ParsedEvents depends
# on -- only the filter is new.
NEW_BODY='RawEvents
| where isnotempty(DeviceId)
| project Timestamp, DeviceId, Level = tostring(Payload.level)
'

echo "changing the KQL body of function/$FN"
"$KUBECTL" -n "$NS" patch function.adx.functional.team/$FN --type=merge \
  -p "$(jq -n --arg b "$NEW_BODY" '{spec:{forProvider:{body:$b}}}')"

# Not "kubectl wait --for=condition=Synced": that condition is still True from
# before the patch, so it returns instantly and proves nothing. Wait for the
# observation itself to catch up.
echo "waiting for the provider to apply it"
observed=""
for _ in $(seq 1 60); do
  observed=$("$KUBECTL" -n "$NS" get function.adx.functional.team/$FN -o jsonpath='{.status.atProvider.body}')
  case "$observed" in
    *"$MARKER"*) break ;;
  esac
  sleep 5
done

case "$observed" in
  *"$MARKER"*) echo "the cluster reports the new body" ;;
  *)
    echo "FAIL: the changed KQL never reached the cluster within 5m." >&2
    echo "expected it to contain: $MARKER" >&2
    echo "read back: $observed" >&2
    "$KUBECTL" -n "$NS" get function.adx.functional.team/$FN -o yaml >&2 || true
    exit 1
    ;;
esac

# The spec change must also be reflected as reconciled, not just written.
"$KUBECTL" -n "$NS" wait function.adx.functional.team/$FN --for=condition=Synced --timeout=2m

# A body that compares unequal after being written would be rewritten on every
# poll. Watching past two intervals catches that.
since=$(date -u +%Y-%m-%dT%H:%M:%SZ)
echo "watching ${QUIET_SECONDS}s for repeated updates"
sleep "$QUIET_SECONDS"

again=$("$KUBECTL" -n "$NS" get events \
  --field-selector reason=UpdatedExternalResource,involvedObject.name=$FN -o json |
  jq -r --arg since "$since" '.items[] | (.eventTime // .lastTimestamp) as $t | select($t != null and $t > $since) | $t')

if [ -n "$again" ]; then
  echo "FAIL: function/$FN was updated again after the change had been applied," >&2
  echo "so the normalized comparison never treats the stored body as equal:" >&2
  echo "$again" >&2
  exit 1
fi

echo "the KQL change was applied once and then left alone"
