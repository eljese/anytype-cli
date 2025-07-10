package auth

import (
	"github.com/spf13/cobra"

	authApiKeyCmd "github.com/anyproto/anytype-cli/cmd/auth/apikey"
	authCreateCmd "github.com/anyproto/anytype-cli/cmd/auth/create"
	authLoginCmd "github.com/anyproto/anytype-cli/cmd/auth/login"
	authLogoutCmd "github.com/anyproto/anytype-cli/cmd/auth/logout"
)

func NewAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate with Anytype",
	}

	authCmd.AddCommand(authLoginCmd.NewLoginCmd())
	authCmd.AddCommand(authLogoutCmd.NewLogoutCmd())
	authCmd.AddCommand(authCreateCmd.NewCreateCmd())
	authCmd.AddCommand(authApiKeyCmd.NewApiKeyCmd())

	return authCmd
}
