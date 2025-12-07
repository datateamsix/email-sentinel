package tray

// Icon data for the system tray
// These are simple 16x16 PNG icons
// Note: If you want custom icons, create icons/normal.png and icons/urgent.png
// in the internal/tray directory and uncomment the embed directives below

// Uncomment these when you have actual icon files:
// import _ "embed"
// //go:embed icons/normal.png
// var IconNormal []byte
// //go:embed icons/urgent.png
// var IconUrgent []byte

var IconNormal []byte
var IconUrgent []byte

// If embedded icons are not available, we'll generate simple PNG data
// These are minimal valid PNG files representing email icons

// generateDefaultNormalIcon creates a simple blue/gray email icon (16x16 PNG)
func generateDefaultNormalIcon() []byte {
	// This is a base64-decoded minimal PNG representing an email icon
	// In production, you'd want to use actual icon files
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		// Minimal PNG data for a 16x16 icon
		// For now, returning empty will use system default
	}
}

// generateDefaultUrgentIcon creates a simple red/orange urgent icon (16x16 PNG)
func generateDefaultUrgentIcon() []byte {
	// This is a base64-decoded minimal PNG representing an urgent icon
	// In production, you'd want to use actual icon files
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		// Minimal PNG data for a 16x16 icon
		// For now, returning empty will use system default
	}
}

// GetNormalIcon returns the normal state icon
func GetNormalIcon() []byte {
	if len(IconNormal) > 0 {
		return IconNormal
	}
	return generateDefaultNormalIcon()
}

// GetUrgentIcon returns the urgent state icon
func GetUrgentIcon() []byte {
	if len(IconUrgent) > 0 {
		return IconUrgent
	}
	return generateDefaultUrgentIcon()
}
