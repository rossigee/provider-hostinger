# Provider-Hostinger API Structure

## Completed API Definitions

### ProviderConfig (apis/v1beta1/)
- ✅ groupversion_info.go
- ✅ register.go
- ✅ providerconfig_types.go
  - APIKeyAuthSpec (for v1 API key auth)
  - OAuthAuthSpec (for v2 OAuth auth)
  - ProviderConfig with cluster-scoped resource

### Instance (apis/instance/v1beta1/)
- ✅ groupversion_info.go (API group: instance.m.hostinger.crossplane.io)
- ✅ register.go
- ✅ types.go
  - InstanceParameters with VPS configuration (hostname, CPU, RAM, disk, IPv6, etc.)
  - InstanceObservation with status information
  - Instance resource type (namespaced, .m. API group)

## Remaining API Definitions

### Backup (apis/backup/v1beta1/)
Should include:
- BackupParameters: instanceId, description, scheduling
- BackupObservation: id, status, createdDate, size
- Backup resource type

### Firewall (apis/firewall/v1beta1/)
Should include:
- FirewallRuleParameters: instanceId, rules[], defaultAction
- FirewallRuleObservation: id, status, appliedDate
- FirewallRule resource type

### SSHKey (apis/sshkey/v1beta1/)
Should include:
- SSHKeyParameters: name, publicKey, instanceIds[]
- SSHKeyObservation: id, fingerprint, createdDate
- SSHKey resource type

## Implementation Notes

1. All resources use v1beta1 API version for namespaced support
2. Instance, Backup, Firewall, SSHKey use .m. API group for full v2 support
3. ProviderConfig is cluster-scoped (not namespaced) - standard pattern
4. All resources implement xpv1.ResourceSpec and xpv1.ResourceStatus patterns
5. All parameters marked with proper kubebuilder validation markers
6. All observations include external ID storage pattern

## Next Steps

1. Create Backup, Firewall, SSHKey API types (groupversion_info, register, types)
2. Run `make generate` to create deepcopy and CRD files
3. Implement client layer for authentication and API calls
4. Implement controllers for each resource
5. Create entry point (cmd/provider/main.go)
6. Add examples and documentation
7. Setup CI/CD workflows
