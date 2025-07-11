package create

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new API key",
		Long:  "Create a new API key for programmatic access to Anytype",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			resp, err := core.CreateAPIKey(name)
			if err != nil {
				return fmt.Errorf("✗ Failed to create API key: %w", err)
			}

			fmt.Println("✓ API key created successfully")
			fmt.Println("ℹ Name:", name)
			fmt.Println("ℹ Key:", resp.AppKey)

			return nil
		},
	}

	return cmd
}
