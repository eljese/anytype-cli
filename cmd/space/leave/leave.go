package leave

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/cmdutil"
	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewLeaveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leave <space-id>",
		Short: "Leave a space",
		Long:  "Leave a space and stop sharing it",
		Args:  cmdutil.ExactArgs(1, "cannot leave space: space-id argument required"),
		RunE: func(cmd *cobra.Command, args []string) error {
			spaceId := args[0]

			if err := core.LeaveSpace(spaceId); err != nil {
				return output.Error("failed to leave space: %w", err)
			}

			output.Success("Successfully sent leave request for space with Id: %s", spaceId)
			return nil
		},
	}

	return cmd
}
