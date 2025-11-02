package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetDataDir(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		wantPath string
	}{
		{
			name:     "with DATA_PATH env",
			envValue: "/custom/data/path",
			wantPath: "/custom/data/path",
		},
		{
			name:     "without DATA_PATH env",
			envValue: "",
			wantPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEnv := os.Getenv("DATA_PATH")
			defer func() {
				os.Setenv("DATA_PATH", originalEnv)
			}()

			if tt.envValue != "" {
				os.Setenv("DATA_PATH", tt.envValue)
			} else {
				os.Unsetenv("DATA_PATH")
			}

			got := GetDataDir()

			if tt.wantPath != "" {
				if got != tt.wantPath {
					t.Errorf("GetDataDir() = %v, want %v", got, tt.wantPath)
				}
			} else {
				if got == "" {
					t.Error("GetDataDir() returned empty path")
				}

				if !strings.HasSuffix(got, "data") {
					t.Errorf("GetDataDir() = %v, expected to end with 'data'", got)
				}
			}
		})
	}
}

func TestGetWorkDir(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		goos     string
		expected string
	}{
		{
			name:     "macOS",
			goos:     "darwin",
			expected: filepath.Join(homeDir, "Library", "Application Support", "anytype"),
		},
		{
			name:     "Windows",
			goos:     "windows",
			expected: filepath.Join(homeDir, "AppData", "Roaming", "anytype"),
		},
		{
			name:     "Linux",
			goos:     "linux",
			expected: filepath.Join(homeDir, ".config", "anytype"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if runtime.GOOS != tt.goos {
				t.Skipf("Skipping test for %s on %s", tt.goos, runtime.GOOS)
			}

			got := GetWorkDir()
			if got != tt.expected {
				t.Errorf("GetWorkDir() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetConfigDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	expected := filepath.Join(homeDir, AnytypeDirName)
	got := GetConfigDir()

	if got != expected {
		t.Errorf("GetConfigDir() = %v, want %v", got, expected)
	}
}

func TestGetConfigFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	expected := filepath.Join(homeDir, AnytypeDirName, ConfigFileName)
	got := GetConfigFilePath()

	if got != expected {
		t.Errorf("GetConfigFilePath() = %v, want %v", got, expected)
	}
}

func TestGetLogsDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	expected := filepath.Join(homeDir, AnytypeDirName, LogsDirName)
	got := GetLogsDir()

	if got != expected {
		t.Errorf("GetLogsDir() = %v, want %v", got, expected)
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"LocalhostIP", LocalhostIP, "127.0.0.1"},
		{"GRPCPort", GRPCPort, "31007"},
		{"GRPCWebPort", GRPCWebPort, "31008"},
		{"APIPort", APIPort, "31009"},
		{"DefaultGRPCAddress", DefaultGRPCAddress, "127.0.0.1:31007"},
		{"DefaultGRPCWebAddress", DefaultGRPCWebAddress, "127.0.0.1:31008"},
		{"DefaultAPIAddress", DefaultAPIAddress, "127.0.0.1:31009"},
		{"GRPCDNSAddress", GRPCDNSAddress, "dns:///127.0.0.1:31007"},
		{"AnytypeDirName", AnytypeDirName, ".anytype"},
		{"ConfigFileName", ConfigFileName, "config.json"},
		{"DataDirName", DataDirName, "data"},
		{"LogsDirName", LogsDirName, "logs"},
		{"AnytypeName", AnytypeName, "anytype"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestGitHubURLs(t *testing.T) {
	baseURL := "https://github.com/anyproto/anytype-cli"

	if GitHubBaseURL != baseURL {
		t.Errorf("GitHubBaseURL = %v, want %v", GitHubBaseURL, baseURL)
	}

	expectedCommitURL := baseURL + "/commit/"
	if GitHubCommitURL != expectedCommitURL {
		t.Errorf("GitHubCommitURL = %v, want %v", GitHubCommitURL, expectedCommitURL)
	}

	expectedReleaseURL := baseURL + "/releases/tag/"
	if GitHubReleaseURL != expectedReleaseURL {
		t.Errorf("GitHubReleaseURL = %v, want %v", GitHubReleaseURL, expectedReleaseURL)
	}
}

func TestAnytypeNetworkAddress(t *testing.T) {
	expected := "N83gJpVd9MuNRZAuJLZ7LiMntTThhPc6DtzWWVjb1M3PouVU"
	if AnytypeNetworkAddress != expected {
		t.Errorf("AnytypeNetworkAddress = %v, want %v", AnytypeNetworkAddress, expected)
	}
}
