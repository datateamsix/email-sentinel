# Linux Quickstart Guide

![Email Sentinel Logo](./images/logo.png)

Complete guide for installing, configuring, and running Email Sentinel on Linux distributions.

**Time Required:** 10-15 minutes
**Skill Level:** Beginner-friendly
**Supported:** Ubuntu, Debian, Fedora, RHEL, Arch, and most modern distributions

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Initial Setup](#initial-setup)
4. [Configure Linux Notifications](#configure-linux-notifications)
5. [Test Notifications](#test-notifications)
6. [Start Monitoring](#start-monitoring)
7. [Set Up Auto-Start](#set-up-auto-start)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software

- **Modern Linux distribution** (Ubuntu 20.04+, Debian 11+, Fedora 35+, etc.)
- **gcc** compiler (for CGO support)
- **Gmail account** to monitor

### Minimal Build Dependencies

**Good news!** As of Email Sentinel v1.0.0, we've migrated to `fyne.io/systray`, which means:
- ✅ **NO GTK3 headers required** (unlike older versions)
- ✅ **NO libayatana-appindicator3 needed**
- ✅ **Only gcc required** for CGO support

### Install Build Dependencies

**Ubuntu/Debian:**

```bash
sudo apt-get update
sudo apt-get install gcc
```

**Fedora/RHEL/CentOS:**

```bash
sudo dnf install gcc
```

**Arch Linux:**

```bash
sudo pacman -S gcc
```

**Alpine Linux:**

```bash
sudo apk add gcc musl-dev
```

---

## Installation

### Option 1: Package Manager (Recommended)

#### Debian/Ubuntu (.deb package)

```bash
# Download latest release
wget https://github.com/datateamsix/email-sentinel/releases/latest/download/email-sentinel_*_amd64.deb

# Install
sudo dpkg -i email-sentinel_*_amd64.deb

# Verify
email-sentinel --version
```

#### Fedora/RHEL/CentOS (.rpm package)

```bash
# Download latest release
wget https://github.com/datateamsix/email-sentinel/releases/latest/download/email-sentinel_*_x86_64.rpm

# Install
sudo rpm -i email-sentinel_*_x86_64.rpm

# Or with dnf
sudo dnf install email-sentinel_*_x86_64.rpm

# Verify
email-sentinel --version
```

#### Arch Linux (AUR)

```bash
# Using yay
yay -S email-sentinel

# Or using paru
paru -S email-sentinel

# Verify
email-sentinel --version
```

### Option 2: Direct Download

```bash
# Download for your architecture
# AMD64 (most common):
wget https://github.com/datateamsix/email-sentinel/releases/latest/download/email-sentinel_linux_amd64.tar.gz

# ARM64 (Raspberry Pi 4, etc.):
wget https://github.com/datateamsix/email-sentinel/releases/latest/download/email-sentinel_linux_arm64.tar.gz

# Extract
tar -xzf email-sentinel_linux_*.tar.gz

# Move to system location
sudo mv email-sentinel /usr/local/bin/

# Make executable
sudo chmod +x /usr/local/bin/email-sentinel

# Verify
email-sentinel --version
```

### Option 3: Build from Source

**Prerequisites:**
- Go 1.22+ ([Download](https://go.dev/dl/))
- gcc compiler

**Install Go (if not already installed):**

```bash
# Ubuntu/Debian
sudo apt-get install golang-go

# Fedora/RHEL
sudo dnf install golang

# Arch Linux
sudo pacman -S go

# Or download from: https://go.dev/dl/
```

**Build Email Sentinel:**

```bash
# Clone repository
git clone https://github.com/datateamsix/email-sentinel.git
cd email-sentinel

# Build with CGO enabled
export CGO_ENABLED=1
go build -o email-sentinel .

# Verify
./email-sentinel --version

# Optional: Install system-wide
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
mkdir -p ~/.config/email-sentinel/

# Copy credentials (adjust path if needed)
cp ~/Downloads/credentials.json ~/.config/email-sentinel/
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
5. Copy authorization code and paste into terminal
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

## Configure Linux Notifications

Email Sentinel uses Linux native notifications via DBus. Configuration varies by desktop environment:

### GNOME (Ubuntu, Fedora Workstation)

**Install notification daemon (usually pre-installed):**

```bash
# Ubuntu/Debian
sudo apt-get install libnotify-bin

# Fedora/RHEL
sudo dnf install libnotify

# Arch Linux
sudo pacman -S libnotify
```

**Test notification system:**

```bash
notify-send "Test" "Notifications are working!"
```

**Enable system tray extension (for tray icon):**

```bash
# Ubuntu/Debian
sudo apt-get install gnome-shell-extension-appindicator

# Enable the extension
gnome-extensions enable appindicatorsupport@rgcjonas.gmail.com
```

**Restart GNOME Shell:**
- Press `Alt+F2`
- Type `r`
- Press `Enter`

### KDE Plasma

**Notifications usually work out of the box.**

**Test notification:**

```bash
notify-send "Test" "Notifications are working!"
```

**System tray:**
- Right-click on system tray
- Settings → Configure System Tray
- Ensure "Status Notifier Items" is enabled

### XFCE

**Install notification daemon:**

```bash
# Ubuntu/Debian
sudo apt-get install xfce4-notifyd

# Fedora/RHEL
sudo dnf install xfce4-notifyd

# Arch Linux
sudo pacman -S xfce4-notifyd
```

**Test notification:**

```bash
notify-send "Test" "Notifications are working!"
```

### i3 / Sway (Tiling Window Managers)

**Install notification daemon:**

```bash
# Install dunst (lightweight notification daemon)
# Ubuntu/Debian
sudo apt-get install dunst

# Fedora/RHEL
sudo dnf install dunst

# Arch Linux
sudo pacman -S dunst
```

**Start dunst:**

```bash
dunst &
```

**Auto-start dunst (add to i3/sway config):**

```bash
# Add to ~/.config/i3/config or ~/.config/sway/config
exec --no-startup-id dunst
```

### Do Not Disturb / Focus Mode

**GNOME:**
```bash
# Disable Do Not Disturb
gsettings set org.gnome.desktop.notifications show-banners true
```

**KDE:**
- System Settings → Notifications → Do Not Disturb → Off

---

## Test Notifications

Before monitoring, verify notifications work:

### Test Desktop Notification

```bash
email-sentinel test desktop
```

**Expected:** Notification appears (top-right on most DEs)
**If not working:** Check notification daemon is running

### Test System Notification Daemon

```bash
# Generic Linux notification test
notify-send "Email Sentinel Test" "If you can see this, notifications work!"
```

### Check Notification Daemon Status

```bash
# Check if notification daemon is running
ps aux | grep -E "dunst|notification"

# For GNOME
systemctl --user status gvfs-udisks2-volume-monitor

# Start notification daemon if needed
# GNOME: Usually automatic
# KDE: Usually automatic
# XFCE: xfce4-notifyd &
# i3/Sway: dunst &
```

### Test System Tray

```bash
email-sentinel start --tray
# Check system tray area for Email Sentinel icon
# Press Ctrl+C to stop
```

**Note:** System tray support varies by desktop environment. GNOME requires the AppIndicator extension (see above).

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

### Option B: System Tray Mode (Recommended for Daily Use)

```bash
email-sentinel start --tray
```

**Features:**
- Runs in background
- Icon in system tray (if desktop environment supports it)
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

### Systemd User Service (Recommended)

**Install auto-start:**

```bash
email-sentinel install
```

**What it does:**
- Creates a systemd user service
- Located at: `~/.config/systemd/user/email-sentinel.service`
- Runs automatically at login
- Starts with system tray enabled

**Verify installation:**

```bash
# Check service status
systemctl --user status email-sentinel
```

**Expected output:**
```
● email-sentinel.service - Email Sentinel
     Loaded: loaded
     Active: active (running)
```

**Enable service (if not already enabled):**

```bash
systemctl --user enable email-sentinel
```

**Start service manually:**

```bash
systemctl --user start email-sentinel
```

**View logs:**

```bash
# Real-time logs
journalctl --user -u email-sentinel -f

# Last 50 lines
journalctl --user -u email-sentinel -n 50
```

**Stop service:**

```bash
systemctl --user stop email-sentinel
```

**Restart service:**

```bash
systemctl --user restart email-sentinel
```

### Desktop Environment Auto-Start (Alternative)

If systemd is not available or you prefer desktop auto-start:

**Create .desktop file:**

```bash
mkdir -p ~/.config/autostart/

cat > ~/.config/autostart/email-sentinel.desktop << 'EOF'
[Desktop Entry]
Type=Application
Name=Email Sentinel
Exec=/usr/local/bin/email-sentinel start --tray
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
EOF
```

**Verify:**

```bash
ls -la ~/.config/autostart/email-sentinel.desktop
```

### Test Auto-Start

**Reboot test:**
1. Restart your computer
2. After login, check system tray for Email Sentinel icon
3. Or run: `systemctl --user status email-sentinel`

**Check running process:**

```bash
ps aux | grep email-sentinel | grep -v grep
```

### Uninstall Auto-Start

```bash
email-sentinel uninstall
```

---

## Troubleshooting

### Notifications Not Appearing

**Check notification daemon:**

```bash
# GNOME
ps aux | grep gvfs-udisks2-volume-monitor

# KDE
ps aux | grep knotify

# XFCE
ps aux | grep xfce4-notifyd

# i3/Sway/Others
ps aux | grep dunst
```

**Install notification daemon (if missing):**

```bash
# Ubuntu/Debian
sudo apt-get install libnotify-bin dunst

# Fedora/RHEL
sudo dnf install libnotify dunst

# Arch Linux
sudo pacman -S libnotify dunst
```

**Test with notify-send:**

```bash
notify-send "Test" "Can you see this?"
```

**Restart notification daemon:**

```bash
# Kill existing daemon
pkill dunst

# Start new instance
dunst &
```

### System Tray Icon Not Appearing

**GNOME - Install AppIndicator extension:**

```bash
sudo apt-get install gnome-shell-extension-appindicator

# Enable extension
gnome-extensions enable appindicatorsupport@rgcjonas.gmail.com

# Restart GNOME Shell (Alt+F2, type 'r', press Enter)
```

**Check if process is running:**

```bash
ps aux | grep email-sentinel | grep -v grep
```

**DBus issues:**

```bash
# Check DBus service
systemctl --user status dbus

# Restart DBus (if needed)
systemctl --user restart dbus
```

### "credentials.json not found"

**Check if file exists:**

```bash
ls -la ~/.config/email-sentinel/credentials.json
```

**If not found, copy credentials:**

```bash
cp ~/Downloads/credentials.json ~/.config/email-sentinel/
```

### Token Expired

**Re-authenticate:**

```bash
email-sentinel init
```

### Build Fails (if building from source)

**Ensure CGO is enabled:**

```bash
export CGO_ENABLED=1
go build -o email-sentinel .
```

**Check Go version:**

```bash
go version
# Must be 1.22 or higher
```

**Install missing dependencies:**

```bash
# Ubuntu/Debian - only gcc needed!
sudo apt-get install gcc

# Fedora/RHEL
sudo dnf install gcc

# Arch Linux
sudo pacman -S gcc
```

**Note:** Unlike older versions, Email Sentinel no longer requires GTK3 development headers (`libgtk-3-dev`) or AppIndicator libraries (`libayatana-appindicator3-dev`). Only `gcc` is needed!

### Systemd Service Not Starting

**Check service status:**

```bash
systemctl --user status email-sentinel
```

**View detailed logs:**

```bash
journalctl --user -u email-sentinel -n 50
```

**Reload systemd configuration:**

```bash
systemctl --user daemon-reload
systemctl --user enable email-sentinel
systemctl --user restart email-sentinel
```

**Check executable path:**

```bash
which email-sentinel
# Should show: /usr/local/bin/email-sentinel or similar
```

### Permission Denied Errors

**Make binary executable:**

```bash
sudo chmod +x /usr/local/bin/email-sentinel
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

# Systemd commands
systemctl --user status email-sentinel
systemctl --user start email-sentinel
systemctl --user stop email-sentinel
systemctl --user restart email-sentinel
journalctl --user -u email-sentinel -f
```

### Config Locations

| File | Location |
|------|----------|
| Config | `~/.config/email-sentinel/config.yaml` |
| Credentials | `~/.config/email-sentinel/credentials.json` |
| Token | `~/.config/email-sentinel/token.json` |
| Priority Rules | `~/.config/email-sentinel/rules.yaml` |
| OTP Rules | `~/.config/email-sentinel/otp_rules.yaml` |
| AI Config | `~/.config/email-sentinel/ai-config.yaml` |
| Database | `~/.config/email-sentinel/history.db` |
| Systemd Service | `~/.config/systemd/user/email-sentinel.service` |

**Open config directory:**

```bash
cd ~/.config/email-sentinel/
ls -la
```

### Environment Variables

Add to `~/.bashrc` or `~/.zshrc`:

```bash
# AI Summaries (optional)
export GEMINI_API_KEY="your-api-key-here"
export ANTHROPIC_API_KEY="your-api-key-here"
export OPENAI_API_KEY="your-api-key-here"

# Reload shell
source ~/.bashrc  # or source ~/.zshrc
```

---

## Desktop Environment Specific Notes

### GNOME (Ubuntu, Fedora Workstation)

- **Tray Icon:** Requires `gnome-shell-extension-appindicator`
- **Notifications:** Built-in, usually works out of the box
- **Do Not Disturb:** Settings → Notifications

### KDE Plasma

- **Tray Icon:** Works natively
- **Notifications:** Built-in, works out of the box
- **Do Not Disturb:** System Settings → Notifications

### XFCE

- **Tray Icon:** Works natively
- **Notifications:** Install `xfce4-notifyd`
- **Settings:** Settings → Notifications

### i3 / Sway

- **Tray Icon:** Configure system tray in i3/sway config
- **Notifications:** Install and start `dunst`
- **Auto-start:** Add to i3/sway config

### Cinnamon (Linux Mint)

- **Tray Icon:** Works natively
- **Notifications:** Built-in, works out of the box

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

**You're all set!** Email Sentinel is now monitoring your Gmail on Linux. 📧✨
