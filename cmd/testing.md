✅ Next Steps

  PowerShell Tip

  In PowerShell, you need .\ before local executables:

  # ❌ Wrong (doesn't work)
  email-sentinel filter list

  # ✅ Correct (works)
  .\email-sentinel.exe filter list

  ---
  Step 1: Add Your First Filter

  Let's create a test filter using your own email address (so you can easily test it):

  # Add a filter for emails from yourself
  .\email-sentinel.exe filter add --name "Self Test" --from "YOUR_EMAIL@gmail.com"

  Replace YOUR_EMAIL@gmail.com with your actual Gmail address!

  Or use interactive mode:
  .\email-sentinel.exe filter add
  Then follow the prompts.

  ---
  Step 2: Verify Filter Was Created

  .\email-sentinel.exe filter list

  Expected output:
  📋 Email Filters (1)
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  [1] Self Test
      From:    YOUR_EMAIL@gmail.com
      Subject: (any)
      Match:   any (OR - either condition triggers)

  ---
  Step 3: Test Desktop Notification

  .\email-sentinel.exe test desktop

  Expected: A notification should pop up saying "Email Sentinel Test"

  If no notification appears:
  - Windows: Settings → Notifications → Enable
  - Turn off Focus Assist

  ---
  Step 4: Test Monitoring Manually

  .\email-sentinel.exe start

  Expected output:
  ✅ Email Sentinel Started
     Monitoring 1 filter(s)
     Polling interval: 45 seconds
     Desktop notifications: enabled

  🔍 Watching for new emails... (Press Ctrl+C to stop)

  [HH:MM:SS] Checked 10 messages, no new matches

  Let it run for ~30 seconds to verify it works, then press Ctrl+C to stop.

  ---
  Step 5: Send Yourself a Test Email

  While start is running:

  1. Open Gmail in browser
  2. Compose new email
  3. To: Your own email address (same one in filter)
  4. Subject: Test for Email Sentinel
  5. Send!

  Within ~45 seconds, you should see:
  - ✅ Terminal shows: 📧 MATCH [Self Test] From: your@email.com | Subject: Test...
  - ✅ Desktop notification pops up

  If that works, press Ctrl+C and continue to Step 6!

  ---
  Step 6: Enable Startup Task

  Now that everything works manually, enable auto-startup:

  # Enable the task
  schtasks /Change /TN "EmailSentinel" /ENABLE

  # Test it runs via the task
  schtasks /Run /TN "EmailSentinel"

  # Wait a few seconds
  timeout 5

  # Verify it's running
  tasklist | findstr email-sentinel

  Expected output:
  email-sentinel.exe            12345 Console                    1     25,000 K

  If you see that → ✅ Success! Auto-startup is working!

  ---
  Step 7: Final Verification

  # Check complete status
  .\email-sentinel.exe status

  Expected:
  📊 Email Sentinel Status
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ✅ Authentication: Configured
     Token: C:\Users\zroda\AppData\Roaming\email-sentinel\token.json

  📋 Filters: 1 configured
     [1] Self Test

  ⚙️  Polling Interval: 45 seconds

  🔔 Notifications:
     Desktop: Enabled
     Mobile: Disabled

  📁 Config File: C:\Users\zroda\AppData\Roaming\email-sentinel\config.yaml

  ---
  🎯 What Happens Now

  ✅ Email Sentinel is installed and will:
  1. Start automatically when you log in to Windows
  2. Monitor Gmail every 45 seconds
  3. Send desktop notifications when emails match your filters
  4. Auto-refresh OAuth token (never expire)
  5. Handle network errors gracefully with retry logic

  ---
  📝 Quick Reference Commands

  # Manage filters
  .\email-sentinel.exe filter add       # Add new filter
  .\email-sentinel.exe filter list      # View all filters
  .\email-sentinel.exe filter edit      # Modify filter
  .\email-sentinel.exe filter remove    # Delete filter

  # Manage startup
  .\email-sentinel.exe uninstall        # Remove from startup
  .\email-sentinel.exe install          # Reinstall

  # Check status
  .\email-sentinel.exe status           # View config
  .\email-sentinel.exe config show      # View settings

  # Manual control
  tasklist | findstr email-sentinel     # Check if running
  taskkill /IM email-sentinel.exe /F    # Stop it
  schtasks /Run /TN "EmailSentinel"     # Start it

  ---
  🚀 Next Steps

  1. Add real filters for important emails:
  .\email-sentinel.exe filter add --name "Boss" --from "boss@company.com"
  .\email-sentinel.exe filter add --name "Job Alerts" --from "linkedin.com,greenhouse.io"
  2. Adjust polling interval (optional):
  .\email-sentinel.exe config set polling 60  # Check every 60 seconds
  3. Setup mobile notifications (optional):
    - Install ntfy app on phone
    - Subscribe to a unique topic
    - Run: .\email-sentinel.exe config set ntfy_topic "your-topic"
    - Run: .\email-sentinel.exe config set mobile true
    - Test: .\email-sentinel.exe test mobile
  4. Reboot to verify auto-start:
    - Restart your computer
    - After login, run: tasklist | findstr email-sentinel
    - Should show it's running automatically!
