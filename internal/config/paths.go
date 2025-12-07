package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const AppName = "email-sentinel"

// ConfigDir returns the OS-appropriate config directory
func ConfigDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
	default: // linux and others
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(baseDir, AppName), nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return dir, nil
}

// TokenPath returns the path where OAuth token should be stored
func TokenPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "token.json"), nil
}
