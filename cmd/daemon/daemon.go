package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/daemon"
)

func NewDaemonCmd() *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Run the Anytype background daemon",
		Long:  "Run the Anytype daemon that manages background tasks (should be run as a system service).",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetString("addr")
			if err != nil {
				return err
			}
			fmt.Println("ℹ Starting daemon on", addr)
			return daemon.StartManager(addr)
		},
	}

	cmd.Flags().StringVar(&addr, "addr", config.DefaultDaemonAddress, "Address for the daemon to listen on")
	return cmd
}
