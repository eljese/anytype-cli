package apikey

import (
	"github.com/spf13/cobra"

	apiKeyCreateCmd "github.com/anyproto/anytype-cli/cmd/auth/apikey/create"
	apiKeyListCmd "github.com/anyproto/anytype-cli/cmd/auth/apikey/list"
	apiKeyRevokeCmd "github.com/anyproto/anytype-cli/cmd/auth/apikey/revoke"
)

// NewApiKeyCmd creates the auth apikey command
func NewApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apikey <command>",
		Short: "Manage API keys",
		Long:  "Create, list, and revoke API keys for programmatic access",
	}

	// Add subcommands
	cmd.AddCommand(apiKeyCreateCmd.NewCreateCmd())
	cmd.AddCommand(apiKeyListCmd.NewListCmd())
	cmd.AddCommand(apiKeyRevokeCmd.NewRevokeCmd())

	return cmd
}
