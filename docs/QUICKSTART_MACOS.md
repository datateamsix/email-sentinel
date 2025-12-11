# macOS Quickstart Guide

![Email Sentinel Logo](./images/logo.png)

Complete guide for installing, configuring, and running Email Sentinel on macOS.

**Time Required:** 10-15 minutes
**Skill Level:** Beginner-friendly

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Initial Setup](#initial-setup)
4. [Configure macOS Notifications](#configure-macos-notifications)
5. [Test Notifications](#test-notifications)
6. [Start Monitoring](#start-monitoring)
7. [Set Up Auto-Start](#set-up-auto-start)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software

- **macOS 11 (Big Sur)** or later
- **Terminal** (pre-installed)
- **Gmail account** to monitor

### Verify Terminal

Open Terminal (Cmd+Space → "Terminal"):

```bash
echo $SHELL
# Should show: /bin/zsh or /bin/bash
```

---

## Installation

### Option 1: Homebrew (Recommended)

Homebrew is the package manager for macOS.

**Install Homebrew (if not already installed):**

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

**Install Email Sentinel:**

```bash
brew tap datateamsix/tap
brew install email-sentinel
```

**Verify installation:**

```bash
email-sentinel --version
# Should show: Email Sentinel vX.X.X
```

### Option 2: Direct Download

**Download from GitHub:**

1. Go to [Releases](https://github.com/datateamsix/email-sentinel/releases/latest)
2. Download the appropriate binary:
   - **Apple Silicon (M1/M2/M3):** `email-sentinel_darwin_arm64.tar.gz`
   - **Intel Mac:** `email-sentinel_darwin_amd64.tar.gz`
   - **Universal:** `email-sentinel_darwin_universal.tar.gz` (works on both)

**Extract and install:**

```bash
# Navigate to Downloads
cd ~/Downloads

# Extract (replace with your downloaded file)
tar -xzf email-sentinel_darwin_universal.tar.gz

# Move to system location
sudo mv email-sentinel /usr/local/bin/

# Verify
email-sentinel --version
```

### Option 3: Build from Source

**Prerequisites:**
- Go 1.22+ ([Download](https://go.dev/dl/))
- Xcode Command Line Tools

**Install Xcode Command Line Tools:**

```bash
xcode-select --install
```

**Build:**

```bash
# Clone repository
git clone https://github.com/datateamsix/email-sentinel.git
cd email-sentinel

# Build
export CGO_ENABLED=1
go build -o email-sentinel .

# Verify
./email-sentinel --version

# Optional: Move to PATH
sudo mv email-sentinel /usr/local/bin/
```

---

## Initial Setup

### Step 1: Gmail API Credentials

Follow the [Gmail API Setup Guide](gmail_api_setup.md) to:
1. Create a Google Cloud project
2. Enable Gmail API
3. Download `credentials.json`

**Place credentials.json:**

```bash
# Create config directory
mkdir -p ~/Library/Application\ Support/email-sentinel/

# Copy credentials (adjust path if needed)
cp ~/Downloads/credentials.json ~/Library/Application\ Support/email-sentinel/
```

### Step 2: Initialize Email Sentinel

```bash
email-sentinel init
```

**What happens:**
1. Opens browser for Google OAuth
2. You'll see "Google hasn't verified this app" - this is normal
3. Click **"Advanced"** → **"Go to Email Sentinel (unsafe)"**
4. Click **"Allow"** to grant permissions
5. Copy authorization code and paste into Terminal
6. Success! Token saved

### Step 3: Create Your First Filter

**Quick test filter (emails from yourself):**

```bash
email-sentinel filter add \
  --name "Self Test" \
  --from "YOUR_EMAIL@gmail.com" \
  --labels "test"
```

**Replace `YOUR_EMAIL@gmail.com` with your actual Gmail address!**

**Real-world example:**

```bash
# Job alerts
email-sentinel filter add \
  --name "Job Opportunities" \
  --from "linkedin.com,greenhouse.io,lever.co" \
  --subject "interview,opportunity,position" \
  --labels "work,career"
```

**Verify filters:**

```bash
email-sentinel filter list
```

---

## Configure macOS Notifications

Email Sentinel uses macOS native notifications. Ensure they're enabled:

### Enable Notifications (macOS Ventura 13+)

1. Open **System Settings** (Apple menu → System Settings)
2. Click **Notifications**
3. Scroll down to find **Terminal** (or the app you're using)
4. Enable:
   - ✅ **Allow notifications**
   - ✅ **Banners** (or Alerts for persistent notifications)
   - ✅ **Play sound for notifications**
   - ✅ **Show in Notification Center**

### Enable Notifications (macOS Monterey 12 and earlier)

1. Open **System Preferences** (Apple menu → System Preferences)
2. Click **Notifications & Focus**
3. Find **Terminal** in the list
4. Set alert style to **Banners** or **Alerts**
5. Enable:
   - ✅ **Allow Notifications**
   - ✅ **Show in Notification Center**
   - ✅ **Play sound for notifications**

### Grant Terminal Full Disk Access (Optional)

If Email Sentinel needs to access certain files:

1. System Settings → **Privacy & Security**
2. Click **Full Disk Access**
3. Click the **+** button
4. Add **Terminal** (or your terminal app)

### Do Not Disturb Settings

Ensure Email Sentinel notifications come through:

1. System Settings → **Focus**
2. Click **Do Not Disturb**
3. Under **Allow notifications:**
   - Add **Terminal** to allowed apps
   - Or set to allow all apps

---

## Test Notifications

Before monitoring, verify notifications work:

### Test Desktop Notification

```bash
email-sentinel test desktop
```

**Expected:** Notification banner appears in top-right corner
**If not working:** Check notification settings above

### Test with Custom Message

```bash
email-sentinel test desktop
```

**Expected notification:**
- **Title:** Email Sentinel Test
- **Message:** If you can see this, desktop notifications are working! ✅

### Test System Tray (Menu Bar)

```bash
email-sentinel start --tray
# Check menu bar (top-right) for Email Sentinel icon
# Press Ctrl+C in Terminal to stop
```

**Features:**
- Icon appears in menu bar
- Click to see recent alerts
- Click alerts to open in Gmail

---

## Start Monitoring

### Option A: Foreground Mode (Recommended for First Time)

```bash
email-sentinel start
```

**What you'll see:**
```
✅ Email Sentinel Started
   Monitoring 1 filter(s)
   Polling interval: 45 seconds
   Desktop notifications: enabled

🔍 Watching for new emails... (Press Ctrl+C to stop)
```

**Send yourself a test email to trigger an alert!**

### Option B: Menu Bar Mode (Recommended for Daily Use)

```bash
email-sentinel start --tray
```

**Features:**
- Runs in background
- Icon in menu bar (top-right)
- Click icon for menu:
  - Recent Alerts (click to open in Gmail)
  - Open History
  - Clear Alerts
  - Quit

**With AI Summaries (optional):**

```bash
# Requires GEMINI_API_KEY environment variable
export GEMINI_API_KEY="your-api-key-here"
email-sentinel start --tray --ai-summary
```

### Option C: Background Daemon

```bash
email-sentinel start --daemon
```

**Check status:**

```bash
email-sentinel status
```

**Stop daemon:**

```bash
email-sentinel stop
```

---

## Set Up Auto-Start

Make Email Sentinel start automatically when you log in.

### Install Auto-Start

```bash
email-sentinel install
```

**What it does:**
- Creates a LaunchAgent plist file
- Located at: `~/Library/LaunchAgents/com.email-sentinel.plist`
- Runs automatically at login
- Starts with menu bar icon enabled

### Verify Installation

**Check if LaunchAgent exists:**

```bash
ls -la ~/Library/LaunchAgents/com.email-sentinel.plist
```

**Check if service is loaded:**

```bash
launchctl list | grep email-sentinel
```

**Expected output:**
```
-    0    com.email-sentinel
```

### Load LaunchAgent Manually (if needed)

```bash
launchctl load ~/Library/LaunchAgents/com.email-sentinel.plist
```

### Unload LaunchAgent

```bash
launchctl unload ~/Library/LaunchAgents/com.email-sentinel.plist
```

### Test Auto-Start

**Reboot test:**
1. Restart your Mac
2. After login, check menu bar for Email Sentinel icon
3. Or run: `email-sentinel status`

**Check running process:**

```bash
ps aux | grep email-sentinel | grep -v grep
```

### Uninstall Auto-Start (if needed)

```bash
email-sentinel uninstall
```

---

## Troubleshooting

### Notifications Not Appearing

**Grant notification permissions:**
1. System Settings → Notifications
2. Find Terminal (or your terminal app)
3. Enable all notification options

**Test notification manually:**

```bash
# Use macOS native notification test
osascript -e 'display notification "Test" with title "Email Sentinel"'
```

**Restart Notification Center:**

```bash
killall NotificationCenter
# NotificationCenter will restart automatically
```

### Menu Bar Icon Not Appearing

**Check if process is running:**

```bash
ps aux | grep email-sentinel | grep -v grep
```

**Restart Email Sentinel:**

```bash
# Stop
pkill email-sentinel

# Start
email-sentinel start --tray
```

**Grant accessibility permissions:**

Some macOS versions require accessibility permissions for menu bar apps:
1. System Settings → Privacy & Security → Accessibility
2. Add Terminal (or your terminal app)

### "credentials.json not found"

**Check if file exists:**

```bash
ls -la ~/Library/Application\ Support/email-sentinel/credentials.json
```

**If not found, copy credentials:**

```bash
cp ~/Downloads/credentials.json ~/Library/Application\ Support/email-sentinel/
```

### Token Expired

**Re-authenticate:**

```bash
email-sentinel init
```

### Build Fails (if building from source)

**Install Xcode Command Line Tools:**

```bash
xcode-select --install
```

**Enable CGO:**

```bash
export CGO_ENABLED=1
go build -o email-sentinel .
```

**Check Go version:**

```bash
go version
# Must be 1.22 or higher
```

### LaunchAgent Not Starting After Reboot

**Check service status:**

```bash
launchctl list | grep email-sentinel
```

**View service logs:**

```bash
tail -f ~/Library/Logs/email-sentinel.log
```

**Re-install LaunchAgent:**

```bash
email-sentinel uninstall
email-sentinel install
```

**Manually load:**

```bash
launchctl load ~/Library/LaunchAgents/com.email-sentinel.plist
```

### Permission Denied Errors

**Make binary executable:**

```bash
chmod +x /usr/local/bin/email-sentinel
```

**Check file ownership:**

```bash
ls -la /usr/local/bin/email-sentinel
```

---

## Quick Reference

### Common Commands

```bash
# Initialize
email-sentinel init

# Add filter
email-sentinel filter add --name "Name" --from "sender.com" --labels "work"

# List filters
email-sentinel filter list

# Test notifications
email-sentinel test desktop

# Start monitoring
email-sentinel start              # Foreground
email-sentinel start --tray       # Menu bar
email-sentinel start --daemon     # Background daemon

# View alerts
email-sentinel alerts
email-sentinel alerts --recent 5

# OTP codes
email-sentinel otp list
email-sentinel otp get

# Check status
email-sentinel status

# Install auto-start
email-sentinel install
```

### Config Locations

| File | Location |
|------|----------|
| Config | `~/Library/Application Support/email-sentinel/config.yaml` |
| Credentials | `~/Library/Application Support/email-sentinel/credentials.json` |
| Token | `~/Library/Application Support/email-sentinel/token.json` |
| Priority Rules | `~/Library/Application Support/email-sentinel/rules.yaml` |
| OTP Rules | `~/Library/Application Support/email-sentinel/otp_rules.yaml` |
| AI Config | `~/Library/Application Support/email-sentinel/ai-config.yaml` |
| Database | `~/Library/Application Support/email-sentinel/history.db` |
| LaunchAgent | `~/Library/LaunchAgents/com.email-sentinel.plist` |

**Open config directory:**

```bash
open ~/Library/Application\ Support/email-sentinel/
```

### Environment Variables

Add to `~/.zshrc` or `~/.bash_profile`:

```bash
# AI Summaries (optional)
export GEMINI_API_KEY="your-api-key-here"
export ANTHROPIC_API_KEY="your-api-key-here"
export OPENAI_API_KEY="your-api-key-here"

# Reload shell
source ~/.zshrc  # or source ~/.bash_profile
```

---

## Next Steps

1. **Add real filters** for job alerts, VIP senders, urgent keywords
2. **Set up mobile notifications** - See [Mobile ntfy Setup](mobile_ntfy_setup.md)
3. **Enable AI summaries** - See [AI Email Summaries](../README.md#-ai-email-summaries)
4. **Configure OTP detection** - See [OTP/2FA Detection](../README.md#-otp2fa-code-detection)

---

## Need Help?

- **Main README:** [../README.md](../README.md)
- **Gmail API Setup:** [gmail_api_setup.md](gmail_api_setup.md)
- **Complete Build Guide:** [build_to_first_alert.md](build_to_first_alert.md)
- **Report Issues:** https://github.com/datateamsix/email-sentinel/issues

**You're all set!** Email Sentinel is now monitoring your Gmail on macOS. 📧✨
