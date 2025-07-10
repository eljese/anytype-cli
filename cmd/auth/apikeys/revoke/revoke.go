package revoke

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/internal"
)

func NewRevokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke <id>",
		Short: "Revoke an API key",
		Long:  "Revoke an API key by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appId := args[0]

			err := internal.RevokeAPIKey(appId)
			if err != nil {
				return fmt.Errorf("✗ Failed to revoke API key: %w", err)
			}

			fmt.Printf("✓ API key with ID '%s' revoked successfully\n", appId)
			return nil
		},
	}

	return cmd
}
