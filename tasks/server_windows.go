//go:build windows
// +build windows

package tasks

import (
	"os"
	"os/exec"
)

func setPlatformSpecificAttrs(cmd *exec.Cmd) {
	// Windows doesn't need Setpgid
}

func terminateProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(os.Interrupt)
}
