# Project configuration
PROJECT_NAME := provider-hostinger
PROJECT_REPO := github.com/rossigee/$(PROJECT_NAME)

# Build configuration
REGISTRY_ORGS = ghcr.io/rossigee
XPKG_REG_ORGS ?= ghcr.io/rossigee
CROSSPLANE_VERSION = 2.0.2
GO_REQUIRED_VERSION ?= 1.25.5
GOLANGCILINT_VERSION ?= 2.7.2

# Images configuration
IMAGES = provider-hostinger

# Crossplane package configuration
XPKGS = provider-hostinger

# Go configuration
GO_SUBDIRS := cmd apis internal
GO_PROJECT := $(PROJECT_REPO)
GO_STATIC_PACKAGES = $(GO_PROJECT)/cmd/provider
GO_LDFLAGS += -X $(GO_PROJECT)/internal/version.Version=$(VERSION)

# Directories
S3_BUCKET_PATH ?= crossplane-releases
HELM_S3_BUCKET_PATH ?= crossplane-releases/helm
PLATFORMS ?= linux_amd64

# Include build system makefiles
-include build/makelib/common.mk
-include build/makelib/output.mk
-include build/makelib/golang.mk
-include build/makelib/k8s_tools.mk
-include build/makelib/imagelight.mk
-include build/makelib/xpkg.mk

# Ensure package metadata exists before build
xpkg.build.provider-hostinger: do.build.images
