# Spikes

Status of the open questions from `tech-implement.md` section 15. "Code" means
the answer was taken from the SDK or Kusto source/documentation and is pinned
by a unit test; "emulator"/"cluster" means it still has to be confirmed by the
integration or e2e suite. Date: 2026-09-07.

| # | Question | Status | Where it matters |
|---|---|---|---|
| S1 | Does `azkustodata` accept `http://localhost` without auth? | **Answered (code).** Yes: no token provider means no trusted-endpoint check; a token provider over http is refused by the SDK. `TestNewEmulatorNoAuth`. | Emulator tests |
| S2 | Exact `error.code` for "already exists" | **Open (emulator).** `kerrors` matches `*AlreadyExists*`, `*Conflict*` codes/types and "already exists" messages; the integration test records the raw code. `.create table` is idempotent anyway (S10). | `kerrors.Classify` |
| S3 | `.show table * policy <p>` for all Tier 1 policies | **Open (emulator), mitigated.** Only `update` is documented. The controller tries the wildcard once per kind and falls back to single `.show` on a permanent error. Integration test logs which path was taken. | Batch observe |
| S4 | Does `.show database schema as json` carry Folder/DocString/CslType? | **Assumed yes (schema fixture).** Parser is tolerant: missing `CslType` falls back to `Type` (.NET name mapping), missing Folder/DocString are empty strings. Confirm with the emulator. | Table observe |
| S5 | Where is a function's `view` flag observable? | **Not observable.** `view` is write-only; documented on the field. | Function |
| S6 | Exact `Role` strings in `.show ... principals` | **Open (cluster).** Matching is by the singular role word, case-insensitive. Only testable against a real cluster (emulator has no auth). | SecurityRole |
| S7 | Which commands run in the emulator? | **Open (emulator).** Integration suite covers Table, Function, policies; external tables/continuous exports are skipped there. | Test matrix |
| S8 | Detect two MRs on the same cluster policy | **Partially.** Owner annotation is set on Create; a second MR is not blocked, documented as user error. | Tier 2 |
| S9 | contrib convention for namespaced groups (`.m.` infix) | **Closed (decision 2026-09-07).** The provider stays under github.com/functional-team with API groups under `functional.team` (`adx.functional.team`, `policy.adx.functional.team`, ...). No crossplane-contrib transfer is planned, so the `.m.` question does not arise. | API groups |
| S10 | `.create table` on an existing table | **Answered (docs).** Returns success without changing the table. | Create race |

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
