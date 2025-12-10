package tray

import (
	"log"
	"os/exec"
	"runtime"
)

// addFilterGUI opens an interactive terminal for adding a new filter
func (app *TrayApp) addFilterGUI() {
	log.Println("📝 Opening Add Filter dialog...")

	go func() {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			// Open cmd with the filter add command
			script := "email-sentinel.exe filter add && pause"
			cmd = exec.Command("cmd", "/c", "start", "cmd", "/k", script)

		case "darwin":
			// macOS - open Terminal with command
			script := "email-sentinel filter add && read -p 'Press Enter to close...'"
			cmd = exec.Command("osascript", "-e", `tell application "Terminal" to do script "`+script+`"`)

		default:
			// Linux - try various terminal emulators
			script := "email-sentinel filter add && read -p 'Press Enter to close...'"
			terminals := []string{"gnome-terminal", "konsole", "xterm"}
			for _, term := range terminals {
				if _, err := exec.LookPath(term); err == nil {
					if term == "gnome-terminal" {
						cmd = exec.Command(term, "--", "bash", "-c", script)
					} else {
						cmd = exec.Command(term, "-e", "bash", "-c", script)
					}
					break
				}
			}
		}

		if cmd != nil {
			if err := cmd.Start(); err != nil {
				log.Printf("❌ Error opening Add Filter dialog: %v", err)
			}
		}
	}()
}

// editFilterGUI opens an interactive terminal for editing a filter
func (app *TrayApp) editFilterGUI() {
	log.Println("✏️ Opening Edit Filter dialog...")

	go func() {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			// Open cmd with filter list and edit command
			script := "email-sentinel.exe filter list && echo. && email-sentinel.exe filter edit && pause"
			cmd = exec.Command("cmd", "/c", "start", "cmd", "/k", script)

		case "darwin":
			// macOS - open Terminal with command
			script := "email-sentinel filter list && echo && email-sentinel filter edit && read -p 'Press Enter to close...'"
			cmd = exec.Command("osascript", "-e", `tell application "Terminal" to do script "`+script+`"`)

		default:
			// Linux - try various terminal emulators
			script := "email-sentinel filter list && echo && email-sentinel filter edit && read -p 'Press Enter to close...'"
			terminals := []string{"gnome-terminal", "konsole", "xterm"}
			for _, term := range terminals {
				if _, err := exec.LookPath(term); err == nil {
					if term == "gnome-terminal" {
						cmd = exec.Command(term, "--", "bash", "-c", script)
					} else {
						cmd = exec.Command(term, "-e", "bash", "-c", script)
					}
					break
				}
			}
		}

		if cmd != nil {
			if err := cmd.Start(); err != nil {
				log.Printf("❌ Error opening Edit Filter dialog: %v", err)
			}
		}
	}()
}
