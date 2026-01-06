# Provider Hostinger

A Crossplane provider for managing Hostinger VPS and cloud services.

## Overview

`provider-hostinger` enables declarative management of Hostinger cloud resources through Kubernetes using [Crossplane](https://crossplane.io/). This provider implements full Crossplane v2 patterns with namespace-scoped resources, allowing multi-tenant deployments and GitOps workflows.

### Supported Resources

- **ProviderConfig** - Provider credentials management (supports both API v1 key and API v2 OAuth)
- **Instance** - VPS instance lifecycle management (create, update, delete)
- **Backup** - Automated and manual backup scheduling
- **FirewallRule** - Network security with inbound/outbound rules
- **SSHKey** - SSH key management for remote access

### Key Features

✅ **Dual Authentication** - Supports Hostinger API v1 (key-based) and v2 (OAuth)
✅ **Namespace-Scoped** - All resources are namespaced for multi-tenant isolation
✅ **Crossplane v2 Compatible** - Modern Crossplane patterns with management policies
✅ **Full CRUD Operations** - Complete lifecycle management of resources
✅ **Kubernetes-Native** - Manage cloud infrastructure as Kubernetes resources
✅ **GitOps Ready** - Works seamlessly with Flux, ArgoCD, and other GitOps tools

## Installation

### Prerequisites

- Kubernetes cluster (v1.24+)
- Crossplane v2.0+ installed
- kubectl configured to access your cluster

### Install the Provider

```bash
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-hostinger
spec:
  package: ghcr.io/rossigee/provider-hostinger:v0.1.0
  packagePullPolicy: IfNotPresent
EOF
```

Wait for the provider to be ready:

```bash
kubectl get providers
# NAME                  INSTALLED   HEALTHY   PACKAGE                                         AGE
# provider-hostinger    True        True      ghcr.io/rossigee/provider-hostinger:v0.1.0     1m
```

## Configuration

### 1. Create Hostinger Credentials Secret

Choose one authentication method:

#### Option A: API v1 Key Authentication

```bash
kubectl create secret generic hostinger-v1-credentials \
  --from-literal=api-key='your-api-key' \
  --from-literal=customer-id='your-customer-id' \
  -n crossplane-system
```

#### Option B: API v2 OAuth Authentication

```bash
kubectl create secret generic hostinger-v2-credentials \
  --from-literal=client-id='your-client-id' \
  --from-literal=client-secret='your-client-secret' \
  -n crossplane-system
```

### 2. Create ProviderConfig

#### For API v1 Key Authentication

```yaml
apiVersion: hostinger.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: crossplane-system
spec:
  credentials:
    apiKeyAuth:
      endpoint: "https://api.hostinger.com/v1"
      apiKeySecretRef:
        name: hostinger-v1-credentials
        key: api-key
      customerIdSecretRef:
        name: hostinger-v1-credentials
        key: customer-id
```

#### For API v2 OAuth Authentication

```yaml
apiVersion: hostinger.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: crossplane-system
spec:
  credentials:
    oauthAuth:
      endpoint: "https://api.hostinger.com/v2"
      tokenEndpoint: "https://auth.hostinger.com/oauth/token"
      clientIdSecretRef:
        name: hostinger-v2-credentials
        key: client-id
      clientSecretSecretRef:
        name: hostinger-v2-credentials
        key: client-secret
```

Apply the configuration:

```bash
kubectl apply -f providerconfig.yaml
```

## Usage Examples

### Create a VPS Instance

```yaml
apiVersion: instance.m.hostinger.crossplane.io/v1beta1
kind: Instance
metadata:
  name: my-vps
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    hostname: "my-vps.example.com"
    osId: "1"           # Ubuntu 22.04
    cpuCount: 2
    ram: 2048           # 2GB
    diskSize: 50        # 50GB
    ipv6Enabled: true
    bandwidth: 1000
  deletionPolicy: Delete
```

Apply:

```bash
kubectl apply -f instance.yaml
```

Check status:

```bash
kubectl get instances
kubectl describe instance my-vps
```

### Configure Firewall Rules

```yaml
apiVersion: firewall.m.hostinger.crossplane.io/v1beta1
kind: FirewallRule
metadata:
  name: my-firewall
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    instanceId: "123456"  # From Instance status
    defaultAction: deny
    rules:
      - port: "22"
        protocol: tcp
        direction: inbound
        action: allow
      - port: "80"
        protocol: tcp
        direction: inbound
        action: allow
      - port: "443"
        protocol: tcp
        direction: inbound
        action: allow
  deletionPolicy: Delete
```

### Schedule Backups

```yaml
apiVersion: backup.m.hostinger.crossplane.io/v1beta1
kind: Backup
metadata:
  name: daily-backup
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    instanceId: "123456"
    description: "Daily backup"
    schedule: daily
  deletionPolicy: Orphan  # Keep backup if resource is deleted
```

### Add SSH Keys

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-ssh-key
  namespace: default
type: Opaque
stringData:
  public-key: |
    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... your-key-here

---
apiVersion: sshkey.m.hostinger.crossplane.io/v1beta1
kind: SSHKey
metadata:
  name: my-ssh-key
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    name: "My SSH Key"
    publicKeySecretRef:
      name: my-ssh-key
      key: public-key
    instanceIds:
      - "123456"
  deletionPolicy: Delete
```

## API Resources

### ProviderConfig

Configuration for authenticating with Hostinger API.

**API Group**: `hostinger.crossplane.io`
**Version**: `v1beta1`
**Scope**: Cluster-scoped

### Instance

VPS instance management.

**API Group**: `instance.m.hostinger.crossplane.io`
**Version**: `v1beta1`
**Scope**: Namespaced
**Kind**: `Instance`

#### Spec.ForProvider

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| hostname | string | Yes | VPS hostname |
| osId | string | Yes | OS template ID |
| cpuCount | int32 | Yes | Number of CPU cores (min: 1) |
| ram | int32 | Yes | RAM in MB (min: 512) |
| diskSize | int32 | Yes | Disk size in GB (min: 10) |
| bandwidth | *int32 | No | Bandwidth in Mbps |
| ipv6Enabled | *bool | No | Enable IPv6 |
| inodes | *int32 | No | Inode limit |
| rootPasswordSecretRef | SecretKeySelector | No | Root password secret reference |

### Backup

Backup scheduling and management.

**API Group**: `backup.m.hostinger.crossplane.io`
**Version**: `v1beta1`
**Scope**: Namespaced
**Kind**: `Backup`

#### Spec.ForProvider

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| instanceId | string | Yes | Target instance ID |
| description | *string | No | Backup description |
| schedule | *BackupScheduleType | No | Schedule: manual, daily, weekly, monthly |

### FirewallRule

Network security rules.

**API Group**: `firewall.m.hostinger.crossplane.io`
**Version**: `v1beta1`
**Scope**: Namespaced
**Kind**: `FirewallRule`

#### Spec.ForProvider

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| instanceId | string | Yes | Target instance ID |
| rules | []FirewallRuleSpec | No | Array of rules |
| defaultAction | *FirewallAction | No | Default action (allow/deny) |

### SSHKey

SSH key management.

**API Group**: `sshkey.m.hostinger.crossplane.io`
**Version**: `v1beta1`
**Scope**: Namespaced
**Kind**: `SSHKey`

#### Spec.ForProvider

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Key name |
| publicKeySecretRef | SecretKeySelector | Yes | Public key secret reference |
| instanceIds | []string | No | Target instance IDs |

## Troubleshooting

### Provider not becoming ready

Check provider logs:

```bash
kubectl logs -n crossplane-system deployment/provider-hostinger
```

Common issues:
- **Connection timeout**: Verify Hostinger API endpoint is accessible
- **Authentication error**: Verify credentials in secret are correct
- **Certificate errors**: Check TLS certificate configuration

### Instance creation fails

Check the instance status:

```bash
kubectl describe instance my-vps
kubectl get instanceclaim my-vps -o yaml
```

Common issues:
- **Quota exceeded**: Hostinger account may have reached resource limits
- **Invalid OS ID**: Verify osId is supported for your account
- **Invalid hostname**: Ensure hostname is properly formatted

### Connection refused

Verify the ProviderConfig:

```bash
kubectl get providerconfig
kubectl describe providerconfig default
```

Check secret exists and has correct keys:

```bash
kubectl get secrets -n crossplane-system hostinger-v1-credentials -o yaml
```

## Development

### Building the Provider

```bash
# Initialize build submodule
git submodule update --init --recursive

# Build and test
make lint
make test
make build

# Create Crossplane package
make xpkg.build

# Publish (requires Docker registry access)
make publish
```

### Testing

```bash
# Run unit tests
make test

# Run linting
make lint

# Run pre-commit validation
make reviewable
```

## Support & Contributing

For issues, feature requests, or contributions:

1. Check [existing issues](https://github.com/rossigee/provider-hostinger/issues)
2. Open a [new issue](https://github.com/rossigee/provider-hostinger/issues/new)
3. Submit a [pull request](https://github.com/rossigee/provider-hostinger/pulls)

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details

## Security

For security vulnerabilities, please email security@golder.tech instead of using the issue tracker.

## Related Resources

- [Crossplane Documentation](https://docs.crossplane.io/)
- [Hostinger API Documentation](https://docs.hostinger.com/)
- [Crossplane Provider Development](https://docs.crossplane.io/v1.15/guides/provider-development/)
