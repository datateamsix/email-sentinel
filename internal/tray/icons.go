package tray

import _ "embed"

// Icon data for the system tray
// These are 16x16 ICO icons embedded at compile time
// Windows systray requires ICO format, not PNG

//go:embed icons/normal.ico
var IconNormal []byte

//go:embed icons/urgent.ico
var IconUrgent []byte

// GetNormalIcon returns the normal state icon
func GetNormalIcon() []byte {
	return IconNormal
}

// GetUrgentIcon returns the urgent state icon
func GetUrgentIcon() []byte {
	return IconUrgent
}
