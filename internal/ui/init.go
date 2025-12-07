package ui

// Package ui provides interactive command-line interface components for Email Sentinel.
//
// This package includes:
//   - Menu system with navigation and submenus
//   - Setup wizard for first-time configuration
//   - Status dashboard for system monitoring
//   - Branding elements (banner, colors, formatting)
//   - Input helpers (confirmations, prompts, spinners)
//   - Keyboard shortcuts and help system
//
// Usage:
//
//	Interactive mode (main menu):
//	  ui.RunInteractiveMenu()
//
//	Setup wizard:
//	  wizard := ui.NewWizard()
//	  wizard.Run()
//
//	Status dashboard:
//	  ui.RunInteractiveDashboard()
//
//	Confirmation prompts:
//	  if ui.Confirm("Delete this filter?") {
//	      // proceed with deletion
//	  }
//
//	Loading spinner:
//	  spinner := ui.NewSpinner("Connecting to Gmail...")
//	  spinner.Start()
//	  // do work
//	  spinner.Stop(true)

// Init initializes the UI package (currently no initialization required)
func Init() {
	// Reserved for future initialization needs
	// Could include:
	// - Terminal capability detection
	// - Color support checking
	// - Custom theme loading
}

// Version returns the application version
func Version() string {
	return AppVersion
}

// Package constants and exports
const (
	// DefaultMenuWidth is the standard width for menu boxes
	DefaultMenuWidth = 63

	// DefaultBannerWidth is the standard width for banners
	DefaultBannerWidth = 58
)
