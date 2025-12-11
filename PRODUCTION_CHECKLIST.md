# Production Readiness Checklist

## Repository Status ✅
- [x] Repository initialized with Git
- [x] Remote configured to GitHub (`origin`)
- [x] Initial commit created (`v1.0.0` tag pushed)
- [x] Secrets removed from history (credentials archived and cleaned)
- [x] `.gitignore` and `.gitattributes` configured
- [x] Repository hygiene files created (CONTRIBUTING.md, SECURITY.md, CODE_OF_CONDUCT.md)

## Build & Release Configuration ✅
- [x] GoReleaser config valid (`.goreleaser.yaml`)
- [x] Explicit builds configured for Linux/Windows (amd64, arm64)
- [x] Archives configured to include README, LICENSE, docs, app-config.yaml
- [x] Checksums generation configured
- [x] Linux packages configured (DEB, RPM, APK via NFPM)
- [x] Scoop (Windows) package manifest configured
- [x] Homebrew (macOS) formula configured
- [x] Docker image publishing configured
- [x] Source archive generation enabled
- [x] Changelog generation configured

## CI/CD Workflows ✅
- [x] **ci.yml** - Runs on every push/PR to main/develop
  - Go fmt validation
  - Go vet linting
  - Unit tests (race detector enabled)
  - Staticcheck (optional)
  - Multi-OS snapshot builds (Ubuntu/macOS/Windows)
  - Artifact uploads
  
- [x] **release.yml** - Runs on tag push (`v*`)
  - Full test suite
  - GoReleaser full release (with publish)
  - Docker Hub login (requires secrets)
  - Artifact uploads

## Secrets & Environment Variables to Configure

### GitHub Secrets Required for Releases:
1. **DOCKER_USERNAME** - Docker Hub username
2. **DOCKER_PASSWORD** - Docker Hub token/password
3. **DOCKER_TOKEN** - (alternative) Docker token
4. **HOMEBREW_TAP_TOKEN** - GitHub token for homebrew-tap repo
5. **SCOOP_BUCKET_TOKEN** - GitHub token for scoop-bucket repo
6. **GITHUB_TOKEN** - Auto-provided by GitHub Actions

### Configuration Steps:
```bash
# 1. Go to GitHub repo Settings → Secrets and variables → Actions
# 2. Add each secret with its corresponding value
# 3. For DOCKER_PASSWORD, use a Docker access token (not password)
# 4. For HOMEBREW_TAP_TOKEN & SCOOP_BUCKET_TOKEN, use GitHub personal access tokens
```

## Files & Configuration Verified ✅
- [x] `main.go` - Entry point exists
- [x] `cmd/` - Command definitions present
- [x] `internal/` - Package structure in place
- [x] `Dockerfile` - Multi-stage build with VERSION arg
- [x] `docker-compose.yml` - Local development setup
- [x] `app-config.yaml` - Application config template
- [x] `README.md` - Updated with image reference
- [x] `images/central-email-warehousing.png` - Diagram included
- [x] `docs/release_plan.md` - Release plan documented
- [x] `.github/workflows/ci.yml` - CI workflow added
- [x] `.github/workflows/release.yml` - Release workflow updated

## Version & Tagging ✅
- [x] First stable release tagged as `v1.0.0`
- [x] Tag includes comprehensive release notes
- [x] Tag message includes features, installation instructions, use cases

## Pre-Release Checklist

Before pushing a new release tag (`v*`), verify:

### Code Quality
```bash
# Run locally:
go fmt ./...
go vet ./...
go test -v ./...
```

### Build Test
```bash
# Test snapshot locally:
goreleaser release --snapshot --clean
```

### Git Status
```bash
# Ensure clean working tree:
git status
git diff main origin/main  # No uncommitted changes
```

### Tag Creation
```bash
# Create annotated tag:
git tag -a v1.1.0 -m "Release v1.1.0 - description"

# Push tag (triggers release workflow):
git push origin v1.1.0
```

## Post-Release Steps

1. **Monitor GitHub Actions** - Check `release.yml` workflow run
2. **Verify GitHub Release** - Download and test binaries
3. **Check Docker Registry** - Verify image pushed to Docker Hub
4. **Validate Scoop Bucket** - Confirm manifest PR or merge
5. **Update Homebrew Tap** - Confirm formula PR or merge
6. **Announce Release** - Update docs, social media, etc.

## Ongoing Maintenance

### Weekly
- [ ] Review and merge dependabot/renovate updates
- [ ] Monitor GitHub Issues

### Per Release
- [ ] Run full CI/CD locally before tagging
- [ ] Update CHANGELOG.md
- [ ] Test release artifacts on multiple platforms
- [ ] Verify package manager installs work

## Deployment Notes

### Docker Deployment
```bash
docker pull datateamsix/email-sentinel:v1.0.0
docker run --rm -it datateamsix/email-sentinel --help
```

### macOS (Homebrew)
```bash
brew tap datateamsix/tap
brew install email-sentinel
```

### Windows (Scoop)
```powershell
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel
```

### Linux (DEB/RPM)
- Download from GitHub Releases
- Install: `sudo dpkg -i email-sentinel_*.deb` or `sudo rpm -i email-sentinel_*.rpm`

---

**Status**: Production-ready ✅
**Last Updated**: December 10, 2025
**Next Release**: v1.0.1 (bug fixes) or v1.1.0 (features)
