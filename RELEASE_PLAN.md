# Email Sentinel v1.0.0 - Public Release Plan

## Current Status

### ✅ Already Complete
- **Core Application**: Fully functional CLI with all features implemented
- **GoReleaser Configuration**: Comprehensive `.goreleaser.yaml` with multi-platform builds
- **GitHub Actions**: CI/CD pipeline configured in `.github/workflows/release.yml`
- **Documentation**: Extensive README.md with examples, guides, and API documentation
- **Packaging**: Linux packaging scripts (DEB, RPM, systemd service files)
- **Docker**: Multi-stage Dockerfile with health checks
- **Multi-Platform Support**: Windows, macOS (Intel + ARM), Linux (amd64, arm64)
- **Package Managers**: Homebrew and Scoop configuration ready

### 🚧 Needs Attention

#### 1. Missing Repository Dependencies
- **Homebrew Tap Repository**: `datateamsix/homebrew-tap` (referenced but not created)
- **Scoop Bucket Repository**: `datateamsix/scoop-bucket` (referenced but not created)

#### 2. GitHub Secrets Configuration
- `HOMEBREW_TAP_TOKEN` - Personal access token for Homebrew tap
- `SCOOP_BUCKET_TOKEN` - Personal access token for Scoop bucket
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_TOKEN` - Docker Hub access token

#### 3. Missing Files for Release
- `ai-config.yaml` needs to be added to archive files in `.goreleaser.yaml`
- `LICENSE` file verification
- Version tagging (no tags exist yet)

#### 4. GitHub Pages Landing Page
- No landing page currently exists
- Need to create static site or enable GitHub Pages

---

## Release Checklist

### Phase 1: Pre-Release Setup (1-2 hours)

#### 1.1 Create Package Manager Repositories

**Homebrew Tap:**
```bash
# 1. Create repository on GitHub: datateamsix/homebrew-tap
# 2. Initialize it with README
# 3. Generate Personal Access Token with 'repo' scope
# 4. Add HOMEBREW_TAP_TOKEN to GitHub Secrets
```

**Scoop Bucket:**
```bash
# 1. Create repository on GitHub: datateamsix/scoop-bucket
# 2. Initialize it with README
# 3. Generate Personal Access Token with 'repo' scope
# 4. Add SCOOP_BUCKET_TOKEN to GitHub Secrets
```

#### 1.2 Docker Hub Setup
```bash
# 1. Create Docker Hub account (if needed): hub.docker.com
# 2. Create repository: datateamsix/email-sentinel
# 3. Generate access token
# 4. Add DOCKER_USERNAME and DOCKER_TOKEN to GitHub Secrets
```

#### 1.3 Update GoReleaser Configuration
```yaml
# Add ai-config.yaml to archives section
archives:
  files:
    - README.md
    - LICENSE
    - docs/**/*
    - otp_rules.yaml
    - rules.yaml
    - ai-config.yaml  # ADD THIS LINE
```

#### 1.4 Verify License File
```bash
# Ensure LICENSE file exists and is up to date
# If using MIT, verify copyright year and owner
```

---

### Phase 2: Landing Page Creation (2-3 hours)

#### 2.1 GitHub Pages Setup

**Option A: Simple One-Page Site (Recommended for Quick Launch)**

Create `docs/index.html`:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Sentinel - Gmail Notification Monitor</title>
    <meta name="description" content="Monitor Gmail and get instant notifications on desktop and mobile">

    <!-- Open Graph / Social Media -->
    <meta property="og:type" content="website">
    <meta property="og:title" content="Email Sentinel">
    <meta property="og:description" content="Monitor Gmail and get instant notifications on desktop and mobile">
    <meta property="og:image" content="https://datateamsix.github.io/email-sentinel/images/logo.png">

    <!-- Styles -->
    <style>
        :root {
            --primary: #4F46E5;
            --secondary: #10B981;
            --dark: #1F2937;
            --light: #F3F4F6;
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            line-height: 1.6;
            color: var(--dark);
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
        }

        header {
            background: linear-gradient(135deg, var(--primary) 0%, #7C3AED 100%);
            color: white;
            padding: 60px 0;
            text-align: center;
        }

        header h1 {
            font-size: 3rem;
            margin-bottom: 20px;
        }

        header p {
            font-size: 1.25rem;
            opacity: 0.9;
            margin-bottom: 30px;
        }

        .cta-buttons {
            display: flex;
            gap: 20px;
            justify-content: center;
            flex-wrap: wrap;
        }

        .btn {
            display: inline-block;
            padding: 12px 30px;
            border-radius: 8px;
            text-decoration: none;
            font-weight: 600;
            transition: transform 0.2s;
        }

        .btn:hover {
            transform: translateY(-2px);
        }

        .btn-primary {
            background: white;
            color: var(--primary);
        }

        .btn-secondary {
            background: rgba(255, 255, 255, 0.2);
            color: white;
            border: 2px solid white;
        }

        .features {
            padding: 80px 0;
            background: white;
        }

        .features h2 {
            text-align: center;
            font-size: 2.5rem;
            margin-bottom: 60px;
        }

        .feature-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 40px;
        }

        .feature-card {
            text-align: center;
            padding: 30px;
            border-radius: 12px;
            background: var(--light);
        }

        .feature-icon {
            font-size: 3rem;
            margin-bottom: 20px;
        }

        .feature-card h3 {
            font-size: 1.5rem;
            margin-bottom: 15px;
        }

        .installation {
            padding: 80px 0;
            background: var(--light);
        }

        .installation h2 {
            text-align: center;
            font-size: 2.5rem;
            margin-bottom: 40px;
        }

        .install-tabs {
            display: flex;
            justify-content: center;
            gap: 10px;
            margin-bottom: 30px;
            flex-wrap: wrap;
        }

        .tab-btn {
            padding: 10px 20px;
            background: white;
            border: 2px solid var(--primary);
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            color: var(--primary);
        }

        .tab-btn.active {
            background: var(--primary);
            color: white;
        }

        .tab-content {
            display: none;
            background: var(--dark);
            color: white;
            padding: 30px;
            border-radius: 12px;
            font-family: 'Monaco', 'Courier New', monospace;
        }

        .tab-content.active {
            display: block;
        }

        pre {
            overflow-x: auto;
            white-space: pre-wrap;
        }

        footer {
            background: var(--dark);
            color: white;
            text-align: center;
            padding: 40px 0;
        }

        .social-links {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-top: 20px;
        }

        .social-links a {
            color: white;
            text-decoration: none;
            font-size: 1.5rem;
        }

        @media (max-width: 768px) {
            header h1 { font-size: 2rem; }
            .features h2, .installation h2 { font-size: 2rem; }
        }
    </style>
</head>
<body>
    <header>
        <div class="container">
            <h1>📧 Email Sentinel</h1>
            <p>Monitor Gmail and get instant notifications on desktop and mobile</p>
            <div class="cta-buttons">
                <a href="https://github.com/datateamsix/email-sentinel/releases/latest" class="btn btn-primary">
                    Download Latest Release
                </a>
                <a href="https://github.com/datateamsix/email-sentinel" class="btn btn-secondary">
                    View on GitHub
                </a>
            </div>
        </div>
    </header>

    <section class="features">
        <div class="container">
            <h2>✨ Features</h2>
            <div class="feature-grid">
                <div class="feature-card">
                    <div class="feature-icon">🔔</div>
                    <h3>Real-Time Alerts</h3>
                    <p>Get instant desktop and mobile notifications when emails match your filters</p>
                </div>

                <div class="feature-card">
                    <div class="feature-icon">🤖</div>
                    <h3>AI Summaries</h3>
                    <p>Optional AI-powered email summaries with questions and action items</p>
                </div>

                <div class="feature-card">
                    <div class="feature-icon">🔐</div>
                    <h3>OTP Detection</h3>
                    <p>Automatically extract and manage 2FA verification codes</p>
                </div>

                <div class="feature-card">
                    <div class="feature-icon">🏷️</div>
                    <h3>Smart Filters</h3>
                    <p>Organize with labels, priority rules, and flexible matching logic</p>
                </div>

                <div class="feature-card">
                    <div class="feature-icon">💻</div>
                    <h3>Cross-Platform</h3>
                    <p>Single binary for Windows, macOS, and Linux</p>
                </div>

                <div class="feature-card">
                    <div class="feature-icon">🔒</div>
                    <h3>Secure & Private</h3>
                    <p>OAuth 2.0 authentication, all data stored locally</p>
                </div>
            </div>
        </div>
    </section>

    <section class="installation">
        <div class="container">
            <h2>📦 Installation</h2>

            <div class="install-tabs">
                <button class="tab-btn active" onclick="showTab('mac')">macOS</button>
                <button class="tab-btn" onclick="showTab('windows')">Windows</button>
                <button class="tab-btn" onclick="showTab('linux')">Linux</button>
                <button class="tab-btn" onclick="showTab('docker')">Docker</button>
            </div>

            <div id="mac" class="tab-content active">
                <pre># Using Homebrew
brew tap datateamsix/tap
brew install email-sentinel

# Verify installation
email-sentinel --version</pre>
            </div>

            <div id="windows" class="tab-content">
                <pre># Using Scoop
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel

# Verify installation
email-sentinel --version</pre>
            </div>

            <div id="linux" class="tab-content">
                <pre># Debian/Ubuntu
wget https://github.com/datateamsix/email-sentinel/releases/latest/download/email-sentinel_*_amd64.deb
sudo dpkg -i email-sentinel_*_amd64.deb

# RHEL/Fedora/CentOS
wget https://github.com/datateamsix/email-sentinel/releases/latest/download/email-sentinel_*_x86_64.rpm
sudo rpm -i email-sentinel_*_x86_64.rpm

# Verify installation
email-sentinel --version</pre>
            </div>

            <div id="docker" class="tab-content">
                <pre># Pull image
docker pull datateamsix/email-sentinel:latest

# Run with mounted config
docker run -v ~/.email-sentinel:/root/.email-sentinel \
  datateamsix/email-sentinel:latest start --daemon</pre>
            </div>
        </div>
    </section>

    <footer>
        <div class="container">
            <p>&copy; 2025 DataTeam Six. Licensed under MIT.</p>
            <div class="social-links">
                <a href="https://github.com/datateamsix/email-sentinel">⭐ Star on GitHub</a>
                <a href="https://github.com/datateamsix/email-sentinel/issues">🐛 Report Issues</a>
            </div>
        </div>
    </footer>

    <script>
        function showTab(tabId) {
            // Hide all tabs
            document.querySelectorAll('.tab-content').forEach(tab => {
                tab.classList.remove('active');
            });
            document.querySelectorAll('.tab-btn').forEach(btn => {
                btn.classList.remove('active');
            });

            // Show selected tab
            document.getElementById(tabId).classList.add('active');
            event.target.classList.add('active');
        }
    </script>
</body>
</html>
```

**Option B: Jekyll/Hugo Static Site (For More Complex Site)**
- Use GitHub Pages with Jekyll theme
- Create blog posts for release announcements
- Add documentation pages

#### 2.2 Enable GitHub Pages
1. Go to repository Settings → Pages
2. Source: Deploy from branch `main` → `/docs` folder
3. Custom domain (optional): email-sentinel.com
4. Enforce HTTPS

---

### Phase 3: First Release (30 minutes)

#### 3.1 Create Release Tag
```bash
# Ensure all changes are committed
git add .
git commit -m "chore: Prepare for v1.0.0 release"
git push origin main

# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0 - First public release

Features:
- Cross-platform email monitoring (Windows, macOS, Linux)
- Real-time desktop and mobile notifications
- AI-powered email summaries (Claude, OpenAI, Gemini)
- OTP/2FA code detection and management
- Smart filters with labels and priority rules
- System tray integration
- Docker support
- Package manager support (Homebrew, Scoop)
"

git push origin v1.0.0
```

#### 3.2 Monitor Release Process
1. GitHub Actions will automatically trigger
2. GoReleaser will build all platform binaries
3. Docker images will be pushed to Docker Hub
4. Homebrew tap and Scoop bucket will be updated
5. GitHub Release will be created with artifacts

#### 3.3 Verify Release
- Check GitHub Releases page
- Download and test binaries for each platform
- Verify Homebrew formula works: `brew install datateamsix/tap/email-sentinel`
- Verify Scoop manifest works: `scoop install datateamsix/email-sentinel`
- Verify Docker image: `docker pull datateamsix/email-sentinel:latest`

---

### Phase 4: Post-Release (1-2 hours)

#### 4.1 Update Documentation
- Add "Latest Release" badge to README.md
- Update installation instructions if needed
- Create CHANGELOG.md from release notes

#### 4.2 Announcement
- Post on GitHub Discussions (create if needed)
- Share on relevant communities (Reddit, Hacker News, etc.)
- Create blog post (optional)

#### 4.3 Monitor Issues
- Watch for bug reports
- Respond to installation issues
- Gather feedback for v1.1.0

---

## Quick Launch Commands

### 1. Set up GitHub Secrets
```bash
# Generate Personal Access Tokens on GitHub:
# Settings → Developer settings → Personal access tokens → Tokens (classic)
# Create token with 'repo' scope

# Add secrets to repository:
# Repository → Settings → Secrets and variables → Actions → New repository secret

# Required secrets:
# - HOMEBREW_TAP_TOKEN
# - SCOOP_BUCKET_TOKEN
# - DOCKER_USERNAME
# - DOCKER_TOKEN
```

### 2. Create Package Repositories
```bash
# On GitHub, create two new repositories:
# 1. datateamsix/homebrew-tap (public)
# 2. datateamsix/scoop-bucket (public)

# Initialize each with a README:
echo "# Homebrew Tap for Email Sentinel" > README.md
git add README.md
git commit -m "Initial commit"
git push origin main
```

### 3. Update .goreleaser.yaml
```bash
# Add ai-config.yaml to archives.files section
# Commit and push the change
```

### 4. Create and Push Tag
```bash
git tag -a v1.0.0 -m "Release v1.0.0 - First public release"
git push origin v1.0.0
```

### 5. Enable GitHub Pages
```bash
# Go to repository Settings → Pages
# Source: main branch → /docs folder
# Save
```

---

## Success Criteria

- ✅ GitHub Release v1.0.0 published with all platform binaries
- ✅ Docker image available on Docker Hub
- ✅ Homebrew formula works: `brew install datateamsix/tap/email-sentinel`
- ✅ Scoop manifest works: `scoop install datateamsix/email-sentinel`
- ✅ Landing page live at https://datateamsix.github.io/email-sentinel/
- ✅ README badges showing release version and download count
- ✅ All CI/CD checks passing

---

## Timeline Estimate

- **Phase 1 (Pre-Release Setup)**: 1-2 hours
- **Phase 2 (Landing Page)**: 2-3 hours
- **Phase 3 (First Release)**: 30 minutes
- **Phase 4 (Post-Release)**: 1-2 hours

**Total**: 5-8 hours for complete public release

---

## Next Version Planning (v1.1.0)

Potential features based on user feedback:
- [ ] Web UI for filter management
- [ ] Outlook/Microsoft Graph API support
- [ ] Native multi-account OAuth support
- [ ] Custom notification sounds
- [ ] Slack/Discord webhook integration
- [ ] Email notification digests
- [ ] Browser extension integration
- [ ] Mobile app (React Native)

---

## Support & Maintenance

### Issue Triage
- Label issues: bug, enhancement, documentation, question
- Prioritize: critical, high, medium, low
- Assign milestones: v1.1.0, v1.2.0, etc.

### Security Policy
- Create SECURITY.md with vulnerability reporting process
- Set up GitHub security advisories
- Monitor dependencies with Dependabot

### Community
- Create CODE_OF_CONDUCT.md
- Create CONTRIBUTING.md
- Set up GitHub Discussions for Q&A
- Consider Discord server for community

---

## Current Repository Status

**Repository**: https://github.com/datateamsix/email-sentinel
**Current State**: Development (no releases yet)
**Code Status**: Production-ready
**Documentation**: Complete
**CI/CD**: Configured but not tested

**Next Steps**: Follow Phase 1 to prepare for first release
