package status

import (
	"context"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			hasMnemonic := false
			if _, err := core.GetStoredMnemonic(); err == nil {
				hasMnemonic = true
			}

			hasToken := false
			token := ""
			if t, err := core.GetStoredToken(); err == nil {
				hasToken = true
				token = t
			}

			// Get account info from config
			configMgr := config.GetConfigManager()
			_ = configMgr.Load()
			cfg := configMgr.Get()
			accountID := cfg.AccountID

			// First check if server is running
			serverRunning := false
			err := core.GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
				_, err := client.AppGetVersion(ctx, &pb.RpcAppGetVersionRequest{})
				return err
			})
			serverRunning = err == nil

			// If server is running and we have a token, we're logged in
			// (server auto-logs in on restart using stored mnemonic)
			isLoggedIn := serverRunning && hasToken

			// Display status based on priority: server -> credentials -> login
			if !serverRunning {
				output.Print("Server is not running. Run 'anytype serve' to start the server.")
				if hasMnemonic || hasToken || accountID != "" {
					output.Print("Credentials are stored in keychain.")
				}
				return nil
			}

			if !hasMnemonic && !hasToken && accountID == "" {
				output.Print("Not authenticated. Run 'anytype auth login' to authenticate.")
				return nil
			}

			output.Print("Anytype")

			if isLoggedIn && accountID != "" {
				output.Print("  ✓ Logged in to account %s (keychain)", accountID)
			} else if hasToken || hasMnemonic {
				output.Print("  ✗ Not logged in (credentials stored in keychain)")
				if !isLoggedIn && hasToken {
					output.Print("    Note: Server is not running or session expired. Run 'anytype serve' to start server.")
				}
			} else {
				output.Print("  ✗ Not logged in")
			}

			output.Print("  - Active session: %v", isLoggedIn)

			if hasMnemonic {
				output.Print("  - Mnemonic: stored")
			}

			if hasToken {
				if len(token) > 8 {
					output.Print("  - Token: %s****", token[:8])
				} else {
					output.Print("  - Token: stored")
				}
			}

			return nil
		},
	}

	return cmd
}
