# Project configuration
PROJECT_NAME := provider-hostinger
PROJECT_REPO := github.com/rossigee/$(PROJECT_NAME)

# Build configuration
REGISTRY_ORGS = ghcr.io/rossigee
XPKG_REG_ORGS ?= ghcr.io/rossigee
CROSSPLANE_VERSION = 2.5.0
GO_REQUIRED_VERSION ?= 1.27.1
GOLANGCILINT_VERSION ?= 2.13.2

# Images configuration
IMAGES = provider-hostinger

# Crossplane package configuration
XPKGS = provider-hostinger
# Override xpkg publish to build all platforms (stock build only builds current arch)
xpkg.release.publish.ghcr.io/rossigee.provider-hostinger:
	@$(foreach plat,$(XPKG_LINUX_PLATFORMS),$(MAKE) xpkg.build.provider-hostinger PLATFORM=$(plat) || exit 1;)
	@$(CROSSPLANE_CLI) xpkg push \
		$(foreach plat,$(XPKG_LINUX_PLATFORMS),--package-files $(XPKG_OUTPUT_DIR)/$(plat)/provider-hostinger-$(VERSION).xpkg ) \
		ghcr.io/rossigee/provider-hostinger:$(VERSION)
	@$(OK) Pushed package ghcr.io/rossigee/provider-hostinger:$(VERSION)

# Go configuration
GO_SUBDIRS := cmd apis internal
GO_PROJECT := $(PROJECT_REPO)
GO_STATIC_PACKAGES = $(GO_PROJECT)/cmd/provider
GO_LDFLAGS += -X $(GO_PROJECT)/internal/version.Version=$(VERSION)

# Directories
S3_BUCKET_PATH ?= crossplane-releases
HELM_S3_BUCKET_PATH ?= crossplane-releases/helm
PLATFORMS ?= linux_amd64 linux_arm64

# Include build system makefiles
-include build/makelib/common.mk
-include build/makelib/output.mk
-include build/makelib/golang.mk
-include build/makelib/k8s_tools.mk
-include build/makelib/imagelight.mk
-include build/makelib/xpkg.mk

# Ensure package metadata exists before build
xpkg.build.provider-hostinger: do.build.images

# Neutralize plain image publish for ghcr (xpkg uses same ref; plain push would clobber package.yaml)
img.release.publish.ghcr.io/rossigee.provider-hostinger:
	@:
