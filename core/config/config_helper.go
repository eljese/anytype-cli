package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func GetAccountIdFromConfig() (string, error) {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	cfg := configMgr.Get()
	if cfg.AccountId == "" {
		return "", fmt.Errorf("no account Id found in config")
	}

	return cfg.AccountId, nil
}

func SetAccountIdToConfig(accountId string) error {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return configMgr.SetAccountId(accountId)
}

func GetTechSpaceIdFromConfig() (string, error) {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	cfg := configMgr.Get()
	if cfg.TechSpaceId == "" {
		return "", fmt.Errorf("no tech space Id found in config")
	}

	return cfg.TechSpaceId, nil
}

func SetTechSpaceIdToConfig(techSpaceId string) error {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return configMgr.SetTechSpaceId(techSpaceId)
}

func LoadStoredConfig() (*Config, error) {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return configMgr.Get(), nil
}

func GetSessionTokenFromConfig() (string, error) {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	cfg := configMgr.Get()
	if cfg.SessionToken == "" {
		return "", fmt.Errorf("no session token found in config")
	}

	return cfg.SessionToken, nil
}

func SetSessionTokenToConfig(token string) error {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return configMgr.SetSessionToken(token)
}

func GetAccountKeyFromConfig() (string, error) {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	cfg := configMgr.Get()
	if cfg.AccountKey == "" {
		return "", fmt.Errorf("no account key found in config")
	}

	return cfg.AccountKey, nil
}

func SetAccountKeyToConfig(accountKey string) error {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return configMgr.SetAccountKey(accountKey)
}

func GetNetworkConfigPathFromConfig() (string, error) {
        configMgr := GetConfigManager()
        if err := configMgr.Load(); err != nil {
                return "", fmt.Errorf("failed to load config: %w", err)
        }

        cfg := configMgr.Get()
        if cfg.NetworkConfigPath != "" {
                return cfg.NetworkConfigPath, nil
        }

        // Try to find client.yml in etc/ relative to current directory or executable
        paths := []string{
                "etc/client.yml",
                "../etc/client.yml",
                "/home/eljese/anytype/etc/client.yml",
        }

        for _, path := range paths {
                if _, err := os.Stat(path); err == nil {
                        return path, nil
                }
        }

        return "", nil
}
func SetNetworkConfigPathToConfig(path string) error {
	configMgr := GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return configMgr.SetNetworkConfigPath(path)
}

func GetNetworkIdFromConfigYAML(configPath string) (string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read network config: %w", err)
	}

	var cfg struct {
		NetworkId string `yaml:"networkId"`
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("failed to parse network config: %w", err)
	}

	if cfg.NetworkId == "" {
		return "", fmt.Errorf("networkId not found in config")
	}

	return cfg.NetworkId, nil
}
