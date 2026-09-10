# Security policy

## Reporting a vulnerability

Please do **not** open a public issue for a security problem.

Use GitHub's private vulnerability reporting on this repository
(*Security* → *Report a vulnerability*), or write to the maintainer listed in
[OWNERS.md](OWNERS.md). Include what you did, what happened and what you
expected; a proof of concept helps but is not required.

This is a single-maintainer project, so expect an acknowledgement within a few
working days rather than within hours. Once a fix is available it ships in a
patch release, and the advisory credits the reporter unless you ask otherwise.

## Supported versions

The project is pre-`v0.1.0`. Only the latest release receives fixes.

## Scope

This provider executes Kusto management commands with the credentials of its
`ProviderConfig` principal. Two properties are part of the security contract
and a break in either is a vulnerability:

- **Secrets never leave the provider.** External table connection strings that
  carry a SAS or an account key are sent obfuscated and must not appear in
  `status`, events, logs or metrics. Rotation is tracked through a hash
  annotation.
- **KQL is executed as written.** Function bodies, update policy queries,
  row level security queries, materialized view queries and continuous export
  queries are sent verbatim, by design. Anyone allowed to create those managed
  resources runs arbitrary KQL as the provider principal, so Kubernetes RBAC on
  those kinds is the control. That is documented behaviour, not a
  vulnerability; a way to inject KQL through a field that is *not* one of those
  is.

Reports about the Azure Data Explorer service itself belong to
[Microsoft](https://msrc.microsoft.com/report).
