# Provider-Hostinger Implementation Status

**Status**: ✅ **Phase 1 Complete** | 🔄 **Phase 2 In Progress** | ⏳ **Phases 3-8 Pending**

## Completed Tasks

### Phase 1: Project Initialization ✅
- [x] Created directory structure (15+ directories)
- [x] Created `.gitmodules` with rossigee/build submodule
- [x] Created `Makefile` with build orchestration
- [x] Created `go.mod` with Crossplane dependencies (Go 1.25)
- [x] Created `package/crossplane.yaml` (CRITICAL - enables Docker image embedding)
- [x] Created `cluster/images/provider-hostinger/Dockerfile` (CRITICAL - ENTRYPOINT pattern)
- [x] Created `VERSION` file (v0.1.0)
- [x] Initialized build submodule: `github.com/rossigee/build`

### Phase 2: API Resource Definitions (PARTIAL) 🔄

#### ✅ Completed
- [x] `apis/init.go` - Root API initialization
- [x] `apis/v1beta1/groupversion_info.go` - ProviderConfig API group definition
- [x] `apis/v1beta1/register.go` - Scheme registration
- [x] `apis/v1beta1/providerconfig_types.go` - Complete ProviderConfig with:
  - APIKeyAuthSpec (v1 API key authentication)
  - OAuthAuthSpec (v2 OAuth authentication)
  - ProviderConfig resource (cluster-scoped)

- [x] `apis/instance/v1beta1/groupversion_info.go` - Instance API group (instance.m.hostinger.crossplane.io)
- [x] `apis/instance/v1beta1/register.go` - Scheme registration
- [x] `apis/instance/v1beta1/types.go` - Complete Instance resource with:
  - InstanceParameters (hostname, CPU, RAM, disk, bandwidth, IPv6, inodes)
  - InstanceObservation (status, IP addresses, dates)
  - Kubebuilder markers for v1beta1 namespaced resources

#### ⏳ Remaining (Same Pattern)
- [ ] `apis/backup/v1beta1/*` - Backup resource (groupversion_info, register, types)
- [ ] `apis/firewall/v1beta1/*` - Firewall rules resource
- [ ] `apis/sshkey/v1beta1/*` - SSH key resource

## Critical Files (DO NOT MODIFY INCORRECTLY)

| File | Status | Purpose | Critical Notes |
|------|--------|---------|-----------------|
| `package/crossplane.yaml` | ✅ | Provider metadata | Must exist (not package.yaml) - enables Docker embedding |
| `cluster/images/provider-hostinger/Dockerfile` | ✅ | Container image | MUST use ENTRYPOINT (not CMD) - required by Crossplane runtime |
| `cmd/provider/main.go` | ⏳ | Entry point | MUST use `os.Getenv()` for paths (not hardcoded) |
| `build/.gitmodules` | ✅ | Build system | Must point to `github.com/rossigee/build` (not crossplane/build) |
| `.gitmodules` | ✅ | Submodule config | Contains build submodule reference |

## Implementation Architecture

```
Provider-Hostinger v1beta1 (Namespace-Scoped Resources)
├── ProviderConfig (Cluster-scoped)
│   ├── v1 API: APIKeyAuthSpec
│   └── v2 API: OAuthAuthSpec
├── Instance (Namespaced, .m. API group)
│   ├── Create/Update/Delete VPS instances
│   └── Fields: hostname, CPU, RAM, disk, IPv6, etc.
├── Backup (Namespaced, .m. API group)
│   ├── Manage instance backups
│   └── Fields: instance ref, scheduling, description
├── Firewall (Namespaced, .m. API group)
│   ├── Manage firewall rules
│   └── Fields: instance ref, rules, default action
└── SSHKey (Namespaced, .m. API group)
    ├── Manage SSH keys
    └── Fields: name, public key, instance attachments
```

## Files Ready for Implementation

### Skeleton Files to Create (Same Pattern as Instance)

**Backup** (apis/backup/v1beta1/):
```go
// groupversion_info.go: Group = "backup.m.hostinger.crossplane.io"
// types.go: BackupParameters (instanceId, description, scheduling)
//           BackupObservation (id, status, createdDate, size)
```

**Firewall** (apis/firewall/v1beta1/):
```go
// groupversion_info.go: Group = "firewall.m.hostinger.crossplane.io"
// types.go: FirewallRuleParameters (instanceId, rules[], defaultAction)
//           FirewallRuleObservation (id, status, appliedDate)
```

**SSHKey** (apis/sshkey/v1beta1/):
```go
// groupversion_info.go: Group = "sshkey.m.hostinger.crossplane.io"
// types.go: SSHKeyParameters (name, publicKey, instanceIds[])
//           SSHKeyObservation (id, fingerprint, createdDate)
```

## Next Implementation Steps

### Step 1: Complete API Definitions (Phase 2)
Create remaining resource type files following the Instance pattern:
```bash
# Create Backup, Firewall, SSHKey API types
# Each needs: groupversion_info.go, register.go, types.go
```

### Step 2: Generate CRDs and Deepcopy
```bash
cd /home/rossg/src/crossplane-providers/provider-hostinger
make generate
```

### Step 3: Implement Client Layer (Phase 3)
- `internal/clients/auth/` - Authentication handlers
- `internal/clients/hostinger.go` - Main API client factory
- `internal/clients/instance/` - Instance client
- `internal/clients/{backup,firewall,sshkey}/` - Other resource clients
- `internal/clients/errors.go` - Error classification

### Step 4: Implement Controllers (Phase 4)
- `internal/controller/hostinger.go` - Controller registration
- `internal/controller/instance/instance.go` - Instance controller
- `internal/controller/{backup,firewall,sshkey}/` - Other controllers

### Step 5: Create Entry Point (Phase 5)
- `cmd/provider/main.go` - Provider entry point with env var configuration

### Step 6: Create Examples & Documentation (Phase 6)
- `examples/instance/` - Instance examples
- `examples/backup/`, `examples/firewall/`, `examples/sshkey/`
- `README.md` - Comprehensive documentation

### Step 7: Setup CI/CD (Phase 7)
- `.github/workflows/ci.yml` - Validation workflow
- `.github/workflows/release.yml` - Publishing workflow

### Step 8: Quality Assurance (Phase 8)
```bash
make lint && make reviewable && make test
make xpkg.build
```

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| Namespaced v1beta1 (no .m.) | Follow provider-cloudflare pattern - simpler, proven |
| Support both v1 + v2 APIs | Flexible for different Hostinger customer setups |
| Separate client per resource | Better testability and separation of concerns |
| Interface-based clients | Enables mocking for unit tests |
| rossigee/build submodule | Upstream crossplane/build is broken |
| ENTRYPOINT in Dockerfile | Required by Crossplane runtime |
| Environment variables | Runtime-determined paths for TLS certificates |

## Quality Checkpoints

### Critical Files Validation
```bash
# ✅ MUST exist:
ls -la package/crossplane.yaml          # Provider metadata
ls -la cluster/images/provider-hostinger/Dockerfile

# ✅ MUST contain:
grep "ENTRYPOINT" cluster/images/provider-hostinger/Dockerfile
grep "os.Getenv" cmd/provider/main.go   # When created

# ✅ MUST NOT contain:
grep "CMD \[" cluster/images/provider-hostinger/Dockerfile  # Should be EMPTY
grep "/tmp/k8s-webhook" cmd/provider/main.go  # Should be EMPTY
```

### Build Verification
```bash
# After completing remaining phases:
make lint                    # Should pass
make reviewable             # Should pass
make test                   # Should pass
make xpkg.build             # Should create xpkg with Docker image
```

## Token Budget Notes

This implementation plan has been structured to be completed in manageable phases:
- **Phase 1**: ✅ Completed (20% complete)
- **Phase 2**: 🔄 In Progress (core APIs designed)
- **Phase 3-8**: ⏳ Ready for implementation

Due to API token constraints, subsequent phases should be implemented in focused sessions, each completing one phase at a time.

## References

- **Plan Document**: `/home/rossg/.claude/plans/gentle-greeting-hennessy.md`
- **Provider Documentation**: ../CLAUDE.md (parent directory)
- **Reference Provider**: ../provider-cloudflare (namespaced resources)
- **Reference Provider**: ../provider-minio (v2 migration patterns)
