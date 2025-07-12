package space

import (
	"github.com/spf13/cobra"

	spaceAutoapproveCmd "github.com/anyproto/anytype-cli/cmd/space/autoapprove"
	spaceJoinCmd "github.com/anyproto/anytype-cli/cmd/space/join"
	spaceLeaveCmd "github.com/anyproto/anytype-cli/cmd/space/leave"
	spaceListCmd "github.com/anyproto/anytype-cli/cmd/space/list"
)

func NewSpaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space <command>",
		Short: "Manage spaces",
	}

	cmd.AddCommand(spaceAutoapproveCmd.NewAutoapproveCmd())
	cmd.AddCommand(spaceJoinCmd.NewJoinCmd())
	cmd.AddCommand(spaceLeaveCmd.NewLeaveCmd())
	cmd.AddCommand(spaceListCmd.NewListCmd())

	return cmd
}
