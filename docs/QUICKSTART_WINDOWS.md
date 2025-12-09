# Windows Quickstart Guide

![Email Sentinel Logo](../images/logo.png)

Complete guide for installing, configuring, and running Email Sentinel on Windows 10/11.

**Time Required:** 10-15 minutes
**Skill Level:** Beginner-friendly

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Initial Setup](#initial-setup)
4. [Configure Windows Notifications](#configure-windows-notifications)
5. [Test Notifications](#test-notifications)
6. [Start Monitoring](#start-monitoring)
7. [Set Up Auto-Start](#set-up-auto-start)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software

- **Windows 10** or **Windows 11**
- **PowerShell** (pre-installed on Windows)
- **Gmail account** to monitor

### Verify PowerShell

Open PowerShell (Win+X → Windows PowerShell):

```powershell
$PSVersionTable.PSVersion
# Should show version 5.1 or higher
```

---

## Installation

### Option 1: Scoop (Recommended)

Scoop is a command-line installer for Windows.

**Install Scoop (if not already installed):**

```powershell
# Run in PowerShell (as regular user, NOT admin)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
```

**Install Email Sentinel:**

```powershell
scoop bucket add datateamsix https://github.com/datateamsix/scoop-bucket
scoop install email-sentinel
```

**Verify installation:**

```powershell
email-sentinel --version
# Should show: Email Sentinel vX.X.X
```

### Option 2: Direct Download

**Download from GitHub:**

1. Go to [Releases](https://github.com/datateamsix/email-sentinel/releases/latest)
2. Download `email-sentinel_windows_amd64.zip`
3. Extract to a folder (e.g., `C:\email-sentinel\`)
4. Add to PATH or run from that folder

**Extract and run:**

```powershell
# Extract (replace with your download path)
Expand-Archive -Path "$env:USERPROFILE\Downloads\email-sentinel_windows_amd64.zip" -DestinationPath "C:\email-sentinel"

# Navigate to folder
cd C:\email-sentinel

# Verify
.\email-sentinel.exe --version
```

### Option 3: Build from Source

**Prerequisites:**
- Go 1.22+ ([Download](https://go.dev/dl/))
- Git ([Download](https://git-scm.com/downloads))

**Build:**

```powershell
# Clone repository
git clone https://github.com/datateamsix/email-sentinel.git
cd email-sentinel

# Build
set CGO_ENABLED=1
go build -o email-sentinel.exe .

# Verify
.\email-sentinel.exe --version
```

---

## Initial Setup

### Step 1: Gmail API Credentials

Follow the [Gmail API Setup Guide](gmail_api_setup.md) to:
1. Create a Google Cloud project
2. Enable Gmail API
3. Download `credentials.json`

**Place credentials.json:**

```powershell
# Create config directory
New-Item -ItemType Directory -Force -Path $env:APPDATA\email-sentinel

# Copy credentials (adjust path if needed)
Copy-Item "$env:USERPROFILE\Downloads\credentials.json" "$env:APPDATA\email-sentinel\"
```

### Step 2: Initialize Email Sentinel

```powershell
email-sentinel init
```

**What happens:**
1. Opens browser for Google OAuth
2. You'll see "Google hasn't verified this app" - this is normal
3. Click **"Advanced"** → **"Go to Email Sentinel (unsafe)"**
4. Click **"Allow"** to grant permissions
5. Copy authorization code and paste into PowerShell
6. Success! Token saved

### Step 3: Create Your First Filter

**Quick test filter (emails from yourself):**

```powershell
email-sentinel filter add `
  --name "Self Test" `
  --from "YOUR_EMAIL@gmail.com" `
  --labels "test"
```

**Replace `YOUR_EMAIL@gmail.com` with your actual Gmail address!**

**Real-world example:**

```powershell
# Job alerts
email-sentinel filter add `
  --name "Job Opportunities" `
  --from "linkedin.com,greenhouse.io,lever.co" `
  --subject "interview,opportunity,position" `
  --labels "work,career"
```

**Verify filters:**

```powershell
email-sentinel filter list
```

---

## Configure Windows Notifications

Email Sentinel uses Windows notifications for alerts. Ensure they're enabled:

### Enable Notifications (Windows 11)

1. Press **Win+I** to open Settings
2. Go to **System** → **Notifications**
3. Ensure **Notifications** toggle is ON
4. Scroll down and ensure notifications are enabled for:
   - **Windows PowerShell** (if running from PowerShell)
   - **Command Prompt** (if running from CMD)
   - Or the specific app if installed via Scoop

### Enable Notifications (Windows 10)

1. Press **Win+I** to open Settings
2. Go to **System** → **Notifications & actions**
3. Ensure **Get notifications from apps and other senders** is ON
4. Scroll down and enable notifications for PowerShell/CMD

### Enable Focus Assist Override

To ensure alerts come through even in Focus mode:

1. Settings → **System** → **Focus assist** (or **Notifications**)
2. Click **Priority only** or **Alarms only**
3. Add **Email Sentinel** to priority list

### Action Center Permissions

Windows Toast notifications appear in Action Center:

1. Settings → **System** → **Notifications**
2. Ensure **Show notifications in action center** is ON
3. Test with: `email-sentinel test toast`

---

## Test Notifications

Before monitoring, verify notifications work:

### Test Desktop Notification

```powershell
email-sentinel test desktop
```

**Expected:** Pop-up notification appears
**If not working:** Check notification settings above

### Test Windows Toast Notification

```powershell
# Normal priority
email-sentinel test toast

# High priority (urgent)
email-sentinel test toast --priority
```

**Expected:** Rich notification in Action Center with clickable link
**Features:**
- Shows sender, subject, preview
- Click to open Gmail link
- Persists in Action Center

### Test System Tray (Optional)

```powershell
email-sentinel start --tray
# Check system tray (bottom-right) for Email Sentinel icon
# Press Ctrl+C to stop
```

---

## Start Monitoring

### Option A: Foreground Mode (Recommended for First Time)

```powershell
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

### Option B: System Tray Mode (Recommended for Daily Use)

```powershell
email-sentinel start --tray
```

**Features:**
- Runs in background
- Icon in system tray (notification area)
- Right-click icon for menu:
  - Recent Alerts (click to open in Gmail)
  - Open History
  - Clear Alerts
  - Quit

**With AI Summaries (optional):**

```powershell
# Requires GEMINI_API_KEY environment variable
email-sentinel start --tray --ai-summary
```

### Option C: Background Daemon

```powershell
email-sentinel start --daemon
```

**Check status:**

```powershell
email-sentinel status
```

**Stop daemon:**

```powershell
email-sentinel stop
```

---

## Set Up Auto-Start

Make Email Sentinel start automatically when Windows boots.

### Install Auto-Start

```powershell
email-sentinel install
```

**What it does:**
- Creates a Windows Task Scheduler task named "EmailSentinel"
- Runs at user logon (no admin required)
- Starts with system tray enabled

### Verify Installation

**Check Task Scheduler:**

```powershell
# Open Task Scheduler GUI
taskschd.msc
```

Or check via PowerShell:

```powershell
Get-ScheduledTask -TaskName "EmailSentinel"
```

**Expected output:**
```
TaskName        State
--------        -----
EmailSentinel   Ready
```

### View Task Details

```powershell
Get-ScheduledTask -TaskName "EmailSentinel" | Format-List
```

### Test Auto-Start

**Manually trigger the task:**

```powershell
Start-ScheduledTask -TaskName "EmailSentinel"
```

**Check if running:**

```powershell
Get-Process | Where-Object {$_.Name -like "*email-sentinel*"}
```

**Reboot test:**
1. Restart your computer
2. After login, check system tray for Email Sentinel icon
3. Or run: `email-sentinel status`

### Uninstall Auto-Start (if needed)

```powershell
email-sentinel uninstall
```

---

## Troubleshooting

### Notifications Not Appearing

**Check Focus Assist:**
```powershell
# Temporarily disable Focus Assist
# Settings → System → Focus assist → Off
```

**Check notification settings:**
```powershell
# Run test
email-sentinel test desktop
email-sentinel test toast

# If no notification, check Settings → System → Notifications
```

**Re-enable notifications for PowerShell:**
1. Settings → System → Notifications
2. Find "Windows PowerShell" or "Windows Terminal"
3. Ensure notifications are enabled

### System Tray Icon Not Appearing

**Show hidden icons:**
1. Right-click taskbar
2. Taskbar settings → **Select which icons appear on the taskbar**
3. Enable "Email Sentinel"

**Restart explorer.exe:**
```powershell
Stop-Process -Name explorer -Force
# Explorer will restart automatically
```

### "credentials.json not found"

**Check if file exists:**

```powershell
Test-Path "$env:APPDATA\email-sentinel\credentials.json"
# Should return: True
```

**If False, copy credentials:**

```powershell
Copy-Item "$env:USERPROFILE\Downloads\credentials.json" "$env:APPDATA\email-sentinel\"
```

### Token Expired

**Re-authenticate:**

```powershell
email-sentinel init
```

### Build Fails (if building from source)

**Enable CGO:**

```powershell
$env:CGO_ENABLED=1
go build -o email-sentinel.exe .
```

**Install MinGW (if needed):**

```powershell
scoop install mingw
```

### Auto-Start Not Working After Reboot

**Check task status:**

```powershell
Get-ScheduledTask -TaskName "EmailSentinel" | Select-Object TaskName,State
```

**Re-install task:**

```powershell
email-sentinel uninstall
email-sentinel install
```

**Manually run task:**

```powershell
Start-ScheduledTask -TaskName "EmailSentinel"

# Wait 5 seconds, then check
Start-Sleep -Seconds 5
Get-Process | Where-Object {$_.Name -like "*email-sentinel*"}
```

---

## Quick Reference

### Common Commands

```powershell
# Initialize
email-sentinel init

# Add filter
email-sentinel filter add --name "Name" --from "sender.com" --labels "work"

# List filters
email-sentinel filter list

# Test notifications
email-sentinel test desktop
email-sentinel test toast

# Start monitoring
email-sentinel start              # Foreground
email-sentinel start --tray       # System tray
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
| Config | `%APPDATA%\email-sentinel\config.yaml` |
| Credentials | `%APPDATA%\email-sentinel\credentials.json` |
| Token | `%APPDATA%\email-sentinel\token.json` |
| Priority Rules | `%APPDATA%\email-sentinel\rules.yaml` |
| OTP Rules | `%APPDATA%\email-sentinel\otp_rules.yaml` |
| AI Config | `%APPDATA%\email-sentinel\ai-config.yaml` |
| Database | `%APPDATA%\email-sentinel\history.db` |

**Open config directory:**

```powershell
explorer $env:APPDATA\email-sentinel
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

**You're all set!** Email Sentinel is now monitoring your Gmail on Windows. 📧✨
