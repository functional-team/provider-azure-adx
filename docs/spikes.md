# Spikes

Status of the open questions from `tech-implement.md` section 15. "Code" means
the answer was taken from the SDK or Kusto source/documentation and is pinned
by a unit test; "emulator"/"cluster" means it still has to be confirmed by the
integration or e2e suite. Date: 2026-09-07.

| # | Question | Status | Where it matters |
|---|---|---|---|
| S1 | Does `azkustodata` accept `http://localhost` without auth? | **Answered (code).** Yes: no token provider means no trusted-endpoint check; a token provider over http is refused by the SDK. `TestNewEmulatorNoAuth`. | Emulator tests |
| S2 | Exact `error.code` for "already exists" | **Answered (emulator, 2026-09-08).** `.create function` on an existing function returns HTTP 400, `code="BadRequest"`, `@type="Kusto.Common.Svc.Exceptions.EntityAlreadyExistsException"`, message "Entity 'X' of kind 'ExpressionFunction' already exists." `kerrors` classifies it as AlreadyExists via the type/message match. | `kerrors.Classify` |
| S3 | `.show table * policy <p>` for all Tier 1 policies | **Answered for retention (emulator, 2026-09-08).** `.show table * policy retention` works; one batch command served 50 resources. Other policy kinds still take the wildcard-then-fallback path until observed. | Batch observe |
| S4 | Does `.show database schema as json` carry Folder/DocString/CslType? | **Answered (real cluster, 2026-09-08).** Yes for tables: `Folder`, `DocString` and per-column `CslType` are all present. But the same response also disproved the fixture's entity-group shape: entities come as a bare array keyed by the group name (`"EntityGroups":{"TelemetrySources":["cluster('c').database('d')"]}`), not as an object with `Name`/`Entities`. Because one bad entity group fails the whole document, this broke **every** Table observe (`cannot unmarshal array into Go struct field Database.Databases.EntityGroups`); `EntityGroup.UnmarshalJSON` now takes both shapes. Also observed: `MaterializedViews` entries carry no `Query` and no `SourceTable` (the fixture claims both) — harmless today because the materialized-view controller reads those from `.show materialized-view`, but the fixture is misleading and should be corrected. | Table observe |
| S5 | Where is a function's `view` flag observable? | **Not observable.** `view` is write-only; documented on the field. | Function |
| S6 | Exact `Role` strings in `.show ... principals` | **Answered (real cluster, 2026-09-09).** Matching by the singular role word, case-insensitive, works: both SecurityRole examples reach `Ready=True Synced=True` in e2e run 34310415774, which requires Observe to find the principal back in `.show ... principals` and match its role — a wrong role string would make the provider re-create the assignment on every reconcile instead. e2e substitutes the service principal for the examples' placeholder principals (`test/e2e/rewrite-examples.sh`). | SecurityRole |
| S7 | Which commands run in the emulator? | **Open (emulator).** Integration suite covers Table, Function, policies; external tables/continuous exports are skipped there. | Test matrix |
| S8 | Detect two MRs on the same cluster policy | **Partially.** Owner annotation is set on Create; a second MR is not blocked, documented as user error. | Tier 2 |
| S9 | contrib convention for namespaced groups (`.m.` infix) | **Closed (decision 2026-09-07).** The provider stays under github.com/functional-team with API groups under `functional.team` (`adx.functional.team`, `policy.adx.functional.team`, ...). No crossplane-contrib transfer is planned, so the `.m.` question does not arise. | API groups |
| S10 | `.create table` on an existing table | **Answered (docs + emulator, 2026-09-08).** Returns success without changing the table. | Create race |
| S11 | Response for `.show table X cslschema` on a missing table | **Answered (emulator, 2026-09-08), open for the real service.** The emulator answers with success and no rows; a real cluster is expected to answer HTTP 400 `EntityNotFound`. All single-entity observes treat "no error, no matching row" and a NotFound error alike, so both shapes work. `kerrors` additionally classifies v1 `Exceptions` frames (HTTP 200 with an error array, see the SDK) by message text as a defensive path. | Observe paths, `kerrors.Classify` |

Additional unverified assumptions collected while implementing (all marked in
code comments):

- `.show functions` echoes `Parameters` in a text form that
  `normalize.Params` canonicalizes (spacing/type aliases). If the echo differs
  structurally, stage 2 hashes still keep functions stable.
- Caching policy `.show` output nests spans as `{"Value": "31.00:00:00"}`; the
  comparison unwraps `Value` wrappers and also accepts plain strings.
- Encoding, row order and multi-database-admins policy JSON shapes: hash-only
  comparison.
- Column layouts of `.show materialized-views`, `.show external tables`,
  `.show continuous-exports`, `.show table T ingestion mappings`,
  `.show entity_groups`, `.show workload_group`: parsed by column name,
  tolerant of extra columns.
