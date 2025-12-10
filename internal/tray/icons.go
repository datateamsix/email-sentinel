package tray

import _ "embed"

// Icon data for the system tray
// These are 16x16 ICO icons embedded at compile time
// Windows systray requires ICO format, not PNG

//go:embed icons/normal.ico
var IconNormal []byte

//go:embed icons/urgent.ico
var IconUrgent []byte

// GetNormalIcon returns the normal state icon (no alerts)
func GetNormalIcon() []byte {
	return IconNormal
}

// GetUrgentIcon returns the urgent state icon (urgent alerts present)
func GetUrgentIcon() []byte {
	return IconUrgent
}

// GetAlertIcon returns the icon for when there are alerts (any priority)
func GetAlertIcon() []byte {
	return IconUrgent
}
