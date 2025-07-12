package auth

import (
	"github.com/spf13/cobra"

	authApiKeyCmd "github.com/anyproto/anytype-cli/cmd/auth/apikey"
	authCreateCmd "github.com/anyproto/anytype-cli/cmd/auth/create"
	authLoginCmd "github.com/anyproto/anytype-cli/cmd/auth/login"
	authLogoutCmd "github.com/anyproto/anytype-cli/cmd/auth/logout"
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate with Anytype",
	}

	cmd.AddCommand(authLoginCmd.NewLoginCmd())
	cmd.AddCommand(authLogoutCmd.NewLogoutCmd())
	cmd.AddCommand(authCreateCmd.NewCreateCmd())
	cmd.AddCommand(authApiKeyCmd.NewApiKeyCmd())

	return cmd
}
