package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// ConfigExists checks if the config file exists
func ConfigExists() bool {
	path, err := ConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// Load reads the config file and unmarshals it into v
func Load(v interface{}) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}

// Save marshals v and writes it to the config file
func Save(v interface{}) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if _, err := EnsureConfigDir(); err != nil {
		return err
	}

	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
