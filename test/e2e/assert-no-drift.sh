#!/usr/bin/env bash
# Asserts that the provider updates only real changes.
#
# Every kind is reconciled once per poll interval. A resource whose desired and
# observed state never compare equal converges once and is then rewritten on
# every pass, forever -- the ingestion mapping did exactly that (it echoed the
# mapping properties in a different shape, so the comparison always saw a
# difference). That bug only surfaced because the drift appeared immediately
# after create; the same defect one poll later would have gone unnoticed.
#
# So: wait until everything has converged, then watch for longer than a poll
# interval and require that no managed resource reports another update. The
# provider records each one as an UpdatedExternalResource event, so this covers
# every kind at once rather than a single sampled resource.
#
# Runs as an uptest post-assert hook. It waits for all managed resources
# itself, so it does not depend on which example it is attached to.
set -euo pipefail

KUBECTL="${KUBECTL:-kubectl}"
# The provider runs with --poll=1m in e2e (test/e2e/runtimeconfig.yaml). Watch
# well past two intervals so a resource that drifts once per poll cannot slip
# through between two observations.
WATCH_SECONDS="${DRIFT_WATCH_SECONDS:-150}"

echo "waiting for all managed resources to converge"
"$KUBECTL" wait managed --all --all-namespaces --for=condition=Ready --timeout=10m
"$KUBECTL" wait managed --all --all-namespaces --for=condition=Synced --timeout=2m

since=$(date -u +%Y-%m-%dT%H:%M:%SZ)
echo "converged at $since, watching ${WATCH_SECONDS}s for updates that should not happen"
sleep "$WATCH_SECONDS"

# Events carry either eventTime (new API) or lastTimestamp (legacy); take
# whichever is set and keep the ones after the convergence mark.
drifted=$("$KUBECTL" get events --all-namespaces \
  --field-selector reason=UpdatedExternalResource -o json |
  jq -r --arg since "$since" '
    .items[]
    | (.eventTime // .lastTimestamp) as $t
    | select($t != null and $t > $since)
    | "\(.involvedObject.kind)/\(.involvedObject.name) at \($t)"' | sort -u)

if [ -n "$drifted" ]; then
  echo "FAIL: these resources were updated again after they had converged," >&2
  echo "which means desired and observed never compare equal for them:" >&2
  echo "$drifted" >&2
  exit 1
fi

echo "no resource was updated after converging"
