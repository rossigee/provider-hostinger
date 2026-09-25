# Getting Started

Guide to getting started with provider-hostinger.

## Installation

Install the provider:

```bash
kubectl crossplane install provider ghcr.io/rossigee/provider-hostinger:v0.2.5
```

## Prerequisites

- Kubernetes cluster with Crossplane v2.5.0 or later installed
- Hostinger API credentials

## Quick Start

1. Create a ProviderConfig:

```yaml
apiVersion: hostinger.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: crossplane-system
spec:
  credentials:
    source: Secret
    secretRef:
      name: hostinger-credentials
      namespace: crossplane-system
```
