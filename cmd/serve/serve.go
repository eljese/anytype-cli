package serve

import (
	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core/output"
	"github.com/anyproto/anytype-cli/core/serviceprogram"
)

func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"start"},
		Short:   "Run the Anytype server",
		Long:    "Run the Anytype server. When running as a system service, it will be managed by the service manager. When running interactively, use Ctrl+C to stop.",
		RunE:    runServer,
	}

	return cmd
}

func runServer(cmd *cobra.Command, args []string) error {
	svcConfig := &service.Config{
		Name:        "anytype",
		DisplayName: "Anytype Server",
		Description: "Anytype personal knowledge management server",
	}

	prg := serviceprogram.New()

	s, err := service.New(prg, svcConfig)
	if err != nil {
		return output.Error("failed to create service: %w", err)
	}

	err = s.Run()
	if err != nil {
		return output.Error("service failed: %w", err)
	}

	return nil
}
