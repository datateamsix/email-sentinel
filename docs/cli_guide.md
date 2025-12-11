# Email Sentinel CLI Guide

![Email Sentinel Logo](../images/logo.png)

A comprehensive guide to using the Email Sentinel command-line interface for monitoring Gmail and receiving real-time notifications.

**Last Updated:** 2025-12-07
**Version:** 1.0

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [Core Concepts](#core-concepts)
3. [Command Reference](#command-reference)
4. [Configuration](#configuration)
5. [Common Workflows](#common-workflows)
6. [Advanced Usage](#advanced-usage)
7. [Troubleshooting](#troubleshooting)

---

## Quick Start

### First-Time Setup (5 minutes)

```bash
# 1. Initialize and authenticate with Gmail
email-sentinel init

# 2. Add your first filter
email-sentinel filter add --name "Important" --from "boss@company.com"

# 3. Test notifications
email-sentinel test desktop

# 4. Start monitoring
email-sentinel start --tray
```

### Daily Usage

```bash
# Start monitoring (with system tray)
email-sentinel start --tray

# View recent alerts
email-sentinel alerts --recent 10

# Check if running
email-sentinel status
```

---

## Core Concepts

### Filters

**Filters** define which emails trigger notifications. Each filter can match on:

- **Sender** (`--from`): Email addresses or domains
- **Subject** (`--subject`): Keywords in the subject line
- **Gmail Scope** (`--scope`): Which Gmail categories to search (NEW)
- **Match Mode** (`--match`): How to combine conditions
  - `any` (OR): Trigger if sender OR subject matches
  - `all` (AND): Trigger only if sender AND subject both match
- **Labels** (`--labels`): Organize filters by category (work, urgent, personal)

**Example:**
```bash
email-sentinel filter add \
  --name "Job Alerts" \
  --from "linkedin.com,greenhouse.io" \
  --subject "interview,opportunity" \
  --scope inbox \
  --match any \
  --labels "work,career"
```

### Gmail Scopes

**Gmail scopes** let you target specific Gmail categories for each filter, reducing unnecessary API calls and improving performance.

**Available Scopes:**

| Scope | Description | Use Case |
|-------|-------------|----------|
| `inbox` | Primary inbox only (default) | Most email filters |
| `primary` | Primary category | Important personal/business emails |
| `social` | Social category | Facebook, Twitter, LinkedIn notifications |
| `promotions` | Promotions category | Marketing emails, newsletters |
| `updates` | Updates category | Receipts, confirmations, statements |
| `forums` | Forums category | Mailing lists, discussion groups |
| `all` | All mail including spam | Broad monitoring |
| `all-except-trash` | Everything except trash | Wide scope excluding deleted |
| `primary+social` | Multiple categories | Combine with `+` separator |

**Examples:**

```bash
# Monitor social category only
email-sentinel filter add \
  --name "Social Notifications" \
  --from "facebook.com,twitter.com" \
  --scope social \
  --labels "social"

# Monitor multiple categories
email-sentinel filter add \
  --name "Updates & Receipts" \
  --subject "receipt,confirmation,order" \
  --scope "primary+updates" \
  --labels "receipts"

# Global scope override (ignores per-filter scopes)
email-sentinel start --tray --search social
```

**How It Works:**

- Each filter searches only its specified Gmail category
- Messages are deduplicated across filters
- More efficient than searching all mail
- Reduces Gmail API quota usage

### Priority Rules

**Priority rules** automatically classify emails as urgent (🔥) or normal (📧).

Configured in `rules.yaml`:
- **Urgent Keywords**: Subject/snippet contains keywords → Priority 1
- **VIP Senders**: Exact email match → Priority 1
- **VIP Domains**: Sender's domain matches → Priority 1

**Location:**
- Windows: `%APPDATA%\email-sentinel\rules.yaml`
- macOS: `~/Library/Application Support/email-sentinel/rules.yaml`
- Linux: `~/.config/email-sentinel/rules.yaml`

### OTP/2FA Detection

**OTP detection** automatically extracts verification codes from emails.

Features:
- 10+ built-in patterns with confidence scoring
- Auto-expiry (default: 5 minutes)
- Clipboard integration with auto-clear
- False positive prevention (rejects sequential/repeating digits)

Configured in `otp_rules.yaml` (same directory as `rules.yaml`)

### Alert History

All email notifications are saved to a SQLite database (`history.db`) with:
- Filter name, sender, subject, preview
- Direct Gmail permalink
- Priority indicator
- Timestamp

Automatic cleanup: Alerts older than midnight are deleted.

---

## Command Reference

### Global Options

All commands support these flags:

| Flag | Description |
|------|-------------|
| `-h, --help` | Show help for any command |
| `--config <path>` | Use custom config file |
| `-v, --verbose` | Enable verbose logging |

---

### Initialization & Authentication

#### `email-sentinel init`

Initialize Email Sentinel and authenticate with Gmail OAuth 2.0.

**Usage:**
```bash
email-sentinel init
```

**What it does:**
1. Searches for `credentials.json` in:
   - Current directory
   - Config directory (`%APPDATA%\email-sentinel` on Windows)
2. Opens browser for Google OAuth
3. Saves authentication token (`token.json`)
4. Displays existing filters and next steps

**First-time users:** See [Gmail API Setup](gmail_api_setup.md) to create `credentials.json`.

**Re-authentication:** Run this command if your token expires or you need to switch accounts.

---

### Filter Management

#### `email-sentinel filter add`

Add a new email filter with optional labels.

**Interactive Mode:**
```bash
email-sentinel filter add
```
Shows existing labels for reuse and guides you through filter creation.

**CLI Mode:**
```bash
email-sentinel filter add \
  --name "Filter Name" \
  --from "sender1.com,sender2@example.com" \
  --subject "keyword1,keyword2" \
  --match any \
  --labels "label1,label2"
```

**Flags:**

| Flag | Short | Required | Description | Example |
|------|-------|----------|-------------|---------|
| `--name` | `-n` | Yes | Filter name | `"Job Alerts"` |
| `--from` | `-f` | No | Sender patterns (comma-separated) | `"linkedin.com,@github.com"` |
| `--subject` | `-s` | No | Subject keywords (comma-separated) | `"urgent,asap"` |
| `--scope` | | No | Gmail scope/category (default: `inbox`) | `social`, `primary+updates` |
| `--match` | `-m` | No | Match mode: `any` or `all` (default: `any`) | `any` |
| `--labels` | `-l` | No | Labels/categories (comma-separated) | `"work,urgent"` |

**Examples:**

```bash
# Match emails from specific sender
email-sentinel filter add --name "Boss" --from "boss@company.com" --labels "work,urgent"

# Match subject keywords only
email-sentinel filter add --name "Urgent" --subject "urgent,asap,critical" --labels "urgent"

# Match sender AND subject (precise matching)
email-sentinel filter add --name "Client Invoices" \
  --from "billing@client.com" \
  --subject "invoice,payment" \
  --match all \
  --labels "work,billing"

# Match multiple senders (OR logic)
email-sentinel filter add --name "Job Sites" \
  --from "linkedin.com,indeed.com,glassdoor.com" \
  --labels "career"
```

**Label Tips:**
- Labels are persistent and suggested when creating new filters
- Labels appear in all notifications (desktop, mobile, toast)
- Use labels to organize filters: `work`, `personal`, `urgent`, `family`, `billing`, etc.

#### `email-sentinel filter list`

Display all configured filters.

**Usage:**
```bash
email-sentinel filter list
```

**Example Output:**
```
📋 Current Filters (3 total)

[1] Job Alerts
    From: linkedin.com, greenhouse.io
    Subject: interview, opportunity
    Match mode: any
    Labels: work, career

[2] Boss
    From: boss@company.com
    Subject: (any)
    Match mode: any
    Labels: work, urgent

[3] Urgent
    From: (any)
    Subject: urgent, asap, critical
    Match mode: any
    Labels: urgent
```

#### `email-sentinel filter edit`

Edit an existing filter.

**Interactive Mode:**
```bash
email-sentinel filter edit
```
Prompts you to select a filter, then walks through editing each field.

**CLI Mode:**
```bash
email-sentinel filter edit "Filter Name"
```

#### `email-sentinel filter remove`

Remove a filter.

**Interactive Mode:**
```bash
email-sentinel filter remove
```
Shows numbered list of filters to remove.

**CLI Mode:**
```bash
email-sentinel filter remove "Filter Name"
```

---

### Monitoring

#### `email-sentinel start`

Start monitoring Gmail for matching emails.

**Usage:**
```bash
# Foreground mode (logs to stdout)
email-sentinel start

# With system tray icon (recommended)
email-sentinel start --tray

# As background daemon
email-sentinel start --daemon
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--tray` | `-t` | Run with system tray icon and menu |
| `--daemon` | `-d` | Run as background daemon (no output) |

**Foreground Mode:**
- Logs appear in terminal
- Press `Ctrl+C` to stop
- Useful for testing and debugging

**System Tray Mode (Recommended):**
- Icon appears in taskbar/menu bar
- Right-click for menu with recent alerts
- Click alerts to open in Gmail
- Quit from tray menu

**Daemon Mode:**
- Runs in background
- No terminal output
- Use `email-sentinel stop` to stop
- Use `email-sentinel status` to check if running

**What happens when started:**
1. Loads filters from `config.yaml`
2. Loads priority rules from `rules.yaml`
3. Loads OTP rules from `otp_rules.yaml`
4. Connects to Gmail API
5. Polls inbox every N seconds (default: 45)
6. Checks for new emails matching filters
7. Sends notifications for matches
8. Extracts OTP codes if detected
9. Saves alerts to `history.db`

#### `email-sentinel stop`

Stop the background daemon.

**Usage:**
```bash
email-sentinel stop
```

Only works if Email Sentinel was started with `--daemon` flag.

#### `email-sentinel status`

Check if Email Sentinel is currently running.

**Usage:**
```bash
email-sentinel status
```

**Example Output:**
```
✅ Email Sentinel is running
   PID: 12345
   Uptime: 2h 34m
   Last check: 15 seconds ago
```

Or:
```
❌ Email Sentinel is not running
```

---

### Alert History

#### `email-sentinel alerts`

View email alert history from the database.

**Usage:**
```bash
# View today's alerts
email-sentinel alerts

# View last N alerts
email-sentinel alerts --recent 10

# View last 5 alerts
email-sentinel alerts --recent 5
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--recent` | `-r` | Show last N alerts (instead of just today) |

**Example Output:**
```
📬 Today's Alerts (3 total)

[1] 🔥 2025-12-07 14:30:15
    Filter: VIP Senders
    Priority: HIGH
    From:   boss@company.com
    Subject: URGENT: Server is down!
    Preview: We need to address this immediately...
    Link:   https://mail.google.com/mail/u/0/#all/abc123

[2] 📧 2025-12-07 13:45:22
    Filter: Job Alerts
    From:   recruiter@linkedin.com
    Subject: New job opportunity
    Preview: We think you'd be a great fit for...
    Link:   https://mail.google.com/mail/u/0/#all/def456

[3] 📧 2025-12-07 09:12:07
    Filter: Urgent
    From:   support@github.com
    Subject: Security alert
    Link:   https://mail.google.com/mail/u/0/#all/ghi789
```

**Priority Indicators:**
- 🔥 = High priority (urgent)
- 📧 = Normal priority

**Clicking Links:**
Copy the Gmail link and paste in browser to open the email directly.

---

### OTP/2FA Management

#### `email-sentinel otp list`

List OTP codes extracted from emails.

**Usage:**
```bash
# List all recent OTP codes
email-sentinel otp list

# List only active (non-expired) codes
email-sentinel otp list --active
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--active` | `-a` | Show only non-expired codes |

**Example Output:**
```
🔐 Active OTP Codes (2 total)

[1] 🔐 849372 (Confidence: 1.00)
    From: noreply@github.com
    Received: 2 minutes ago
    Expires: in 3 minutes
    Source: subject

[2] 🔐 482716 (Confidence: 0.95)
    From: accounts.google.com
    Received: 4 minutes ago
    Expires: in 1 minute
    Source: body
```

**Confidence Scores:**
- `1.00` = Perfect match (e.g., "Your code is 123456")
- `0.90-0.99` = Very likely OTP
- `0.70-0.89` = Likely OTP
- Below `0.70` = Rejected (configurable in `otp_rules.yaml`)

#### `email-sentinel otp get`

Get the most recent OTP code and copy to clipboard.

**Usage:**
```bash
email-sentinel otp get
```

**Example Output:**
```
🔐 Most Recent OTP: 849372

✅ Copied to clipboard
   From: noreply@github.com
   Received: 2 minutes ago
   Expires: in 3 minutes

⚠️  Clipboard will auto-clear in 2 minutes
```

**Security:**
- Clipboard auto-clears after 2 minutes (configurable)
- Only shows non-expired codes

#### `email-sentinel otp clear`

Remove expired OTP codes from the database.

**Usage:**
```bash
email-sentinel otp clear
```

**Example Output:**
```
🗑️  Cleared 5 expired OTP codes
```

**Auto-Cleanup:**
Email Sentinel automatically clears expired codes during monitoring. Manual clearing is optional.

#### `email-sentinel otp test`

Test OTP detection on sample text.

**Usage:**
```bash
email-sentinel otp test "Your GitHub verification code is 849372"
```

**Example Output:**
```
Testing text: Your GitHub verification code is 849372

✅ OTP Detected: 849372
Confidence: 1.00
Pattern: your_code_is
Source: text
```

**Use Cases:**
- Test if Email Sentinel will detect codes from specific services
- Verify custom patterns work correctly
- Debug OTP detection issues

---

### Testing

#### `email-sentinel test desktop`

Test desktop notification system.

**Usage:**
```bash
email-sentinel test desktop
```

**What it does:**
Sends a test notification using OS-native notification system.

**Expected Result:**
Notification pops up with:
- Title: "Email Sentinel Test"
- Message: "If you can see this, desktop notifications are working! ✅"

**Troubleshooting:**
- **Windows:** Settings → System → Notifications → Enable
- **macOS:** System Preferences → Notifications → Terminal → Allow
- **Linux:** Install `libnotify-bin` or `dunst`

#### `email-sentinel test toast`

Test Windows toast notification (Windows only).

**Usage:**
```bash
# Test normal priority
email-sentinel test toast

# Test high priority
email-sentinel test toast --priority
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--priority` | Test high-priority urgent notification |

**What it does:**
Sends a rich toast notification to Windows Action Center.

**Expected Result:**
- Toast appears in Action Center
- Shows subject, sender, preview
- Includes clickable Gmail link
- High-priority toasts show 🔥 icon

#### `email-sentinel test mobile`

Test mobile push notification via ntfy.sh.

**Usage:**
```bash
email-sentinel test mobile
```

**Prerequisites:**
1. Mobile notifications enabled in config
2. ntfy topic configured
3. ntfy app installed on phone

**Setup:** See [Mobile ntfy Setup Guide](mobile_ntfy_setup.md)

**Expected Result:**
Push notification arrives on phone with test message.

#### `email-sentinel test filter`

Test if a filter matches specific email metadata.

**Usage:**
```bash
email-sentinel test filter "Filter Name" "sender@example.com" "subject text"
```

**Example:**
```bash
email-sentinel test filter "Job Alerts" "recruiter@linkedin.com" "New opportunity"
```

**Output:**
```
✅ Filter "Job Alerts" matches
   Sender matched: recruiter@linkedin.com
   Subject matched: opportunity
   Match mode: any
```

Or:
```
❌ Filter "Job Alerts" does not match
   Sender: No match
   Subject: No match
```

**Use Case:**
Debug why alerts aren't triggering for specific emails.

---

### Configuration

#### `email-sentinel config show`

Display current configuration.

**Usage:**
```bash
email-sentinel config show
```

**Example Output:**
```
📝 Current Configuration

Config file: C:\Users\username\AppData\Roaming\email-sentinel\config.yaml

Polling interval: 45 seconds
Desktop notifications: enabled
Mobile notifications: enabled
  ntfy topic: my-secret-topic-x7k2

Filters: 3
Priority rules: loaded
OTP detection: enabled
```

#### `email-sentinel config set`

Modify configuration values.

**Usage:**
```bash
# Set polling interval (seconds)
email-sentinel config set polling 30

# Enable/disable mobile notifications
email-sentinel config set mobile true
email-sentinel config set mobile false

# Set ntfy topic
email-sentinel config set ntfy_topic "your-topic-name"

# Enable/disable desktop notifications
email-sentinel config set desktop true
email-sentinel config set desktop false
```

**Configuration Keys:**

| Key | Type | Description | Default |
|-----|------|-------------|---------|
| `polling` | int | Seconds between Gmail checks | `45` |
| `desktop` | bool | Enable desktop notifications | `true` |
| `mobile` | bool | Enable mobile push notifications | `false` |
| `ntfy_topic` | string | ntfy.sh topic name | `""` |

**Notes:**
- Changes take effect on next `start`
- Restart Email Sentinel after config changes
- Direct YAML editing also supported

---

### Auto-Start

#### `email-sentinel install`

Install Email Sentinel to run automatically on system startup.

**Usage:**
```bash
# Install auto-start
email-sentinel install

# Preview what would be installed (dry-run)
email-sentinel install --show
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--show` | Preview installation without making changes |

**What it does:**

**Windows:**
- Creates Task Scheduler task named "EmailSentinel"
- Trigger: At logon
- Action: Run `email-sentinel.exe start --tray`

**macOS:**
- Creates LaunchAgent: `~/Library/LaunchAgents/com.email-sentinel.plist`
- Runs on login

**Linux:**
- Creates systemd user service: `~/.config/systemd/user/email-sentinel.service`
- Enabled for auto-start

**Verification:**

**Windows:**
```bash
# Check task exists
Get-ScheduledTask -TaskName "EmailSentinel"

# View in GUI
taskschd.msc
```

**macOS:**
```bash
launchctl list | grep email-sentinel
```

**Linux:**
```bash
systemctl --user status email-sentinel
```

#### `email-sentinel uninstall`

Remove Email Sentinel from automatic startup.

**Usage:**
```bash
email-sentinel uninstall
```

**What it does:**
- **Windows:** Deletes Task Scheduler task
- **macOS:** Removes LaunchAgent plist
- **Linux:** Disables and removes systemd service

---

### Interactive Mode

#### `email-sentinel` (no args)

Launch interactive menu system.

**Usage:**
```bash
email-sentinel
```

**Features:**
- First-run setup wizard
- Visual menu navigation
- Real-time status dashboard
- Guided filter creation
- Notification configuration
- OTP/2FA setup

**Main Menu:**
```
╔═══════════════════════════════════════════════════════════╗
║                        MAIN MENU                          ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  [1] 🚀 Start Monitoring      Start watching for emails   ║
║  [2] 📋 Manage Filters        Add, edit, remove filters   ║
║  [3] 🔔 Notifications         Configure alerts            ║
║  [4] 📊 Status & History      View alerts and status      ║
║  [5] ⚙️  Settings             Configure app settings      ║
║  [6] 🔧 Setup Wizard          Re-run initial setup        ║
║                                                           ║
║  [q] Exit                                                 ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
```

**Keyboard Shortcuts:**
- `1-9`: Select menu item
- `Enter`: Confirm
- `b`: Back
- `q`: Quit
- `r`: Refresh (in dashboard)
- `?`: Help

---

## Configuration

### Config File Locations

| OS | Path |
|----|------|
| Windows | `%APPDATA%\email-sentinel\config.yaml` |
| macOS | `~/Library/Application Support/email-sentinel/config.yaml` |
| Linux | `~/.config/email-sentinel/config.yaml` |

### Config File Structure

**config.yaml:**
```yaml
polling_interval: 45  # Seconds between Gmail checks

filters:
  - name: "Job Opportunities"
    from:
      - "@linkedin.com"
      - "greenhouse.io"
    subject:
      - "interview"
      - "opportunity"
    match: any  # "any" (OR) or "all" (AND)
    labels:
      - work
      - career

notifications:
  desktop: true
  mobile:
    enabled: true
    ntfy_topic: "your-secret-topic-name"
```

**rules.yaml:**
```yaml
priority_rules:
  urgent_keywords:
    - urgent
    - asap
    - action required
    - deadline

  vip_senders:
    - boss@company.com

  vip_domains:
    - importantclient.com

notification_settings:
  quiet_hours_start: "22:00"
  quiet_hours_end: "08:00"
```

**otp_rules.yaml:**
```yaml
enabled: true
expiry_duration: "5m"
confidence_threshold: 0.7
auto_copy_to_clipboard: false
clipboard_auto_clear: "2m"

trusted_otp_senders:
  - accounts.google.com
  - noreply@github.com
  - amazon.com

custom_patterns:
  - name: "custom_code"
    regex: "Code:\\s*([A-Z0-9]{6})"
    confidence: 0.8
```

### Database Files

| File | Purpose | Location |
|------|---------|----------|
| `history.db` | Alert history and OTP codes | Same as config directory |
| `seen.db` | Track processed emails | Same as config directory |
| `token.json` | Gmail OAuth token | Same as config directory |

---

## Common Workflows

### Daily Monitoring Workflow

```bash
# Morning: Start monitoring with tray
email-sentinel start --tray

# Throughout day: Receive notifications automatically

# Evening: Check what you received
email-sentinel alerts

# View OTP codes if needed
email-sentinel otp list --active
```

### Job Search Workflow

```bash
# Add job-related filters
email-sentinel filter add \
  --name "LinkedIn Jobs" \
  --from "linkedin.com" \
  --subject "opportunity,interview,application" \
  --labels "career"

email-sentinel filter add \
  --name "Job Boards" \
  --from "indeed.com,glassdoor.com,ziprecruiter.com" \
  --labels "career"

# Start monitoring
email-sentinel start --tray

# Check job alerts
email-sentinel alerts | grep career
```

### VIP Email Workflow

```bash
# Add VIP filters
email-sentinel filter add \
  --name "Boss" \
  --from "boss@company.com" \
  --labels "work,urgent"

email-sentinel filter add \
  --name "Important Client" \
  --from "client@bigcompany.com" \
  --labels "work,urgent"

# Configure priority rules
# Edit: %APPDATA%\email-sentinel\rules.yaml
# Add to vip_senders:
#   - boss@company.com
#   - client@bigcompany.com

# Test notifications
email-sentinel test desktop
email-sentinel test toast --priority

# Start with tray
email-sentinel start --tray
```

### Multi-Account Monitoring

See [Central Email Setup Guide](central_email_setup.md) for forwarding emails from multiple accounts.

**Quick summary:**
1. Create central Gmail "collector" account
2. Set up forwarding rules in other accounts
3. Monitor central account with Email Sentinel
4. Forwarded emails retain original sender/subject

### OTP/2FA Workflow

```bash
# Configure OTP detection (one-time)
# Edit: %APPDATA%\email-sentinel\otp_rules.yaml
# Ensure enabled: true

# Start monitoring
email-sentinel start --tray

# When you request OTP (e.g., GitHub login):
# 1. Email arrives with code
# 2. Email Sentinel extracts code automatically
# 3. Check with:
email-sentinel otp get

# Code is copied to clipboard, paste into login form
```

---

## Advanced Usage

### Custom Polling Interval

Faster polling = more responsive, but uses more API quota.

```bash
# Set to 30 seconds (faster)
email-sentinel config set polling 30

# Set to 60 seconds (slower, saves quota)
email-sentinel config set polling 60
```

**Gmail API Quota:**
- Free tier: 1 billion quota units/day
- Each poll ≈ 5-10 units
- Even 10-second polling won't exceed quota

### Filter Match Modes

**ANY (OR logic):**
```bash
email-sentinel filter add \
  --name "Work Emails" \
  --from "company.com" \
  --subject "urgent,meeting" \
  --match any
```
Triggers if:
- Email from `company.com` OR
- Subject contains "urgent" or "meeting"

**ALL (AND logic):**
```bash
email-sentinel filter add \
  --name "Urgent Work" \
  --from "company.com" \
  --subject "urgent" \
  --match all
```
Triggers only if:
- Email from `company.com` AND
- Subject contains "urgent"

### Label Organization

Create a labeling system for your filters:

```bash
# Work-related
email-sentinel filter add --name "Boss" --from "boss@company.com" --labels "work,urgent"
email-sentinel filter add --name "Team" --from "team@company.com" --labels "work"
email-sentinel filter add --name "Clients" --from "client.com" --labels "work,clients"

# Personal
email-sentinel filter add --name "Family" --from "mom@example.com,dad@example.com" --labels "personal,family"
email-sentinel filter add --name "Friends" --from "friend1@example.com" --labels "personal,social"

# Financial
email-sentinel filter add --name "Bank" --from "bank.com" --labels "finance,important"
email-sentinel filter add --name "Invoices" --subject "invoice,payment" --labels "finance,billing"

# Career
email-sentinel filter add --name "LinkedIn" --from "linkedin.com" --labels "career,networking"
email-sentinel filter add --name "Job Boards" --from "indeed.com,glassdoor.com" --labels "career,jobs"
```

Labels appear in all notifications for easy categorization.

### Custom OTP Patterns

Add custom patterns for services Email Sentinel doesn't recognize:

**Edit `otp_rules.yaml`:**
```yaml
custom_patterns:
  # Bank codes like "Your Bank Code: 123456"
  - name: "bank_code"
    regex: "Bank Code:\\s*([0-9]{6})"
    confidence: 0.9

  # Alphanumeric codes like "Code: ABC123"
  - name: "alpha_code"
    regex: "Code:\\s*([A-Z0-9]{6})"
    confidence: 0.85

  # 8-digit codes
  - name: "long_code"
    regex: "verification code is\\s*([0-9]{8})"
    confidence: 0.9
```

**Test patterns:**
```bash
email-sentinel otp test "Your Bank Code: 849372"
```

### Running Multiple Instances

Monitor multiple Gmail accounts by running separate instances:

**Instance 1 (Personal):**
```bash
# Use custom config directory
email-sentinel --config ~/.email-sentinel-personal/config.yaml start --tray
```

**Instance 2 (Work):**
```bash
email-sentinel --config ~/.email-sentinel-work/config.yaml start --tray
```

Each instance needs separate:
- Config directory
- `credentials.json`
- `token.json`

### Scripting & Automation

Email Sentinel is scriptable for automation:

**PowerShell (Windows):**
```powershell
# Auto-restart if crashed
while ($true) {
    .\email-sentinel.exe start --tray
    Start-Sleep -Seconds 5
}

# Check status and restart if needed
$status = .\email-sentinel.exe status
if ($status -notlike "*running*") {
    .\email-sentinel.exe start --daemon
}

# Get recent alerts programmatically
$alerts = .\email-sentinel.exe alerts --recent 5 | Out-String
if ($alerts -like "*boss@company.com*") {
    # Do something
}
```

**Bash (macOS/Linux):**
```bash
# Auto-restart if crashed
while true; do
    ./email-sentinel start --tray
    sleep 5
done

# Check and start if not running
./email-sentinel status || ./email-sentinel start --daemon

# Parse alert output
./email-sentinel alerts --recent 5 | grep "boss@company.com"
```

---

## Troubleshooting

### Authentication Issues

**Problem: "credentials.json not found"**

**Solution:**
```bash
# Check current directory
ls credentials.json

# Check config directory
ls %APPDATA%\email-sentinel\credentials.json  # Windows
ls ~/Library/Application\ Support/email-sentinel/credentials.json  # macOS
ls ~/.config/email-sentinel/credentials.json  # Linux

# Copy to config directory
cp credentials.json %APPDATA%\email-sentinel\  # Windows
cp credentials.json ~/Library/Application\ Support/email-sentinel/  # macOS
cp credentials.json ~/.config/email-sentinel/  # Linux
```

**Problem: "Token has expired"**

**Solution:**
```bash
# Re-authenticate
email-sentinel init
```

**Problem: "Access blocked" during OAuth**

**Solution:**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. APIs & Services → OAuth consent screen
3. Test users → Add Users
4. Enter your Gmail address
5. Retry `email-sentinel init`

### Notification Issues

**Problem: Desktop notifications not appearing**

**Windows:**
```powershell
# Test notification
.\email-sentinel.exe test desktop

# Check Windows settings:
# Settings → System → Notifications → Enable notifications
```

**macOS:**
```bash
# Test notification
./email-sentinel test desktop

# Check macOS settings:
# System Preferences → Notifications → Terminal → Allow Notifications
```

**Linux:**
```bash
# Install notification daemon
sudo apt install libnotify-bin  # Debian/Ubuntu
sudo dnf install libnotify       # Fedora

# Test
./email-sentinel test desktop
```

**Problem: Mobile notifications not working**

**Solution:**
```bash
# Verify config
email-sentinel config show

# Should show:
#   Mobile notifications: enabled
#   ntfy topic: <your-topic>

# If not configured:
email-sentinel config set mobile true
email-sentinel config set ntfy_topic "your-topic-name"

# Test
email-sentinel test mobile

# Check phone:
# - ntfy app installed?
# - Subscribed to correct topic?
# - Notifications enabled for ntfy app?
```

### Filter Issues

**Problem: No alerts showing up**

**Check 1: Is Email Sentinel running?**
```bash
email-sentinel status
```

**Check 2: Do filters exist?**
```bash
email-sentinel filter list
```

**Check 3: Test filter matching**
```bash
email-sentinel test filter "Filter Name" "sender@example.com" "subject text"
```

**Check 4: Check polling interval**
```bash
email-sentinel config show
# Wait at least polling_interval seconds for new emails
```

**Check 5: Email already seen?**
Email Sentinel only alerts on NEW emails. Send a fresh email to test.

### OTP Issues

**Problem: OTP codes not detected**

**Solution:**
```bash
# Test detection on sample text
email-sentinel otp test "Your verification code is 482716"

# If not detected, check confidence threshold
# Edit otp_rules.yaml:
confidence_threshold: 0.5  # Lower = more sensitive

# Check if OTP detection enabled
# Edit otp_rules.yaml:
enabled: true
```

**Problem: False positives (invoice numbers detected as OTP)**

**Solution:**
Email Sentinel rejects sequential (123456) and repeating (111111) digits by default.

```yaml
# Edit otp_rules.yaml:
confidence_threshold: 0.8  # Higher = fewer false positives
```

### Performance Issues

**Problem: High CPU usage**

**Solution:**
```bash
# Increase polling interval
email-sentinel config set polling 60

# Reduce number of filters (combine similar ones)
email-sentinel filter list
email-sentinel filter remove "Unused Filter"
```

**Problem: Gmail API quota exceeded**

Very rare (1 billion units/day), but if it happens:
```bash
# Increase polling interval dramatically
email-sentinel config set polling 300  # 5 minutes
```

### Database Issues

**Problem: Alert history corrupted**

**Solution:**
```bash
# Backup database
cp %APPDATA%\email-sentinel\history.db history.db.bak  # Windows

# Delete and restart (history will be recreated)
del %APPDATA%\email-sentinel\history.db  # Windows
email-sentinel start
```

**Problem: Seen messages database too large**

**Solution:**
```bash
# Clear seen messages (will re-process recent emails)
del %APPDATA%\email-sentinel\seen.db  # Windows
rm ~/.config/email-sentinel/seen.db   # Linux/macOS
```

---

## Quick Reference

### Essential Commands

```bash
# Setup
email-sentinel init

# Filters
email-sentinel filter add --name "Name" --from "sender.com"
email-sentinel filter list

# Start/Stop
email-sentinel start --tray
email-sentinel stop
email-sentinel status

# Alerts
email-sentinel alerts
email-sentinel alerts --recent 10

# OTP
email-sentinel otp list --active
email-sentinel otp get

# Config
email-sentinel config show
email-sentinel config set polling 45

# Testing
email-sentinel test desktop
email-sentinel test toast
email-sentinel test mobile

# Auto-start
email-sentinel install
email-sentinel uninstall
```

### File Locations (Windows)

```
%APPDATA%\email-sentinel\
├── config.yaml          # Main configuration
├── rules.yaml           # Priority rules
├── otp_rules.yaml       # OTP detection config
├── history.db           # Alert history & OTP codes
├── seen.db              # Processed emails
├── token.json           # Gmail OAuth token
└── credentials.json     # Gmail API credentials
```

### File Locations (macOS)

```
~/Library/Application Support/email-sentinel/
├── config.yaml
├── rules.yaml
├── otp_rules.yaml
├── history.db
├── seen.db
├── token.json
└── credentials.json
```

### File Locations (Linux)

```
~/.config/email-sentinel/
├── config.yaml
├── rules.yaml
├── otp_rules.yaml
├── history.db
├── seen.db
├── token.json
└── credentials.json
```

---

## Additional Resources

- **Complete Setup Guide:** [build_to_first_alert.md](build_to_first_alert.md)
- **Gmail API Setup:** [gmail_api_setup.md](gmail_api_setup.md)
- **Mobile Notifications:** [mobile_ntfy_setup.md](mobile_ntfy_setup.md)
- **Multi-Account Setup:** [central_email_setup.md](central_email_setup.md)
- **Main README:** [../README.md](../README.md)
- **GitHub Issues:** https://github.com/datateamsix/email-sentinel/issues

---

**Last Updated:** 2025-12-07
**Email Sentinel Version:** 1.0
