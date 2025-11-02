package core

import (
	"errors"
	"os"
	"testing"

	"github.com/anyproto/anytype-cli/core/config"
)

func setupTestHome(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "anytype-keyring-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
	t.Setenv("HOME", tempDir)
	return tempDir
}

func TestKeyringFallback(t *testing.T) {
	setupTestHome(t)

	keyringUnavailable = true
	defer func() { keyringUnavailable = false }()

	testAccountKey := "test-account-key-12345"

	savedToKeyring, err := SaveAccountKey(testAccountKey)
	if err != nil {
		t.Fatalf("SaveAccountKey failed: %v", err)
	}
	if savedToKeyring {
		t.Error("Expected SaveAccountKey to return false (config fallback) when keyring unavailable, got true")
	}

	configMgr := config.GetConfigManager()
	err = configMgr.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	cfg := configMgr.Get()
	if cfg.AccountKey != testAccountKey {
		t.Errorf("Expected account key %q, got %q", testAccountKey, cfg.AccountKey)
	}

	retrievedKey, err := GetStoredAccountKey()
	if err != nil {
		t.Fatalf("GetStoredAccountKey failed: %v", err)
	}

	if retrievedKey != testAccountKey {
		t.Errorf("Expected retrieved key %q, got %q", testAccountKey, retrievedKey)
	}

	err = DeleteStoredAccountKey()
	if err != nil {
		t.Fatalf("DeleteStoredAccountKey failed: %v", err)
	}

	err = configMgr.Load()
	if err != nil {
		t.Fatalf("Failed to load config after delete: %v", err)
	}

	cfg = configMgr.Get()
	if cfg.AccountKey != "" {
		t.Errorf("Expected empty account key after delete, got %q", cfg.AccountKey)
	}

	_, err = GetStoredAccountKey()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestTokenFallback(t *testing.T) {
	setupTestHome(t)

	keyringUnavailable = true
	defer func() { keyringUnavailable = false }()

	testToken := "test-session-token-67890"

	savedToKeyring, err := SaveToken(testToken)
	if err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}
	if savedToKeyring {
		t.Error("Expected SaveToken to return false (config fallback) when keyring unavailable, got true")
	}

	configMgr := config.GetConfigManager()
	err = configMgr.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	cfg := configMgr.Get()
	if cfg.SessionToken != testToken {
		t.Errorf("Expected session token %q, got %q", testToken, cfg.SessionToken)
	}

	retrievedToken, err := GetStoredToken()
	if err != nil {
		t.Fatalf("GetStoredToken failed: %v", err)
	}

	if retrievedToken != testToken {
		t.Errorf("Expected retrieved token %q, got %q", testToken, retrievedToken)
	}

	err = DeleteStoredToken()
	if err != nil {
		t.Fatalf("DeleteStoredToken failed: %v", err)
	}

	err = configMgr.Load()
	if err != nil {
		t.Fatalf("Failed to load config after delete: %v", err)
	}

	cfg = configMgr.Get()
	if cfg.SessionToken != "" {
		t.Errorf("Expected empty session token after delete, got %q", cfg.SessionToken)
	}

	_, err = GetStoredToken()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestConfigFilePermissions(t *testing.T) {
	setupTestHome(t)

	keyringUnavailable = true
	defer func() { keyringUnavailable = false }()

	testAccountKey := "test-account-key-permissions"
	savedToKeyring, err := SaveAccountKey(testAccountKey)
	if err != nil {
		t.Fatalf("SaveAccountKey failed: %v", err)
	}
	if savedToKeyring {
		t.Error("Expected SaveAccountKey to return false (config fallback) when keyring unavailable, got true")
	}

	configMgr := config.GetConfigManager()
	configPath := configMgr.GetFilePath()
	if configPath == "" {
		t.Skip("Could not determine config file path")
	}

	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	mode := fileInfo.Mode().Perm()
	expectedMode := os.FileMode(0600)
	if mode != expectedMode {
		t.Errorf("Expected file permissions %v, got %v", expectedMode, mode)
	}
}

func TestEmptyCredentialRetrieval(t *testing.T) {
	setupTestHome(t)

	keyringUnavailable = true
	defer func() { keyringUnavailable = false }()

	configMgr := config.GetConfigManager()
	err := configMgr.Delete()
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to delete config: %v", err)
	}

	_, err = GetStoredAccountKey()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound for non-existent account key, got %v", err)
	}

	_, err = GetStoredToken()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound for non-existent token, got %v", err)
	}
}
