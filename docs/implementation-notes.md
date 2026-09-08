# Implementation notes

Deviations from `tech-implement.md` and decisions that were only possible once
the code met the real SDK and generators. Date: 2026-09-07.

## SDK

- **`azkustodata.Statement` is not an interface.** In v1.2.2 it is a type alias
  for `*kql.Builder`, and `kql.New` only accepts compile-time constants. Dynamic
  commands go through `kql.New("").AddUnsafe(text)`. This is encapsulated in
  `internal/clients/kusto/cmd.Command.Statement()`; the rest of the code never
  touches the SDK statement type.
- **Spike S1 is answered by the SDK source.** `azkustodata.NewConn` refuses a
  token provider over `http://` and skips the trusted-endpoint validation when
  no credentials are configured, so the emulator works with
  `credentials.source: None` and nothing else. A unit test
  (`TestNewEmulatorNoAuth`) pins this.
- **`Mgmt` returns `query/v1.Dataset`.** Values arrive as pointers to Go types
  (`*bool`, `*int64`, `*time.Duration`, `[]byte` for dynamic). The adapter in
  `client.go` flattens them into the provider's own `Result`/`Row` model so
  domain packages and tests never import the SDK.
- **Errors**: `errors.GetKustoError` uses type assertions and does not unwrap.
  `kerrors` uses `errors.As` instead and decodes the REST body itself
  (`error.code`, `@type`, `@message`, `@permanent`, `innererror.code`).
- **Timespans overflow `time.Duration`.** Kusto's "unlimited" retention is
  1,000,000 days, more than `time.Duration` can hold. `internal/timespan` keeps
  .NET ticks (100 ns, `int64`) and formats both the `.NET` (`d.hh:mm:ss`) and
  the KQL literal form.

## crossplane-runtime

- **Blocked guardrails are returned as errors from `Observe`.** The plan wanted
  `Synced=False` with a reason while returning `ResourceUpToDate: true`. The
  managed reconciler overwrites `Synced` with `ReconcileSuccess` after a
  successful Observe, so that is not possible. Returning a `kerrors.Blocked`
  error from Observe gives `Synced=False`, reason `ReconcileError`, message
  `<Reason>: <text>` (for example `UnsupportedColumnTypeChange: column X ...`),
  and Update is never called. The reconciler requeues with backoff, which is
  the documented "visible, not loud" behaviour.
- **Annotations after Update are not persisted by the reconciler**, only after
  Create. Kinds that record hash annotations after an update call
  `base.PersistAnnotations`, which merge-patches a copy of the object and copies
  the new `resourceVersion` back so the in-memory status set during Observe is
  not lost and the following status update does not conflict.
- **Create must not set status.** The reconciler reverts everything but
  annotations after Create. The materialized view operation id therefore lives
  in the `adx.functional.team/operation-id` annotation, mirrored into
  `status.atProvider` by Observe.
- **Initializers**: `base.Register` installs `SpecNameAsExternalName` before the
  standard `NameAsExternalName`, so `spec.forProvider.name` wins over
  `metadata.name`.

## Code generation

- **Interface assertions live in `interfaces_test.go` files.** angryjet type
  checks the package before generating the very methods the assertions
  reference. Test files are invisible to it, `go vet`/`go test` still enforce
  the assertions.
- **Hand-written `ResolveReferences` files carry `//go:build !angryjet`.** They
  pass the managed resource to crossplane-runtime, which again needs generated
  methods. `apis/generate.go` runs angryjet with `GOFLAGS=-tags=angryjet` so
  those files are excluded during generation only.
- **No interfaces in API packages.** controller-gen cannot generate deepcopy
  for interface types; the `Policy` interface used by the generic policy
  controller lives in `internal/controller/policy`.
- **Policy kinds were emitted from a table** (17 kinds, identical shape). The
  generated files are ordinary source and are maintained by hand from here on.

## Kusto behaviour built into the code

- `.create table` is idempotent on an existing table (spike S10): the provider
  still tolerates an `AlreadyExists` classification on every Create.
- Batch observe for table policies uses `.show table * policy <name>`. If a
  cluster rejects the wildcard for a policy type (spike S3), the first
  `Permanent` error flips a per-kind flag and the kind falls back to single
  `.show` commands for the rest of the process lifetime.
- Policies whose `.show` JSON shape could not be verified (`encoding`,
  `roworder`, `multidatabaseadmins`) are compared through the stage 2 hash only
  (`Def.HashOnly`). They converge after the first write and detect drift, but
  cannot report field level diffs.
- Function bodies: the A6 example (`{ StormEvents | take myLimit}` vs
  `{ StormEvents | take myLimit }`) is already equal after stage 1, because the
  difference is whitespace next to the braces. Interior whitespace is never
  collapsed, so stage 2 stays mandatory.

## Hosting and API groups

- Decided on 2026-09-07 (evening): the project stays with functional.team.
  Module path `github.com/functional-team/provider-azure-adx`, packages at
  `ghcr.io/functional-team/provider-azure-adx`, API groups `adx.functional.team`,
  `policy.adx.functional.team`, `security.adx.functional.team`,
  `cluster.adx.functional.team`. Provider annotations use the same prefix
  (`adx.functional.team/applied-hash`, ...). No crossplane-contrib transfer is
  planned; a later move would be a group rename and therefore a breaking change.

## Operational

- Provider flags added beyond the template: `--poll` (default 10m),
  `--observe-cache-ttl` (60s), `--disable-observe-cache`,
  `--kusto-commands-per-second` (5), `--kusto-max-inflight` (4),
  `--kusto-client-idle-timeout` (30m).
- Metrics: `adx_commands_total{kind,op,class}`,
  `adx_command_duration_seconds{kind,op}`, `adx_commands_throttled_total{kind}`,
  `adx_observe_cache_requests_total{section,result}`.
- Commands are logged at debug level in redacted form only; `h@'...'` and
  `h"..."` literals are replaced by placeholders.
