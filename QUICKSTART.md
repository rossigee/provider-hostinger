# Provider-Hostinger Quickstart

## Current Status

**✅ Phase 1 Complete** - Project Initialization & Critical Build System Files
**🔄 Phase 2 In Progress** - API Resource Definitions (Core resources complete)

## What's Been Created

### Critical Infrastructure (Foundation Ready)
```
✅ .gitmodules                          - Build submodule reference
✅ Makefile                             - Build system orchestration
✅ go.mod                               - Go 1.27.1 dependencies
✅ package/crossplane.yaml              - CRITICAL: Provider metadata
✅ cluster/images/provider-hostinger/Dockerfile - CRITICAL: ENTRYPOINT pattern
✅ VERSION                              - v0.2.5 release version
✅ build/                               - rossigee/build submodule initialized
```

### API Definitions (Partial)
```
✅ apis/init.go                         - Root API initialization
✅ apis/v1beta1/*                       - ProviderConfig resource
  ├─ groupversion_info.go
  ├─ register.go
  └─ providerconfig_types.go            - v1 API key + v2 OAuth support

✅ apis/instance/v1beta1/*              - VPS Instance resource
  ├─ groupversion_info.go
  ├─ register.go
  └─ types.go                           - Complete Instance definition
```

### Documentation
```
✅ IMPLEMENTATION_STATUS.md             - Detailed progress report
✅ API_STRUCTURE.md                     - API architecture overview
```

## Architecture Overview

```
Provider-Hostinger
├── ProviderConfig (v1beta1, cluster-scoped)
│   ├── API v1: APIKeyAuthSpec (endpoint, api_key, customer_id)
│   └── API v2: OAuthAuthSpec (endpoint, client_id, client_secret, token_endpoint)
│
├── Instance (v1beta1, namespaced, .m. group)
│   ├── Spec: hostname, osId, cpuCount, ram, diskSize, bandwidth, ipv6Enabled
│   └── Status: id, status, ipAddress, ipv6Address, creationDate
│
├── Backup (v1beta1, namespaced, .m. group) [TO BE CREATED]
├── Firewall (v1beta1, namespaced, .m. group) [TO BE CREATED]
└── SSHKey (v1beta1, namespaced, .m. group) [TO BE CREATED]
```

## Next Implementation (Phase 2 Completion)

### Step 1: Complete API Definitions
Create remaining resource types (Backup, Firewall, SSHKey) following the Instance pattern:

```bash
# Each resource needs 3 files:
# apis/<resource>/v1beta1/groupversion_info.go    (API group definition)
# apis/<resource>/v1beta1/register.go             (Scheme registration)
# apis/<resource>/v1beta1/types.go                (Resource type definitions)
```

**Pattern to follow** (see `apis/instance/v1beta1/types.go`):
- Parameters struct (spec.forProvider fields)
- Observation struct (status.atProvider fields)
- Main resource struct with kubebuilder markers
- List type struct

### Step 2: Generate CRDs
```bash
cd /home/rossg/src/crossplane-providers/provider-hostinger
make generate
```

This will:
- Create `zz_generated.*.go` files (deepcopy, etc.)
- Generate CRDs in `config/crd/`
- Update package metadata

### Step 3: Implement Client Layer (Phase 3)
Structure:
```
internal/clients/
├── auth/
│   ├── authenticator.go        # Base interface
│   ├── v1keyauth.go            # API v1 key auth
│   └── v2oauthauth.go          # API v2 OAuth auth
├── hostinger.go                # Main API client factory
├── instance/
│   ├── instance.go             # Instance client
│   ├── interfaces.go           # Client interface
│   └── client_test.go
├── {backup,firewall,sshkey}/   # Other resource clients
└── errors.go                   # Error classification
```

### Step 4: Implement Controllers (Phase 4)
Structure:
```
internal/controller/
├── hostinger.go                # Controller registration
└── instance/
    ├── instance.go             # Instance controller
    └── instance_test.go
```

Controller pattern:
- Connector: Creates API clients from ProviderConfig
- External: Implements CRUD operations (Observe, Create, Update, Delete)
- Setup: Registers controller with manager

### Step 5: Create Entry Point (Phase 5)
```
cmd/provider/main.go            # Provider binary entry point
```

Key requirements:
- Use `os.Getenv("WEBHOOK_TLS_CERT_DIR")` for certificate paths
- Support `LEADER_ELECT` environment variable
- Register all APIs and controllers
- Configure webhook server (if needed)

### Step 6: Add Examples & Documentation (Phase 6)
```
examples/
├── instance/
│   ├── providerconfig-v1.yaml  # API v1 example
│   ├── providerconfig-v2.yaml  # API v2 example
│   └── instance.yaml           # Instance example
├── backup/
├── firewall/
└── sshkey/

README.md                        # Comprehensive documentation
```

### Step 7: Setup CI/CD (Phase 7)
```
.github/workflows/
├── ci.yml                      # Validation (lint, test, security)
└── release.yml                 # Publishing (Docker + xpkg)
```

Key pattern: **"CI validates, Release publishes"**
- CI runs on push/PR (no registry publishing)
- Release runs on tags (publishes to ghcr.io/rossigee)

### Step 8: Quality Assurance (Phase 8)
```bash
# Verify critical files
ls -la package/crossplane.yaml
grep "ENTRYPOINT" cluster/images/provider-hostinger/Dockerfile
grep "os.Getenv" cmd/provider/main.go

# Build and test
make lint && make reviewable && make test
make xpkg.build

# Verify Docker image embedding
tar -tf _output/xpkg/linux_amd64/provider-hostinger-*.xpkg | grep manifest.json
```

## Critical Success Factors

| Requirement | Status | Details |
|-------------|--------|---------|
| `package/crossplane.yaml` | ✅ | Enables Docker image embedding in xpkg |
| `ENTRYPOINT` in Dockerfile | ✅ | Required by Crossplane runtime (not CMD) |
| `rossigee/build` submodule | ✅ | Upstream `crossplane/build` is broken |
| Environment variables | ⏳ | Will be implemented in cmd/provider/main.go |
| API group naming | ✅ | Using .m. for namespaced resources |
| Error classification | ⏳ | Will implement in clients/errors.go |

## File Reference

| File | Purpose | Created | Notes |
|------|---------|---------|-------|
| `.gitmodules` | Submodule config | ✅ | Points to github.com/rossigee/build |
| `Makefile` | Build orchestration | ✅ | Includes all standard targets |
| `go.mod` | Dependencies | ✅ | Go 1.27.1, Crossplane 2.5.0 |
| `package/crossplane.yaml` | Provider metadata | ✅ | **CRITICAL** - NOT package.yaml |
| `Dockerfile` | Container image | ✅ | **CRITICAL** - Uses ENTRYPOINT |
| `apis/v1beta1/*` | ProviderConfig | ✅ | v1 key + v2 OAuth support |
| `apis/instance/v1beta1/*` | Instance resource | ✅ | Full VPS instance definition |
| `cmd/provider/main.go` | Entry point | ⏳ | Needs environment var config |
| `internal/clients/*` | API clients | ⏳ | Authentication + resource clients |
| `internal/controller/*` | Controllers | ⏳ | Resource reconciliation logic |

## Environment

- **Go Version**: 1.27.1
- **Crossplane Version**: 2.5.0
- **Build System**: github.com/rossigee/build
- **Registry**: ghcr.io/rossigee/provider-hostinger
- **API Pattern**: v1beta1 (namespaced) with .m. API groups
- **Current Directory**: `/home/rossg/src/crossplane-providers/provider-hostinger`

## References

- **Full Plan**: `/home/rossg/.claude/plans/gentle-greeting-hennessy.md`
- **Status Report**: `./IMPLEMENTATION_STATUS.md` (this directory)
- **API Structure**: `./API_STRUCTURE.md` (this directory)
- **Reference Providers**:
  - `../provider-cloudflare` - Namespaced v1beta1 resources
  - `../provider-minio` - v2 migration patterns
- **Parent Documentation**: `../CLAUDE.md` - Critical build system requirements

## Quick Commands

```bash
# Navigate to provider directory
cd /home/rossg/src/crossplane-providers/provider-hostinger

# After completing remaining API definitions:
make generate          # Generate CRDs and deepcopy

# During development:
make lint              # Code linting
make test              # Unit tests
make build             # Build binary
make docker.build      # Build Docker image

# Final validation:
make reviewable        # Full pre-commit validation
make xpkg.build        # Build Crossplane package

# Before publishing:
make publish           # Complete build + publish workflow
```

---

**Status**: Ready for Phase 2 Completion (Backup, Firewall, SSHKey API definitions)

**Next Session**: Create remaining 3 resource types, then `make generate`, then move to Phase 3 (Client Layer)
