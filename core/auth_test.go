package core

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetDefaultDataPath(t *testing.T) {
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

			got := getDefaultDataPath()

			if tt.wantPath != "" {
				if got != tt.wantPath {
					t.Errorf("getDefaultDataPath() = %v, want %v", got, tt.wantPath)
				}
			} else {
				if got == "" {
					t.Error("getDefaultDataPath() returned empty path")
				}

				if !strings.HasSuffix(got, "data") {
					t.Errorf("getDefaultDataPath() = %v, expected to end with 'data'", got)
				}
			}
		})
	}
}

func TestGetDefaultWorkDir(t *testing.T) {
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

			got := getDefaultWorkDir()
			if got != tt.expected {
				t.Errorf("getDefaultWorkDir() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateAccountKey(t *testing.T) {
	tests := []struct {
		name        string
		accountKey  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid 64-byte account key (real example)",
			accountKey: "bNYSkBlOzNMKpDupAgL3g31Hnq7JpeX45O6MCpUqNdt16Avbgy5T5oQECKvAoy3+E4wHGPpCRCVWZQQCXRh7xw==",
			wantErr:    false,
		},
		{
			name:       "valid 32-byte account key (minimum)",
			accountKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			wantErr:    false,
		},
		{
			name:        "empty account key",
			accountKey:  "",
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "mnemonic instead of account key",
			accountKey:  "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12",
			wantErr:     true,
			errContains: "appears to be a mnemonic phrase",
		},
		{
			name:        "not valid base64",
			accountKey:  "this-is-not-valid-base64!!!",
			wantErr:     true,
			errContains: "must be valid base64",
		},
		{
			name:        "base64 but insufficient key material",
			accountKey:  "QUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQQ==", // 30 bytes decoded
			wantErr:     true,
			errContains: "insufficient key material",
		},
		{
			name:        "very short base64",
			accountKey:  "YWJj", // "abc" decoded (3 bytes)
			wantErr:     true,
			errContains: "insufficient key material",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccountKey(tt.accountKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccountKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateAccountKey() error = %q, want to contain %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}
