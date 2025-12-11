# 📚 Email Sentinel Documentation

Complete documentation for Email Sentinel - a lightweight, smart real-time notification management system that protects your focus and productivity.

---

## 🚀 Getting Started

**New to Email Sentinel?** Start here:

### Quick Installation Guides

Choose your platform for a complete walkthrough:

| Platform | Guide | Time | Features |
|----------|-------|------|----------|
| 🪟 **Windows** | [Windows Quickstart](QUICKSTART_WINDOWS.md) | 10-15 min | Task Scheduler, Toast Notifications, System Tray |
| 🍎 **macOS** | [macOS Quickstart](QUICKSTART_MACOS.md) | 10-15 min | LaunchAgent, Menu Bar, Notification Center |
| 🐧 **Linux** | [Linux Quickstart](QUICKSTART_LINUX.md) | 10-15 min | systemd, Desktop Notifications |

### Universal Setup

- **[Build to First Alert](BUILD_TO_FIRST_ALERT.md)** - Step-by-step guide for any platform (15-20 min)
- **[Gmail API Setup](GMAIL_API_SETUP.md)** - Required Google Cloud configuration (5 min)

---

## 📖 Feature Guides

### Core Features

- **[CLI Command Reference](CLI_GUIDE.md)** - Complete command documentation
- **Filter Management** - See CLI Guide: [Filter Commands](CLI_GUIDE.md#filters)
- **Notification System** - Desktop, Mobile, Toast notifications
- **Alert History** - View and manage past alerts

### Advanced Features

- **[Gmail Category Scopes](../README.md#gmail-category-scopes)** - Target specific Gmail categories (primary, social, promotions)
- **[AI Email Summaries](AI_ARCHITECTURE.md)** - Automatic email summarization with Claude/GPT/Gemini
- **[OTP/2FA Detection](../README.md#otp2fa-code-detection)** - Extract verification codes from emails
- **[Priority Rules](../README.md#priority-rules)** - Auto-classify urgent emails
- **[Filter Labels](../README.md#filter-labels)** - Organize filters by category
- **[System Tray](TRAY_SYSTEM_ARCHITECTURE.md)** - Background app architecture

### Multi-Account Setup

- **[Central Email Setup](CENTRAL_EMAIL_SETUP.md)** - Monitor multiple accounts via email forwarding
  - Gmail, Outlook, Yahoo, ProtonMail, etc.
  - Single OAuth authentication
  - Retain original sender information

---

## 🔧 Configuration

### Configuration Files

All configuration stored in platform-specific directories:

| Platform | Location |
|----------|----------|
| **Windows** | `%APPDATA%\email-sentinel\` |
| **macOS** | `~/Library/Application Support/email-sentinel/` |
| **Linux** | `~/.config/email-sentinel/` |

### Key Configuration Files

- **`app-config.yaml`** - Main configuration (polling, notifications, AI, priority rules)
- **`config.yaml`** - Filter definitions
- **`otp_rules.yaml`** - OTP/2FA detection settings
- **`credentials.json`** - Gmail API OAuth credentials (from Google Cloud)
- **`token.json`** - OAuth access token (auto-generated)
- **`history.db`** - SQLite database for alert history

---

## 💡 Common Use Cases

### Job Search Monitoring
Monitor LinkedIn, Greenhouse, Lever for interview opportunities.
- [Example Filter](../README.md#job-search-monitoring)

### VIP Sender Alerts
Never miss emails from your boss, clients, or important contacts.
- [Example Filter](../README.md#vip-sender-alerts)

### Urgent Keyword Detection
Get alerted for emails containing "urgent", "asap", "deadline", etc.
- [Example Filter](../README.md#urgent-keyword-monitoring)

### Social Media Digest
Monitor social notifications from Facebook, Twitter, LinkedIn.
- [Example Filter](../README.md#social-media-digest)

### Newsletter Tracking
Track newsletters without cluttering your primary inbox.
- [Example Filter](../README.md#newsletter-monitoring)

---

## 🏗️ Architecture & Design

### Technical Documentation

- **[System Architecture](../README.md#architecture)** - High-level design
- **[System Tray Architecture](TRAY_SYSTEM_ARCHITECTURE.md)** - Background app design
- **[AI Architecture](AI_ARCHITECTURE.md)** - Email summarization pipeline
- **[Configuration Migration](CONFIG_MIGRATION_GUIDE.md)** - Upgrading between versions

### Data Flow

1. Gmail API polling (configurable interval)
2. Filter matching with Gmail scope support
3. Priority evaluation
4. OTP detection
5. AI summary generation (optional)
6. Multi-channel notifications
7. Alert storage in SQLite

---

## 🛠️ Troubleshooting

### Common Issues

**OAuth & Authentication:**
- "App not verified" warning - [Solution](../README.md#app-not-verified-during-oauth)
- "Access blocked" error - [Solution](../README.md#access-blocked-error)
- Token expired - [Solution](../README.md#token-expired)

**Notifications:**
- Not appearing - [Windows](QUICKSTART_WINDOWS.md#notifications-not-appearing) | [macOS](QUICKSTART_MACOS.md#troubleshooting) | [Linux](QUICKSTART_LINUX.md#troubleshooting)
- System tray icon missing - [Solution](../README.md#system-tray-icon-missing)

**Build & Installation:**
- Build fails - [Solution](../README.md#build-fails-from-source)
- CGO issues - [Windows Guide](QUICKSTART_WINDOWS.md#build-fails-if-building-from-source)

---

## 📱 Platform-Specific Features

### Windows
- **Windows Toast Notifications** - Rich notifications in Action Center
- **Task Scheduler Auto-Start** - Launch on Windows boot
- **System Tray Integration** - Background monitoring

[Windows Quickstart →](QUICKSTART_WINDOWS.md)

### macOS
- **Menu Bar App** - Native macOS menu bar integration
- **Notification Center** - macOS native notifications
- **LaunchAgent Auto-Start** - Launch on login

[macOS Quickstart →](QUICKSTART_MACOS.md)

### Linux
- **systemd Service** - Run as user service
- **Desktop Notifications** - GNOME/KDE/Xfce support
- **System Tray** - Desktop environment integration

[Linux Quickstart →](QUICKSTART_LINUX.md)

---

## 🤝 Contributing

Want to contribute? Check out:

- **[Contributing Guidelines](../README.md#contributing)** - How to contribute
- **[Architecture Docs](TRAY_SYSTEM_ARCHITECTURE.md)** - Understand the codebase
- **[GitHub Issues](https://github.com/datateamsix/email-sentinel/issues)** - Report bugs or request features

---

## 📞 Support

- **Main README:** [../README.md](../README.md)
- **Issues:** https://github.com/datateamsix/email-sentinel/issues
- **Discussions:** https://github.com/datateamsix/email-sentinel/discussions

---

## 📋 Quick Reference

### Common Commands

```bash
# Initialize
email-sentinel init

# Add filter with Gmail scope
email-sentinel filter add --name "Name" --from "sender.com" --scope inbox --labels "work"

# Start monitoring
email-sentinel start --tray

# View alerts
email-sentinel alerts

# Get OTP code
email-sentinel otp get

# Install auto-start
email-sentinel install
```

### Links by Topic

**Setup:**
- [Gmail API Setup](GMAIL_API_SETUP.md)
- [Windows Setup](QUICKSTART_WINDOWS.md)
- [macOS Setup](QUICKSTART_MACOS.md)
- [Linux Setup](QUICKSTART_LINUX.md)

**Features:**
- [CLI Commands](CLI_GUIDE.md)
- [Gmail Scopes](../README.md#gmail-category-scopes)
- [AI Summaries](AI_ARCHITECTURE.md)
- [OTP Detection](../README.md#otp2fa-code-detection)
- [Multi-Account](CENTRAL_EMAIL_SETUP.md)

**Advanced:**
- [System Tray](TRAY_SYSTEM_ARCHITECTURE.md)
- [Architecture](../README.md#architecture)
- [Configuration](CONFIG_MIGRATION_GUIDE.md)

---

**Last Updated:** December 2025
**Version:** 1.0
