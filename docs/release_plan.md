# Release Plan — Email Sentinel v1.0.0

This document contains the release execution plan and checklist that was previously embedded in `.goreleaser.yaml`.

## Status
- Ready for Release
- Timeline: 4-6 hours total
- Date: December 2025

## Pre-Flight Checklist (high level)
- Core application fully functional
- Documentation streamlined and professional
- Gmail scope filtering feature implemented
- Top 5 use cases documented
- GoReleaser configuration exists
- GitHub Actions workflow configured
- Dockerfile created
- Cross-platform builds configured
- README.md polished

## Outstanding Items
- Homebrew Tap repository (datateamsix/homebrew-tap)
- Scoop Bucket repository (datateamsix/scoop-bucket)
- GitHub Secrets configuration
- Update .goreleaser.yaml (add app-config.yaml)
- Docker Hub setup and test
- GitHub Pages landing page
- Create v1.0.0 git tag
- Test all package managers

## Execution Plan (summary)
Phase 1: Repository setup — create Homebrew tap and Scoop bucket repositories, add tokens and secrets.

Phase 2: Update configuration — add `app-config.yaml` to archives, verify extra files exist, validate GoReleaser config.

Phase 3: CI and release — run `goreleaser check`, perform a dry-run snapshot, then tag and release.

## Notes
- Keep sensitive files out of the repo and remove any archived credentials before pushing.
- Use `goreleaser check` and a dry-run (`--snapshot --skip-publish`) to validate configuration before publishing.

