# 🔐 Gmail API & Google Cloud Platform Setup

![Email Sentinel Logo](../images/logo.png)

Complete step-by-step guide to setting up Gmail API access for Email Sentinel.

**Time Required:** ~10 minutes
**Cost:** FREE (uses Google Cloud free tier)

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Step 1: Create Google Cloud Project](#step-1-create-google-cloud-project)
4. [Step 2: Enable Gmail API](#step-2-enable-gmail-api)
5. [Step 3: Configure OAuth Consent Screen](#step-3-configure-oauth-consent-screen)
6. [Step 4: Create OAuth Credentials](#step-4-create-oauth-credentials)
7. [Step 5: Download Credentials](#step-5-download-credentials)
8. [Verification](#verification)
9. [Troubleshooting](#troubleshooting)
10. [Security & Privacy](#security--privacy)

---

## Overview

Email Sentinel needs permission to read your Gmail inbox. This is done securely through Google's OAuth 2.0 system.

**What we'll set up:**
- Google Cloud Platform (GCP) project
- Gmail API access
- OAuth 2.0 credentials
- Read-only permissions (Email Sentinel NEVER modifies or deletes emails)

**Important:**
- ✅ Your credentials stay on your computer
- ✅ Email Sentinel only reads emails (readonly access)
- ✅ No emails are sent to external servers
- ✅ All processing happens locally on your machine

---

## Prerequisites

### Required
- **Gmail account** - The account you want to monitor
- **Google account** - Can be the same as your Gmail account

### Costs
- **FREE** - Gmail API is free for personal use
- **Quota:** 1 billion quota units per day (way more than you'll ever use)
- **No credit card required**

---

## Step 1: Create Google Cloud Project

### 1.1 Go to Google Cloud Console

Open your browser and navigate to:
**https://console.cloud.google.com/**

**First time?** You may see a welcome screen - click "Get started" or "Try for free"

### 1.2 Accept Terms of Service

If prompted:
- Read the Terms of Service
- Check the box to agree
- Click "Agree and Continue"

### 1.3 Create a New Project

1. Click the **project dropdown** at the top (next to "Google Cloud")
   - It might say "Select a project" or show an existing project name

2. Click **"NEW PROJECT"** (top right of the dialog)

3. Fill in project details:
   - **Project name:** `Email Sentinel`
   - **Organization:** Leave as "No organization" (unless you have one)
   - **Location:** Leave as default

4. Click **"CREATE"**

5. Wait a few seconds for the project to be created

### 1.4 Select Your New Project

1. Click the **project dropdown** again
2. Select **"Email Sentinel"** from the list
3. The top bar should now show "Email Sentinel" as the active project

---

## Step 2: Enable Gmail API

### 2.1 Navigate to API Library

**Option A: Direct Link**
- Go to: https://console.cloud.google.com/apis/library

**Option B: Using Menu**
1. Click ☰ (hamburger menu) in top-left
2. Navigate to: **APIs & Services** → **Library**

### 2.2 Search for Gmail API

1. In the search box, type: **"Gmail API"**
2. Click on **"Gmail API"** from the results
   - Look for the one with the Gmail icon (red/white envelope)

### 2.3 Enable the API

1. Click the blue **"ENABLE"** button
2. Wait 10-20 seconds for the API to be enabled
3. You'll see "API enabled" notification

**Important:** Make sure you're still on the "Email Sentinel" project!

---

## Step 3: Configure OAuth Consent Screen

This screen determines what users see when authorizing Email Sentinel.

### 3.1 Navigate to OAuth Consent Screen

1. Click ☰ (hamburger menu)
2. Go to: **APIs & Services** → **OAuth consent screen**

**Or use direct link:** https://console.cloud.google.com/apis/credentials/consent

### 3.2 Select User Type

You'll see two options:

**Select: "External"**
- Allows you to use any Gmail account
- You can add yourself as a test user
- **Don't worry about verification** - for personal use, you don't need it

Click **"CREATE"**

### 3.3 Fill in App Information

**Page 1: OAuth consent screen**

Fill in these fields:

1. **App name:** `Email Sentinel`

2. **User support email:** Your email address (select from dropdown)

3. **App logo:** Leave blank (optional)

4. **Application home page:** Leave blank (optional)

5. **Application privacy policy link:** Leave blank (optional)

6. **Application terms of service link:** Leave blank (optional)

7. **Authorized domains:** Leave blank

8. **Developer contact information:**
   - Enter your email address

Click **"SAVE AND CONTINUE"**

### 3.4 Configure Scopes

**Page 2: Scopes**

This defines what Email Sentinel can access.

1. Click **"ADD OR REMOVE SCOPES"**

2. In the filter box, search for: **"Gmail API"**

3. Find and check this scope:
   - ✅ **`.../auth/gmail.readonly`**
   - Description: "Read all resources and their metadata—no write operations."

4. Click **"UPDATE"**

5. Verify the scope appears in "Your sensitive scopes" table

6. Click **"SAVE AND CONTINUE"**

**Important:** Only select `gmail.readonly` - this ensures Email Sentinel can NEVER modify or delete your emails.

### 3.5 Add Test Users

**Page 3: Test users**

Since the app is "External" and unverified, you must add yourself as a test user.

1. Click **"+ ADD USERS"**

2. Enter your Gmail address (the one you want to monitor)
   - Example: `your.email@gmail.com`

3. Click **"ADD"**

4. Verify your email appears in the test users list

5. Click **"SAVE AND CONTINUE"**

### 3.6 Review Summary

**Page 4: Summary**

1. Review all the information
2. Everything should look correct
3. Click **"BACK TO DASHBOARD"**

**✅ OAuth Consent Screen is now configured!**

---

## Step 4: Create OAuth Credentials

Now we'll create the credentials file that Email Sentinel will use.

### 4.1 Navigate to Credentials

1. Click ☰ (hamburger menu)
2. Go to: **APIs & Services** → **Credentials**

**Or use direct link:** https://console.cloud.google.com/apis/credentials

### 4.2 Create Credentials

1. Click **"+ CREATE CREDENTIALS"** (top of page)

2. Select **"OAuth client ID"** from the dropdown

### 4.3 Configure OAuth Client

If prompted to configure the consent screen:
- Click **"CONFIGURE CONSENT SCREEN"**
- You've already done this, so click **"BACK"**

1. **Application type:** Select **"Desktop app"**
   - This is important! Don't select "Web application"

2. **Name:** `Email Sentinel Desktop`
   - Or any name you prefer

3. Click **"CREATE"**

### 4.4 Success Screen

You'll see a popup: **"OAuth client created"**

- Shows your Client ID
- Shows your Client secret
- **Don't worry about copying these now**

Click **"OK"** to close the popup

---

## Step 5: Download Credentials

### 5.1 Find Your OAuth Client

You should now see your credentials listed:

**OAuth 2.0 Client IDs** section shows:
- Name: "Email Sentinel Desktop"
- Type: "Desktop"
- Created: [date]

### 5.2 Download JSON

1. On the right side of your credential row, click the **⬇ Download** icon
   - It looks like a download arrow

2. A JSON file will download:
   - File name: `client_secret_XXXXX.apps.googleusercontent.com.json`

3. **Rename this file to:** `credentials.json`
   - This is important! Email Sentinel looks for `credentials.json`

### 5.3 Save the File

Save `credentials.json` to a safe location. You'll need it for Email Sentinel.

**Recommended locations:**
- Desktop (temporary)
- Email Sentinel project directory
- Config directory (see build guide)

**⚠️ Keep this file safe!**
- Don't share it publicly
- Don't commit it to Git
- It contains your OAuth client secret

---

## Verification

Let's verify everything is set up correctly.

### Check 1: Gmail API is Enabled

1. Go to: **APIs & Services** → **Enabled APIs & services**
2. You should see **"Gmail API"** in the list
3. Status should be "Enabled"

### Check 2: OAuth Consent Screen Configured

1. Go to: **APIs & Services** → **OAuth consent screen**
2. You should see:
   - User Type: External
   - Status: Testing (or Published)
   - Test users: Your email address

### Check 3: Credentials Downloaded

1. Go to: **APIs & Services** → **Credentials**
2. You should see your OAuth 2.0 Client ID
3. You should have `credentials.json` file downloaded

### Check 4: File Contents

Open `credentials.json` in a text editor. It should look like:

```json
{
  "installed": {
    "client_id": "XXXXX.apps.googleusercontent.com",
    "project_id": "email-sentinel-...",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    ...
  }
}
```

If it looks like this, you're good! ✅

---

## Troubleshooting

### Problem: "Gmail API not found"

**Solution:**
1. Make sure you selected the correct project
2. Check project dropdown at top - should say "Email Sentinel"
3. Try disabling and re-enabling the Gmail API

### Problem: "Access blocked: This app isn't verified"

**This is NORMAL for personal apps!**

**Solution:**
1. When you see this screen, click **"Advanced"**
2. Click **"Go to Email Sentinel (unsafe)"**
3. This is safe - you created the app yourself!

**Why this happens:**
- Google shows this warning for all unverified apps
- Verification costs $100-$200 and takes weeks
- For personal use, you don't need verification
- Adding yourself as a "test user" allows you to bypass this

### Problem: "Test user not found"

**Solution:**
1. Go to **OAuth consent screen**
2. Scroll to **Test users**
3. Click **"+ ADD USERS"**
4. Add your Gmail address
5. Click **"SAVE"**

### Problem: Wrong credential type downloaded

**Symptoms:**
- File contains "web" instead of "installed"
- Email Sentinel says "Invalid credentials"

**Solution:**
1. You created "Web application" instead of "Desktop app"
2. Go to **APIs & Services** → **Credentials**
3. Delete the wrong credential
4. Create new one, select **"Desktop app"**
5. Download again

### Problem: Can't find downloaded file

**Windows (PowerShell):**
```powershell
# Check Downloads folder
Get-ChildItem $env:USERPROFILE\Downloads\client_secret*.json
```

**Windows (CMD):**
```cmd
dir %USERPROFILE%\Downloads\client_secret*.json
```

**macOS/Linux (Bash):**
```bash
ls ~/Downloads/client_secret*.json
```

---

## Security & Privacy

### What Email Sentinel Can Access

With `gmail.readonly` scope:
- ✅ **Read** email messages and metadata
- ✅ **Search** emails
- ✅ **List** labels and threads
- ❌ **CANNOT** send emails
- ❌ **CANNOT** delete emails
- ❌ **CANNOT** modify emails
- ❌ **CANNOT** move emails to trash

### How Credentials are Used

1. **credentials.json** - Contains OAuth client ID and secret
   - Stored locally on your computer
   - Used to initiate OAuth flow

2. **token.json** - Contains your access token (created after first auth)
   - Stored locally in config directory
   - Used for API requests
   - Auto-refreshes when expired

### Data Storage

- **All data stays local** - Nothing is sent to external servers
- **No cloud storage** - Alert history is in local SQLite database
- **No analytics** - Email Sentinel doesn't phone home

### Revoking Access

If you want to revoke Email Sentinel's access:

1. Go to: https://myaccount.google.com/permissions
2. Find "Email Sentinel"
3. Click "Remove Access"

To re-authorize, run: `email-sentinel init`

### Best Practices

1. ✅ **Never share credentials.json publicly**
2. ✅ **Add credentials.json to .gitignore**
3. ✅ **Use a separate Google Cloud project for each app**
4. ✅ **Regularly review authorized apps:** https://myaccount.google.com/permissions
5. ✅ **Use read-only scopes when possible**

---

## Quota Limits

Gmail API has generous quotas for personal use:

**Daily Quota:**
- **1,000,000,000** quota units per day
- **Typical usage:** ~1,000 units per day with Email Sentinel
- **You'll never hit the limit** with normal usage

**Per-User Rate Limits:**
- **250 quota units per second**
- Email Sentinel polls every 45 seconds by default
- Well within rate limits

**What counts as quota:**
- Each API request uses quota units
- Reading a message: 5 units
- Listing messages: 5 units

**Monitoring your quota:**
1. Go to: **APIs & Services** → **Dashboard**
2. Click on **Gmail API**
3. View usage charts

---

## Cost Information

**Gmail API is FREE for:**
- Personal use
- Up to 1 billion quota units per day
- No credit card required

**When you might pay:**
- Only if you exceed quota limits (extremely unlikely)
- Exceeding limits would require 1000s of API requests per second

**Email Sentinel's typical usage:**
- ~2,000 API requests per day
- Well within free tier
- **Cost: $0.00**

---

## Next Steps

Now that you have `credentials.json`, continue with the main setup:

1. **Build Email Sentinel:** See [Build to First Alert Guide](build_to_first_alert.md)
2. **Place credentials.json** in project directory
3. **Run:** `email-sentinel init`
4. **Complete OAuth flow** in browser

---

## Additional Resources

**Google Documentation:**
- Gmail API Overview: https://developers.google.com/gmail/api/guides
- OAuth 2.0: https://developers.google.com/identity/protocols/oauth2
- API Quotas: https://developers.google.com/gmail/api/reference/quota

**Google Cloud Console:**
- Dashboard: https://console.cloud.google.com/
- APIs & Services: https://console.cloud.google.com/apis/dashboard
- Credentials: https://console.cloud.google.com/apis/credentials

**Email Sentinel:**
- Main README: See [README.md](../README.md)
- Complete Setup: See [build_to_first_alert.md](build_to_first_alert.md)

---

## Quick Reference Checklist

- [ ] Create Google Cloud project: "Email Sentinel"
- [ ] Enable Gmail API
- [ ] Configure OAuth consent screen (External)
- [ ] Add scope: `gmail.readonly`
- [ ] Add yourself as test user
- [ ] Create OAuth credentials (Desktop app)
- [ ] Download credentials as `credentials.json`
- [ ] Save file securely
- [ ] Proceed to Email Sentinel setup

**✅ You're ready to use Email Sentinel with Gmail API!**
