# Email Sentinel - Deployment Guide

This document provides comprehensive instructions for deploying, releasing, and distributing Email Sentinel.

## Table of Contents

- [Overview](#overview)
- [Release Process](#release-process)
- [Package Manager Setup](#package-manager-setup)
- [GitHub Secrets Configuration](#github-secrets-configuration)
- [Testing Releases](#testing-releases)
- [Docker Deployment](#docker-deployment)
- [Troubleshooting](#troubleshooting)

## Overview

Email Sentinel uses **GoReleaser** for automated multi-platform builds and releases. When you push a version tag (e.g., `v1.0.0`), GitHub Actions automatically:

1. Builds binaries for Windows, macOS, and Linux (multiple architectures)
2. Creates installers (.deb, .rpm packages)
3. Generates checksums and archives
4. Publishes to GitHub Releases
5. Updates Homebrew tap and Scoop bucket
6. Builds and pushes Docker image

## Release Process

### 1. Prepare for Release

```bash
# Ensure you're on main branch with latest changes
git checkout main
git pull origin main

# Run tests
go test ./...

# Verify build works locally
go build -o email-sentinel .
```

### 2. Create Version Tag

Use semantic versioning (`vMAJOR.MINOR.PATCH`):

```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release v1.0.0: Initial public release"

# Push tag to trigger release workflow
git push origin v1.0.0
```

### 3. Monitor Release

1. Go to **Actions** tab in GitHub repository
2. Watch the "Release" workflow progress
3. Once complete, check the **Releases** page
4. Verify all artifacts are present:
   - Windows ZIP archives
   - macOS tar.gz archives (universal binary)
   - Linux tar.gz archives
   - .deb packages
   - .rpm packages
   - Checksums file
   - Source code archives

### 4. Announce Release

After successful release:

1. Edit release notes on GitHub (auto-generated, can be customized)
2. Announce on social media, Discord, or relevant channels
3. Update documentation if needed

## Package Manager Setup

### Homebrew Tap (macOS/Linux)

#### Initial Setup

1. Create a new GitHub repository: `datateamsix/homebrew-tap`

2. Add GitHub secret `HOMEBREW_TAP_TOKEN`:
   - Go to GitHub Settings → Developer settings → Personal access tokens
   - Generate new token (classic) with `repo` scope
   - Add to repository secrets as `HOMEBREW_TAP_TOKEN`

#### User Installation

```bash
brew tap datateamsix/tap
brew install email-sentinel

# Update
brew upgrade email-sentinel
```

### Scoop Bucket (Windows)

#### Initial Setup

1. Create a new GitHub repository: `datateamsix/scoop-bucket`

2. Add GitHub secret `SCOOP_BUCKET_TOKEN`:
   - Use the same PAT as Homebrew or create a new one
   - Add to repository secrets as `SCOOP_BUCKET_TOKEN`

#### User Installation

```powershell
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel

# Update
scoop update email-sentinel
```

### Linux Packages

DEB and RPM packages are automatically built and attached to GitHub Releases.

#### Debian/Ubuntu

```bash
# Download from GitHub Releases
wget https://github.com/datateamsix/email-sentinel/releases/download/v1.0.0/email-sentinel_1.0.0_amd64.deb

# Install
sudo dpkg -i email-sentinel_1.0.0_amd64.deb

# If dependencies missing:
sudo apt-get install -f
```

#### RHEL/Fedora/CentOS

```bash
# Download from GitHub Releases
wget https://github.com/datateamsix/email-sentinel/releases/download/v1.0.0/email-sentinel_1.0.0_x86_64.rpm

# Install
sudo rpm -i email-sentinel_1.0.0_x86_64.rpm

# Or with yum/dnf for dependency resolution:
sudo yum install email-sentinel_1.0.0_x86_64.rpm
```

## GitHub Secrets Configuration

Required secrets for automated releases:

### GITHUB_TOKEN (Automatic)

Automatically provided by GitHub Actions. No setup needed.

### HOMEBREW_TAP_TOKEN

**Purpose**: Update Homebrew formula automatically

**Setup**:
1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Select scopes: `repo` (Full control of private repositories)
4. Copy token
5. Go to repository Settings → Secrets and variables → Actions
6. New repository secret: `HOMEBREW_TAP_TOKEN`

### SCOOP_BUCKET_TOKEN

**Purpose**: Update Scoop manifest automatically

**Setup**: Same as HOMEBREW_TAP_TOKEN (can use same token or create separate)

### DOCKER_USERNAME and DOCKER_TOKEN

**Purpose**: Push Docker images to Docker Hub

**Setup**:
1. Create Docker Hub account if needed
2. Go to Docker Hub → Account Settings → Security
3. Create new access token
4. Add secrets:
   - `DOCKER_USERNAME`: Your Docker Hub username
   - `DOCKER_TOKEN`: Access token from step 3

### Optional: GPG_PRIVATE_KEY

**Purpose**: Sign releases cryptographically

**Setup**:
```bash
# Generate GPG key
gpg --full-generate-key

# Export private key
gpg --export-secret-keys --armor YOUR_KEY_ID > private.key

# Add to GitHub secrets (paste entire contents)
```

Then uncomment the `signs` section in `.goreleaser.yaml`.

## Testing Releases

### Local Testing with GoReleaser

Test releases locally before pushing tags:

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Validate configuration
goreleaser check

# Build snapshot (no publish, no tag required)
goreleaser release --snapshot --clean

# Check output in dist/ directory
ls -lah dist/
```

### Test Installation

After snapshot build:

```bash
# Extract and test binary
cd dist
tar -xzf email-sentinel_*_linux_amd64.tar.gz
./email-sentinel --version

# Test DEB package (Ubuntu/Debian)
sudo dpkg -i email-sentinel_*_amd64.deb
email-sentinel --version
sudo dpkg -r email-sentinel

# Test RPM package (Fedora/RHEL)
sudo rpm -i email-sentinel_*_x86_64.rpm
email-sentinel --version
sudo rpm -e email-sentinel
```

### Test Docker Image

```bash
# Build Docker image locally
docker build -t email-sentinel:test .

# Run
docker run --rm email-sentinel:test --version

# Test with config
docker run -v $(pwd)/config:/home/sentinel/.config/email-sentinel \
  email-sentinel:test init
```

## Docker Deployment

### Docker Hub (Automated)

Releases automatically push to Docker Hub when you configure secrets:

- `DOCKER_USERNAME`
- `DOCKER_TOKEN`

Images are tagged as:
- `datateamsix/email-sentinel:latest`
- `datateamsix/email-sentinel:v1.0.0`
- `datateamsix/email-sentinel:1`
- `datateamsix/email-sentinel:1.0`

### Manual Docker Build

```bash
# Build
docker build -t datateamsix/email-sentinel:latest .

# Tag specific version
docker tag datateamsix/email-sentinel:latest datateamsix/email-sentinel:v1.0.0

# Push
docker push datateamsix/email-sentinel:latest
docker push datateamsix/email-sentinel:v1.0.0
```

### Docker Compose Deployment

```bash
# Clone repository or download docker-compose.yml
git clone https://github.com/datateamsix/email-sentinel.git
cd email-sentinel

# Create config directory
mkdir -p config

# Copy credentials.json to config/ directory
cp /path/to/credentials.json config/

# Start service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop service
docker-compose down
```

## Version Injection

Version information is injected at build time via ldflags:

```bash
# Manual build with version
go build -ldflags="\
  -X github.com/datateamsix/email-sentinel/internal/ui.AppVersion=1.0.0 \
  -X github.com/datateamsix/email-sentinel/internal/ui.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X github.com/datateamsix/email-sentinel/internal/ui.GitCommit=$(git rev-parse --short HEAD)" \
  -o email-sentinel .

# Verify
./email-sentinel --version
```

GoReleaser handles this automatically using git tags.

## Troubleshooting

### Release Workflow Fails

**Check:**
1. Go to Actions tab → Failed workflow → View logs
2. Common issues:
   - Missing GitHub secrets
   - Invalid `.goreleaser.yaml` syntax
   - Test failures
   - Network issues

**Fix:**
```bash
# Test locally first
goreleaser check
goreleaser release --snapshot --clean
```

### Homebrew Formula Update Fails

**Check:**
1. Verify `HOMEBREW_TAP_TOKEN` has correct permissions
2. Ensure `datateamsix/homebrew-tap` repository exists
3. Check workflow logs for error details

**Fix:**
```bash
# Manually update formula in homebrew-tap repo
# Edit Formula/email-sentinel.rb with new version and checksum
```

### Docker Image Not Pushed

**Check:**
1. Verify `DOCKER_USERNAME` and `DOCKER_TOKEN` are set
2. Check Docker Hub repository exists: `datateamsix/email-sentinel`
3. Verify token has write permissions

**Fix:**
```bash
# Login to Docker Hub manually
docker login

# Push manually
docker push datateamsix/email-sentinel:latest
```

### Package Installation Fails

**DEB Package:**
```bash
# Check dependencies
dpkg -I email-sentinel_*.deb

# Force install dependencies
sudo apt-get install -f
```

**RPM Package:**
```bash
# Check dependencies
rpm -qpR email-sentinel_*.rpm

# Install with dependency resolution
sudo yum install email-sentinel_*.rpm
```

### Binary "Not Signed" Warnings

**macOS:**
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine email-sentinel

# Or allow in System Preferences → Security & Privacy
```

**Windows:**
```powershell
# Right-click → Properties → Unblock
# Or run as administrator to install
```

**Long-term fix**: Purchase code signing certificates (see Sprint 5 in architecture doc)

## Release Checklist

Before releasing a new version:

- [ ] All tests passing locally: `go test ./...`
- [ ] Documentation updated (README, CHANGELOG)
- [ ] Version number follows semantic versioning
- [ ] GitHub secrets configured (HOMEBREW, SCOOP, DOCKER)
- [ ] Test snapshot build: `goreleaser release --snapshot --clean`
- [ ] Create and push tag: `git tag v1.0.0 && git push origin v1.0.0`
- [ ] Monitor GitHub Actions workflow
- [ ] Verify release artifacts on GitHub Releases page
- [ ] Test package installation (at least one platform)
- [ ] Update release notes with highlights
- [ ] Announce release

## Resources

- [GoReleaser Documentation](https://goreleaser.com)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Scoop Manifests](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests)
- [Docker Build Best Practices](https://docs.docker.com/develop/dev-best-practices/)

## Support

For deployment issues:
- Open an issue: https://github.com/datateamsix/email-sentinel/issues
- Check discussions: https://github.com/datateamsix/email-sentinel/discussions
