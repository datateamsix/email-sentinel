# Email Sentinel v1.0.0 - Production Release Checklist

**Release Manager:** Production Engineering Team
**Target Release Date:** December 2025
**Status:** Ready for Production Release

---

## Executive Summary

This checklist provides a comprehensive guide for publishing Email Sentinel v1.0.0 to production. All critical components have been reviewed and are production-ready.

### ✅ Pre-Release Status
- [x] Core application fully functional
- [x] Documentation complete and professional
- [x] GoReleaser configuration validated
- [x] GitHub Actions workflows configured
- [x] Cross-platform builds configured
- [x] Docker image configuration ready
- [x] No sensitive credentials in repository
- [x] README.md polished and comprehensive
- [x] LICENSE file present (MIT)
- [x] app-config.yaml included in release archives

### 🔧 Configuration Review Results

**GoReleaser Configuration (.goreleaser.yaml)**
- ✅ Configuration validated with `goreleaser check`
- ⚠️ Deprecation warnings (non-blocking):
  - `dockers` → migrate to `dockers_v2` (future enhancement)
  - `brews` → migrate to `homebrew_casks` (future enhancement)
- ✅ Builds configured for:
  - Linux (amd64, arm64) - CGO disabled
  - Windows (amd64, arm64) - CGO disabled
  - macOS (amd64, arm64) - CGO enabled (for systray support)
- ✅ Packaging configured:
  - Homebrew formula (datateamsix/homebrew-tap)
  - Scoop manifest (datateamsix/scoop-bucket)
  - DEB/RPM/APK packages
  - Docker images (DockerHub)
  - Source archives

**GitHub Actions Workflows**
- ✅ `.github/workflows/release.yml` - Production release workflow
- ✅ `.github/workflows/ci.yml` - Continuous integration
- ✅ Runs on macOS (required for CGO macOS builds)
- ✅ Full test suite execution before release

---

## Phase 1: GitHub Repository Setup

### 1.1 Create Package Manager Repositories

**Homebrew Tap Repository**
```bash
# On GitHub, create new repository:
# Name: homebrew-tap
# Owner: datateamsix
# Visibility: Public
# Initialize: README only
```

**Scoop Bucket Repository**
```bash
# On GitHub, create new repository:
# Name: scoop-bucket
# Owner: datateamsix
# Visibility: Public
# Initialize: README only
```

### 1.2 Configure GitHub Secrets

Navigate to: `Settings → Secrets and variables → Actions → New repository secret`

**Required Secrets:**

| Secret Name | Description | How to Get |
|------------|-------------|------------|
| `HOMEBREW_TAP_TOKEN` | GitHub Personal Access Token for Homebrew Tap | Create at https://github.com/settings/tokens → Fine-grained token → Permissions: `Contents: Read/Write` → Repositories: `datateamsix/homebrew-tap` |
| `SCOOP_BUCKET_TOKEN` | GitHub Personal Access Token for Scoop Bucket | Create at https://github.com/settings/tokens → Fine-grained token → Permissions: `Contents: Read/Write` → Repositories: `datateamsix/scoop-bucket` |
| `DOCKER_USERNAME` | DockerHub username | Your DockerHub account username |
| `DOCKER_TOKEN` | DockerHub access token | Create at https://hub.docker.com/settings/security → New Access Token |
| `GITHUB_TOKEN` | GitHub Actions token | Automatically provided by GitHub Actions |

**Token Permissions Required:**
- Homebrew/Scoop tokens: `Contents: Read/Write` only
- Docker token: `Read, Write, Delete` permissions
- Expiration: Set to 1 year, calendar reminder to rotate

### 1.3 Verify DockerHub Setup

```bash
# Ensure DockerHub repository exists
# Repository: datateamsix/email-sentinel
# Visibility: Public
# Auto-build: Disabled (handled by GoReleaser)
```

---

## Phase 2: Pre-Release Validation

### 2.1 Run Local Validation

```powershell
# Validate GoReleaser configuration
goreleaser check

# Expected output:
#   • checking path=.goreleaser.yaml
#   • 1 configuration file(s) validated
#   • thanks for using GoReleaser!
```

### 2.2 Security Audit

```bash
# Verify no credentials committed
git log --all --full-history -- credentials.json
git log --all --full-history -- token.json
git log --all --full-history -- .env

# Should return empty - if not, use git-filter-repo to remove
```

### 2.3 Test Suite Execution

```bash
# Run full test suite
go test -v -race ./...

# Run linters
go vet ./...
gofmt -s -l .

# Expected: All tests pass, no lint errors
```

### 2.4 Documentation Review

**Files to verify:**
- [x] README.md - Complete with all features documented
- [x] docs/QUICKSTART_WINDOWS.md
- [x] docs/QUICKSTART_MACOS.md
- [x] docs/QUICKSTART_LINUX.md
- [x] docs/CLI_GUIDE.md
- [x] docs/mobile_ntfy_setup.md
- [x] docs/GMAIL_API_GCP_SETUP.md
- [x] LICENSE (MIT)
- [x] app-config.yaml (production template)

---

## Phase 3: Release Execution

### 3.1 Pre-Release Git Preparation

```bash
# Ensure working directory is clean
git status
# Expected: "nothing to commit, working tree clean"

# Ensure on main branch
git branch --show-current
# Expected: "main"

# Pull latest changes
git pull origin main

# Verify recent commits
git log --oneline -5
```

### 3.2 Create Git Tag (v1.0.0)

```bash
# Create annotated tag with release notes
git tag -a v1.0.0 -m "Release v1.0.0 - Initial Production Release

🎉 Email Sentinel v1.0.0 - Initial Production Release

Features:
- Gmail API monitoring with OAuth 2.0
- Smart filtering by sender, subject, and Gmail categories
- Desktop notifications (Windows/macOS/Linux)
- Mobile notifications via ntfy.sh
- OTP/2FA code extraction and clipboard support
- Digital account tracking (subscriptions, trials)
- AI email summaries (Gemini, Claude, OpenAI)
- System tray integration
- Filter expiration with grace period
- Priority rules (urgent keywords, VIP senders/domains)
- Alert history with Gmail links
- Cross-platform builds (Windows, macOS, Linux)

Installation:
- Homebrew (macOS): brew tap datateamsix/tap && brew install email-sentinel
- Scoop (Windows): scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket && scoop install email-sentinel
- DEB/RPM packages for Linux
- Docker: docker pull datateamsix/email-sentinel
- Source builds supported

Documentation: https://github.com/datateamsix/email-sentinel
License: MIT"

# Verify tag created
git tag -l -n9 v1.0.0

# Push tag to GitHub (this triggers release workflow)
git push origin v1.0.0
```

### 3.3 Monitor Release Workflow

```bash
# Navigate to GitHub Actions
# URL: https://github.com/datateamsix/email-sentinel/actions

# Monitor "Release" workflow
# Triggered by: tag push v1.0.0
# Runner: macos-latest (required for CGO builds)

# Expected stages:
# 1. Checkout code ✓
# 2. Fetch all tags ✓
# 3. Set up Go ✓
# 4. Run tests ✓
# 5. Log in to Docker Hub ✓
# 6. Run GoReleaser ✓
# 7. Upload artifacts ✓
# 8. Release notification ✓

# Estimated duration: 10-15 minutes
```

### 3.4 Verify Release Artifacts

**GitHub Release Page**
```
URL: https://github.com/datateamsix/email-sentinel/releases/tag/v1.0.0

Expected artifacts:
✅ email-sentinel_1.0.0_Windows_amd64.zip
✅ email-sentinel_1.0.0_Windows_arm64.zip
✅ email-sentinel_1.0.0_Linux_amd64.tar.gz
✅ email-sentinel_1.0.0_Linux_arm64.tar.gz
✅ email-sentinel_1.0.0_macOS_amd64.tar.gz
✅ email-sentinel_1.0.0_macOS_arm64.tar.gz
✅ email-sentinel_1.0.0.deb
✅ email-sentinel_1.0.0.rpm
✅ email-sentinel_1.0.0.apk
✅ email-sentinel_1.0.0_source.tar.gz
✅ checksums.txt
✅ Release notes (auto-generated from commits)
```

**Homebrew Tap**
```bash
# Verify formula published
# URL: https://github.com/datateamsix/homebrew-tap

# Check commits by goreleaserbot
# File: Formula/email-sentinel.rb

# Test installation
brew tap datateamsix/tap
brew install email-sentinel
email-sentinel --version
# Expected: email-sentinel version 1.0.0
```

**Scoop Bucket**
```powershell
# Verify manifest published
# URL: https://github.com/datateamsix/scoop-bucket

# Check commits by goreleaserbot
# File: email-sentinel.json

# Test installation
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel
email-sentinel --version
# Expected: email-sentinel version 1.0.0
```

**Docker Hub**
```bash
# Verify images published
# URL: https://hub.docker.com/r/datateamsix/email-sentinel

# Expected tags:
# - latest
# - 1.0.0
# - 1
# - 1.0

# Test image
docker pull datateamsix/email-sentinel:latest
docker run --rm datateamsix/email-sentinel:latest --version
# Expected: email-sentinel version 1.0.0
```

---

## Phase 4: Post-Release Validation

### 4.1 Installation Testing

**macOS (Homebrew)**
```bash
brew tap datateamsix/tap
brew install email-sentinel
email-sentinel --version
email-sentinel --help
```

**Windows (Scoop)**
```powershell
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel
email-sentinel --version
email-sentinel --help
```

**Linux (DEB - Ubuntu/Debian)**
```bash
# Download from releases
wget https://github.com/datateamsix/email-sentinel/releases/download/v1.0.0/email-sentinel_1.0.0_amd64.deb
sudo dpkg -i email-sentinel_1.0.0_amd64.deb
email-sentinel --version
```

**Linux (RPM - Fedora/RHEL)**
```bash
# Download from releases
wget https://github.com/datateamsix/email-sentinel/releases/download/v1.0.0/email-sentinel_1.0.0_amd64.rpm
sudo rpm -i email-sentinel_1.0.0_amd64.rpm
email-sentinel --version
```

**Docker**
```bash
docker pull datateamsix/email-sentinel:1.0.0
docker run --rm datateamsix/email-sentinel:1.0.0 --version
```

### 4.2 Functional Smoke Tests

**Basic functionality test:**
```bash
# 1. Version check
email-sentinel --version

# 2. Help command
email-sentinel --help

# 3. Config show (should show defaults)
email-sentinel config show

# 4. Filter list (should be empty initially)
email-sentinel filter list

# All commands should execute without errors
```

### 4.3 Update Documentation Links

**Verify working links in README.md:**
- [x] Release badges point to v1.0.0
- [x] Download links point to latest release
- [x] Documentation links resolve correctly
- [x] Installation instructions accurate

---

## Phase 5: Communication & Monitoring

### 5.1 Release Announcement

**GitHub Discussions**
```markdown
Title: 🎉 Email Sentinel v1.0.0 Released!

We're excited to announce the first production release of Email Sentinel!

Email Sentinel is a cross-platform CLI tool that monitors Gmail and sends real-time desktop and mobile notifications when emails match custom filters.

🔥 Key Features:
- Smart email filtering by sender, subject, and Gmail categories
- Desktop + mobile notifications
- OTP/2FA code extraction
- Digital account tracking (subscriptions, trials)
- AI email summaries (optional)
- System tray integration
- Cross-platform (Windows, macOS, Linux)

📦 Installation:
- Homebrew (macOS): brew tap datateamsix/tap && brew install email-sentinel
- Scoop (Windows): scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket && scoop install email-sentinel
- Linux: DEB/RPM packages available
- Docker: docker pull datateamsix/email-sentinel

📚 Documentation: https://github.com/datateamsix/email-sentinel
⬇️ Download: https://github.com/datateamsix/email-sentinel/releases/tag/v1.0.0

Thank you to everyone who contributed!
```

### 5.2 Monitoring Plan

**First 24 Hours:**
- Monitor GitHub Issues for installation problems
- Watch GitHub Discussions for questions
- Check download counts on release page
- Monitor DockerHub pull statistics

**First Week:**
- Review installation success rates
- Collect feedback on documentation clarity
- Monitor for bug reports
- Track platform-specific issues

---

## Phase 6: Rollback Plan (If Needed)

### 6.1 Emergency Rollback Procedure

If critical issues are discovered:

```bash
# 1. Mark release as pre-release
# GitHub UI: Edit release → Check "This is a pre-release"

# 2. Add warning to README
git checkout main
# Add warning banner to README.md
git commit -m "docs: add v1.0.0 known issues warning"
git push origin main

# 3. Delete broken tag (if necessary)
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# 4. Create hotfix release (v1.0.1)
# Fix issues, test thoroughly, tag v1.0.1
```

### 6.2 Known Issues Template

```markdown
## Known Issues in v1.0.0

### Critical
- None identified

### Non-Critical
- macOS builds require Xcode Command Line Tools for source builds
- Windows local snapshot builds require clang (resolved in GitHub Actions)

### Workarounds
- Use package managers (Homebrew, Scoop) for simplified installation
```

---

## Appendix A: Required GitHub Secrets

### Creating Fine-Grained Tokens

**Homebrew Tap Token:**
1. Go to: https://github.com/settings/tokens?type=beta
2. Click "Generate new token"
3. Token name: `HOMEBREW_TAP_TOKEN (email-sentinel)`
4. Expiration: 1 year
5. Repository access: Only select repositories → `datateamsix/homebrew-tap`
6. Permissions:
   - Contents: Read and write
7. Click "Generate token"
8. Copy token and save to GitHub Secrets as `HOMEBREW_TAP_TOKEN`

**Scoop Bucket Token:**
1. Go to: https://github.com/settings/tokens?type=beta
2. Click "Generate new token"
3. Token name: `SCOOP_BUCKET_TOKEN (email-sentinel)`
4. Expiration: 1 year
5. Repository access: Only select repositories → `datateamsix/scoop-bucket`
6. Permissions:
   - Contents: Read and write
7. Click "Generate token"
8. Copy token and save to GitHub Secrets as `SCOOP_BUCKET_TOKEN`

**Docker Hub Token:**
1. Go to: https://hub.docker.com/settings/security
2. Click "New Access Token"
3. Description: `email-sentinel-goreleaser`
4. Permissions: Read, Write, Delete
5. Click "Generate"
6. Copy token and save to GitHub Secrets as `DOCKER_TOKEN`

---

## Appendix B: Quick Reference Commands

### Tag and Release
```bash
# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Delete tag (if needed)
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```

### Local Testing
```bash
# Validate config
goreleaser check

# Dry run (snapshot build)
goreleaser build --snapshot --clean --skip=validate

# Full dry run
goreleaser release --snapshot --skip=publish --clean
```

### Installation Testing
```bash
# Homebrew
brew tap datateamsix/tap
brew install email-sentinel

# Scoop
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel

# Docker
docker pull datateamsix/email-sentinel:latest
```

---

## Appendix C: Deprecation Warnings Resolution (Future)

**Current Warnings (Non-Blocking):**

1. **`dockers` → `dockers_v2`**
   ```yaml
   # Future migration (not required for v1.0.0)
   # Replace dockers: with dockers_v2:
   # See: https://goreleaser.com/deprecations#dockers
   ```

2. **`brews` → `homebrew_casks`**
   ```yaml
   # Future migration (not required for v1.0.0)
   # Note: This only applies if shipping a macOS .app bundle
   # CLI tools can continue using `brews`
   # See: https://goreleaser.com/deprecations#brews
   ```

**Action:** Schedule for v1.1.0 release cycle

---

## Final Checklist

Before executing Phase 3 (tagging v1.0.0):

- [ ] All Phase 1 repositories created
- [ ] All Phase 1 GitHub secrets configured
- [ ] All Phase 2 validations passed
- [ ] README badges point to correct repository
- [ ] DockerHub repository exists and is public
- [ ] Team notification sent about impending release
- [ ] Production Engineering sign-off obtained

**Sign-off:**
- [ ] Technical Lead: _______________
- [ ] Production Engineering: _______________
- [ ] Quality Assurance: _______________

---

**Document Version:** 1.0
**Last Updated:** 2025-12-11
**Next Review:** After v1.0.0 release completion
