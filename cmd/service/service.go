package service

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core/output"
	"github.com/anyproto/anytype-cli/core/serviceprogram"
)

// getService creates a service instance with our standard configuration
func getService() (service.Service, error) {
	options := service.KeyValue{
		"UserService": true,
	}

	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		logDir := filepath.Join(homeDir, ".anytype", "logs")
		if err := os.MkdirAll(logDir, 0755); err == nil {
			options["LogDirectory"] = logDir
		}
	}

	svcConfig := &service.Config{
		Name:        "anytype",
		DisplayName: "Anytype",
		Description: "Anytype",
		Arguments:   []string{"serve"},
		Option:      options,
	}

	prg := serviceprogram.New()
	return service.New(prg, svcConfig)
}

func NewServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage anytype as a user service",
		Long:  "Install, uninstall, start, stop, and check status of anytype running as a user service.",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "install",
			Short: "Install as a user service",
			RunE:  installService,
		},
		&cobra.Command{
			Use:   "uninstall",
			Short: "Uninstall the user service",
			RunE:  uninstallService,
		},
		&cobra.Command{
			Use:   "start",
			Short: "Start the service",
			RunE:  startService,
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the service",
			RunE:  stopService,
		},
		&cobra.Command{
			Use:   "restart",
			Short: "Restart the service",
			RunE:  restartService,
		},
		&cobra.Command{
			Use:   "status",
			Short: "Check service status",
			RunE:  statusService,
		},
	)

	return cmd
}

func installService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	err = s.Install()
	if err != nil {
		return output.Error("failed to install service: %w", err)
	}

	output.Success("anytype service installed successfully")
	output.Print("\nTo manage the service:")
	output.Print("  Start:   anytype service start")
	output.Print("  Stop:    anytype service stop")
	output.Print("  Restart: anytype service restart")
	output.Print("  Status:  anytype service status")

	return nil
}

func uninstallService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	err = s.Uninstall()
	if err != nil {
		return output.Error("failed to uninstall service: %w", err)
	}

	output.Success("anytype service uninstalled successfully")
	return nil
}

func startService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	// Check if service is installed first
	_, err = s.Status()
	if err != nil && errors.Is(err, service.ErrNotInstalled) {
		output.Warning("anytype service is not installed")
		output.Info("Run 'anytype service install' to install it first")
		return nil
	}

	err = s.Start()
	if err != nil {
		return output.Error("failed to start service: %w", err)
	}

	output.Success("anytype service started")
	return nil
}

func stopService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	// Check if service is installed first
	_, err = s.Status()
	if err != nil && errors.Is(err, service.ErrNotInstalled) {
		output.Warning("anytype service is not installed")
		output.Info("Run 'anytype service install' to install it first")
		return nil
	}

	err = s.Stop()
	if err != nil {
		return output.Error("failed to stop service: %w", err)
	}

	output.Success("anytype service stopped")
	return nil
}

func restartService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	// Check if service is installed first
	_, err = s.Status()
	if err != nil && errors.Is(err, service.ErrNotInstalled) {
		output.Warning("anytype service is not installed")
		output.Info("Run 'anytype service install' to install it first")
		return nil
	}

	err = s.Restart()
	if err != nil {
		return output.Error("failed to restart service: %w", err)
	}

	output.Success("anytype service restarted")
	return nil
}

func statusService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	status, err := s.Status()
	if err != nil {
		if errors.Is(err, service.ErrNotInstalled) {
			output.Info("anytype service is not installed")
			output.Info("Run 'anytype service install' to install it")
			return nil
		}
		return output.Error("failed to get service status: %w", err)
	}

	switch status {
	case service.StatusRunning:
		output.Success("anytype service is running")
	case service.StatusStopped:
		output.Info("anytype service is stopped")
		output.Info("Run 'anytype service start' to start it")
	default:
		output.Info("anytype service status: %v", status)
	}

	return nil
}
