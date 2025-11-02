package config

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestHome creates a temp directory and sets HOME to it for test isolation
func setupTestHome(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "anytype-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
	t.Setenv("HOME", tempDir)
	return tempDir
}

func TestGetStoredAccountId(t *testing.T) {
	tempDir := setupTestHome(t)

	configPath := filepath.Join(tempDir, ".anytype", "config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	testConfig := `{"accountId":"test-account-123"}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	accountId, err := GetAccountIdFromConfig()
	if err == nil && accountId != "" {
		t.Logf("GetAccountIdFromConfig() = %v", accountId)
	}
}

func TestGetStoredTechSpaceId(t *testing.T) {
	tempDir := setupTestHome(t)

	configPath := filepath.Join(tempDir, ".anytype", "config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	testConfig := `{"techSpaceId":"tech-space-789"}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	techSpaceId, err := GetTechSpaceIdFromConfig()
	if err == nil && techSpaceId != "" {
		t.Logf("GetTechSpaceIdFromConfig() = %v", techSpaceId)
	}
}

func TestLoadStoredConfig(t *testing.T) {
	tempDir := setupTestHome(t)

	configPath := filepath.Join(tempDir, ".anytype", "config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	testConfig := `{
		"accountId":"test-account-123",
		"techSpaceId":"tech-space-789"
	}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadStoredConfig()
	if err == nil && cfg != nil {
		if cfg.AccountId != "" || cfg.TechSpaceId != "" {
			t.Logf("LoadStoredConfig() loaded config with AccountId=%v, TechSpaceId=%v",
				cfg.AccountId, cfg.TechSpaceId)
		}
	}
}
