package core

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		gitState string
		want     string
	}{
		{
			name:     "with version",
			version:  "1.0.0",
			gitState: "",
			want:     "1.0.0",
		},
		{
			name:     "empty version returns dev",
			version:  "",
			gitState: "",
			want:     "dev",
		},
		{
			name:     "dirty git state",
			version:  "1.0.0",
			gitState: "dirty",
			want:     "1.0.0-dirty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalVersion := Version
			originalGitState := GitState
			defer func() {
				Version = originalVersion
				GitState = originalGitState
			}()

			Version = tt.version
			GitState = tt.gitState

			got := GetVersion()
			if got != tt.want {
				t.Errorf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetReleaseURL(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		want    string
	}{
		{
			name:    "release version",
			version: "1.0.0",
			commit:  "abc123",
			want:    "https://github.com/anyproto/anytype-cli/releases/tag/1.0.0",
		},
		{
			name:    "dev version with commit",
			version: "",
			commit:  "abc123",
			want:    "https://github.com/anyproto/anytype-cli/commit/abc123",
		},
		{
			name:    "pre-release version",
			version: "1.0.0-beta",
			commit:  "abc123",
			want:    "https://github.com/anyproto/anytype-cli/commit/abc123",
		},
		{
			name:    "no version or commit",
			version: "",
			commit:  "",
			want:    "https://github.com/anyproto/anytype-cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalVersion := Version
			originalCommit := Commit
			defer func() {
				Version = originalVersion
				Commit = originalCommit
			}()

			Version = tt.version
			Commit = tt.commit

			got := GetReleaseURL()
			if got != tt.want {
				t.Errorf("GetReleaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVersionBrief(t *testing.T) {
	originalVersion := Version
	originalBuildTime := BuildTime
	originalCommit := Commit
	defer func() {
		Version = originalVersion
		BuildTime = originalBuildTime
		Commit = originalCommit
	}()

	Version = "1.0.0"
	BuildTime = "2024-01-01"
	Commit = "abc123"

	got := GetVersionBrief()

	if !strings.Contains(got, "anytype-cli 1.0.0") {
		t.Errorf("GetVersionBrief() missing version: %v", got)
	}
	if !strings.Contains(got, "2024-01-01") {
		t.Errorf("GetVersionBrief() missing build time: %v", got)
	}
	if !strings.Contains(got, "github.com/anyproto/anytype-cli") {
		t.Errorf("GetVersionBrief() missing URL: %v", got)
	}
}

func TestGetVersionVerbose(t *testing.T) {
	originalVersion := Version
	originalCommit := Commit
	originalBuildTime := BuildTime
	defer func() {
		Version = originalVersion
		Commit = originalCommit
		BuildTime = originalBuildTime
	}()

	Version = "1.0.0"
	Commit = "abc123"
	BuildTime = "2024-01-01"

	got := GetVersionVerbose()

	expectedParts := []string{
		"anytype-cli 1.0.0",
		"Commit: abc123",
		"Built: 2024-01-01",
		"Go:",
		"OS/Arch:",
		"Heart:",
		"URL:",
	}

	for _, part := range expectedParts {
		if !strings.Contains(got, part) {
			t.Errorf("GetVersionVerbose() missing %q: %v", part, got)
		}
	}
}
