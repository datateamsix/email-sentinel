/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install email-sentinel to run automatically on startup",
	Long: `Install email-sentinel as a system service/startup app.

This configures email-sentinel to start automatically when you log in to your computer.
After installation, email-sentinel will monitor your Gmail in the background.

Platform Support:
  • Windows:  Creates a Task Scheduler task
  • macOS:    Creates a launchd agent
  • Linux:    Creates a systemd user service

Requirements:
  • Must have completed 'email-sentinel init'  (OAuth configured)
  • Must have at least one filter configured

Examples:
  email-sentinel install          # Install for current user
  email-sentinel install --show   # Show what would be installed`,
	Run: runInstall,
}

var showOnly bool

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&showOnly, "show", false, "Show installation steps without installing")
}

func runInstall(cmd *cobra.Command, args []string) {
	fmt.Println("📦 Email Sentinel - Auto-Startup Installation")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")

	// Get executable path
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("❌ Error: Could not determine executable path: %v\n", err)
		os.Exit(1)
	}

	// Resolve symlinks
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		fmt.Printf("❌ Error: Could not resolve executable path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Executable: %s\n", exePath)
	fmt.Printf("Platform:   %s\n", runtime.GOOS)
	fmt.Println("")

	switch runtime.GOOS {
	case "windows":
		installWindows(exePath, showOnly)
	case "darwin":
		installMacOS(exePath, showOnly)
	case "linux":
		installLinux(exePath, showOnly)
	default:
		fmt.Printf("❌ Unsupported platform: %s\n", runtime.GOOS)
		os.Exit(1)
	}
}

func installWindows(exePath string, showOnly bool) {
	taskName := "EmailSentinel"
	xmlPath := filepath.Join(os.TempDir(), "email-sentinel-task.xml")

	// Create Task Scheduler XML
	xmlContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Description>Email Sentinel - Gmail notification monitor</Description>
  </RegistrationInfo>
  <Triggers>
    <LogonTrigger>
      <Enabled>true</Enabled>
    </LogonTrigger>
  </Triggers>
  <Principals>
    <Principal>
      <LogonType>InteractiveToken</LogonType>
      <RunLevel>LeastPrivilege</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>true</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>true</RunOnlyIfNetworkAvailable>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions>
    <Exec>
      <Command>%s</Command>
      <Arguments>start --tray</Arguments>
    </Exec>
  </Actions>
</Task>`, exePath)

	if showOnly {
		fmt.Println("📋 Installation Preview (Windows):")
		fmt.Println("")
		fmt.Println("1. Task Scheduler task will be created:")
		fmt.Printf("   Name: %s\n", taskName)
		fmt.Printf("   Executable: %s start --tray\n", exePath)
		fmt.Println("   Trigger: At logon")
		fmt.Println("")
		fmt.Println("2. Command that would be run:")
		fmt.Printf("   schtasks /Create /XML \"%s\" /TN \"%s\" /F\n", xmlPath, taskName)
		fmt.Println("")
		fmt.Println("Run without --show to perform installation")
		return
	}

	// Write XML file
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0644); err != nil {
		fmt.Printf("❌ Error creating task XML: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(xmlPath)

	// Create task using schtasks
	cmd := exec.Command("schtasks", "/Create", "/XML", xmlPath, "/TN", taskName, "/F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Error creating task: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		fmt.Println("")
		fmt.Println("Manual installation:")
		fmt.Println("1. Open Task Scheduler (taskschd.msc)")
		fmt.Println("2. Create Basic Task")
		fmt.Printf("3. Trigger: When I log on\n")
		fmt.Printf("4. Action: Start a program\n")
		fmt.Printf("5. Program: %s\n", exePath)
		fmt.Printf("6. Arguments: start\n")
		os.Exit(1)
	}

	fmt.Println("✅ Successfully installed!")
	fmt.Println("")
	fmt.Println("Email Sentinel will now start automatically when you log in.")
	fmt.Println("")
	fmt.Println("To manage:")
	fmt.Println("  • View: Task Scheduler (taskschd.msc)")
	fmt.Println("  • Disable: schtasks /Change /TN \"EmailSentinel\" /DISABLE")
	fmt.Println("  • Remove: schtasks /Delete /TN \"EmailSentinel\" /F")
	fmt.Println("")
	fmt.Println("To start now:")
	fmt.Printf("  %s start\n", exePath)
}

func installMacOS(exePath string, showOnly bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("❌ Error: Could not determine home directory: %v\n", err)
		os.Exit(1)
	}

	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	plistPath := filepath.Join(launchAgentsDir, "com.datateamsix.email-sentinel.plist")

	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.datateamsix.email-sentinel</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>start</string>
        <string>--tray</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <dict>
        <key>NetworkState</key>
        <true/>
    </dict>
    <key>StandardOutPath</key>
    <string>%s/Library/Logs/email-sentinel.log</string>
    <key>StandardErrorPath</key>
    <string>%s/Library/Logs/email-sentinel-error.log</string>
</dict>
</plist>`, exePath, homeDir, homeDir)

	if showOnly {
		fmt.Println("📋 Installation Preview (macOS):")
		fmt.Println("")
		fmt.Println("1. LaunchAgent plist will be created:")
		fmt.Printf("   %s\n", plistPath)
		fmt.Println("")
		fmt.Println("2. Commands that would be run:")
		fmt.Printf("   launchctl unload %s (if exists)\n", plistPath)
		fmt.Printf("   launchctl load %s\n", plistPath)
		fmt.Println("")
		fmt.Println("Run without --show to perform installation")
		return
	}

	// Ensure LaunchAgents directory exists
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		fmt.Printf("❌ Error creating LaunchAgents directory: %v\n", err)
		os.Exit(1)
	}

	// Write plist file
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		fmt.Printf("❌ Error creating plist file: %v\n", err)
		os.Exit(1)
	}

	// Unload if already loaded (ignore errors)
	exec.Command("launchctl", "unload", plistPath).Run()

	// Load the agent
	cmd := exec.Command("launchctl", "load", plistPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Error loading LaunchAgent: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		os.Exit(1)
	}

	fmt.Println("✅ Successfully installed!")
	fmt.Println("")
	fmt.Println("Email Sentinel will now start automatically when you log in.")
	fmt.Println("")
	fmt.Printf("Configuration: %s\n", plistPath)
	fmt.Printf("Logs: ~/Library/Logs/email-sentinel.log\n")
	fmt.Println("")
	fmt.Println("To manage:")
	fmt.Printf("  • Stop:    launchctl unload %s\n", plistPath)
	fmt.Printf("  • Start:   launchctl load %s\n", plistPath)
	fmt.Printf("  • Remove:  rm %s\n", plistPath)
}

func installLinux(exePath string, showOnly bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("❌ Error: Could not determine home directory: %v\n", err)
		os.Exit(1)
	}

	systemdDir := filepath.Join(homeDir, ".config", "systemd", "user")
	servicePath := filepath.Join(systemdDir, "email-sentinel.service")

	serviceContent := fmt.Sprintf(`[Unit]
Description=Email Sentinel - Gmail notification monitor
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=%s start --tray
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
`, exePath)

	if showOnly {
		fmt.Println("📋 Installation Preview (Linux):")
		fmt.Println("")
		fmt.Println("1. Systemd user service will be created:")
		fmt.Printf("   %s\n", servicePath)
		fmt.Println("")
		fmt.Println("2. Commands that would be run:")
		fmt.Println("   systemctl --user daemon-reload")
		fmt.Println("   systemctl --user enable email-sentinel")
		fmt.Println("   systemctl --user start email-sentinel")
		fmt.Println("")
		fmt.Println("Run without --show to perform installation")
		return
	}

	// Ensure systemd user directory exists
	if err := os.MkdirAll(systemdDir, 0755); err != nil {
		fmt.Printf("❌ Error creating systemd directory: %v\n", err)
		os.Exit(1)
	}

	// Write service file
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		fmt.Printf("❌ Error creating service file: %v\n", err)
		os.Exit(1)
	}

	// Reload systemd
	if err := exec.Command("systemctl", "--user", "daemon-reload").Run(); err != nil {
		fmt.Printf("❌ Error reloading systemd: %v\n", err)
		os.Exit(1)
	}

	// Enable service
	if err := exec.Command("systemctl", "--user", "enable", "email-sentinel").Run(); err != nil {
		fmt.Printf("❌ Error enabling service: %v\n", err)
		os.Exit(1)
	}

	// Start service
	cmd := exec.Command("systemctl", "--user", "start", "email-sentinel")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Error starting service: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		fmt.Println("")
		fmt.Println("Check logs with: journalctl --user -u email-sentinel -f")
		os.Exit(1)
	}

	fmt.Println("✅ Successfully installed and started!")
	fmt.Println("")
	fmt.Println("Email Sentinel is now running and will start automatically on boot.")
	fmt.Println("")
	fmt.Printf("Configuration: %s\n", servicePath)
	fmt.Println("")
	fmt.Println("To manage:")
	fmt.Println("  • Status:  systemctl --user status email-sentinel")
	fmt.Println("  • Stop:    systemctl --user stop email-sentinel")
	fmt.Println("  • Restart: systemctl --user restart email-sentinel")
	fmt.Println("  • Logs:    journalctl --user -u email-sentinel -f")
	fmt.Println("  • Disable: systemctl --user disable email-sentinel")
}

// Uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove email-sentinel from startup",
	Long:  `Remove email-sentinel from automatic startup configuration.`,
	Run:   runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) {
	fmt.Println("🗑️  Email Sentinel - Uninstall")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")

	switch runtime.GOOS {
	case "windows":
		uninstallWindows()
	case "darwin":
		uninstallMacOS()
	case "linux":
		uninstallLinux()
	default:
		fmt.Printf("❌ Unsupported platform: %s\n", runtime.GOOS)
		os.Exit(1)
	}
}

func uninstallWindows() {
	cmd := exec.Command("schtasks", "/Delete", "/TN", "EmailSentinel", "/F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "cannot find") {
			fmt.Println("⚠️  Email Sentinel is not installed")
		} else {
			fmt.Printf("❌ Error: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
		}
		return
	}

	fmt.Println("✅ Email Sentinel removed from startup")
}

func uninstallMacOS() {
	homeDir, _ := os.UserHomeDir()
	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents", "com.datateamsix.email-sentinel.plist")

	// Unload
	exec.Command("launchctl", "unload", plistPath).Run()

	// Remove file
	if err := os.Remove(plistPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("⚠️  Email Sentinel is not installed")
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
		return
	}

	fmt.Println("✅ Email Sentinel removed from startup")
}

func uninstallLinux() {
	// Stop service
	exec.Command("systemctl", "--user", "stop", "email-sentinel").Run()

	// Disable service
	exec.Command("systemctl", "--user", "disable", "email-sentinel").Run()

	// Remove service file
	homeDir, _ := os.UserHomeDir()
	servicePath := filepath.Join(homeDir, ".config", "systemd", "user", "email-sentinel.service")

	if err := os.Remove(servicePath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("⚠️  Email Sentinel is not installed")
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
		return
	}

	// Reload systemd
	exec.Command("systemctl", "--user", "daemon-reload").Run()

	fmt.Println("✅ Email Sentinel removed from startup")
}
