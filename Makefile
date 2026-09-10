# ====================================================================================
# Setup Project
PROJECT_NAME := provider-azure-adx
PROJECT_REPO := github.com/functional-team/$(PROJECT_NAME)

PLATFORMS ?= linux_amd64 linux_arm64
-include build/makelib/common.mk

# ====================================================================================
# Setup Output

-include build/makelib/output.mk

# ====================================================================================
# Setup Go

NPROCS ?= 1
GO_TEST_PARALLEL := $(shell echo $$(( $(NPROCS) / 2 )))
GO_STATIC_PACKAGES = $(GO_PROJECT)/cmd/provider
GO_LDFLAGS += -X $(GO_PROJECT)/internal/version.Version=$(VERSION)
GO_SUBDIRS += cmd internal apis
# Integration tests live under test/ and are guarded by the "integration" build tag.
GO_INTEGRATION_TEST_PACKAGES = $(GO_PROJECT)/test/integration
GO111MODULE = on
GOLANGCILINT_VERSION = 2.12.2
-include build/makelib/golang.mk

# ====================================================================================
# Setup Kubernetes tools

-include build/makelib/k8s_tools.mk

# ====================================================================================
# Setup Images

IMAGES = provider-azure-adx
-include build/makelib/imagelight.mk

# ====================================================================================
# Setup XPKG

# Images are published to ghcr.io. The release workflow additionally pushes to
# xpkg.upbound.io/functional-team when the UPBOUND_PUBLISH variable is set --
# that registry is the one the Upbound Marketplace lists, and it renders the
# API reference from the CRDs plus the annotations in package/crossplane.yaml.
# The API groups stay under functional.team either way; registry path and API
# group are unrelated.
XPKG_REG_ORGS ?= ghcr.io/functional-team
XPKG_REG_ORGS_NO_PROMOTE ?= ghcr.io/functional-team
XPKGS = provider-azure-adx
# The examples are embedded in the package (--examples-root) and are what the
# Marketplace shows as the manifest for each kind. Three of them carry uptest
# hook annotations, which are test scaffolding and point at files a user does
# not have, so the package gets a staged copy without them. Not
# XPKG_CLEANUP_EXAMPLES_ENABLED: that strips the same annotations by
# round-tripping every file through a YAML parser, losing all comments and the
# key order. uptest keeps reading examples/ directly, so what the e2e suite
# tests stays the manifests we ship.
# Per platform, because build.all runs the platforms in parallel (make -j2) and
# each sub-make stages into this directory before calling xpkg build. Sharing
# one directory means one platform wipes it while the other is reading it:
# v0.1.0-rc.5 died on "failed to parse examples: yaml: line 1: did not find
# expected node content" for linux_arm64 while linux_amd64 built fine, and
# rc.4 had the same race and got away with it.
XPKG_EXAMPLES_DIR = $(OUTPUT_DIR)/examples/$(PLATFORM)
-include build/makelib/xpkg.mk

xpkg.examples:
	@$(INFO) staging examples for the package
	@rm -rf $(XPKG_EXAMPLES_DIR)
	@cd $(ROOT_DIR)/examples && find . -name '*.yaml' | while read -r f; do \
		mkdir -p "$(XPKG_EXAMPLES_DIR)/$$(dirname "$$f")"; \
		awk -f $(ROOT_DIR)/hack/strip-uptest-annotations.awk "$$f" > "$(XPKG_EXAMPLES_DIR)/$$f" || exit 1; \
	done || $(FAIL)
	@$(OK) staging examples for the package

.PHONY: xpkg.examples

# NOTE(hasheddan): we force image building to happen prior to xpkg build so that
# we ensure image is present in daemon.
xpkg.build.provider-azure-adx: do.build.images xpkg.examples

fallthrough: submodules
	@echo Initial setup complete. Running make again . . .
	@make

# ====================================================================================
# Tests

# Integration tests run against the Kusto emulator (testcontainers) or, when
# ADX_TEST_CLUSTER_URI is set, against a real development cluster.
test-integration:
	@$(INFO) go test integration (tag: integration)
	@CGO_ENABLED=0 $(GOHOST) test -tags integration -count=1 -timeout 30m ./test/integration/... || $(FAIL)
	@$(OK) go test integration

# ====================================================================================
# End-to-end tests (uptest against a real cluster)
#
# `make e2e` builds the provider, creates a kind cluster with Crossplane
# (controlplane.mk), deploys the freshly built package into it
# (local.xpkg.mk) and runs uptest over the examples. The ProviderConfig and
# credentials are created by test/e2e/setup.sh from ADX_E2E_* variables, see
# test/e2e/README.md.
UPTEST_LOCAL_DEPLOY_TARGET = local.xpkg.deploy.provider.$(PROJECT_NAME)
# calloutpolicy is excluded because the dev cluster does not store user rules
# at all: ".alter cluster policy callout" succeeds, and ".show" keeps
# returning exactly the 21 built-in rules -- including for the example rule
# taken verbatim from Kusto's own documentation. Whether that is a provider
# defect, a command that needs .alter-merge, or a cluster that refuses user
# callouts is open (docs/spikes.md, S12); until it is settled the kind cannot
# be verified here, and hiding the drift behind a hash comparison would only
# make the provider claim success it has not achieved.
# clustermanagedidentitypolicy names a second identity for DataConnection that
# is not attached to the dev cluster, and queryaccelerationpolicy applies to an
# external table named ExportsDelta that no example creates -- query
# acceleration only works on Delta Lake external tables, so it needs delta data
# in storage, not just a container.
# function.yaml's parameters value "(limit:long = 100)" collides with
# chainsaw/uptest's own templating: any string starting with "(" and ending
# with ")" is parsed as an embedded JMESPath expression, and "limit:long =
# 100" isn't valid JMESPath. No escape/opt-out exists upstream (checked
# chainsaw's templating docs). The Function kind is still exercised, by the
# parameterless functions the update-policy and row-level-security examples
# ship, and by the emulator integration suite.
UPTEST_INPUT_MANIFESTS ?= $(shell find examples -name '*.yaml' \
	-not -path 'examples/provider/*' \
	-not -path 'examples/composition/*' \
	-not -name 'calloutpolicy.yaml' \
	-not -name 'clustermanagedidentitypolicy.yaml' \
	-not -name 'queryaccelerationpolicy.yaml' \
	-not -name 'function.yaml' \
	| sort | tr '\n' ',' | sed 's/,$$//')
UPTEST_SETUP_SCRIPT ?= test/e2e/setup.sh
# Runs the provider with --poll=1m so slow drift becomes observable within a
# run; see the file for why.
DRC_FILE ?= test/e2e/runtimeconfig.yaml
UPTEST_DEFAULT_TIMEOUT ?= 1800s
CROSSPLANE_VERSION ?= 2.0.2
-include build/makelib/local.xpkg.mk
-include build/makelib/controlplane.mk
-include build/makelib/uptest.mk

# Update the submodules, such as the common build scripts.
submodules:
	@git submodule sync
	@git submodule update --init --recursive

go.cachedir:
	@go env GOCACHE

go.mod.cachedir:
	@go env GOMODCACHE

# NOTE(hasheddan): we must ensure up is installed in tool cache prior to build
# as including the k8s_tools machinery prior to the xpkg machinery sets UP to
# point to tool cache.
build.init: $(CROSSPLANE_CLI)

# This is for running out-of-cluster locally, and is for convenience. Running
# this make target will print out the command which was used. For more control,
# try running the binary directly with different arguments.
run: go.build
	@$(INFO) Running Crossplane locally out-of-cluster . . .
	@# To see other arguments that can be provided, run the command with --help instead
	$(GO_OUT_DIR)/provider --debug

dev: $(KIND) $(KUBECTL)
	@$(INFO) Creating kind cluster
	@$(KIND) create cluster --name=$(PROJECT_NAME)-dev
	@$(KUBECTL) cluster-info --context kind-$(PROJECT_NAME)-dev
	@$(INFO) Installing provider-azure-adx CRDs
	@$(KUBECTL) apply -R -f package/crds
	@$(INFO) Starting provider-azure-adx controllers
	@$(GO) run cmd/provider/main.go --debug

dev-clean: $(KIND) $(KUBECTL)
	@$(INFO) Deleting kind cluster
	@$(KIND) delete cluster --name=$(PROJECT_NAME)-dev

# Generate the CRD reference documentation under docs/api.
crddoc:
	@$(INFO) generating CRD reference docs
	@$(GO) run fybrik.io/crdoc@v0.6.4 --resources package/crds --output docs/api/README.md || $(FAIL)
	@$(OK) generating CRD reference docs

.PHONY: submodules fallthrough test-integration run dev dev-clean crddoc

define CROSSPLANE_MAKE_HELP
Crossplane Targets:
    submodules            Update the submodules, such as the common build scripts.
    run                   Run crossplane locally, out-of-cluster. Useful for development.
    test-integration      Run integration tests against the Kusto emulator.
    e2e                   Run uptest end-to-end tests against a real cluster.
    crddoc                Regenerate docs/api from the CRDs.

endef
# The reason CROSSPLANE_MAKE_HELP is used instead of CROSSPLANE_HELP is because the crossplane
# binary will try to use CROSSPLANE_HELP if it is set, and this is for something different.
export CROSSPLANE_MAKE_HELP

crossplane.help:
	@echo "$$CROSSPLANE_MAKE_HELP"

help-special: crossplane.help

.PHONY: crossplane.help help-special
