# Provider-Hostinger: Comprehensive Project Outline

## 🎯 Project Summary

A fresh Crossplane v2 provider for Hostinger VPS management, based on provider-cloudflare patterns with v1 (API key) + v2 (OAuth) authentication support. Namespace-scoped resources using .m. API groups for full v2 compatibility.

**Location**: `/home/rossg/src/crossplane-providers/provider-hostinger`
**Status**: 🟢 Phase 1 Complete | 🟡 Phase 2 In Progress | 🔵 Phases 3-8 Ready

---

## 📊 Implementation Progress

### Summary Statistics
- **Files Created**: 17 (5 Go files, 4 Markdown docs, 1 Dockerfile, 1 Makefile, 1 go.mod, 5 config)
- **Lines of Code**: ~500+ (Go) + documentation
- **Completion**: ~20% (Phase 1-2 partial)
- **Critical Files**: 100% Complete ✅

### Phase Breakdown
```
Phase 1: Project Initialization        ✅ 100% Complete (7/7 items)
Phase 2: API Resource Definitions      🔄  50% Complete (5/9 items)
Phase 3: Client Layer                  ⏳   0% Complete (ready for implementation)
Phase 4: Controllers                   ⏳   0% Complete (ready for implementation)
Phase 5: Entry Point                   ⏳   0% Complete (ready for implementation)
Phase 6: Examples & Documentation      ⏳   0% Complete (ready for implementation)
Phase 7: CI/CD Workflows               ⏳   0% Complete (ready for implementation)
Phase 8: Quality Assurance             ⏳   0% Complete (ready for implementation)
```

---

## 📂 Project Structure (What's Been Created)

```
provider-hostinger/
├── ✅ .gitmodules                          (rossigee/build submodule)
├── ✅ Makefile                             (build orchestration - 40 lines)
├── ✅ VERSION                              (v0.1.0)
├── ✅ go.mod                               (Go 1.25 dependencies)
├── ✅ Dockerfile                           (CRITICAL: ENTRYPOINT pattern)
├── ✅ QUICKSTART.md                        (This phase's guide)
├── ✅ IMPLEMENTATION_STATUS.md             (Detailed progress)
├── ✅ API_STRUCTURE.md                     (API architecture)
├── ✅ PROJECT_OUTLINE.md                   (This file)
│
├── build/                                  (rossigee/build git submodule)
│   ├── makelib/                            (build system makefiles)
│   └── ...
│
├── apis/
│   ├── ✅ init.go                          (root API initialization)
│   ├── v1beta1/                            (ProviderConfig API)
│   │   ├── ✅ groupversion_info.go
│   │   ├── ✅ register.go
│   │   └── ✅ providerconfig_types.go      (v1 key + v2 OAuth)
│   └── instance/v1beta1/                   (VPS Instance resource)
│       ├── ✅ groupversion_info.go
│       ├── ✅ register.go
│       └── ✅ types.go                     (Complete definition)
│   ├── backup/v1beta1/                     (TO CREATE)
│   ├── firewall/v1beta1/                   (TO CREATE)
│   └── sshkey/v1beta1/                     (TO CREATE)
│
├── cluster/images/provider-hostinger/
│   └── ✅ Dockerfile                       (ENTRYPOINT-based container)
│
├── cmd/
│   └── provider/
│       └── main.go                         (TO CREATE)
│
├── config/
│   ├── crd/                                (generated CRDs - TO CREATE)
│   ├── manager/
│   │   ├── config.yaml
│   │   └── manager.yaml
│   └── provider/
│       └── provider.yaml
│
├── examples/
│   ├── instance/
│   │   ├── providerconfig.yaml             (TO CREATE)
│   │   └── instance.yaml                   (TO CREATE)
│   ├── backup/
│   ├── firewall/
│   └── sshkey/
│
├── internal/
│   ├── clients/
│   │   ├── auth/
│   │   │   ├── authenticator.go            (TO CREATE)
│   │   │   ├── v1keyauth.go                (TO CREATE)
│   │   │   └── v2oauthauth.go              (TO CREATE)
│   │   ├── instance/
│   │   │   ├── instance.go                 (TO CREATE)
│   │   │   └── interfaces.go               (TO CREATE)
│   │   ├── {backup,firewall,sshkey}/       (TO CREATE)
│   │   ├── hostinger.go                    (TO CREATE)
│   │   └── errors.go                       (TO CREATE)
│   ├── controller/
│   │   ├── hostinger.go                    (TO CREATE)
│   │   └── instance/
│   │       └── instance.go                 (TO CREATE)
│   ├── metrics/
│   │   └── metrics.go                      (TO CREATE)
│   └── version/
│       └── version.go                      (TO CREATE)
│
├── package/
│   ├── ✅ crossplane.yaml                  (CRITICAL: provider metadata)
│   └── crds/                               (generated - TO CREATE)
│
├── test/
│   ├── unit/
│   │   └── controller_test.go              (TO CREATE)
│   └── e2e/
│       └── provider_test.go                (TO CREATE)
│
├── .github/
│   └── workflows/
│       ├── ci.yml                          (TO CREATE)
│       └── release.yml                     (TO CREATE)
│
└── README.md                               (TO CREATE)
```

---

## 🔑 Critical Files Status

### ✅ COMPLETED (MUST NOT MODIFY INCORRECTLY)

| File | Status | Why Critical |
|------|--------|--------------|
| `package/crossplane.yaml` | ✅ Ready | Enables Docker image embedding in .xpkg |
| `cluster/images/provider-hostinger/Dockerfile` | ✅ Ready | ENTRYPOINT pattern required by Crossplane |
| `.gitmodules` | ✅ Ready | Points to rossigee/build (upstream is broken) |
| `Makefile` | ✅ Ready | Build orchestration for all targets |
| `go.mod` | ✅ Ready | Go 1.25 with Crossplane 1.21.0 |

### ⏳ TO CREATE (NEXT PHASES)

| File | Phase | Purpose |
|------|-------|---------|
| `cmd/provider/main.go` | 5 | Entry point with env var config |
| `internal/clients/auth/*` | 3 | Authentication handlers |
| `internal/controller/*` | 4 | Resource controllers |
| `.github/workflows/ci.yml` | 7 | Validation workflow |
| `examples/*.yaml` | 6 | Usage examples |
| `README.md` | 6 | Documentation |

---

## 🚀 API Resources Defined

### ProviderConfig (Cluster-Scoped)
- **Location**: `apis/v1beta1/`
- **Status**: ✅ Complete
- **Features**:
  - API v1: APIKeyAuthSpec (endpoint, api_key, customer_id)
  - API v2: OAuthAuthSpec (endpoint, client_id, client_secret, token_endpoint)
  - Secret references for sensitive data

### Instance (Namespaced, .m. API Group)
- **Location**: `apis/instance/v1beta1/`
- **Status**: ✅ Complete
- **Features**:
  - Create, read, update, delete VPS instances
  - Parameters: hostname, osId, cpuCount, ram, diskSize, bandwidth, IPv6, inodes
  - Observations: id, status, ipAddress, ipv6Address, creationDate, expirationDate

### Backup (Namespaced, .m. API Group) [TO CREATE]
- **Location**: `apis/backup/v1beta1/`
- **Structure**:
  - Parameters: instanceId, description, scheduling
  - Observations: id, status, createdDate, size

### Firewall (Namespaced, .m. API Group) [TO CREATE]
- **Location**: `apis/firewall/v1beta1/`
- **Structure**:
  - Parameters: instanceId, rules[], defaultAction
  - Observations: id, status, appliedDate

### SSHKey (Namespaced, .m. API Group) [TO CREATE]
- **Location**: `apis/sshkey/v1beta1/`
- **Structure**:
  - Parameters: name, publicKey, instanceIds[]
  - Observations: id, fingerprint, createdDate

---

## 📋 Comprehensive To-Do List

### ✅ COMPLETED ITEMS (7)
1. [x] Create directory structure
2. [x] Create .gitmodules with build submodule
3. [x] Create Makefile with build system
4. [x] Create go.mod with Crossplane dependencies
5. [x] Create package/crossplane.yaml (CRITICAL)
6. [x] Create Dockerfile with ENTRYPOINT (CRITICAL)
7. [x] Initialize build submodule

### 🔄 IN PROGRESS (Phase 2 - API Definitions)
8. [x] Create ProviderConfig API (v1beta1)
9. [x] Create Instance resource (v1beta1, .m. group)
10. [ ] Create Backup resource (same pattern)
11. [ ] Create Firewall resource (same pattern)
12. [ ] Create SSHKey resource (same pattern)
13. [ ] Run `make generate` to create CRDs

### ⏳ PENDING ITEMS (27)
14. [ ] Phase 3: Create authentication handlers
15. [ ] Phase 3: Create Hostinger API client factory
16. [ ] Phase 3: Create instance client
17. [ ] Phase 3: Create other resource clients
18. [ ] Phase 3: Create error classification
19. [ ] Phase 4: Create controller registration
20. [ ] Phase 4: Create Instance controller
21. [ ] Phase 4: Create other controllers
22. [ ] Phase 5: Create cmd/provider/main.go
23. [ ] Phase 6: Create ProviderConfig examples
24. [ ] Phase 6: Create Instance examples
25. [ ] Phase 6: Create other examples
26. [ ] Phase 6: Write README.md
27. [ ] Phase 7: Create ci.yml workflow
28. [ ] Phase 7: Create release.yml workflow
29. [ ] Phase 8: Run make lint
30. [ ] Phase 8: Run make reviewable
31. [ ] Phase 8: Run make test
32. [ ] Phase 8: Run make xpkg.build
33. [ ] Phase 8: Verify Docker image embedding
34. [ ] Phase 8: Test deployment
35. [ ] Documentation review and finalization

---

## 🎓 Key Design Patterns

### API Design Pattern
```go
// All resources follow this pattern:
type <Resource>Parameters struct {
    // spec.forProvider fields
}

type <Resource>Observation struct {
    // status.atProvider fields
}

type <Resource>Spec struct {
    xpv1.ResourceSpec `json:",inline"`
    ForProvider       <Resource>Parameters `json:"forProvider"`
}

type <Resource>Status struct {
    xpv1.ResourceStatus `json:",inline"`
    AtProvider          <Resource>Observation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,hostinger}
type <Resource> struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   <Resource>Spec   `json:"spec,omitempty"`
    Status <Resource>Status `json:"status,omitempty"`
}
```

### Controller Pattern
```
Connector (creates API clients from ProviderConfig)
    ↓
External (implements managed CRUD interface)
    ├─ Observe (check existence + state)
    ├─ Create (provision resource)
    ├─ Update (modify resource)
    └─ Delete (deprovision resource)
```

### Client Pattern
```
Interface-based clients for testability
    ↓
Separate client per resource type
    ↓
Helper functions (GenerateObservation, UpToDate, LateInitialize)
    ↓
Error classification (IsNotFound, IsInvalidConfig, etc.)
```

---

## 📚 Documentation References

| Document | Purpose | Location |
|----------|---------|----------|
| QUICKSTART.md | Phase-by-phase implementation guide | ./QUICKSTART.md |
| IMPLEMENTATION_STATUS.md | Detailed progress tracking | ./IMPLEMENTATION_STATUS.md |
| API_STRUCTURE.md | API architecture overview | ./API_STRUCTURE.md |
| PROJECT_OUTLINE.md | This file - comprehensive overview | ./PROJECT_OUTLINE.md |
| Implementation Plan | Full scope and design decisions | `/home/rossg/.claude/plans/gentle-greeting-hennessy.md` |
| Parent Documentation | Critical build system requirements | `../CLAUDE.md` |

---

## 🔍 Verification Commands

### Check Critical Files
```bash
# Must exist and be correct:
ls -la package/crossplane.yaml                          # Provider metadata
grep "ENTRYPOINT" cluster/images/provider-hostinger/Dockerfile
grep "github.com/rossigee/build" .gitmodules

# Must NOT contain:
grep "CMD \[" cluster/images/provider-hostinger/Dockerfile   # Should fail
grep "/tmp/k8s-webhook" cmd/provider/main.go                 # Will be checked later
```

### Current Project Statistics
```bash
find . -type f -name "*.go" | wc -l           # Count Go files
wc -l **/*.go                                  # Lines of code
du -sh .                                       # Total size
```

---

## 🛠️ Build System Commands

```bash
# Navigate to project
cd /home/rossg/src/crossplane-providers/provider-hostinger

# After Phase 2 completion:
make generate           # Generate CRDs and deepcopy

# During development:
make lint              # Code linting
make test              # Unit tests
make build             # Build binary
make docker.build      # Build Docker image

# Full validation:
make reviewable        # generate + lint + test
make xpkg.build        # Build Crossplane package

# Publishing:
make publish           # Complete workflow (version tag required)
```

---

## 📊 Success Metrics

| Metric | Target | Status |
|--------|--------|--------|
| Directory structure complete | ✅ 15+ directories | ✅ Done |
| Critical build files | ✅ 5 files | ✅ Done |
| API definitions | ✅ 5 resources | 🔄 2/5 complete |
| Client layer | ✅ Complete | ⏳ Pending |
| Controllers | ✅ Complete | ⏳ Pending |
| Entry point | ✅ Functional | ⏳ Pending |
| Examples | ✅ All resources | ⏳ Pending |
| CI/CD workflows | ✅ Both files | ⏳ Pending |
| Make targets | ✅ All working | ⏳ Test after Phase 2 |
| Docker image embedding | ✅ Verified | ⏳ Test after completion |

---

## 🎯 Next Immediate Steps

1. **Phase 2 Completion** (THIS SESSION)
   - Create Backup, Firewall, SSHKey API types
   - Run `make generate`
   - Verify CRDs generated

2. **Phase 3** (NEXT SESSION)
   - Implement authentication handlers
   - Implement API client factory
   - Implement resource clients

3. **Phase 4** (NEXT SESSION)
   - Implement controllers
   - Test controller logic

4. **Phase 5+** (SUBSEQUENT SESSIONS)
   - Create entry point
   - Add examples
   - Setup CI/CD
   - Quality assurance

---

## 📞 Important Contacts

**Current Implementation**:
- Status: 🟢 Phase 1 Complete, 🟡 Phase 2 In Progress
- Location: `/home/rossg/src/crossplane-providers/provider-hostinger`
- Git Branch: master
- Plan File: `/home/rossg/.claude/plans/gentle-greeting-hennessy.md`

**Parent Repository**:
- Location: `/home/rossg/src/crossplane-providers/`
- Contains 14+ other providers as references
- Build system: `github.com/rossigee/build`

---

**Last Updated**: 2025-01-06
**Next Review**: After Phase 2 completion
**Status**: ✅ Ready for Phase 2 Completion → Phase 3 Implementation
