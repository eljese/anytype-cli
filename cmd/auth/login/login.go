package login

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/daemon"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to your Anytype vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			mnemonic, _ := cmd.Flags().GetString("mnemonic")
			rootPath, _ := cmd.Flags().GetString("path")
			apiAddr, _ := cmd.Flags().GetString("api-addr")

			statusResp, err := daemon.SendTaskStatus("server")
			if err != nil || statusResp.Status != "running" {
				return fmt.Errorf("server is not running")
			}

			if err := core.Login(mnemonic, rootPath, apiAddr); err != nil {
				return fmt.Errorf("✗ Failed to log in: %w", err)
			}
			fmt.Println("✓ Successfully logged in")
			return nil

		},
	}

	cmd.Flags().String("mnemonic", "", "Provide mnemonic (12 words) for authentication")
	cmd.Flags().String("path", "", "Provide custom root path for wallet recovery")
	cmd.Flags().String("api-addr", "", fmt.Sprintf("API listen address (default: %s)", config.DefaultAPIAddress))

	return cmd
}
