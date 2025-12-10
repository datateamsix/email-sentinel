/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package appconfig

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestMigrationFromLegacyConfigs tests the migration from old config files
func TestMigrationFromLegacyConfigs(t *testing.T) {
	// Create a temporary directory for test configs
	tempDir := t.TempDir()

	// Create old ai-config.yaml
	aiConfig := `ai_summary:
  enabled: true
  provider: "gemini"
  api:
    gemini:
      model: "gemini-2.0-flash-exp"
      max_tokens: 1000
      temperature: 0.3
    claude:
      model: "claude-3-5-haiku-20241022"
      max_tokens: 1000
      temperature: 0.3
    openai:
      model: "gpt-4o-mini"
      max_tokens: 1000
      temperature: 0.3
  behavior:
    enable_cache: true
  rate_limit:
    max_per_hour: 60
    max_per_day: 1500
  prompt:
    system: "You are an AI assistant."
`
	if err := os.WriteFile(filepath.Join(tempDir, "ai-config.yaml"), []byte(aiConfig), 0600); err != nil {
		t.Fatalf("Failed to create test ai-config.yaml: %v", err)
	}

	// Create old rules.yaml
	rulesConfig := `priority_rules:
  urgent_keywords:
    - urgent
    - asap
    - critical
  vip_senders:
    - boss@company.com
  vip_domains:
    - costar.com
notification_settings:
  quiet_hours_start: "22:00"
  quiet_hours_end: "08:00"
  weekend_mode: "quiet"
`
	if err := os.WriteFile(filepath.Join(tempDir, "rules.yaml"), []byte(rulesConfig), 0600); err != nil {
		t.Fatalf("Failed to create test rules.yaml: %v", err)
	}

	// Create old otp_rules.yaml
	otpConfig := `enabled: true
expiry_duration: "5m"
auto_copy_to_clipboard: true
clipboard_auto_clear: "2m"
custom_patterns:
  - name: "custom_code"
    regex: "\\b\\d{6}\\b"
    confidence: 0.9
trusted_otp_senders:
  - noreply@github.com
  - accounts.google.com
`
	if err := os.WriteFile(filepath.Join(tempDir, "otp_rules.yaml"), []byte(otpConfig), 0600); err != nil {
		t.Fatalf("Failed to create test otp_rules.yaml: %v", err)
	}

	// Test migration by calling migrateAIConfig, migrateRules, migrateOTPRules
	appCfg := DefaultConfig()

	// Migrate AI config
	if err := migrateAIConfig(filepath.Join(tempDir, "ai-config.yaml"), appCfg); err != nil {
		t.Fatalf("Failed to migrate ai-config.yaml: %v", err)
	}

	// Verify AI config migration
	if !appCfg.AISummary.Enabled {
		t.Error("Expected AISummary.Enabled to be true")
	}
	if appCfg.AISummary.Provider != "gemini" {
		t.Errorf("Expected provider 'gemini', got '%s'", appCfg.AISummary.Provider)
	}
	if appCfg.AISummary.Providers.Gemini.Model != "gemini-2.0-flash-exp" {
		t.Errorf("Expected Gemini model 'gemini-2.0-flash-exp', got '%s'", appCfg.AISummary.Providers.Gemini.Model)
	}
	if !appCfg.AISummary.Cache.Enabled {
		t.Error("Expected Cache.Enabled to be true")
	}
	if appCfg.AISummary.Prompt.System != "You are an AI assistant." {
		t.Error("Prompt.System was not migrated correctly")
	}

	// Migrate rules
	if err := migrateRules(filepath.Join(tempDir, "rules.yaml"), appCfg); err != nil {
		t.Fatalf("Failed to migrate rules.yaml: %v", err)
	}

	// Verify rules migration
	if len(appCfg.Priority.UrgentKeywords) != 3 {
		t.Errorf("Expected 3 urgent keywords, got %d", len(appCfg.Priority.UrgentKeywords))
	}
	if appCfg.Priority.UrgentKeywords[0] != "urgent" {
		t.Error("UrgentKeywords not migrated correctly")
	}
	if len(appCfg.Priority.VIPSenders) != 1 || appCfg.Priority.VIPSenders[0] != "boss@company.com" {
		t.Error("VIPSenders not migrated correctly")
	}
	if appCfg.Notifications.QuietHours.Start != "22:00" {
		t.Errorf("Expected QuietHours.Start '22:00', got '%s'", appCfg.Notifications.QuietHours.Start)
	}
	if appCfg.Notifications.WeekendMode != "quiet" {
		t.Errorf("Expected WeekendMode 'quiet', got '%s'", appCfg.Notifications.WeekendMode)
	}

	// Migrate OTP rules
	if err := migrateOTPRules(filepath.Join(tempDir, "otp_rules.yaml"), appCfg); err != nil {
		t.Fatalf("Failed to migrate otp_rules.yaml: %v", err)
	}

	// Verify OTP migration
	if !appCfg.OTP.Enabled {
		t.Error("Expected OTP.Enabled to be true")
	}
	if appCfg.OTP.ExpiryDuration != "5m" {
		t.Errorf("Expected ExpiryDuration '5m', got '%s'", appCfg.OTP.ExpiryDuration)
	}
	if !appCfg.OTP.Clipboard.AutoCopy {
		t.Error("Expected Clipboard.AutoCopy to be true")
	}
	if appCfg.OTP.Clipboard.ClearAfter != "2m" {
		t.Errorf("Expected ClearAfter '2m', got '%s'", appCfg.OTP.Clipboard.ClearAfter)
	}
	if len(appCfg.OTP.CustomPatterns) != 1 {
		t.Errorf("Expected 1 custom pattern, got %d", len(appCfg.OTP.CustomPatterns))
	}
	if len(appCfg.OTP.TrustedSenders) != 2 {
		t.Errorf("Expected 2 trusted senders, got %d", len(appCfg.OTP.TrustedSenders))
	}

	t.Log("✅ All migrations successful!")
}

// TestSaveAndLoad tests saving and loading the unified config
func TestSaveAndLoad(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "app-config.yaml")

	// Create a test config
	cfg := DefaultConfig()
	cfg.AISummary.Enabled = true
	cfg.AISummary.Provider = "claude"
	cfg.Priority.UrgentKeywords = []string{"test", "urgent"}

	// Marshal and save
	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Load it back
	loadedData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var loadedCfg AppConfig
	if err := yaml.Unmarshal(loadedData, &loadedCfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify
	if !loadedCfg.AISummary.Enabled {
		t.Error("Expected AISummary.Enabled to be true")
	}
	if loadedCfg.AISummary.Provider != "claude" {
		t.Errorf("Expected provider 'claude', got '%s'", loadedCfg.AISummary.Provider)
	}
	if len(loadedCfg.Priority.UrgentKeywords) != 2 {
		t.Errorf("Expected 2 keywords, got %d", len(loadedCfg.Priority.UrgentKeywords))
	}

	t.Log("✅ Save and load test successful!")
}
