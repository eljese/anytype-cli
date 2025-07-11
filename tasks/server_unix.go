//go:build !windows
// +build !windows

package tasks

import (
	"os/exec"
	"syscall"
)

func setPlatformSpecificAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func terminateProcess(pid int) error {
	return syscall.Kill(-pid, syscall.SIGTERM)
}
