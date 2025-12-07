# 📧 Email Sentinel

A cross-platform CLI tool that monitors your Gmail inbox and sends real-time notifications when emails match your custom filters. Get instant alerts on desktop and mobile without constantly checking your inbox.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)](https://github.com/yourusername/email-sentinel)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 🎯 Problem Statement

Waiting for important emails (job opportunities, client responses, urgent messages) means constantly checking your inbox — a major productivity killer. Email Sentinel solves this by:

- **Monitoring Gmail** in the background via the Gmail API
- **Filtering emails** by sender address or subject line keywords
- **Pushing notifications** to your desktop and mobile device instantly

## ✨ Features

### Core Functionality
- **Cross-Platform**: Single binary for Windows, macOS, and Linux
- **Flexible Filters**: Match by sender, subject keywords, or both
- **AND/OR Logic**: Configure whether all conditions must match or any condition triggers
- **Smart Priority Rules**: YAML-based rules engine for automatic priority classification
- **Low Resource**: Lightweight polling with configurable intervals
- **Secure**: OAuth 2.0 authentication, credentials stored locally
- **No Cost**: Uses Gmail API free tier (1B quota units/day)

### Notifications
- **Desktop Notifications**: Native OS notifications (Windows, macOS, Linux)
- **Windows Toast Notifications**: Rich, clickable notifications in Action Center with Gmail links
- **Mobile Push**: Free push notifications to iPhone/Android via [ntfy.sh](https://ntfy.sh)
- **System Tray**: Background app with tray icon showing recent alerts (Windows/macOS/Linux)

### Alert Management
- **Alert History**: SQLite database stores all alerts with automatic daily cleanup
- **View Past Alerts**: Review missed notifications with `email-sentinel alerts`
- **Direct Gmail Links**: Click any alert to open the email directly in Gmail
- **Priority Indicators**: Visual distinction between normal and urgent emails (🔥 vs 📧)

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Commands](#commands)
- [Filter Examples](#filter-examples)
- [Priority Rules](#priority-rules)
- [Alert History](#alert-history)
- [System Tray](#system-tray)
- [Mobile Notifications](#mobile-notifications)
- [Architecture](#architecture)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## 🔧 Prerequisites

### Required

1. **Go 1.22+** - [Download](https://go.dev/dl/)
2. **Google Cloud Platform Account** - [Console](https://console.cloud.google.com/)
3. **Gmail Account** - The account you want to monitor

### Google Cloud Setup

1. Create a new project in [Google Cloud Console](https://console.cloud.google.com/)
2. Enable the **Gmail API**:
   - Navigate to: APIs & Services → Library
   - Search for "Gmail API" → Enable
3. Configure **OAuth Consent Screen**:
   - Navigate to: APIs & Services → OAuth consent screen
   - User type: External
   - App name: `email-sentinel`
   - Scopes: Add `https://www.googleapis.com/auth/gmail.readonly`
   - Test users: Add your Gmail address
4. Create **OAuth Credentials**:
   - Navigate to: APIs & Services → Credentials
   - Create Credentials → OAuth Client ID
   - Application type: Desktop app
   - Download JSON → Save as `credentials.json`

## 📦 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/email-sentinel.git
cd email-sentinel

# Build for your platform
go build -o email-sentinel .

# Or build for all platforms
./scripts/build-all.sh
```

### Cross-Compilation

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o email-sentinel.exe .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o email-sentinel-mac .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o email-sentinel-mac-arm .

# Linux
GOOS=linux GOARCH=amd64 go build -o email-sentinel-linux .
```

## 🚀 Quick Start

```bash
# 1. Place credentials.json in project root or config directory

# 2. Initialize and authenticate
./email-sentinel init

# 3. Add a filter
./email-sentinel filter add --name "Job Alerts" --from "linkedin.com,greenhouse.io"

# 4. Test notifications
./email-sentinel test desktop
./email-sentinel test toast          # Windows only - test Action Center

# 5. Start monitoring with system tray
./email-sentinel start --tray        # Recommended for background use

# 6. View alert history
./email-sentinel alerts

# 7. Install auto-startup (optional)
./email-sentinel install
```

**For detailed setup instructions**, see **[Complete Setup Guide](docs/build_to_first_alert.md)** - comprehensive guide for Windows, macOS, and Linux.

**Want to monitor multiple email accounts?** See **[Central Email Setup](docs/central_email_setup.md)** - forward emails from all accounts to one central Gmail.

**Having issues on Windows?** See **[WINDOWS_INSTALL_TROUBLESHOOTING.md](WINDOWS_INSTALL_TROUBLESHOOTING.md)**

## ⚙️ Configuration

### Config File Location

| OS      | Path                                                    |
|---------|---------------------------------------------------------|
| Windows | `%APPDATA%\email-sentinel\config.yaml`                  |
| macOS   | `~/Library/Application Support/email-sentinel/config.yaml` |
| Linux   | `~/.config/email-sentinel/config.yaml`                  |

### Config Structure

```yaml
polling_interval: 45  # Seconds between Gmail checks

filters:
  - name: "Job Opportunities"
    from:
      - "@linkedin.com"
      - "greenhouse.io"
      - "lever.co"
    subject:
      - "interview"
      - "application"
      - "opportunity"
    match: any  # "any" (OR) or "all" (AND)

  - name: "Client Emails"
    from:
      - "client@company.com"
    subject: []
    match: any

notifications:
  desktop: true
  mobile:
    enabled: true
    ntfy_topic: "your-secret-topic-name"
```

### Filter Match Modes

| Mode  | Behavior                                          | Use Case                              |
|-------|---------------------------------------------------|---------------------------------------|
| `any` | Triggers if sender OR subject matches             | Broader matching, more notifications  |
| `all` | Triggers only if sender AND subject both match    | Precise matching, fewer notifications |

## 📖 Commands

### `init`

Initialize email-sentinel with Gmail OAuth authentication.

```bash
email-sentinel init
```

- Locates `credentials.json`
- Opens browser for Google OAuth
- Saves authentication token
- Displays existing filters and next steps

### `filter add`

Add a new email filter.

```bash
# Interactive mode
email-sentinel filter add

# With flags
email-sentinel filter add \
  --name "Job Alerts" \
  --from "linkedin.com,greenhouse.io" \
  --subject "interview,opportunity" \
  --match any
```

| Flag        | Short | Description                           |
|-------------|-------|---------------------------------------|
| `--name`    | `-n`  | Filter name (required)                |
| `--from`    | `-f`  | Sender patterns, comma-separated      |
| `--subject` | `-s`  | Subject patterns, comma-separated     |
| `--match`   | `-m`  | Match mode: `any` or `all` (default: `any`) |

### `filter list`

Display all configured filters.

```bash
email-sentinel filter list
```

### `filter edit`

Edit an existing filter.

```bash
# Interactive selection
email-sentinel filter edit

# By name
email-sentinel filter edit "Job Alerts"
```

### `filter remove`

Remove a filter.

```bash
# Interactive selection
email-sentinel filter remove

# By name
email-sentinel filter remove "Job Alerts"
```

### `start`

Start monitoring Gmail for matching emails.

```bash
# Foreground (logs to stdout)
email-sentinel start

# With system tray icon (recommended)
email-sentinel start --tray

# As daemon (background)
email-sentinel start --daemon
```

**Flags:**
- `--tray` / `-t` - Run with system tray icon
- `--daemon` / `-d` - Run as background daemon

### `stop`

Stop the background daemon.

```bash
email-sentinel stop
```

### `status`

Check if email-sentinel is running.

```bash
email-sentinel status
```

### `alerts`

View email alert history.

```bash
# View all alerts from today
email-sentinel alerts

# View last N alerts
email-sentinel alerts --recent 5
```

Shows:
- Priority indicators (🔥 for urgent, 📧 for normal)
- Timestamp, sender, subject
- Gmail permalink for each alert

### `test`

Test notification systems.

```bash
# Test desktop notification
email-sentinel test desktop

# Test Windows toast notification
email-sentinel test toast
email-sentinel test toast --priority    # Test urgent notification

# Test mobile notification
email-sentinel test mobile

# Test filter matching
email-sentinel test filter "Filter Name" "sender@example.com" "subject text"
```

### `config`

View or modify configuration.

```bash
# Show current config
email-sentinel config show

# Set polling interval
email-sentinel config set polling 30
```

### `install`

Install email-sentinel to run automatically on startup.

```bash
# Install auto-startup
email-sentinel install

# Preview what would be installed
email-sentinel install --show
```

**Platform support:**
- **Windows**: Creates Task Scheduler task (runs at logon)
- **macOS**: Creates LaunchAgent plist
- **Linux**: Creates systemd user service

### `uninstall`

Remove email-sentinel from automatic startup.

```bash
email-sentinel uninstall
```

## 📝 Filter Examples

### Job Search

```bash
email-sentinel filter add \
  --name "Job Applications" \
  --from "linkedin.com,greenhouse.io,lever.co,indeed.com" \
  --subject "interview,application,opportunity,position" \
  --match any
```

### Specific Sender

```bash
email-sentinel filter add \
  --name "From Boss" \
  --from "boss@company.com"
```

### Subject Keywords Only

```bash
email-sentinel filter add \
  --name "Urgent Emails" \
  --subject "urgent,asap,emergency,critical"
```

### Precise Match (AND)

```bash
email-sentinel filter add \
  --name "Client Invoices" \
  --from "billing@client.com" \
  --subject "invoice,payment" \
  --match all
```

## 🎯 Priority Rules

Email Sentinel includes a smart rules engine that automatically classifies emails as normal or high priority.

### Configuration

Priority rules are stored in `rules.yaml` (auto-created on first run):

**Location:**
- Windows: `%APPDATA%\email-sentinel\rules.yaml`
- macOS: `~/Library/Application Support/email-sentinel/rules.yaml`
- Linux: `~/.config/email-sentinel/rules.yaml`

### Example Rules

```yaml
priority_rules:
  # Keywords that mark emails as urgent
  urgent_keywords:
    - urgent
    - asap
    - action required
    - deadline
    - invoice
    - payment
    - critical

  # Specific email addresses that are always high priority
  vip_senders:
    - boss@company.com
    - ceo@company.com

  # Entire domains that are always high priority
  vip_domains:
    - importantclient.com
    - criticalvendor.com

notification_settings:
  # Future: Quiet hours, weekend mode
  quiet_hours_start: "22:00"
  quiet_hours_end: "08:00"
```

### How It Works

When an email matches a filter, the priority engine evaluates:

1. **Urgent Keywords** - Subject/snippet contains keywords → Priority 1
2. **VIP Senders** - Exact email match → Priority 1
3. **VIP Domains** - Sender's domain matches → Priority 1
4. Otherwise → Priority 0 (normal)

**Priority Indicators:**
- 🔥 High priority (1) - Red icon, urgent notification sound
- 📧 Normal priority (0) - Standard icon and sound

## 📜 Alert History

All email notifications are automatically saved to a local SQLite database.

### View Alerts

```bash
# View today's alerts
email-sentinel alerts

# View last 5 alerts
email-sentinel alerts --recent 5
```

### Example Output

```
📬 Today's Alerts (3 total)

[1] 🔥 2025-12-06 14:30:15
    Filter: VIP Senders
    Priority: HIGH
    From:   boss@company.com
    Subject: URGENT: Server is down!
    Preview: We need to address this immediately...
    Link:   https://mail.google.com/mail/u/0/#all/abc123

[2] 📧 2025-12-06 13:45:22
    Filter: Job Alerts
    From:   recruiter@linkedin.com
    Subject: New job opportunity
    Link:   https://mail.google.com/mail/u/0/#all/def456
```

### Features

- **Persistent Storage** - Survives restarts, view missed alerts
- **Clickable Links** - Direct Gmail permalinks
- **Auto-Cleanup** - Alerts older than midnight are deleted automatically
- **Priority Indicators** - Visual distinction between urgent and normal

## 📱 System Tray

Run Email Sentinel as a background app with a system tray icon.

### Start with Tray

```bash
email-sentinel start --tray
```

### Features

**Tray Icon:**
- Normal state: Default email icon
- Urgent alert: Icon switches to red/orange
- Tooltip shows monitoring status

**Menu Options:**
- **Recent Alerts** - Submenu showing last 10 alerts
  - Click any alert to open in Gmail
  - Shows time, priority icon, and subject
- **Open History** - View all alerts in terminal
- **Quit** - Stop monitoring and exit

### Example Menu

```
Email Sentinel
├─ Recent Alerts
│  ├─ 🔥 [14:30] URGENT: Server is down!
│  ├─ 📧 [14:15] Weekly team meeting
│  ├─ 📧 [13:45] New feature request
│  └─ ...
├─ ────────────
├─ Open History
├─ ────────────
└─ Quit
```

### Behavior

- **New Alert Arrives** → Icon flashes urgent color for 5 seconds
- **Click Alert** → Opens Gmail in browser
- **Auto-Refresh** → Menu updates every 30 seconds
- **Background Mode** → No terminal window needed

### Platform Support

- ✅ Windows - System tray (notification area)
- ✅ macOS - Menu bar
- ✅ Linux - System tray (GNOME, KDE, etc.)

## 📱 Mobile Notifications

Email Sentinel uses [ntfy.sh](https://ntfy.sh) for free mobile push notifications.

### Setup

1. **Install ntfy app**:
   - [iOS App Store](https://apps.apple.com/app/ntfy/id1625396347)
   - [Android Play Store](https://play.google.com/store/apps/details?id=io.heckel.ntfy)

2. **Subscribe to your topic**:
   - Open ntfy app
   - Subscribe to a unique topic name (e.g., `michaels-email-alerts-x7k2`)
   - Use a random/unique name for privacy

3. **Configure email-sentinel**:
   ```bash
   email-sentinel config set ntfy_topic "your-topic-name"
   email-sentinel config set mobile_enabled true
   ```

### How It Works

```
[Email Matches Filter]
        ↓
[HTTP POST to ntfy.sh]
        ↓
[Push to your phone]
```

No account required. Free. Open source.

## 🏗️ Architecture

### Project Structure

```
email-sentinel/
├── cmd/                      # CLI commands (Cobra)
│   ├── root.go               # Base command + global flags
│   ├── init.go               # OAuth initialization
│   ├── start.go              # Start watcher daemon
│   ├── stop.go               # Stop daemon
│   ├── status.go             # Check daemon status
│   ├── config.go             # Config management
│   ├── filter.go             # Filter parent command
│   ├── add.go                # filter add
│   ├── list.go               # filter list
│   ├── edit.go               # filter edit
│   └── remove.go             # filter remove
│
├── internal/                 # Private packages
│   ├── config/
│   │   ├── paths.go          # OS-specific config paths
│   │   └── manager.go        # Config load/save
│   │
│   ├── filter/
│   │   ├── types.go          # Filter struct definitions
│   │   └── engine.go         # Filter matching logic
│   │
│   ├── gmail/
│   │   ├── auth.go           # OAuth flow
│   │   ├── client.go         # Gmail API wrapper
│   │   └── message.go        # Message parsing
│   │
│   ├── notify/
│   │   ├── desktop.go        # OS notifications
│   │   └── mobile.go         # ntfy.sh integration
│   │
│   └── state/
│       └── seen.go           # Track processed messages
│
├── scripts/
│   ├── build-all.sh          # Cross-compilation script
│   └── install-startup.sh    # Register as startup app
│
├── main.go                   # Entry point
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── .gitignore
```

### Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                         EMAIL SENTINEL                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────────┐    │
│   │ Gmail API   │───▶│ Filter      │───▶│ Notifications   │    │
│   │ (polling)   │    │ Engine      │    │                 │    │
│   └─────────────┘    └─────────────┘    │ ├─ Desktop      │    │
│         │                   │           │ └─ Mobile       │    │
│         │                   │           └─────────────────┘    │
│         ▼                   ▼                                  │
│   ┌─────────────┐    ┌─────────────┐                           │
│   │ OAuth Token │    │ Config YAML │                           │
│   └─────────────┘    └─────────────┘                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Key Dependencies

| Package                           | Purpose                    |
|-----------------------------------|----------------------------|
| `github.com/spf13/cobra`          | CLI framework              |
| `github.com/spf13/viper`          | Config management          |
| `google.golang.org/api/gmail/v1`  | Gmail API client           |
| `golang.org/x/oauth2/google`      | OAuth 2.0                  |
| `github.com/gen2brain/beeep`      | Desktop notifications      |
| `gopkg.in/yaml.v3`                | YAML parsing               |

## 🛠️ Development

### Setup

```bash
# Clone
git clone https://github.com/yourusername/email-sentinel.git
cd email-sentinel

# Install dependencies
go mod download

# Run locally
go run . --help

# Build
go build -o email-sentinel .

# Run tests
go test ./...
```

### Adding a New Command

```bash
# Install cobra-cli
go install github.com/spf13/cobra-cli@latest

# Add command
cobra-cli add mycommand

# Add subcommand
cobra-cli add subcommand -p parentCmd
```

### Code Style

```bash
# Format
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run

# Vet
go vet ./...
```

## ❓ Troubleshooting

### "App not verified" warning during OAuth

This is expected for personal/development apps. Click:
1. **Advanced**
2. **Go to email-sentinel (unsafe)**

Your app only accesses your own Gmail — this is safe.

### "Access blocked" error

Add your email as a test user:
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. APIs & Services → OAuth consent screen
3. Test users → Add Users
4. Enter your Gmail address

### Token expired

Re-run initialization:
```bash
email-sentinel init
```

### Notifications not appearing (Windows)

Ensure notifications are enabled in Windows Settings:
- Settings → System → Notifications
- Enable notifications for your terminal/app

### Notifications not appearing (macOS)

Grant notification permissions:
- System Preferences → Notifications
- Find email-sentinel or your terminal app
- Enable "Allow Notifications"

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Multi-Account Monitoring

Email Sentinel supports monitoring multiple email accounts through a **central collector inbox** approach:

- Set up email forwarding from all your accounts (Gmail, Outlook, Yahoo, iCloud, etc.)
- Monitor one central Gmail account with Email Sentinel
- Forwarded messages retain original sender and subject metadata
- See detailed setup guide: **[Central Email Setup](docs/central_email_setup.md)**

### Ideas for Contribution

- [ ] Support for Outlook/Microsoft Graph API
- [ ] Native multi-account OAuth support
- [ ] Web UI for filter management
- [ ] Email notification summaries/digests
- [ ] Custom notification sounds
- [ ] Slack/Discord webhook integration
- [ ] Filter import/export

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [ntfy.sh](https://ntfy.sh) - Free push notification service
- [Google Gmail API](https://developers.google.com/gmail/api) - Email access

---

**Note for AI Agents**: This project uses Go modules. The module path is defined in `go.mod`. All internal packages use the `internal/` directory convention, preventing external imports. Configuration is stored in OS-appropriate locations (see `internal/config/paths.go`). The filter matching logic is in `internal/filter/engine.go`. OAuth tokens are stored separately from app credentials for security.