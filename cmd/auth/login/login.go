package login

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewLoginCmd() *cobra.Command {
	var accountKey string
	var rootPath string
	var listenAddress string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to your bot account",
		Long:  "Authenticate using your account key to access your Anytype bot account and stored data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := core.LoginBot(accountKey, rootPath, listenAddress); err != nil {
				return output.Error("failed to log in: %w", err)
			}
			output.Success("Successfully logged in")
			return nil

		},
	}

	cmd.Flags().StringVar(&accountKey, "account-key", "", "Account key for authentication")
	cmd.Flags().StringVar(&rootPath, "path", "", "Root path for account data")
	cmd.Flags().StringVar(&listenAddress, "listen-address", config.DefaultAPIAddress, "API listen address in `host:port` format")

	return cmd
}
