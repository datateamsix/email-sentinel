# 📱 Mobile Push Notifications Setup (ntfy.sh)

![Email Sentinel Logo](../images/logo.png)

Complete guide to setting up free mobile push notifications for Email Sentinel using [ntfy.sh](https://ntfy.sh).

---

## Table of Contents
1. [What is ntfy.sh?](#what-is-ntfysh)
2. [Download the App](#download-the-app)
3. [Create Your Topic](#create-your-topic)
4. [Configure Email Sentinel](#configure-email-sentinel)
5. [Test Your Setup](#test-your-setup)
6. [Troubleshooting](#troubleshooting)

---

## What is ntfy.sh?

**ntfy.sh** is a free, open-source service that sends push notifications to your phone without requiring an account.

**Key Features:**
- ✅ **100% Free** - No account needed, no credit card
- ✅ **Privacy-Focused** - Messages are deleted after 12 hours
- ✅ **Cross-Platform** - Works on iPhone and Android
- ✅ **Simple** - Just pick a unique topic name
- ✅ **Reliable** - Used by thousands of users worldwide

**How it works:**
1. You subscribe to a unique "topic" in the ntfy app
2. Email Sentinel sends notifications to that topic
3. Your phone receives instant push notifications

---

## Download the App

### iPhone (iOS)

1. Open the **App Store**
2. Search for **"ntfy"**
3. Download the app by **Philipp Heckel**
   - Direct link: https://apps.apple.com/app/ntfy/id1625396347
4. Open the app after installation

### Android

1. Open the **Google Play Store**
2. Search for **"ntfy"**
3. Download the app by **Philipp Heckel**
   - Direct link: https://play.google.com/store/apps/details?id=io.heckel.ntfy
4. Open the app after installation

---

## Create Your Topic

A "topic" is like a channel that only you know about. Anyone can publish to it, so make it unique!

### Step 1: Choose a Unique Topic Name

**Good topic names:**
- `michaels-email-alerts-x7k2p9`
- `johndoe-gmail-monitor-2025`
- `mywork-alerts-secret123`

**Bad topic names:**
- `email` (too common, others might use it)
- `alerts` (not unique)
- `test` (way too common)

**Tips:**
- Include your name or initials
- Add random numbers/letters at the end
- Make it hard to guess
- No spaces, use hyphens `-` or underscores `_`

### Step 2: Subscribe in the ntfy App

#### On iPhone:
1. Open the **ntfy** app
2. Tap the **"+"** button (bottom right)
3. Enter your unique topic name
   - Example: `your-name-email-alerts-abc123`
4. Tap **"Subscribe"**
5. You'll see your topic in the list

#### On Android:
1. Open the **ntfy** app
2. Tap the **"+"** button (bottom right)
3. Under **"Subscribe to topic"**, enter your unique topic name
   - Example: `your-name-email-alerts-abc123`
4. Tap **"Subscribe"**
5. You'll see your topic in the list

### Step 3: Test Your Topic (Optional)

Before connecting Email Sentinel, test that notifications work:

1. Open a web browser on your computer
2. Go to: `https://ntfy.sh`
3. In the **"Publish"** box:
   - **Topic:** Enter your topic name
   - **Message:** Type "Test from web"
4. Click **"Send"**
5. Check your phone - you should get a notification!

---

## Configure Email Sentinel

Now connect Email Sentinel to your ntfy topic.

### Step 1: Enable Mobile Notifications

**Windows (PowerShell/CMD):**
```powershell
.\email-sentinel.exe config set mobile true
```

**macOS/Linux (Bash):**
```bash
./email-sentinel config set mobile true
```

### Step 2: Set Your Topic

**Windows (PowerShell/CMD):**
```powershell
.\email-sentinel.exe config set ntfy_topic "your-topic-name-here"
```

**macOS/Linux (Bash):**
```bash
./email-sentinel config set ntfy_topic "your-topic-name-here"
```

**Example:**
```powershell
# Windows
.\email-sentinel.exe config set ntfy_topic "michaels-email-alerts-x7k2p9"

# macOS/Linux
./email-sentinel config set ntfy_topic "michaels-email-alerts-x7k2p9"
```

### Step 3: Verify Configuration

**Windows (PowerShell/CMD):**
```powershell
.\email-sentinel.exe config show
```

**macOS/Linux (Bash):**
```bash
./email-sentinel config show
```

You should see:
```yaml
notifications:
  desktop: true
  mobile:
    enabled: true
    ntfy_topic: your-topic-name-here
```

---

## Test Your Setup

### Test from Email Sentinel

**Windows (PowerShell/CMD):**
```powershell
.\email-sentinel.exe test mobile
```

**macOS/Linux (Bash):**
```bash
./email-sentinel test mobile
```

**Expected output:**
```
📱 Sending test mobile notification...

Sending to topic: your-topic-name-here

✅ Test notification sent!

Check your phone for a notification from ntfy.sh
```

**On your phone:**
- Notification should appear within 1-2 seconds
- Title: "Email Sentinel Test"
- Message: "If you can see this on your phone, mobile notifications are working! ✅"

### If It Works:
🎉 **Success!** You're all set. Mobile notifications will now be sent when emails match your filters.

### If It Doesn't Work:
See [Troubleshooting](#troubleshooting) below.

---

## Using Mobile Notifications

Once configured, Email Sentinel will automatically send mobile notifications when emails match your filters.

### What You'll Receive

When an email matches a filter, you'll get a push notification with:

**Normal Priority Email:**
- **Title:** "📧 Email Match: [Filter Name]"
- **Message:** "From: sender@example.com | Subject: Email subject"

**High Priority Email (VIP/Urgent):**
- **Title:** "🔥 URGENT: [Filter Name]"
- **Message:** "From: boss@company.com | Subject: Important message"

### Notification Features

- **Instant delivery** - Usually arrives in 1-2 seconds
- **Works anywhere** - As long as your phone has internet
- **No battery drain** - ntfy uses efficient push notifications
- **Silent hours** - Configure quiet hours in `rules.yaml` (coming soon)

---

## Troubleshooting

### Problem: No notification received

**Check 1: Phone has internet connection**
- Try opening a website on your phone
- ntfy needs internet (WiFi or cellular data)

**Check 2: ntfy app is running**
- Open the ntfy app
- Verify your topic is still subscribed
- Check if any notifications are shown in the app

**Check 3: Notification permissions**

**iPhone:**
1. Settings → Notifications → ntfy
2. Enable "Allow Notifications"
3. Set alert style to "Alerts" or "Banners"

**Android:**
1. Settings → Apps → ntfy → Notifications
2. Enable "All ntfy notifications"
3. Check "Importance" is set to "High"

**Check 4: Topic name is correct**

**Windows (PowerShell/CMD):**
```powershell
.\email-sentinel.exe config show
```

**macOS/Linux (Bash):**
```bash
./email-sentinel config show
```

Verify the `ntfy_topic` matches exactly what you subscribed to in the app.

**Check 5: Test manually**
1. Open browser: https://ntfy.sh
2. Send a test message to your topic
3. If this works but Email Sentinel doesn't, the issue is with Email Sentinel config

### Problem: Delayed notifications

**Cause:** Phone's battery saver or data saver mode

**Solution:**
1. Add ntfy to battery optimization exemptions
   - **Android:** Settings → Battery → Battery optimization → ntfy → Don't optimize
   - **iPhone:** No action needed (iOS handles this automatically)

### Problem: Topic name was guessed by someone else

**Symptoms:**
- Receiving random notifications
- Notifications that aren't from Email Sentinel

**Solution:**
1. Create a NEW topic with more random characters
   - Example: `myalerts-x9k2p7q5m3`
2. Update Email Sentinel:

   **Windows:**
   ```powershell
   .\email-sentinel.exe config set ntfy_topic "myalerts-x9k2p7q5m3"
   ```

   **macOS/Linux:**
   ```bash
   ./email-sentinel config set ntfy_topic "myalerts-x9k2p7q5m3"
   ```
3. Unsubscribe from old topic in ntfy app

### Problem: Too many notifications

**Option 1: Adjust filter rules**
- Make filters more specific
- Use `--match all` instead of `--match any`

**Option 2: Disable mobile notifications temporarily**

**Windows:**
```powershell
.\email-sentinel.exe config set mobile false
```

**macOS/Linux:**
```bash
./email-sentinel config set mobile false
```

**Option 3: Use priority rules**
- Configure `rules.yaml` to only notify for urgent emails
- Future feature: weekend_mode and quiet_hours

---

## Advanced Features

### Self-Hosted ntfy Server (Optional)

If you want complete privacy, you can run your own ntfy server:

1. Follow instructions at: https://docs.ntfy.sh/install/
2. Change the ntfy URL in Email Sentinel config:

   **Windows:**
   ```powershell
   # This is a future feature - not yet implemented
   .\email-sentinel.exe config set ntfy_url "https://your-server.com"
   ```

   **macOS/Linux:**
   ```bash
   # This is a future feature - not yet implemented
   ./email-sentinel config set ntfy_url "https://your-server.com"
   ```

### Multiple Devices

You can receive notifications on multiple phones:

1. Install ntfy app on second device
2. Subscribe to the SAME topic name
3. Both devices will receive notifications

### Mute Specific Topics

In the ntfy app:
1. Long-press on your topic
2. Tap "Mute notifications"
3. Choose duration (1 hour, 1 day, forever)

---

## Privacy & Security

### What data is sent?

Email Sentinel sends to ntfy:
- Filter name that matched
- Email sender (From address)
- Email subject line

**NOT sent:**
- Email body content
- Attachments
- Your Gmail credentials

### How long is data stored?

- **ntfy.sh default:** Messages are deleted after **12 hours**
- **In transit:** Messages are sent over HTTPS (encrypted)
- **On your phone:** Stored until you clear notifications

### Can others see my notifications?

- **Only if they guess your topic name**
- This is why choosing a unique, random topic is important
- Consider topics like: `user-randomstring-20chars`

### Security Best Practices

1. ✅ Use a long, random topic name
2. ✅ Don't share your topic name
3. ✅ Avoid using personal information in topic names
4. ✅ If you suspect your topic is compromised, create a new one

---

## Need Help?

**ntfy Documentation:** https://docs.ntfy.sh

**Email Sentinel Issues:** https://github.com/datateamsix/email-sentinel/issues

**Test mobile notifications:**

**Windows:**
```powershell
.\email-sentinel.exe test mobile
```

**macOS/Linux:**
```bash
./email-sentinel test mobile
```

---

## Quick Reference

**Windows (PowerShell/CMD):**
```powershell
# Enable mobile notifications
.\email-sentinel.exe config set mobile true

# Set your ntfy topic
.\email-sentinel.exe config set ntfy_topic "your-unique-topic"

# Test mobile notifications
.\email-sentinel.exe test mobile

# Check configuration
.\email-sentinel.exe config show

# Disable mobile notifications
.\email-sentinel.exe config set mobile false
```

**macOS/Linux (Bash):**
```bash
# Enable mobile notifications
./email-sentinel config set mobile true

# Set your ntfy topic
./email-sentinel config set ntfy_topic "your-unique-topic"

# Test mobile notifications
./email-sentinel test mobile

# Check configuration
./email-sentinel config show

# Disable mobile notifications
./email-sentinel config set mobile false
```

---

**You're all set!** 📱 You'll now receive instant push notifications on your phone when important emails arrive.
