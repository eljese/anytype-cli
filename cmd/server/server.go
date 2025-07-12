package server

import (
	"github.com/spf13/cobra"

	serverStartCmd "github.com/anyproto/anytype-cli/cmd/server/start"
	serverStatusCmd "github.com/anyproto/anytype-cli/cmd/server/status"
	serverStopCmd "github.com/anyproto/anytype-cli/cmd/server/stop"
)

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <command>",
		Short: "Manage the Anytype local server",
	}

	cmd.AddCommand(serverStartCmd.NewStartCmd())
	cmd.AddCommand(serverStopCmd.NewStopCmd())
	cmd.AddCommand(serverStatusCmd.NewStatusCmd())

	return cmd
}
