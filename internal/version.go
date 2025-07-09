package internal

import (
	"fmt"
	"runtime"
)

// Set via ldflags during build
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
	GitState  = "unknown"
)

func GetVersionVerbose() string {
	return fmt.Sprintf("Anytype CLI %s\nCommit: %s\nBuilt: %s\nGo: %s\nOS/Arch: %s/%s\nURL: %s",
		GetVersion(), Commit, BuildTime, runtime.Version(), runtime.GOOS, runtime.GOARCH, GetReleaseURL())
}

func GetVersionBrief() string {
	return fmt.Sprintf("anytype-cli version %s (%s)\n%s", GetVersion(), BuildTime, GetReleaseURL())
}

func GetVersion() string {
	if GitState == "dirty" {
		return Version + "-dirty"
	}
	return Version
}

func GetReleaseURL() string {
	if GitState == "dirty" || Version == "v0.0.0" || Version == "dev" {
		return "https://github.com/anyproto/anytype-cli/commit/" + Commit
	}
	return "https://github.com/anyproto/anytype-cli/releases/tag/" + Version
}
