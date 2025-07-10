package tasks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/anyproto/anytype-cli/internal/config"
)

// ServerTask is a background task that starts the server process.
// It spawns the server and waits until the given context is canceled.
func ServerTask(ctx context.Context) error {
	grpcPort := config.GRPCPort
	grpcWebPort := config.GRPCWebPort

	cmd := exec.Command("./dist/grpc-server")
	cmd.Env = append(os.Environ(),
		config.EnvGRPCAddr+"="+config.DefaultBindAddress+":"+grpcPort,
		config.EnvGRPCWebAddr+"="+config.DefaultBindAddress+":"+grpcWebPort,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		return <-done
	case err := <-done:
		return err
	}
}
