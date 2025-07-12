package config

import (
	"github.com/spf13/cobra"

	configGetCmd "github.com/anyproto/anytype-cli/cmd/config/get"
	configResetCmd "github.com/anyproto/anytype-cli/cmd/config/reset"
	configSetCmd "github.com/anyproto/anytype-cli/cmd/config/set"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage configuration",
		Long:  `Manage Anytype CLI configuration settings`,
	}

	cmd.AddCommand(configGetCmd.NewGetCmd())
	cmd.AddCommand(configSetCmd.NewSetCmd())
	cmd.AddCommand(configResetCmd.NewResetCmd())

	return cmd
}
