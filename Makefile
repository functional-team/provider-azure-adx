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

# Images are published to ghcr.io under the personal account until the project
# moves to crossplane-contrib (see docs/tech-implement.md, section 5.1).
XPKG_REG_ORGS ?= ghcr.io/functional-team
XPKG_REG_ORGS_NO_PROMOTE ?= ghcr.io/functional-team
XPKGS = provider-azure-adx
-include build/makelib/xpkg.mk

# NOTE(hasheddan): we force image building to happen prior to xpkg build so that
# we ensure image is present in daemon.
xpkg.build.provider-azure-adx: do.build.images

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
	-not -name 'clustermanagedidentitypolicy.yaml' \
	-not -name 'queryaccelerationpolicy.yaml' \
	-not -name 'function.yaml' \
	| sort | tr '\n' ',' | sed 's/,$$//')
UPTEST_SETUP_SCRIPT ?= test/e2e/setup.sh
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
