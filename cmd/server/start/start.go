package start

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/daemon"
)

func NewStartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Anytype local server",
		RunE: func(cmd *cobra.Command, args []string) error {
			rootPath, _ := cmd.Flags().GetString("path")
			apiAddr, _ := cmd.Flags().GetString("api-addr")

			resp, err := daemon.SendTaskStart("server", nil)
			if err != nil {
				return fmt.Errorf("failed to start server: %w", err)
			}

			displayAddr := apiAddr
			if displayAddr == "" {
				displayAddr = config.DefaultAPIAddress
			}

			fmt.Printf("✓ Server started successfully via daemon. Response: %s, API Address: %s\n", resp.Status, displayAddr)
			time.Sleep(2 * time.Second)

			mnemonic, err := core.GetStoredMnemonic()
			if err == nil && mnemonic != "" {
				fmt.Println("ℹ Keychain mnemonic found. Attempting to login...")
				if err := core.LoginAccount(mnemonic, rootPath, apiAddr); err != nil {
					fmt.Println("✗ Failed to login using keychain mnemonic:", err)
				} else {
					fmt.Println("✓ Successfully logged in using keychain mnemonic.")
				}
			} else {
				fmt.Println("ℹ No keychain mnemonic found. Please login using 'anytype auth login' or create a new account with 'anytype auth create'.")
			}
			return nil
		},
	}

	startCmd.Flags().String("path", "", "Custom root path for wallet recovery")
	startCmd.Flags().String("api-addr", "", fmt.Sprintf("API listen address (default: %s)", config.DefaultAPIAddress))

	return startCmd
}
