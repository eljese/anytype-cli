package logout

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out and clear stored credentials",
		Long:  "End your current session and remove stored authentication tokens and account key from the system keychain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := core.Logout(); err != nil {
				return output.Error("failed to log out: %w", err)
			}
			output.Success("Successfully logged out. Stored credentials removed.")
			return nil
		},
	}

	return cmd
}
