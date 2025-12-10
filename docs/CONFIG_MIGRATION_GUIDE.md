# Configuration Migration Guide

## Overview

Email Sentinel has migrated to a **unified configuration file** (`app-config.yaml`) that consolidates all settings into one easy-to-manage file.

**Previous (v0.x):**
- `ai-config.yaml` - AI summary settings
- `rules.yaml` - Priority rules and notifications
- `otp_rules.yaml` - OTP/2FA detection settings

**Current (v1.0+):**
- `app-config.yaml` - **All settings in one file**

## Automatic Migration

Email Sentinel **automatically migrates** your existing configuration when you:

1. Run any command (e.g., `email-sentinel start`)
2. Or explicitly run: `email-sentinel config migrate`

### What Happens During Migration

1. ✅ Detects old config files (`ai-config.yaml`, `rules.yaml`, `otp_rules.yaml`)
2. ✅ Merges all settings into new `app-config.yaml`
3. ✅ Preserves all your custom settings
4. ✅ Keeps old files as backup (does not delete them)
5. ✅ Shows migration summary

### Migration Example

```bash
$ email-sentinel config migrate

🔄 Starting configuration migration...
📦 Migrating from separate config files to unified app-config.yaml...
📦 Migrated from: [ai-config.yaml rules.yaml otp_rules.yaml]
✅ Successfully migrated to app-config.yaml

✅ Configuration loaded successfully!

📋 Configuration Summary:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 Monitoring:
   Polling Interval: 45 seconds
   WAL Mode: true
   Cleanup Interval: 1h

🤖 AI Summary:
   Enabled: true
   Provider: gemini
   Cache Enabled: true

⚡ Priority Rules:
   Urgent Keywords: 35 configured
   VIP Senders: 3 configured
   VIP Domains: 3 configured

🔐 OTP Detection:
   Enabled: true
   Expiry Duration: 5m
   Trusted Senders: 14 configured
   Auto-copy to Clipboard: false

🔔 Notifications:
   Desktop: true
   Mobile: false
   Weekend Mode: normal
   Quiet Hours: 22:00 - 08:00

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📁 Config file: C:\Users\YourName\AppData\Roaming\email-sentinel\app-config.yaml

💡 Tip: You can now edit app-config.yaml to customize your settings
   The old config files are kept as backup and not deleted
```

## Field Mapping

### AI Config Migration

| Old (ai-config.yaml) | New (app-config.yaml) |
|---------------------|----------------------|
| `ai_summary.enabled` | `ai_summary.enabled` |
| `ai_summary.provider` | `ai_summary.provider` |
| `ai_summary.api.claude.*` | `ai_summary.providers.claude.*` |
| `ai_summary.api.openai.*` | `ai_summary.providers.openai.*` |
| `ai_summary.api.gemini.*` | `ai_summary.providers.gemini.*` |
| `ai_summary.behavior.enable_cache` | `ai_summary.cache.enabled` |
| `ai_summary.rate_limit.max_per_hour` | `ai_summary.providers.*.rate_limit.requests_per_minute` (converted) |
| `ai_summary.rate_limit.max_per_day` | `ai_summary.providers.*.rate_limit.requests_per_day` |

### Rules Migration

| Old (rules.yaml) | New (app-config.yaml) |
|------------------|----------------------|
| `priority_rules.urgent_keywords` | `priority.urgent_keywords` |
| `priority_rules.vip_senders` | `priority.vip_senders` |
| `priority_rules.vip_domains` | `priority.vip_domains` |
| `notification_settings.quiet_hours_start` | `notifications.quiet_hours.start` |
| `notification_settings.quiet_hours_end` | `notifications.quiet_hours.end` |
| `notification_settings.weekend_mode` | `notifications.weekend_mode` |

### OTP Rules Migration

| Old (otp_rules.yaml) | New (app-config.yaml) |
|---------------------|----------------------|
| `enabled` | `otp.enabled` |
| `expiry_duration` | `otp.expiry_duration` |
| `auto_copy_to_clipboard` | `otp.clipboard.auto_copy` |
| `clipboard_auto_clear` | `otp.clipboard.clear_after` |
| `custom_patterns` | `otp.custom_patterns` |
| `trusted_otp_senders` | `otp.trusted_senders` |

## New Features in app-config.yaml

The unified config adds several new features not available in the old separate files:

### 1. **Monitoring Settings**
```yaml
monitoring:
  polling_interval: 45
  database:
    wal_mode: true
    cleanup_interval: "1h"
```

### 2. **OTP Trusted Domains**
```yaml
otp:
  trusted_domains:
    - amazon.com
    - paypal.com
    - stripe.com
```

### 3. **Enhanced Notification Controls**
```yaml
notifications:
  desktop:
    enabled: true
    duration: 10
    sound: true
  mobile:
    enabled: false
    topic: ""
    server: "https://ntfy.sh"
    priority: 4
  quiet_hours:
    start: "22:00"
    end: "08:00"
    allow_urgent: true
  weekend_mode: normal
```

### 4. **AI Cache Settings**
```yaml
ai_summary:
  cache:
    enabled: true
    ttl: "24h"
    max_size: 1000
```

### 5. **OTP Trigger Phrases**
```yaml
otp:
  trigger_phrases:
    - verification code
    - security code
    - one-time password
```

## Manual Migration (If Needed)

If automatic migration doesn't work or you want to manually migrate:

1. **Backup your old configs:**
   ```bash
   cp ai-config.yaml ai-config.yaml.backup
   cp rules.yaml rules.yaml.backup
   cp otp_rules.yaml otp_rules.yaml.backup
   ```

2. **Copy the template:**
   ```bash
   cp app-config.yaml.example app-config.yaml
   ```

3. **Edit and merge settings** from your old files into `app-config.yaml`

4. **Test the new config:**
   ```bash
   email-sentinel config migrate
   ```

## Troubleshooting

### Migration Failed

If you see `⚠️ Migration failed`:

1. **Check file permissions** - Ensure old config files are readable
2. **Validate YAML syntax** - Use a YAML validator on old config files
3. **Check logs** - Look for specific error messages

### Old Config Files Not Detected

If migration says "no legacy config files found":

1. **Check config directory:**
   - Windows: `%APPDATA%\email-sentinel\`
   - macOS: `~/Library/Application Support/email-sentinel/`
   - Linux: `~/.config/email-sentinel/`

2. **Verify file names** are exactly:
   - `ai-config.yaml`
   - `rules.yaml`
   - `otp_rules.yaml`

### Settings Not Migrated

If some settings didn't migrate:

1. **Run migration again:**
   ```bash
   email-sentinel config migrate
   ```

2. **Check the output** for warnings about specific files

3. **Manually add missing settings** to `app-config.yaml`

## After Migration

### What to Do Next

1. ✅ **Review migrated config:**
   ```bash
   email-sentinel config migrate
   ```

2. ✅ **Edit if needed:**
   - Windows: `%APPDATA%\email-sentinel\app-config.yaml`
   - macOS: `~/Library/Application Support/email-sentinel/app-config.yaml`
   - Linux: `~/.config/email-sentinel/app-config.yaml`

3. ✅ **Test your setup:**
   ```bash
   email-sentinel start
   ```

### Should I Delete Old Config Files?

**No, keep them as backup!**

Email Sentinel automatically uses `app-config.yaml` if it exists, ignoring the old files. You can delete them later once you're confident the migration worked correctly.

### Using the Old Files

If you want to go back to old config files:
1. Delete or rename `app-config.yaml`
2. Email Sentinel will attempt migration again on next run

## FAQ

**Q: Will this break my existing setup?**
A: No, migration is automatic and preserves all settings. Old files remain as backup.

**Q: Can I still use the old config files?**
A: No, once `app-config.yaml` exists, Email Sentinel uses only that file.

**Q: What if I have custom settings not in the migration?**
A: Manually add them to `app-config.yaml` using the documented structure.

**Q: Does this affect my filters?**
A: No, filters remain in `filters.yaml` and are unaffected by this migration.

**Q: Are API keys migrated?**
A: No, API keys are loaded from environment variables (`GEMINI_API_KEY`, `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`), not stored in config files.

## Getting Help

If you encounter issues during migration:

1. Check this guide thoroughly
2. Run `email-sentinel config migrate` to see detailed output
3. Report issues: https://github.com/datateamsix/email-sentinel/issues

---

**Updated:** 2025-12-10
**Version:** 1.0.0
