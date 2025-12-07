# 📨 Setting Up Your Central Email Account (Collector Inbox)

Email Sentinel works best when all important email flows into one central Gmail account.

You can use:
- A brand-new Gmail account created just for alerts
- Your existing Gmail account
- A work Gmail (if allowed)

**Forwarded messages retain their original sender, subject, and metadata** — allowing Email Sentinel to classify and notify you instantly.

---

## 🧭 Why Use a Central Gmail Inbox?

Instead of connecting multiple accounts, you route specific messages from your other inboxes into a single "collector inbox."

### Benefits:

- ✅ Only one account to authenticate with Email Sentinel
- ✅ Unified monitoring of all important messages
- ✅ No multi-OAuth complexity
- ✅ Faster polling and fewer API rate issues
- ✅ Simpler configuration for end-users

---

## 🏁 Step 1 — Choose or Create Your Central Gmail Account

You can either create a dedicated Gmail account or use an existing one.

### Option A — Create a new Gmail account (recommended)

1. Visit: https://accounts.google.com/signup
2. Create a new account (e.g., `my.alerts.center@gmail.com`)
3. Log in at least once

**Why this is recommended:**
- Keeps your primary inbox separate and clean

### Option B — Use an existing Gmail account

This works fine if you prefer to consolidate everything into your existing inbox.

**Note:** Email Sentinel only reads messages that match your filters.

---

## 🧭 Step 2 — Connect Your Central Gmail to Email Sentinel

Run:

```bash
email-sentinel init
```

This opens a Google OAuth prompt in your browser.

Once authenticated, Sentinel is now authorized to read the collector inbox.

---

## 📨 Step 3 — Forward Emails From Other Accounts

Below are instructions for forwarding from:
- Gmail
- Outlook / Hotmail / Office 365
- ProtonMail
- Yahoo Mail
- iCloud
- Custom domains (Google Workspace, cPanel, Namecheap, Cloudflare, etc.)

---

### 1️⃣ Forwarding From Another Gmail Account

#### Step 1 — Open Forwarding Settings

1. Go to the Gmail account you want to forward from
2. Click ⚙️ **Settings** → **See all settings**
3. Open the **Forwarding and POP/IMAP** tab
4. Click **Add a forwarding address**

#### Step 2 — Add Your Central Gmail

Enter the collector Gmail (e.g., `my.alerts.center@gmail.com`).

#### Step 3 — Confirm the Verification Code

- Google sends a verification code to the central inbox
- Copy it back into the forwarding setup page

#### Step 4 — Set Up Forwarding Rules

To avoid forwarding every single email:

1. Go to **Settings** → **Filters and Blocked Addresses**
2. Click **Create a new filter**
3. Choose criteria (sender, subject, labels, etc.)
4. Choose **Forward it to:** `<your central Gmail>`

#### Recommended Filters

Forward:
- Job alerts
- Important senders
- Billing/invoices
- Client messages
- Password reset or security emails

---

### 2️⃣ Forwarding From Outlook / Hotmail / Office 365

#### Step 1 — Open Outlook Settings

1. Go to: https://outlook.live.com
2. Click ⚙️ **Settings** → **View all Outlook settings**

#### Step 2 — Create a Rule

1. Navigate to **Mail** → **Rules**
2. Click **Add new rule**
3. Name it (e.g., "Forward to Sentinel")

#### Step 3 — Set Conditions

Choose:
- "From" sensitivity (specific senders)
- "Subject includes"
- "Has attachment"
- etc.

#### Step 4 — Add Action

Select:
- ➡️ **Forward to:** `<your central Gmail>`

#### Step 5 — Save

Your Outlook account will now forward matching messages automatically.

---

### 3️⃣ Forwarding From ProtonMail

ProtonMail does not support standard forwarding for encrypted messages, but you can do automatic forwarding with **ProtonMail Bridge** or Proton's new **"Auto-Forwarding" Beta**.

#### Option A — Auto-Forwarding (ProtonMail Professional / Visionary / Proton Unlimited)

1. Go to **Settings** → **Messages and drafts** → **Auto-forward**
2. Add your central Gmail address
3. Choose conditions or forward-all
4. Save and confirm

#### Option B — ProtonMail Bridge (Desktop App)

If you run Bridge:

1. Configure your ProtonMail account in Bridge
2. Add it as an account in a local email client (Apple Mail, Outlook, Thunderbird)
3. Set up forwarding rules inside the mail client to auto-forward matching messages

**This works reliably but requires the Bridge app running in background.**

---

### 4️⃣ Forwarding From Yahoo Mail

#### Step 1 — Enable Forwarding

1. Go to: https://mail.yahoo.com
2. Click **Settings** → **More settings**
3. Navigate to **Mailboxes**
4. Under your mailbox, find **Forwarding**

#### Step 2 — Add Central Gmail

1. Enter `<your central Gmail>`
2. Verify via confirmation email

#### Step 3 — Add Filters (Optional)

Yahoo allows basic filtering via:
- Filters
- Blocked addresses
- Subject rules

Forward only what you need.

---

### 5️⃣ Forwarding From iCloud Mail

iCloud forwarding is straightforward.

#### Step 1 — Open iCloud Mail Settings

1. Visit: https://icloud.com/mail
2. Click the ⚙️ **Gear icon**
3. Select **Preferences** → **General**

#### Step 2 — Set Forwarding

Check:
- **"Forward my email to:"**
- Enter `<your central Gmail>`

#### Step 3 — Optional: Exclude Junk

Enable:
- ✅ "Hide my email"
- ✅ "Forward only messages that pass filtering"

---

### 6️⃣ Forwarding From Custom Domains (Google Workspace, Namecheap, Cloudflare, etc.)

#### Google Workspace Admin

1. Admin console → **Apps** → **Google Workspace** → **Gmail**
2. **Routing** → **Configure forwarding rule**
3. Choose conditions
4. Forward to `<your central Gmail>`

#### Namecheap / cPanel

Navigate to:
1. **Email** → **Forwarders** → **Add Forwarder**
2. Forward `support@mydomain.com`
3. Send to `<your central Gmail>`

#### Cloudflare Email Routing

1. Cloudflare dashboard → **Email** → **Routes**
2. Add routing rule:
   - **Match:** `*@yourdomain.com` or specific addresses
   - **Forward to:** `<your central Gmail>`

**Cloudflare is extremely reliable and recommended.**

---

## 🔧 Recommended Forwarding Strategy

Most users will want:

- ✅ All job alerts → forward
- ✅ All billing/invoice emails → forward
- ✅ VIP contacts → forward
- ✅ Any email with keywords: "urgent", "asap", "invoice", "deadline", "payment"

**Let each account's own filter system decide what matters.**

Your central inbox becomes the single point that Email Sentinel monitors.

---

## 🔍 Step 4 — Test Forwarding

Send yourself an email from another account with a subject like:

```
Forwarding Test — From <YourService>
```

### Confirm:

- ✅ It arrives in the central Gmail inbox
- ✅ Email Sentinel processes it
- ✅ Desktop/mobile notifications fire

---

## 🛡️ Privacy & Security Notes

- ✅ Email Sentinel only reads the central Gmail you authorize
- ✅ Forwarding rules give you fine-grained outbound privacy control
- ✅ You can stop forwarding anytime
- ✅ Your original inboxes remain untouched
- ✅ No passwords are shared or stored for other accounts

---

## 🎉 You're All Set!

Your central Gmail inbox is now configured as a unified collection point for all important emails across your accounts.

**Next Steps:**
1. Add filters in Email Sentinel: `email-sentinel filter add`
2. Start monitoring: `email-sentinel start --tray`
3. View alert history: `email-sentinel alerts`

For more help, see:
- [Complete Setup Guide](build_to_first_alert.md)
- [Gmail API Setup](gmail_api_setup.md)
- [Mobile Notifications](mobile_ntfy_setup.md)
