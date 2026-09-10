# Destructive changes and guardrails

Crossplane applies the desired state immediately. This page lists what the
provider does for every change that could lose data or access, and how you see
it.

| Change | Behaviour | Signal |
|---|---|---|
| Column added to `Table.columns` | `.alter-merge table` (Merge) or `.alter table` (Replace) | `Synced=True` |
| Column removed from the spec, `schemaUpdateMode: Merge` (default) | Ignored. The column and its data stay in the cluster. | `status.atProvider.driftColumns` lists it, `Synced=True` |
| Column removed from the spec, `schemaUpdateMode: Replace` | `.alter table` drops the column **and its data**. | none beyond the event; this is the opt-in |
| Column order changed | Ignored in Merge, applied in Replace | |
| Column type changed | **Never applied.** No command is sent. | `Synced=False`, message starts with `UnsupportedColumnTypeChange` |
| Table/Function/MaterializedView renamed (annotation or `spec.forProvider.name`) | Rejected by validation once set. If the annotation is edited by hand the old entity is orphaned and a new one created. | validation error |
| `MaterializedView.sourceTable`, `backfill`, `effectiveDateTime` changed | Rejected by validation (create-only). | validation error |
| `MaterializedView.query` changed | `.alter materialized-view`. Kusto rejects changes to group-by expressions or column types. | on rejection `Synced=False`, message starts with `MaterializedViewAlterRejected`; nothing is dropped |
| Managed resource deleted (`Table`, `MaterializedView`, `ExternalTable`, ...) | Crossplane default: the entity is dropped, tables and views with all data. | set `spec.deletionPolicy: Orphan` or exclude `Delete` from `managementPolicies` to keep it |
| Policy managed resource deleted | `.delete ... policy`: the entity falls back to the inherited policy. No data is touched. | |
| `SecurityRole` in `Authoritative` mode | Principals not in the spec are removed on every reconcile, including manually added ones. Deleting the resource runs `.set ... none`. | choose `mode: Additive` for shared roles |
| `SecurityRole` in `Additive` mode | Only the spec's principals are added; deleting the resource drops only those. | |
| External table connection string with a secret | Sent obfuscated, never shown. A rotated Secret triggers an update. | `Synced=False` briefly during the update |
| Kusto throttles the cluster (HTTP 429) | Backoff and retry, counted in `adx_commands_throttled_total`. | raise `--observe-cache-ttl`/`--poll`, lower `--kusto-commands-per-second` |

There is deliberately no plan/apply mode. For a preview, create the resource with
`spec.managementPolicies: ["Observe"]` and read `status.atProvider`, then widen
the policies.

## What the callout policy can and cannot promise

A cluster answers `.show cluster policy callout` with its 19 immutable built-in
rules alongside the ones a `CalloutPolicy` manages, and marks none of them --
every entry carries only `CalloutType`, `CalloutUriRegex` and `CanCall`
(verified against a real cluster on 2026-09-10). Built-in and managed rules are
therefore indistinguishable in the response.

That makes one guarantee impossible and leaves the other intact:

- **Not possible:** "no rules exist besides the managed ones". A rule added by
  hand is indistinguishable from a built-in and is ignored, the same way
  `SecurityRole` in `Additive` mode leaves foreign principals alone.
- **Still holds:** every managed rule must be present and unchanged. A managed
  rule that was altered or removed in the cluster is reported as drift and
  restored.

Requiring the lists to match exactly is not a stricter alternative, it simply
never worked: with the built-ins present the comparison could never be equal, so
the policy was rewritten on every poll interval, forever.
