package space

import (
	"github.com/spf13/cobra"

	spaceAutoapproveCmd "github.com/anyproto/anytype-cli/cmd/space/autoapprove"
	spaceJoinCmd "github.com/anyproto/anytype-cli/cmd/space/join"
	spaceLeaveCmd "github.com/anyproto/anytype-cli/cmd/space/leave"
)

func NewSpaceCmd() *cobra.Command {
	spaceCmd := &cobra.Command{
		Use:   "space <command>",
		Short: "Manage spaces",
	}

	spaceCmd.AddCommand(spaceAutoapproveCmd.NewAutoapproveCmd())
	spaceCmd.AddCommand(spaceJoinCmd.NewJoinCmd())
	spaceCmd.AddCommand(spaceLeaveCmd.NewLeaveCmd())

	return spaceCmd
}
