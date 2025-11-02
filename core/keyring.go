package core

import (
	"github.com/zalando/go-keyring"
)

const (
	keyringService        = "anytype-cli"
	keyringTokenUser      = "session-token"
	keyringBotAccountUser = "account-key"
)

func SaveToken(token string) error {
	return keyring.Set(keyringService, keyringTokenUser, token)
}

func GetStoredToken() (string, error) {
	return keyring.Get(keyringService, keyringTokenUser)
}

func DeleteStoredToken() error {
	return keyring.Delete(keyringService, keyringTokenUser)
}

func SaveAccountKey(accountKey string) error {
	return keyring.Set(keyringService, keyringBotAccountUser, accountKey)
}

func GetStoredAccountKey() (string, error) {
	return keyring.Get(keyringService, keyringBotAccountUser)
}

func DeleteStoredAccountKey() error {
	return keyring.Delete(keyringService, keyringBotAccountUser)
}
