package service

import (
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
		DisplayName: "Anytype Server",
		Description: "Anytype personal knowledge management server",
		Arguments:   []string{"serve"},
		Option:      options,
	}

	prg := serviceprogram.New()
	return service.New(prg, svcConfig)
}

func NewServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage Anytype as a system service",
		Long:  "Install, uninstall, and control Anytype as a system service.",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "install",
			Short: "Install Anytype as a system service",
			RunE:  installService,
		},
		&cobra.Command{
			Use:   "uninstall",
			Short: "Uninstall the Anytype system service",
			RunE:  uninstallService,
		},
		&cobra.Command{
			Use:   "start",
			Short: "Start the Anytype service",
			RunE:  startService,
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the Anytype service",
			RunE:  stopService,
		},
		&cobra.Command{
			Use:   "restart",
			Short: "Restart the Anytype service",
			RunE:  restartService,
		},
		&cobra.Command{
			Use:   "status",
			Short: "Check the status of the Anytype service",
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

	output.Success("Anytype service installed successfully")
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

	output.Success("Anytype service uninstalled successfully")
	return nil
}

func startService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	err = s.Start()
	if err != nil {
		return output.Error("failed to start service: %w", err)
	}

	output.Success("Anytype service started")
	return nil
}

func stopService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	err = s.Stop()
	if err != nil {
		return output.Error("failed to stop service: %w", err)
	}

	output.Success("Anytype service stopped")
	return nil
}

func restartService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	err = s.Restart()
	if err != nil {
		return output.Error("failed to restart service: %w", err)
	}

	output.Success("Anytype service restarted")
	return nil
}

func statusService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	status, err := s.Status()
	if err != nil {
		if err == service.ErrNotInstalled {
			output.Info("Anytype service is not installed")
			output.Info("Run 'anytype service install' to install it")
			return nil
		}
		return output.Error("failed to get service status: %w", err)
	}

	switch status {
	case service.StatusRunning:
		output.Success("Anytype service is running")
	case service.StatusStopped:
		output.Info("Anytype service is stopped")
		output.Info("Run 'anytype service start' to start it")
	default:
		output.Info("Anytype service status: %v", status)
	}

	return nil
}
