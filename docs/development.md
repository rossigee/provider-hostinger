# Development

Guide for developing the provider.

## Prerequisites

- Go 1.27.1 or later
- Kubernetes cluster
- Crossplane v2.5.0 or later
- kubectl configured

## Building

```bash
make build
```

## Testing

```bash
make test
```

## Running Locally

```bash
make run
```

## Release Process

1. Set `VERSION` to the exact release version.
2. Update `CHANGELOG.md` and current installation references.
3. Create the exact `vMAJOR.MINOR.PATCH` tag at the current `origin/master` commit.
4. The tag-only release workflow builds `linux_amd64` and `linux_arm64` xpkg files, publishes the version and `latest` tags, verifies equal digests and both architectures, and creates the GitHub Release.
