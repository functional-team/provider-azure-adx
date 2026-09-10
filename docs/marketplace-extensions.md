# Marketplace extension assets

`extensions/` holds the assets the Upbound Marketplace renders next to the
package. They are *not* part of the xpkg the build produces; the release
workflow appends them to the pushed package version with

    up alpha xpkg append --extensions-root=./extensions xpkg.upbound.io/functional-team/provider-azure-adx:<version>

That command is an alpha feature of the `up` CLI (>= v0.39.0) and only applies
to the Upbound registry. The GHCR copy of the package carries no extensions.

| Path | Rendered as |
|---|---|
| `extensions/icons/icon.svg` | the package icon on the Marketplace listing page |

The layout is the one
[Upbound documents](https://docs.upbound.io/manuals/marketplace/packages/#add-documentation-icons-and-other-assets-to-your-package),
which also allows `readme/`, `docs/`, `release-notes/` and `sbom/`. Nothing
else lives under `extensions/`, this file included: the documented tree has no
file at its root, and `xpkg append` is alpha, so the directory holds only what
the tool expects to find.

Assets appear on the listing page of the package version. Until the repository
is published to the Marketplace there is no listing page, and nothing renders
anywhere -- a version in state ACCEPTED is, in Upbound's words, "available for
publishing to the Marketplace, but not yet visible to others".

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
