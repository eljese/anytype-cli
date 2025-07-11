package logout

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
)

func NewLogoutCmd() *cobra.Command {
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out and remove stored credentials from keychain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := core.Logout(); err != nil {
				return fmt.Errorf("✗ Failed to log out: %w", err)
			}
			fmt.Println("✓ Successfully logged out. Stored credentials removed.")
			return nil
		},
	}

	return logoutCmd
}
