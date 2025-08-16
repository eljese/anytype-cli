package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetStoredAccountID(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "anytype-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	testConfig := `{"accountId":"test-account-123"}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	accountID, err := GetStoredAccountID()
	if err == nil && accountID != "" {
		t.Logf("GetStoredAccountID() = %v", accountID)
	}
}

func TestGetStoredTechSpaceID(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "anytype-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, ".anytype", "config.json")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	testConfig := `{"techSpaceId":"tech-space-789"}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	techSpaceID, err := GetStoredTechSpaceID()
	if err == nil && techSpaceID != "" {
		t.Logf("GetStoredTechSpaceID() = %v", techSpaceID)
	}
}

func TestLoadStoredConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "anytype-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, ".anytype", "config.json")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	testConfig := `{
		"accountId":"test-account-123",
		"techSpaceId":"tech-space-789"
	}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	cfg, err := LoadStoredConfig()
	if err == nil && cfg != nil {
		if cfg.AccountID != "" || cfg.TechSpaceID != "" {
			t.Logf("LoadStoredConfig() loaded config with AccountID=%v, TechSpaceID=%v",
				cfg.AccountID, cfg.TechSpaceID)
		}
	}
}
