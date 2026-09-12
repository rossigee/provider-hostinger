# Provider Hostinger Documentation

A Crossplane v2 provider for managing Hostinger cloud resources. All managed resources are namespaced (`*.m.hostinger.crossplane.io/v1beta1`) with multi-tenancy support.

## Resource Documentation

### Compute

| Resource | API Group | Description |
|----------|-----------|-------------|
| Instance | `instance.m.hostinger.crossplane.io/v1beta1` | Virtual instances |

### Networking

| Resource | API Group | Description |
|----------|-----------|-------------|
| Firewall | `firewall.m.hostinger.crossplane.io/v1beta1` | Firewall rules |
| SSHKey | `sshkey.m.hostinger.crossplane.io/v1beta1` | SSH keys |

### Backup

| Resource | API Group | Description |
|----------|-----------|-------------|
| Backup | `backup.m.hostinger.crossplane.io/v1beta1` | Backup management |

### Provider

| Resource | API Group | Description |
|----------|-----------|-------------|
| ProviderConfig | `hostinger.m.crossplane.io/v1beta1` | API credentials (cluster-scoped) |

## API Coverage Gaps

Hostinger API surface not yet modeled: DNS zone/record management, container/Kubernetes plans, mailbox/email accounts, domain registration, snapshots vs backups scheduling, firewall templates, and usage/billing queries.
