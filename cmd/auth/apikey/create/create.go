package create

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/cmdutil"
	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new API key",
		Long:  "Create a new API key for programmatic access to Anytype",
		Args:  cmdutil.ExactArgs(1, "cannot create API key: name argument required"),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			resp, err := core.CreateAPIKey(name)
			if err != nil {
				return output.Error("failed to create API key: %w", err)
			}

			output.Success("API key created successfully")
			output.Info("Name: %s", name)
			output.Info("Key: %s", resp.AppKey)

			return nil
		},
	}

	return cmd
}
