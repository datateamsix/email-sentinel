/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the email-sentinel daemon",
	Long: `Stop the email-sentinel background daemon if it's running.

Note: Currently, daemon mode requires manual process management.
If you started email-sentinel in the foreground, use Ctrl+C to stop it.

For true daemon mode, consider using a process manager like:
- Windows: nssm (Non-Sucking Service Manager)
- macOS: launchd
- Linux: systemd

Example:
  email-sentinel stop`,
	Run: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop(cmd *cobra.Command, args []string) {
	fmt.Println("⏹️  Stop Command")
	fmt.Println("")
	fmt.Println("Currently, daemon management is not fully implemented.")
	fmt.Println("")
	fmt.Println("To stop email-sentinel:")
	fmt.Println("  • If running in foreground: Press Ctrl+C")
	fmt.Println("  • If running as daemon: Manually kill the process")
	fmt.Println("")
	fmt.Println("Future versions will include full daemon management.")
}
