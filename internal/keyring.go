package internal

import "github.com/zalando/go-keyring"

const (
	keyringService      = "anytype-cli"
	keyringMnemonicUser = "mnemonic"
	keyringTokenUser    = "session-token"
)

func SaveMnemonic(mnemonic string) error {
	return keyring.Set(keyringService, keyringMnemonicUser, mnemonic)
}

func GetStoredMnemonic() (string, error) {
	return keyring.Get(keyringService, keyringMnemonicUser)
}

func DeleteStoredMnemonic() error {
	return keyring.Delete(keyringService, keyringMnemonicUser)
}

func SaveToken(token string) error {
	return keyring.Set(keyringService, keyringTokenUser, token)
}

func GetStoredToken() (string, error) {
	return keyring.Get(keyringService, keyringTokenUser)
}

func DeleteStoredToken() error {
	return keyring.Delete(keyringService, keyringTokenUser)
}
