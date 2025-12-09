# Migration Guide: getlantern/systray → fyne.io/systray

## Quick Reference for Afternoon Implementation

**Estimated Time:** 1 hour
**Risk Level:** Low
**Rollback:** Easy

---

## Step 1: Backup & Branch (5 minutes)

```bash
cd C:\Users\zroda\Desktop\email-sentinel-cli

# Create backup branch
git checkout -b backup/pre-fyne-migration

# Create feature branch
git checkout -b feature/migrate-to-fyne-systray
```

---

## Step 2: Update Dependencies (5 minutes)

```bash
# Remove old dependency
go get github.com/getlantern/systray@none

# Add fyne.io/systray
go get fyne.io/systray@latest

# Clean up
go mod tidy

# Verify
go list -m fyne.io/systray
# Expected: fyne.io/systray v1.11.0 or newer
```

---

## Step 3: Update Code (5 minutes)

### File: `internal/tray/tray.go`

**Line 13 - Change import:**

```diff
import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/datateamsix/email-sentinel/internal/storage"
-	"github.com/getlantern/systray"
+	"fyne.io/systray"
)
```

**That's it! API is compatible, no other changes needed.**

---

## Step 4: Build & Test (30 minutes)

### Build

```bash
# Windows
set CGO_ENABLED=1
go build -o email-sentinel.exe .

# Verify build
.\email-sentinel.exe --version
```

### Test Suite

#### Test 1: Basic Tray Launch
```bash
.\email-sentinel.exe start --tray
```
**Expected:** Tray icon appears in system tray

#### Test 2: Menu Functionality
1. Click tray icon
2. Verify "Recent Alerts" menu
3. Check "Open History" option
4. Check "Quit" option

**Expected:** All menus visible and clickable

#### Test 3: Alert Update
```bash
# In another terminal, trigger test notification
.\email-sentinel.exe test desktop
```
**Expected:** Tray refreshes, shows new alert

#### Test 4: AI Summary Integration
```bash
.\email-sentinel.exe start --tray --ai-summary
```
**Expected:** Tray starts with AI enabled

#### Test 5: OTP Display
```bash
.\email-sentinel.exe otp test "Your verification code is 849372"
```
**Expected:** OTP code detected and can be viewed in tray

#### Test 6: Gmail Link Click
1. Start tray
2. Click on an alert in "Recent Alerts"
3. **Expected:** Browser opens Gmail link

#### Test 7: Icon State Change
1. Create high-priority alert
2. **Expected:** Icon changes to urgent (red/orange)
3. Wait 5 seconds
4. **Expected:** Icon returns to normal

#### Test 8: Cleanup Interval
```bash
.\email-sentinel.exe start --tray --cleanup-interval 1
```
**Expected:** Auto-cleanup runs every minute (check logs)

---

## Step 5: Cross-Platform Testing (If Available)

### Linux Testing (Ubuntu/Debian)

```bash
# Should NOT require GTK3 headers anymore
export CGO_ENABLED=1
go build -o email-sentinel .

# Start tray
./email-sentinel start --tray
```

**Expected:** Builds without needing `libgtk-3-dev` or `libayatana-appindicator3-dev`

### macOS Testing

```bash
export CGO_ENABLED=1
go build -o email-sentinel .

./email-sentinel start --tray
```

**Expected:** Icon appears in macOS menu bar

---

## Step 6: Performance Check (5 minutes)

### Memory Usage

```bash
# Windows
tasklist /FI "IMAGENAME eq email-sentinel.exe" /FO LIST

# Expected: ~15-25 MB (similar to before)
```

### CPU Usage

```bash
# Monitor during normal operation
# Expected: <0.1% idle, <2% during alert updates
```

---

## Step 7: Documentation Update (10 minutes)

### Update README.md

**Find section:** "## 🔧 Prerequisites"

**Before:**
```markdown
### Linux Dependencies

```bash
sudo apt-get install gcc libgtk-3-dev libayatana-appindicator3-dev
```
```

**After:**
```markdown
### Linux Dependencies

```bash
# Only gcc required for CGO
sudo apt-get install gcc
```
```

### Update docs/build_to_first_alert.md

Search for Linux dependency instructions and remove GTK3 references.

---

## Rollback Procedure (If Needed)

```bash
# Switch back to backup branch
git checkout backup/pre-fyne-migration

# Rebuild with old version
go build -o email-sentinel.exe .

# Test
.\email-sentinel.exe start --tray
```

---

## Success Criteria

- [x] Build completes without errors
- [x] Tray icon appears on Windows
- [x] Menus are functional
- [x] Recent alerts populate correctly
- [x] Click alerts → Opens Gmail
- [x] Icon changes on urgent alerts
- [x] AI summaries display in tooltips
- [x] OTP codes show with 🔐 icon
- [x] Quit from menu exits cleanly
- [x] No memory leaks after 10 minutes
- [x] Performance is comparable

**If all checked:** Ready to commit!

---

## Commit & Push

```bash
git add .
git commit -m "feat: Migrate to fyne.io/systray for simplified Linux deployment

- Remove GTK3 dependency requirement on Linux
- Update to fyne.io/systray v1.11.0
- Maintain full API compatibility
- Simplify build instructions for Linux users

Closes #[issue-number-if-applicable]"

git push origin feature/migrate-to-fyne-systray
```

---

## Post-Migration Validation

### Run Full Test Suite

```bash
# Start monitoring
.\email-sentinel.exe start --tray

# Open another terminal
.\email-sentinel.exe test desktop
.\email-sentinel.exe test toast
.\email-sentinel.exe otp test "Code: 123456"
.\email-sentinel.exe alerts

# Check tray responsiveness
# Check logs for errors
# Monitor for 15 minutes
```

### Stress Test

```bash
# Run for extended period
.\email-sentinel.exe start --tray --ai-summary --cleanup-interval 60

# Leave running overnight (optional)
# Check memory usage next morning
# Should be stable with no leaks
```

---

## Troubleshooting

### Issue: Build fails with CGO errors

```bash
# Verify CGO is enabled
go env CGO_ENABLED
# Should output: 1

# If not
set CGO_ENABLED=1  # Windows
export CGO_ENABLED=1  # Linux/macOS
```

### Issue: Tray icon not appearing

**Windows:**
- Check taskbar settings → "Select which icons appear on taskbar"
- Enable "Email Sentinel" if hidden

**Linux:**
- Install GNOME Shell extension: `gnome-shell-extension-appindicator`
- Restart GNOME Shell: `Alt+F2`, type `r`, press Enter

**macOS:**
- Check menu bar spacing settings
- Restart the application

### Issue: Menu not updating

```bash
# Check database
.\email-sentinel.exe alerts

# If empty, create test alert
.\email-sentinel.exe test desktop

# Verify tray refresh logic
# Check logs for update errors
```

---

## Additional Resources

- **Full Architecture Doc:** `docs/tray_system_architecture.md`
- **fyne.io/systray Docs:** https://pkg.go.dev/fyne.io/systray
- **GitHub Repo:** https://github.com/fyne-io/systray
- **API Reference:** https://pkg.go.dev/fyne.io/systray#section-documentation

---

## Notes

- API is 99% compatible with getlantern/systray
- Only import statement needs changing
- No logic rewrites required
- GTK3 dependency removed on Linux
- All features remain functional

**Good luck with the migration! 🚀**
