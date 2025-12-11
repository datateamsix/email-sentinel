# 📧 Email Sentinel

![Email Sentinel Gopher](docs/images/go-night.svg)

**Stop wasting time digging for emails in your inbox. Get notified on the important stuff as it arrives.**

Email notifications are blunt instruments. Email Sentinel is your simple, no-frills workhorse scalpel for cutting through the BS and noise of email.

[![Latest Release](https://img.shields.io/github/v/release/datateamsix/email-sentinel)](https://github.com/datateamsix/email-sentinel/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/datateamsix/email-sentinel/total)](https://github.com/datateamsix/email-sentinel/releases)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)](https://github.com/datateamsix/email-sentinel)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## Why Email Sentinel?

**The Problem**: You're waiting for important emails—job responses, client messages, urgent alerts—so you compulsively check your inbox every 5 minutes or waste time digging through spam and noise. This destroys your focus and productivity.

**The Solution**: Email Sentinel monitors Gmail silently in the background and sends **instant notifications** only for emails matching your specific filters. Work in peace, knowing you won't miss what matters.

**5 Killer Use Cases:**
- 🎯 **Job Hunting** - Never miss interview invites or recruiter messages
- 🔐 **OTP Codes** - Auto-extract 2FA codes, copy to clipboard instantly
- 👔 **VIP Alerts** - Get notified the second your boss or key clients email
- 🚨 **Spam Rescue** - Catch important emails misfiled in spam/promotions
- ⚡ **Urgent Keywords** - Auto-detect "ASAP", "urgent", deadlines across all senders

---

## ✨ What It Does

- 📬 **Monitors Gmail** via API (no IMAP polling)
- 🎯 **Smart Filtering** by sender, subject, or both
- 🔔 **Desktop + Mobile Notifications** (Windows/macOS/Linux + ntfy.sh)
- 🤖 **AI Email Summaries** (optional, with Claude/GPT/Gemini)
- 🔐 **OTP/2FA Code Extraction** (copy codes instantly)
- 📊 **Alert History** with Gmail links
- 🪟 **System Tray App** (runs in background)
- 📱 **Gmail Category Scopes** (inbox, social, promotions, etc.)
- 🏷️ **Filter Labels** (organize by work, urgent, personal)
- ⚡ **Priority Rules** (auto-detect urgent emails)

---

## 🚀 Quick Start

### 1. Install

**Windows (Scoop):**
```powershell
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel
```

**macOS (Homebrew):**
```bash
brew tap datateamsix/tap
brew install email-sentinel
```

**Linux / From Source:**
```bash
# Download from releases
wget https://github.com/datateamsix/email-sentinel/releases/latest/download/email-sentinel_linux_amd64.tar.gz
tar -xzf email-sentinel_linux_amd64.tar.gz
sudo mv email-sentinel /usr/local/bin/

# Or build from source
git clone https://github.com/datateamsix/email-sentinel.git
cd email-sentinel
go build -o email-sentinel .
```

### 2. Set Up Gmail API

You need Gmail API credentials (one-time setup, ~5 minutes):

**📖 Follow:** [Gmail API Setup Guide](docs/GMAIL_API_SETUP.md)

Short version:
1. Create project at [Google Cloud Console](https://console.cloud.google.com/)
2. Enable Gmail API
3. Create OAuth Desktop credentials
4. Download `credentials.json`

### 3. Initialize

```bash
# Place credentials.json in config directory or current folder
email-sentinel init
# Opens browser for OAuth → Grant permissions → Done!
```

### 4. Add Filter

```bash
# Example: Get notified for job opportunities
email-sentinel filter add \
  --name "Job Alerts" \
  --from "linkedin.com,greenhouse.io,lever.co" \
  --subject "interview,opportunity,position" \
  --scope inbox \
  --labels "work,career"
```

### 5. Start Monitoring

```bash
# Run with system tray (recommended)
email-sentinel start --tray

# Or run in foreground
email-sentinel start
```

**That's it!** Send yourself a test email and watch the notification appear. ✨

---

## 📚 Documentation

### Getting Started Guides

Choose your platform for step-by-step setup:

- **[Windows Quickstart](docs/QUICKSTART_WINDOWS.md)** - Complete Windows 10/11 guide with Task Scheduler auto-start
- **[macOS Quickstart](docs/QUICKSTART_MACOS.md)** - Menu bar app setup with LaunchAgent auto-start
- **[Linux Quickstart](docs/QUICKSTART_LINUX.md)** - Desktop environment setup with systemd auto-start
- **[Complete Setup Guide](docs/BUILD_TO_FIRST_ALERT.md)** - Universal guide from build to first alert

### Feature Guides

- **[CLI Command Reference](docs/CLI_GUIDE.md)** - All commands explained
- **[Central Email Setup](docs/CENTRAL_EMAIL_SETUP.md)** - Monitor multiple accounts via forwarding
- **[Mobile Notifications (ntfy)](docs/mobile_ntfy_setup.md)** - Free push notifications to phone
- **[AI Architecture](docs/AI_ARCHITECTURE.md)** - How AI summaries work
- **[System Tray Architecture](docs/TRAY_SYSTEM_ARCHITECTURE.md)** - Background app design

---

## 🎯 Top 5 Use Cases

### 1. 🎯 Job Hunting & Recruiter Alerts

**The Problem:** Miss interview invites or recruiter messages buried in your inbox.

**The Solution:** Get instant alerts when recruiters reach out or job applications progress.

```bash
email-sentinel filter add \
  --name "Job Opportunities" \
  --from "linkedin.com,greenhouse.io,lever.co,indeed.com,glassdoor.com" \
  --subject "interview,opportunity,position,application,recruiter" \
  --scope inbox \
  --labels "career,urgent"

# Also check promotions where some recruiting emails land
email-sentinel filter add \
  --name "Job Promotions" \
  --from "linkedin.com,indeed.com,ziprecruiter.com" \
  --scope promotions \
  --labels "career"
```

**What you'll catch:**
- Interview invitations
- Application status updates
- Direct recruiter messages
- Job match notifications
- Networking opportunities

---

### 2. 🔐 OTP & 2FA Verification Codes

**The Problem:** Constantly switching between email and apps to copy verification codes.

**The Solution:** Auto-extract OTP codes, copy to clipboard, view in tray.

```bash
# OTP detection works automatically with any filter
# Emails from common OTP senders are auto-detected
email-sentinel filter add \
  --name "Verification Codes" \
  --from "noreply@github.com,accounts.google.com,amazon.com,paypal.com" \
  --subject "verification,code,otp,2fa,authenticate" \
  --scope "inbox+updates" \
  --labels "otp,security"

# View all OTP codes
email-sentinel otp list --active

# Copy latest code to clipboard
email-sentinel otp get
```

**What you'll catch:**
- GitHub verification codes
- Google account codes
- Banking 2FA codes
- Amazon, PayPal, Microsoft codes
- Any 6-digit verification codes

---

### 3. 👔 Boss & VIP Contact Alerts

**The Problem:** Miss urgent emails from your boss, clients, or important contacts.

**The Solution:** Never miss messages from people who matter most.

```bash
# Your boss
email-sentinel filter add \
  --name "Boss" \
  --from "boss@company.com" \
  --scope inbox \
  --labels "work,vip,urgent"

# Key clients (add to app-config.yaml for priority 🔥)
email-sentinel filter add \
  --name "Top Clients" \
  --from "client1@company.com,client2@agency.com,ceo@bigclient.com" \
  --scope inbox \
  --labels "work,vip,billing"

# Important domains (entire company)
email-sentinel filter add \
  --name "Strategic Partner" \
  --from "@importantpartner.com" \
  --scope inbox \
  --labels "work,vip"
```

**Pro tip:** Add VIP senders to `app-config.yaml` for automatic urgent (🔥) priority:
```yaml
priority:
  vip_senders:
    - boss@company.com
    - ceo@bigclient.com
  vip_domains:
    - importantpartner.com
```

---

### 4. 🚨 Spam Rescue & Important Email Recovery

**The Problem:** Important emails sometimes land in spam/promotions and you miss them.

**The Solution:** Monitor spam and promotions for specific important senders.

```bash
# Check spam for specific contacts
email-sentinel filter add \
  --name "Spam Check - Important Senders" \
  --from "important@freelancer.com,contractor@startup.com" \
  --scope "all-except-trash" \
  --labels "spam-rescue,important"

# Monitor promotions for bills/receipts
email-sentinel filter add \
  --name "Bills in Promotions" \
  --from "billing@utilities.com,noreply@bank.com" \
  --subject "bill,invoice,statement,payment" \
  --scope "promotions+updates" \
  --labels "billing,important"

# Catch forwarded emails that might be misfiled
email-sentinel filter add \
  --name "Forwarded Important" \
  --subject "fwd:,forwarded" \
  --from "assistant@company.com" \
  --scope "all-except-trash" \
  --labels "forwarded"
```

**What you'll catch:**
- Important emails misfiled by Gmail
- Bills/invoices in promotions
- Forwarded messages
- Contractor/freelancer emails
- Small business communications

---

### 5. ⚡ Urgent Keywords & Time-Sensitive Alerts

**The Problem:** Miss emails marked "urgent", "asap", or with deadlines.

**The Solution:** Auto-detect urgency keywords across all senders.

```bash
# Urgent keywords in any email
email-sentinel filter add \
  --name "Urgent Keywords" \
  --subject "urgent,asap,emergency,critical,deadline,action required,time sensitive" \
  --scope "inbox+primary" \
  --labels "urgent"

# Work-related urgent items
email-sentinel filter add \
  --name "Work Deadlines" \
  --from "@company.com,@client.com" \
  --subject "deadline,due,eod,asap,urgent" \
  --scope inbox \
  --labels "work,urgent,deadline"

# Payment & billing urgency
email-sentinel filter add \
  --name "Payment Deadlines" \
  --subject "overdue,payment due,final notice,reminder" \
  --scope "inbox+updates+promotions" \
  --labels "billing,urgent"
```

**Pro tip:** Add urgent keywords to `app-config.yaml` for automatic priority:
```yaml
priority:
  urgent_keywords:
    - urgent
    - asap
    - deadline
    - action required
    - emergency
    - time sensitive
    - final notice
```

**What you'll catch:**
- Deadline reminders
- Emergency notifications
- Payment due notices
- Time-sensitive requests
- Critical updates

---

## 🔥 Advanced Features

### Gmail Category Scopes

Target specific Gmail categories when creating filters:

```bash
# Monitor only primary inbox
email-sentinel filter add --name "Important" --from "ceo@company.com" --scope primary

# Monitor social category
email-sentinel filter add --name "Social" --from "facebook.com" --scope social

# Monitor multiple categories
email-sentinel filter add --name "Updates" --subject "receipt,confirmation" --scope "primary+updates"
```

**Available scopes:**
- `inbox` - Primary inbox (default)
- `primary` - Primary category
- `social` - Social category
- `promotions` - Promotions category
- `updates` - Updates category (receipts, confirmations)
- `forums` - Forums category (mailing lists)
- `all` - All mail including spam
- `all-except-trash` - Everything except trash
- `primary+social` - Combine multiple with `+`

**Override all filters globally:**
```bash
# Search only social category for all filters
email-sentinel start --tray --search social
```

### AI Email Summaries

Get instant AI-powered summaries of emails (FREE with Gemini):

```bash
# 1. Get free API key
# Visit: https://makersuite.google.com/app/apikey

# 2. Set environment variable
export GEMINI_API_KEY="your-api-key"  # macOS/Linux
setx GEMINI_API_KEY "your-api-key"    # Windows

# 3. Enable in config
# Edit: app-config.yaml, set ai_summary.enabled: true

# 4. Start with AI
email-sentinel start --tray --ai-summary
```

Summaries include:
- 📝 Concise overview
- ❓ Questions asked
- ✅ Action items

Supports: **Gemini** (free), **Claude**, **OpenAI GPT**

### OTP/2FA Code Detection

Automatically extract verification codes from emails:

```bash
# View recent codes
email-sentinel otp list

# Get latest code (copies to clipboard)
email-sentinel otp get

# Test detection
email-sentinel otp test "Your verification code is 849372"
```

Supports 15+ services including Gmail, GitHub, Amazon, PayPal, banking, and more.

### Priority Rules

Auto-detect urgent emails with `app-config.yaml`:

```yaml
priority:
  urgent_keywords:
    - urgent
    - asap
    - action required
    - deadline
    - critical

  vip_senders:
    - boss@company.com
    - ceo@company.com

  vip_domains:
    - importantclient.com
```

Urgent emails show 🔥 icon in notifications.

### Filter Labels

Organize filters by category:

```bash
# Labels appear in all notifications
email-sentinel filter add \
  --name "Client Emails" \
  --from "client@company.com" \
  --labels "work,urgent,billing"

# List shows: 🏷️ work, urgent, billing
```

Labels help you:
- Categorize alerts at a glance
- Filter notifications by type
- Organize your workflow

### Alert History

View past alerts with clickable Gmail links:

```bash
# View today's alerts
email-sentinel alerts

# View last 10 alerts
email-sentinel alerts --recent 10

# Clear all alerts
email-sentinel alerts clear
```

### System Tray Features

Run with `--tray` for background operation:

```bash
email-sentinel start --tray
```

**Features:**
- 📧 Icon shows alert status (mailbox with flag)
- 📋 Recent Alerts submenu (click to open in Gmail)
- 🔥 Red icon flash for urgent alerts
- ⚙️ Manage filters from tray menu
- 📊 View alert history
- 🗑️ Clear alerts

### Multi-Account Monitoring

Monitor multiple email accounts by forwarding to one central Gmail:

**📖 See:** [Central Email Setup Guide](docs/CENTRAL_EMAIL_SETUP.md)

```
┌──────────────┐
│  Work Gmail  │──┐
└──────────────┘  │
┌──────────────┐  │  ┌──────────────────┐       ┌──────────────────┐
│   Outlook    │──┼─▶│  Central Gmail   │──────▶│ Email Sentinel   │
└──────────────┘  │  └──────────────────┘       └──────────────────┘
┌──────────────┐  │
│ Personal GMail│──┘
└──────────────┘
```

Benefits:
- Single OAuth authentication
- Monitor unlimited accounts
- Retain original sender info
- Multi-stage filtering (email client → Gmail → Email Sentinel)

---

## 📖 Command Reference

### Filters

```bash
# Add filter
email-sentinel filter add [--name] [--from] [--subject] [--scope] [--labels] [--match]

# List filters
email-sentinel filter list

# Edit filter
email-sentinel filter edit [name]

# Remove filter
email-sentinel filter remove [name]
```

### Monitoring

```bash
# Start (foreground)
email-sentinel start

# Start with system tray
email-sentinel start --tray

# Start with AI summaries
email-sentinel start --tray --ai-summary

# Start with global scope override
email-sentinel start --search social

# Start as daemon
email-sentinel start --daemon

# Check status
email-sentinel status

# Stop daemon
email-sentinel stop
```

### Alerts & OTP

```bash
# View alerts
email-sentinel alerts [--recent N]

# Clear alerts
email-sentinel alerts clear

# List OTP codes
email-sentinel otp list [--active]

# Get latest OTP
email-sentinel otp get

# Clear expired OTPs
email-sentinel otp clear

# Test OTP extraction
email-sentinel otp test "text"
```

### Configuration

```bash
# Show config
email-sentinel config show

# Set polling interval
email-sentinel config set polling 30

# View config file location
# Windows: %APPDATA%\email-sentinel\
# macOS: ~/Library/Application Support/email-sentinel/
# Linux: ~/.config/email-sentinel/
```

### Testing

```bash
# Test desktop notification
email-sentinel test desktop

# Test Windows toast
email-sentinel test toast

# Test toast with priority
email-sentinel test toast --priority

# Test mobile (ntfy)
email-sentinel test mobile

# Test filter matching
email-sentinel test filter "Filter Name" "sender@example.com" "subject text"
```

### Auto-Start

```bash
# Install auto-start
email-sentinel install

# Preview install
email-sentinel install --show

# Uninstall auto-start
email-sentinel uninstall
```

---

## ⚙️ Configuration Files

All config files are in the platform-specific config directory:

| Platform | Location |
|----------|----------|
| **Windows** | `%APPDATA%\email-sentinel\` |
| **macOS** | `~/Library/Application Support/email-sentinel/` |
| **Linux** | `~/.config/email-sentinel/` |

### Main Configuration (`app-config.yaml`)

Unified configuration for all features:

```yaml
# Polling interval (seconds)
polling_interval: 45

# Filter configuration
filters:
  enabled: true

# Notifications
notifications:
  desktop:
    enabled: true
  mobile:
    enabled: false
    ntfy_topic: ""
  quiet_hours:
    enabled: false
    start: "22:00"
    end: "08:00"
  weekend_mode: false

# Priority rules
priority:
  urgent_keywords:
    - urgent
    - asap
    - critical
  vip_senders:
    - boss@company.com
  vip_domains:
    - importantclient.com

# AI summaries (optional)
ai_summary:
  enabled: false
  provider: gemini  # gemini, claude, openai
  providers:
    gemini:
      model: gemini-1.5-flash
      max_tokens: 1024
      temperature: 0.3

# OTP settings
otp:
  enabled: true
  expiry_duration: "5m"
  trusted_senders:
    - accounts.google.com
    - noreply@github.com
    - amazon.com
    - paypal.com
  trusted_domains:
    - amazon.com
    - paypal.com
```

**Note**: All settings are now unified in `app-config.yaml`. If upgrading from an older version with separate config files (`rules.yaml`, `otp_rules.yaml`, `ai-config.yaml`), run `email-sentinel config migrate` to automatically convert to the new format.

---

## 🔧 Architecture

### How It Works

```
┌─────────────────────────────────────────────────────────┐
│                    EMAIL SENTINEL                       │
├─────────────────────────────────────────────────────────┤
│                                                         │
│   ┌──────────┐    ┌──────────┐    ┌──────────────┐    │
│   │ Gmail API│───▶│  Filter  │───▶│Notifications │    │
│   │ (poll)   │    │  Engine  │    │  • Desktop   │    │
│   └──────────┘    └──────────┘    │  • Mobile    │    │
│         │              │           │  • Toast     │    │
│         │              │           │  • Tray      │    │
│         ▼              ▼           └──────────────┘    │
│   ┌──────────┐    ┌──────────┐                         │
│   │  OAuth   │    │ Config   │                         │
│   │  Token   │    │  YAML    │                         │
│   └──────────┘    └──────────┘                         │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Key Components

- **Gmail Client** (`internal/gmail/`) - API wrapper with auto-refresh
- **Filter Engine** (`internal/filter/`) - Match logic with Gmail scope support
- **Notification System** (`internal/notify/`) - Desktop, mobile, toast
- **System Tray** (`internal/tray/`) - Background app with menu
- **Alert Storage** (`internal/storage/`) - SQLite database
- **AI Service** (`internal/ai/`) - Optional summaries
- **OTP Detector** (`internal/otp/`) - Code extraction
- **Priority Rules** (`internal/rules/`) - Urgency classification

### Data Flow

1. **Polling**: Gmail API fetches recent messages (per-filter scope)
2. **Deduplication**: Messages checked against seen state
3. **Filtering**: Each message tested against all filters
4. **Priority**: Rules engine evaluates urgency
5. **AI Summary** (optional): Generate async summary
6. **OTP Detection**: Extract verification codes
7. **Notifications**: Send desktop, mobile, toast alerts
8. **Storage**: Save to SQLite with Gmail links
9. **Tray Update**: Refresh recent alerts menu

---

## 🛠️ Troubleshooting

### "App not verified" during OAuth

**This is normal for personal apps.**

1. Click **"Advanced"**
2. Click **"Go to email-sentinel (unsafe)"**
3. Click **"Allow"**

Your app only accesses your own Gmail - this is safe.

### "Access blocked" error

Add your email as a test user:

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. APIs & Services → OAuth consent screen
3. Test users → Add Users
4. Enter your Gmail address

### Token expired

Re-authenticate:
```bash
email-sentinel init
```

### Notifications not appearing

**Windows:**
- Settings → System → Notifications → Enable
- Disable Focus Assist temporarily
- Enable for PowerShell/CMD/Terminal

**macOS:**
- System Preferences → Notifications
- Find email-sentinel or Terminal
- Enable "Allow Notifications"

**Linux:**
- Check notification daemon (notify-send)
- Desktop environment notification settings

### System tray icon missing

**Windows:**
- Taskbar settings → Select which icons appear
- Enable "Email Sentinel"
- Restart explorer.exe

**macOS:**
- Icon appears in menu bar (top-right)
- Check for icon with mailbox symbol

**Linux:**
- Ensure system tray support (GNOME Shell Extension, etc.)

### Build fails (from source)

**Enable CGO:**
```bash
# Windows
set CGO_ENABLED=1

# macOS/Linux
export CGO_ENABLED=1

go build -o email-sentinel .
```

**Install build tools:**
- Windows: `scoop install mingw`
- macOS: Install Xcode Command Line Tools
- Linux: `sudo apt install build-essential`

---

## 🤝 Contributing

Contributions welcome! Areas for improvement:

- [ ] Native multi-account OAuth support
- [ ] Outlook/Microsoft Graph API
- [ ] Web UI for filter management
- [ ] Slack/Discord webhook integration
- [ ] Custom notification sounds
- [ ] Email digest/summary reports
- [ ] Filter import/export
- [ ] IMAP fallback support

**To contribute:**

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [ntfy.sh](https://ntfy.sh) - Free push notification service
- [Google Gmail API](https://developers.google.com/gmail/api) - Email access
- [Fyne Systray](https://github.com/fyne-io/systray) - Cross-platform system tray

---

## 📞 Support

- **Documentation:** [docs/](docs/)
- **Issues:** https://github.com/datateamsix/email-sentinel/issues
- **Discussions:** https://github.com/datateamsix/email-sentinel/discussions

---

**Built with ❤️ by [DataTeamSix](https://github.com/datateamsix)**

**Protect your productivity. Never miss what matters.** 📧✨
