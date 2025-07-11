package leave

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
)

func NewLeaveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leave <space-id>",
		Short: "Leave a space",
		Long:  "Leave a space and stop sharing it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spaceID := args[0]

			if err := core.LeaveSpace(spaceID); err != nil {
				return fmt.Errorf("failed to leave space: %w", err)
			}

			fmt.Printf("Successfully sent leave request for space with ID: %s\n", spaceID)
			return nil
		},
	}

	return cmd
}
