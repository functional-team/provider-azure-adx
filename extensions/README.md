# Marketplace extension assets

Assets the Upbound Marketplace renders next to the package. They are *not* part
of the xpkg the build produces; the release workflow appends them to the pushed
package version with

    up alpha xpkg append --extensions-root=./extensions xpkg.upbound.io/functional-team/provider-azure-adx:<version>

That command is an alpha feature of the `up` CLI (>= v0.39.0) and only applies
to the Upbound registry. The GHCR copy of the package carries no extensions.

| Path | Rendered as |
|---|---|
| `icons/icon.svg` | the package icon in the repository view and the Marketplace listing |

The overview text is *not* here: it lives in the `meta.crossplane.io/readme`
annotation in `package/crossplane.yaml`, so that it travels with the package on
every registry rather than only on xpkg.upbound.io.

## The icon

Drawn for this project (a data store plus a lens, in blue). Deliberately not
Microsoft's Azure Data Explorer icon: Microsoft permits its product icons in
"architectural diagrams, training materials, or documentation" and states
"Don't use Microsoft product icons to represent your product or service"
(https://learn.microsoft.com/azure/architecture/icons/). A marketplace avatar
for a third-party provider is exactly that use.
