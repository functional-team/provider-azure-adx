#!/usr/bin/env bash
# Shared body of the update checks. Sourced by the post-assert hooks; not
# meant to be run on its own.
#
# uptest has an update step of its own, but it cannot be used here: its
# template patches the resource without --namespace, while every managed
# resource of this provider is namespaced, so the patch never finds the object
# (its assert and delete templates do pass the namespace -- the update one is
# the odd one out). Its retry loop then spins forever, because it increments
# the attempt counter with the non-POSIX "((attempt++))" while running under
# /usr/bin/sh, so the log repeats "attempt 1/10" until chainsaw kills the step
# at the timeout. Observed in e2e run 34452246023.
#
# check_update <resource> <namespace> <patch json> <jsonpath> <marker>
#
# Asserts that the patch reaches the cluster. That it is then written only
# once is not checked here but by assert-no-drift.sh, which watches every
# resource after they have all converged -- these two among them. Waiting for
# quiet in each hook as well cost 300s and pushed the apply phase past its
# budget, for an assertion already covered.
check_update() {
  local resource="$1" ns="$2" patch="$3" jsonpath="$4" marker="$5"
  local kubectl="${KUBECTL:-kubectl}"

  echo "patching $resource"
  "$kubectl" -n "$ns" patch "$resource" --type=merge -p "$patch"

  # Not "kubectl wait --for=condition=Synced": that condition is still True
  # from before the patch, so it returns instantly and proves nothing. Wait for
  # the observation itself to catch up.
  echo "waiting for the provider to apply it"
  local observed=""
  for _ in $(seq 1 60); do
    observed=$("$kubectl" -n "$ns" get "$resource" -o "jsonpath=$jsonpath")
    case "$observed" in
      *"$marker"*) break ;;
    esac
    sleep 5
  done

  case "$observed" in
    *"$marker"*) echo "the cluster reports the change" ;;
    *)
      echo "FAIL: the change never reached the cluster within 5m." >&2
      echo "expected $jsonpath to contain: $marker" >&2
      echo "read back: $observed" >&2
      "$kubectl" -n "$ns" get "$resource" -o yaml >&2 || true
      return 1
      ;;
  esac

  "$kubectl" -n "$ns" wait "$resource" --for=condition=Synced --timeout=2m
  echo "the change reached the cluster"
}
