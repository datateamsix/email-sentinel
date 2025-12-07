# 🚀 Complete Setup Guide: Build to First Alert

This guide walks you through the complete process from building Email Sentinel to receiving your first email notification.

**Time Required:** ~15-20 minutes
**Skill Level:** Beginner-friendly

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Step 1: Set Up Central Email Account (Optional)](#step-1-set-up-central-email-account-optional)
3. [Step 2: Set Up Google Cloud & Gmail API](#step-2-set-up-google-cloud--gmail-api)
4. [Step 3: Build Email Sentinel](#step-3-build-email-sentinel)
5. [Step 4: Initialize & Authenticate](#step-4-initialize--authenticate)
6. [Step 5: Create Your First Filter](#step-5-create-your-first-filter)
7. [Step 6: Set Up Priority Rules](#step-6-set-up-priority-rules)
8. [Step 7: Test Notifications](#step-7-test-notifications)
9. [Step 8: Start Monitoring](#step-8-start-monitoring)
10. [Step 9: Trigger Your First Alert](#step-9-trigger-your-first-alert)
11. [Step 10: View Alert History](#step-10-view-alert-history)
12. [Step 11: Set Up Auto-Start](#step-11-set-up-auto-start)
13. [Troubleshooting](#troubleshooting)

---

## Prerequisites

Before starting, ensure you have:

### Required Software
- [ ] **Go 1.22+** installed ([Download](https://go.dev/dl/))
- [ ] **Git** installed ([Download](https://git-scm.com/downloads))
- [ ] **Gmail account** to monitor

### Required Accounts
- [ ] **Google Cloud Platform account** (free tier is sufficient)
  - Create at: https://console.cloud.google.com/

### Verification

**Windows (PowerShell):**
```powershell
# Check Go version
go version
# Should show: go version go1.22 or higher

# Check Git
git --version
# Should show: git version 2.x.x
```

**macOS/Linux (Bash):**
```bash
# Check Go version
go version
# Should show: go version go1.22 or higher

# Check Git
git --version
# Should show: git version 2.x.x
```

---

## Step 1: Set Up Central Email Account (Optional)

**🎯 Recommended for users monitoring multiple email accounts**

If you want to monitor emails from multiple accounts (personal Gmail, work email, Outlook, etc.), you can set up a **central "collector" Gmail account** that receives forwarded emails from all your other accounts.

### Why Use This Approach?

- ✅ Monitor multiple email accounts with a single Email Sentinel setup
- ✅ Only authenticate once (no multi-OAuth complexity)
- ✅ Forwarded messages retain original sender and subject metadata
- ✅ Simpler configuration and faster polling
- ✅ Multi Layer/Stage Filtering of Email 

### Quick Setup:

1. **Create or choose** a Gmail account to act as your central inbox
2. **Set up forwarding rules** in your other email accounts (Gmail, Outlook, Yahoo, etc.)
3. **Forward only important emails** using filters (job alerts, VIP senders, urgent keywords)

**See detailed guide:** [Central Email Setup (Collector Inbox)](central_email_setup.md)

### Skip This Step If:

- You only have one Gmail account to monitor
- You prefer to authenticate each account separately (future feature)

---

## Step 2: Set Up Google Cloud & Gmail API

**See detailed guide:** [Gmail API Setup](gmail_api_setup.md)

**Quick summary:**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project: "Email Sentinel"
3. Enable Gmail API
4. Configure OAuth Consent Screen (External)
5. Add yourself as test user
6. Create OAuth 2.0 credentials (Desktop app)
7. Download credentials as `credentials.json`

**Result:** You should have a file named `credentials.json`

---

## Step 3: Build Email Sentinel

### Option A: Clone from GitHub

**Windows (PowerShell):**
```powershell
# Clone the repository
git clone https://github.com/datateamsix/email-sentinel.git
cd email-sentinel

# Build the binary
go build -o email-sentinel.exe .
```

**macOS/Linux (Bash):**
```bash
# Clone the repository
git clone https://github.com/datateamsix/email-sentinel.git
cd email-sentinel

# Build the binary
go build -o email-sentinel .
```

### Option B: Build from Source (if you have the code locally)

**Windows (PowerShell):**
```powershell
# Navigate to the project directory
cd C:\path\to\email-sentinel

# Install dependencies
go mod download

# Build the binary
go build -o email-sentinel.exe .
```

**macOS/Linux (Bash):**
```bash
# Navigate to the project directory
cd /path/to/email-sentinel

# Install dependencies
go mod download

# Build the binary
go build -o email-sentinel .
```

### Verify Build

**Windows (PowerShell/CMD):**
```powershell
.\email-sentinel.exe --help
```

**macOS/Linux (Bash):**
```bash
./email-sentinel --help
```

**Expected output:**
```
Email Sentinel - Real-time Gmail monitoring with custom filters
...
Available Commands:
  alerts      View today's email alerts
  config      View and modify configuration
  filter      Manage email filters
  ...
```

---

## Step 4: Initialize & Authenticate

### Place credentials.json

Move your downloaded `credentials.json` to one of these locations:

**Option 1: Project directory (recommended for testing)**

**Windows (PowerShell):**
```powershell
# Copy from Downloads to current directory
Copy-Item $env:USERPROFILE\Downloads\credentials.json .
```

**Windows (CMD):**
```cmd
copy %USERPROFILE%\Downloads\credentials.json .
```

**macOS/Linux (Bash):**
```bash
cp ~/Downloads/credentials.json .
```

**Option 2: Config directory (recommended for production)**

**Windows (PowerShell):**
```powershell
# Create directory if it doesn't exist
New-Item -ItemType Directory -Force -Path $env:APPDATA\email-sentinel
# Copy credentials
Copy-Item credentials.json $env:APPDATA\email-sentinel\
```

**Windows (CMD):**
```cmd
mkdir %APPDATA%\email-sentinel
copy credentials.json %APPDATA%\email-sentinel\
```

**macOS (Bash):**
```bash
mkdir -p ~/Library/Application\ Support/email-sentinel/
cp credentials.json ~/Library/Application\ Support/email-sentinel/
```

**Linux (Bash):**
```bash
mkdir -p ~/.config/email-sentinel/
cp credentials.json ~/.config/email-sentinel/
```

### Run Initialization

**Windows (PowerShell/CMD):**
```powershell
.\email-sentinel.exe init
```

**macOS/Linux (Bash):**
```bash
./email-sentinel init
```

**What happens:**
1. Email Sentinel finds `credentials.json`
2. Opens browser to Google OAuth consent screen
3. You'll see: "Google hasn't verified this app"
   - Click **"Advanced"**
   - Click **"Go to Email Sentinel (unsafe)"**
   - This is normal for personal apps!
4. Click **"Allow"** to grant permissions
5. You'll see an authorization code
6. **Copy the code** and paste into terminal
7. Success! Token saved

**Expected output:**
```
🔐 Initializing Email Sentinel...

Found credentials: ./credentials.json

Opening browser for OAuth authentication...

Enter authorization code: [paste code here]

✅ Successfully authenticated!
   Token saved to: C:\Users\...\email-sentinel\token.json

Next steps:
  1. Add filters: email-sentinel filter add
  2. Start monitoring: email-sentinel start
```

---

## Step 5: Create Your First Filter

Let's create a simple test filter that will definitely trigger.

### Add a Self-Test Filter

**Windows (PowerShell/CMD):**
```powershell
.\email-sentinel.exe filter add `
  --name "Self Test" `
  --from "YOUR_EMAIL@gmail.com" `
  --labels "test"
```

**macOS/Linux (Bash):**
```bash
./email-sentinel filter add \
  --name "Self Test" \
  --from "YOUR_EMAIL@gmail.com" \
  --labels "test"
```

**Replace `YOUR_EMAIL@gmail.com` with your actual Gmail address!**

**Example (Windows):**
```powershell
.\email-sentinel.exe filter add --name "Self Test" --from "john.doe@gmail.com" --labels "test"
```

**Example (macOS/Linux):**
```bash
./email-sentinel filter add --name "Self Test" --from "john.doe@gmail.com" --labels "test"
```

**Note about Labels:**
- Labels help organize your filters (e.g., "work", "personal", "urgent")
- Labels appear in ALL notifications (desktop, mobile, toast)
- Once created, labels are saved and suggested when creating new filters
- Use `--labels "work,urgent"` for multiple labels (comma-separated)

### Verify Filter

**Windows:**
```powershell
.\email-sentinel.exe filter list
```

**macOS/Linux:**
```bash
./email-sentinel filter list
```

**Expected output:**
```
📋 Current Filters (1 total)

[1] Self Test
    From: john.doe@gmail.com
    Subject: (any)
    Match mode: any
```

### Add Real-World Filters (Optional)

**Windows (PowerShell/CMD):**
```powershell
# Job alerts
.\email-sentinel.exe filter add `
  --name "Job Opportunities" `
  --from "linkedin.com,greenhouse.io,lever.co" `
  --subject "interview,opportunity,position"

# Important people
.\email-sentinel.exe filter add `
  --name "VIP Emails" `
  --from "boss@company.com,client@important.com"

# Urgent keywords
.\email-sentinel.exe filter add `
  --name "Urgent" `
  --subject "urgent,asap,critical,emergency"
```

**macOS/Linux (Bash):**
```bash
# Job alerts
./email-sentinel filter add \
  --name "Job Opportunities" \
  --from "linkedin.com,greenhouse.io,lever.co" \
  --subject "interview,opportunity,position"

# Important people
./email-sentinel filter add \
  --name "VIP Emails" \
  --from "boss@company.com,client@important.com"

# Urgent keywords
./email-sentinel filter add \
  --name "Urgent" \
  --subject "urgent,asap,critical,emergency"
```

---

## Step 6: Set Up Priority Rules

Priority rules automatically mark emails as urgent (🔥) or normal (📧).

### Locate rules.yaml

The file will be auto-created on first `start`, but you can create it now:

**Windows:** `%APPDATA%\email-sentinel\rules.yaml`
**macOS:** `~/Library/Application Support/email-sentinel/rules.yaml`
**Linux:** `~/.config/email-sentinel/rules.yaml`

### Edit Priority Rules

**Windows (PowerShell):**
```powershell
notepad $env:APPDATA\email-sentinel\rules.yaml
```

**Windows (CMD):**
```cmd
notepad %APPDATA%\email-sentinel\rules.yaml
```

**macOS:**
```bash
open ~/Library/Application\ Support/email-sentinel/rules.yaml
```

**Linux:**
```bash
nano ~/.config/email-sentinel/rules.yaml
# or
vim ~/.config/email-sentinel/rules.yaml
```

### Example Configuration

```yaml
priority_rules:
  urgent_keywords:
    - urgent
    - asap
    - action required
    - deadline
    - invoice
    - critical
    - emergency

  vip_senders:
    - boss@company.com
    - ceo@company.com

  vip_domains:
    - importantclient.com

notification_settings:
  quiet_hours_start: ""
  quiet_hours_end: ""
  weekend_mode: normal
```

**Save the file.**

---

## Step 7: Test Notifications

Before starting monitoring, verify notifications work.

### Test Desktop Notifications

**Windows:**
```powershell
.\email-sentinel.exe test desktop
```

**macOS/Linux:**
```bash
./email-sentinel test desktop
```

**Expected:**
- Notification pops up on your screen
- Title: "Email Sentinel Test"
- Message: "If you can see this, desktop notifications are working! ✅"

**If no notification appears:**
- **Windows:** Settings → System → Notifications → Enable
- **macOS:** System Preferences → Notifications → Terminal → Allow
- **Linux:** Install `libnotify-bin` or `dunst`

### Test Windows Toast Notifications (Windows Only)

**Windows (PowerShell/CMD):**
```powershell
# Test normal priority
.\email-sentinel.exe test toast

# Test high priority
.\email-sentinel.exe test toast --priority
```

**Expected:**
- Rich notification in Windows Action Center
- Clickable link
- Shows subject, sender, preview

### Test Mobile Notifications (Optional)

**First:** Set up ntfy.sh
**See:** [Mobile ntfy Setup Guide](mobile_ntfy_setup.md)

**Windows:**
```powershell
# Configure ntfy
.\email-sentinel.exe config set mobile true
.\email-sentinel.exe config set ntfy_topic "your-unique-topic"

# Test
.\email-sentinel.exe test mobile
```

**macOS/Linux:**
```bash
# Configure ntfy
./email-sentinel config set mobile true
./email-sentinel config set ntfy_topic "your-unique-topic"

# Test
./email-sentinel test mobile
```

---

## Step 8: Start Monitoring

Now start Email Sentinel to monitor your inbox.

### Option A: Foreground Mode (Recommended for First Time)

**Windows:**
```powershell
.\email-sentinel.exe start
```

**macOS/Linux:**
```bash
./email-sentinel start
```

**What you'll see:**
```
✅ Email Sentinel Started
   Monitoring 1 filter(s)
   Polling interval: 45 seconds
   Desktop notifications: enabled

🔍 Watching for new emails... (Press Ctrl+C to stop)

[14:23:05] Checked 10 messages, no new matches
[14:23:50] Checked 10 messages, no new matches
...
```

### Option B: With System Tray (Recommended for Daily Use)

**Windows:**
```powershell
.\email-sentinel.exe start --tray
```

**macOS/Linux:**
```bash
./email-sentinel start --tray
```

**What happens:**
- System tray icon appears in taskbar
- Right-click icon for menu
- Monitoring runs in background
- Recent alerts shown in tray menu

---

## Step 9: Trigger Your First Alert

Now let's trigger an alert by sending yourself an email.

### Send Test Email

**Method 1: Gmail Web Interface**
1. Open Gmail in browser
2. Click **"Compose"**
3. **To:** YOUR_EMAIL@gmail.com (same email you're monitoring)
4. **Subject:** Test Alert
5. **Body:** This is a test email
6. Click **"Send"**

**Method 2: From Another Email Account**
- Send an email to your Gmail address
- From address should match your filter

### Wait for Detection

- **Default polling interval:** 45 seconds
- **Check terminal output:**

**Expected output:**
```
[14:25:15] Checked 10 messages, no new matches
[14:26:00] Checked 10 messages, no new matches
📧 MATCH [Self Test] From: john.doe@gmail.com | Subject: Test Alert
[14:26:45] Checked 10 messages, no new matches
```

### Verify Notifications

You should receive:
1. ✅ **Desktop notification** (pop-up)
2. ✅ **Windows toast** (Action Center) - if on Windows
3. ✅ **Mobile push** (phone) - if configured
4. ✅ **System tray update** (if using --tray)

**Notification shows:**
- Filter name: "Self Test"
- Sender: your-email@gmail.com
- Subject: Test Alert
- (Windows toast includes snippet + Gmail link)

---

## Step 10: View Alert History

All alerts are saved to a SQLite database.

### View Today's Alerts

**Windows:**
```powershell
.\email-sentinel.exe alerts
```

**macOS/Linux:**
```bash
./email-sentinel alerts
```

**Expected output:**
```
📬 Today's Alerts (1 total)

[1] 📧 2025-12-06 14:26:00
    Filter: Self Test
    From:   john.doe@gmail.com
    Subject: Test Alert
    Preview: This is a test email
    Link:   https://mail.google.com/mail/u/0/#all/abc123xyz
```

### View Last 5 Alerts

**Windows:**
```powershell
.\email-sentinel.exe alerts --recent 5
```

**macOS/Linux:**
```bash
./email-sentinel alerts --recent 5
```

### Click to Open Email

- Copy the **Link** from the output
- Paste into browser
- Gmail opens directly to that email

---

## Step 11: Set Up Auto-Start

Make Email Sentinel start automatically when your computer boots.

### Install Auto-Start

**Windows:**
```powershell
.\email-sentinel.exe install
```

**macOS/Linux:**
```bash
./email-sentinel install
```

**What it does:**
- **Windows:** Creates Task Scheduler task
- **macOS:** Creates LaunchAgent
- **Linux:** Creates systemd service

### Verify Installation

**Windows (PowerShell):**
```powershell
# Check Task Scheduler
taskschd.msc

# OR check if task exists
Get-ScheduledTask -TaskName "EmailSentinel"
```

**macOS:**
```bash
launchctl list | grep email-sentinel
```

**Linux:**
```bash
systemctl --user status email-sentinel
```

### Reboot Test

1. Restart your computer
2. After login, check if Email Sentinel is running:

**Windows (PowerShell):**
```powershell
# Check if process exists
Get-Process | Where-Object {$_.Name -like "*email-sentinel*"}

# OR using tasklist
tasklist | findstr email-sentinel
```

**macOS/Linux:**
```bash
# Check process
ps aux | grep email-sentinel

# Or check system tray for icon
```

### Uninstall Auto-Start (if needed)

**Windows:**
```powershell
.\email-sentinel.exe uninstall
```

**macOS/Linux:**
```bash
./email-sentinel uninstall
```

---

## Congratulations! 🎉

You've successfully set up Email Sentinel from build to first alert!

### What You've Accomplished:

- ✅ Built Email Sentinel from source
- ✅ Configured Gmail API authentication
- ✅ Created email filters
- ✅ Set up priority rules
- ✅ Tested all notification systems
- ✅ Received your first alert
- ✅ Viewed alert history
- ✅ Configured auto-start

### Next Steps:

1. **Add Real Filters** - Replace the test filter with actual use cases

   **Windows:**
   ```powershell
   .\email-sentinel.exe filter add --name "Job Alerts" --from "linkedin.com,indeed.com"
   ```

   **macOS/Linux:**
   ```bash
   ./email-sentinel filter add --name "Job Alerts" --from "linkedin.com,indeed.com"
   ```

2. **Customize Priority Rules** - Edit `rules.yaml` with your VIPs
   ```yaml
   vip_senders:
     - boss@company.com
   ```

3. **Run with System Tray** - Use `--tray` for background mode

   **Windows:**
   ```powershell
   .\email-sentinel.exe start --tray
   ```

   **macOS/Linux:**
   ```bash
   ./email-sentinel start --tray
   ```

4. **Set Up Mobile Notifications** - Follow [mobile_ntfy_setup.md](mobile_ntfy_setup.md)

---

## Troubleshooting

### Problem: "credentials.json not found"

**Solution:**

**Windows (PowerShell):**
```powershell
# Check if file exists in current directory
Test-Path credentials.json

# List files
Get-ChildItem credentials.json
```

**Windows (CMD):**
```cmd
dir credentials.json
```

**macOS/Linux:**
```bash
# Check if file exists
ls credentials.json

# If not, copy it to the project directory
```

### Problem: "Token has expired"

**Solution:**

**Windows:**
```powershell
# Re-run authentication
.\email-sentinel.exe init
```

**macOS/Linux:**
```bash
# Re-run authentication
./email-sentinel init
```

### Problem: No alerts showing up

**Check 1: Is monitoring running?**

**Windows:**
```powershell
.\email-sentinel.exe status
```

**macOS/Linux:**
```bash
./email-sentinel status
```

**Check 2: Do filters match?**

**Windows:**
```powershell
.\email-sentinel.exe filter list

# Test filter matching
.\email-sentinel.exe test filter "Filter Name" "sender@example.com" "subject text"
```

**macOS/Linux:**
```bash
./email-sentinel filter list

# Test filter matching
./email-sentinel test filter "Filter Name" "sender@example.com" "subject text"
```

**Check 3: Check seen messages**
- Email might have already been marked as seen
- Try sending a NEW email

**Check 4: Check polling interval**

**Windows:**
```powershell
.\email-sentinel.exe config show
```

**macOS/Linux:**
```bash
./email-sentinel config show
```
- Default is 45 seconds - wait at least this long

### Problem: Notifications not appearing

**Windows:**
```powershell
# Test each notification type
.\email-sentinel.exe test desktop
.\email-sentinel.exe test toast
.\email-sentinel.exe test mobile
```

**macOS/Linux:**
```bash
# Test each notification type
./email-sentinel test desktop
./email-sentinel test mobile
```

See notification-specific troubleshooting in each test command output.

### Problem: Build fails

**Check Go version:**
```bash
go version
# Must be 1.22 or higher
```

**Clean and rebuild:**
```bash
go clean
go mod tidy
go build -o email-sentinel.exe .
```

---

## Quick Reference Commands

**Windows (PowerShell/CMD):**
```powershell
# Initialize
.\email-sentinel.exe init

# Add filter
.\email-sentinel.exe filter add --name "Name" --from "sender.com"

# List filters
.\email-sentinel.exe filter list

# Test notifications
.\email-sentinel.exe test desktop
.\email-sentinel.exe test toast

# Start monitoring
.\email-sentinel.exe start              # Foreground
.\email-sentinel.exe start --tray       # With system tray

# View alerts
.\email-sentinel.exe alerts

# Check status
.\email-sentinel.exe status
.\email-sentinel.exe config show

# Install auto-start
.\email-sentinel.exe install
```

**macOS/Linux (Bash):**
```bash
# Initialize
./email-sentinel init

# Add filter
./email-sentinel filter add --name "Name" --from "sender.com"

# List filters
./email-sentinel filter list

# Test notifications
./email-sentinel test desktop

# Start monitoring
./email-sentinel start              # Foreground
./email-sentinel start --tray       # With system tray

# View alerts
./email-sentinel alerts

# Check status
./email-sentinel status
./email-sentinel config show

# Install auto-start
./email-sentinel install
```

---

## Need More Help?

- **Gmail API Setup:** See [gmail_api_setup.md](gmail_api_setup.md)
- **Mobile Notifications:** See [mobile_ntfy_setup.md](mobile_ntfy_setup.md)
- **Report Issues:** https://github.com/datateamsix/email-sentinel/issues

**You're all set!** Email Sentinel is now monitoring your Gmail and you'll receive instant notifications for important emails. 📧✨
