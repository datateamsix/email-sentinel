# Quick Release Commands - Email Sentinel v1.0.0

**⚠️ WARNING: Only execute these commands after completing all items in PRODUCTION_RELEASE_CHECKLIST.md**

---

## Pre-Flight Check

```bash
# 1. Verify you're on main branch
git branch --show-current
# Expected: main

# 2. Ensure working directory is clean
git status
# Expected: "nothing to commit, working tree clean"

# 3. Pull latest changes
git pull origin main

# 4. Verify GoReleaser config
goreleaser check
# Expected: "1 configuration file(s) validated"

# 5. Run tests
go test -v ./...
# Expected: All tests pass
```

---

## Create and Push v1.0.0 Tag

```bash
# Create annotated tag with release message
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

# Verify tag created correctly
git tag -l -n9 v1.0.0

# Push tag to GitHub (THIS TRIGGERS THE RELEASE WORKFLOW!)
git push origin v1.0.0
```

---

## Monitor Release

```bash
# After pushing the tag, monitor the GitHub Actions workflow:
# URL: https://github.com/datateamsix/email-sentinel/actions

# The release workflow will:
# 1. Run tests
# 2. Build binaries for all platforms
# 3. Create DEB/RPM/APK packages
# 4. Build and push Docker images
# 5. Update Homebrew tap
# 6. Update Scoop bucket
# 7. Create GitHub release with artifacts
# 8. Generate changelog

# Estimated duration: 10-15 minutes
```

---

## Verify Release

After the workflow completes, verify all artifacts:

```bash
# 1. Check GitHub Release
# URL: https://github.com/datateamsix/email-sentinel/releases/tag/v1.0.0
# Verify all archives, packages, and checksums are present

# 2. Test Homebrew installation (macOS)
brew tap datateamsix/tap
brew install email-sentinel
email-sentinel --version

# 3. Test Scoop installation (Windows)
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel
email-sentinel --version

# 4. Test Docker image
docker pull datateamsix/email-sentinel:1.0.0
docker run --rm datateamsix/email-sentinel:1.0.0 --version
```

---

## Rollback (Emergency Only)

If critical issues are discovered immediately after release:

```bash
# 1. Mark release as pre-release (via GitHub UI)
# Edit release → Check "This is a pre-release" → Save

# 2. Delete tag locally
git tag -d v1.0.0

# 3. Delete tag remotely
git push origin :refs/tags/v1.0.0

# 4. Fix issues, then re-tag as v1.0.1
git tag -a v1.0.1 -m "Release v1.0.1 - Hotfix"
git push origin v1.0.1
```

---

## Post-Release Actions

```bash
# 1. Update release notes (via GitHub UI)
# Add any additional context, known issues, or migration guides

# 2. Create GitHub Discussion announcement
# Navigate to: https://github.com/datateamsix/email-sentinel/discussions
# Create announcement post

# 3. Monitor issues and discussions
# Respond to installation issues promptly
# Document common problems in FAQ
```

---

## Next Steps (v1.0.1, v1.1.0, etc.)

For future releases:

```bash
# Update version in future commits
# Tag next version
git tag -a v1.0.1 -m "Release v1.0.1 - Bug fixes"
git push origin v1.0.1

# Or for feature release
git tag -a v1.1.0 -m "Release v1.1.0 - New features"
git push origin v1.1.0
```

---

**Last Updated:** 2025-12-11
**For:** Email Sentinel v1.0.0 Initial Production Release
