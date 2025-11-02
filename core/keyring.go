package core

import (
	"errors"

	"github.com/anyproto/anytype-cli/core/config"
	"github.com/zalando/go-keyring"
)

const (
	keyringService        = "anytype-cli"
	keyringTokenUser      = "session-token"
	keyringAccountKeyUser = "account-key"
)

var (
	keyringUnavailable = false
	ErrNotFound        = errors.New("credentials not found")
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

// SaveToken saves the session token to the keyring if available, otherwise to the config file.
// Returns true if saved to keyring, false if saved to config file.
func SaveToken(token string) (bool, error) {
	if isKeyringAvailable() {
		return true, keyring.Set(keyringService, keyringTokenUser, token)
	}

	return false, config.SetSessionTokenToConfig(token)
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
		return "", ErrNotFound
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

// SaveAccountKey saves the account key to the keyring if available, otherwise to the config file.
// Returns true if saved to keyring, false if saved to config file.
func SaveAccountKey(accountKey string) (bool, error) {
	if isKeyringAvailable() {
		return true, keyring.Set(keyringService, keyringAccountKeyUser, accountKey)
	}

	return false, config.SetAccountKeyToConfig(accountKey)
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
		return "", ErrNotFound
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
