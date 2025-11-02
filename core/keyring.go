package core

import (
	"errors"

	"github.com/anyproto/anytype-cli/core/config"
	"github.com/zalando/go-keyring"

	"github.com/anyproto/anytype-cli/core/output"
)

const (
	keyringService        = "anytype-cli"
	keyringTokenUser      = "session-token"
	keyringAccountKeyUser = "account-key"
)

var (
	keyringUnavailable = false
)

// isKeyringAvailable checks if the keyring is accessible
func isKeyringAvailable() bool {
	if keyringUnavailable {
		return false
	}

	err := keyring.Set(keyringService, "test", "test")
	if err != nil {
		keyringUnavailable = true
		return false
	}
	_ = keyring.Delete(keyringService, "test")
	return true
}

func SaveToken(token string) error {
	if isKeyringAvailable() {
		return keyring.Set(keyringService, keyringTokenUser, token)
	}

	output.Warning("System keyring unavailable (requires D-Bus on Linux, Keychain on macOS, Credential Manager on Windows)")
	output.Warning("Storing session token in config file: %s (insecure)", config.GetConfigManager().GetFilePath())

	return config.SetSessionTokenToConfig(token)
}

func GetStoredToken() (string, error) {
	if isKeyringAvailable() {
		token, err := keyring.Get(keyringService, keyringTokenUser)
		if err == nil {
			return token, nil
		}
		if !errors.Is(err, keyring.ErrNotFound) {
			keyringUnavailable = true
		}
	}

	token, _ := config.GetSessionTokenFromConfig()
	if token == "" {
		return "", keyring.ErrNotFound
	}
	return token, nil
}

func DeleteStoredToken() error {
	var keyringErr error
	if isKeyringAvailable() {
		keyringErr = keyring.Delete(keyringService, keyringTokenUser)
	}

	configErr := config.SetSessionTokenToConfig("")

	if keyringErr != nil && configErr != nil {
		return configErr
	}
	return nil
}

func SaveAccountKey(accountKey string) error {
	if isKeyringAvailable() {
		return keyring.Set(keyringService, keyringAccountKeyUser, accountKey)
	}

	output.Warning("System keyring unavailable (requires D-Bus on Linux, Keychain on macOS, Credential Manager on Windows)")
	output.Warning("Storing account key in config file: %s (insecure)", config.GetConfigManager().GetFilePath())

	return config.SetAccountKeyToConfig(accountKey)
}

func GetStoredAccountKey() (string, error) {
	if isKeyringAvailable() {
		key, err := keyring.Get(keyringService, keyringAccountKeyUser)
		if err == nil {
			return key, nil
		}
		if !errors.Is(err, keyring.ErrNotFound) {
			keyringUnavailable = true
		}
	}

	key, _ := config.GetAccountKeyFromConfig()
	if key == "" {
		return "", keyring.ErrNotFound
	}
	return key, nil
}

func DeleteStoredAccountKey() error {
	var keyringErr error
	if isKeyringAvailable() {
		keyringErr = keyring.Delete(keyringService, keyringAccountKeyUser)
	}

	configErr := config.SetAccountKeyToConfig("")

	if keyringErr != nil && configErr != nil {
		return configErr
	}
	return nil
}
