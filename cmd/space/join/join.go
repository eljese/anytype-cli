package join

import (
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/cmdutil"
	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewJoinCmd() *cobra.Command {
	var (
		networkId     string
		inviteCid     string
		inviteFileKey string
	)

	cmd := &cobra.Command{
		Use:   "join <invite-link>",
		Short: "Join a space",
		Long:  "Join a space using an invite link (https://invite.any.coop/...)",
		Args:  cmdutil.ExactArgs(1, "cannot join space: invite-link argument required"),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]
			var spaceId string

			if networkId == "" {
				networkId = config.AnytypeNetworkAddress
			}

			if strings.HasPrefix(input, "https://invite.any.coop/") {
				u, err := url.Parse(input)
				if err != nil {
					return output.Error("invalid invite link: %w", err)
				}

				path := strings.TrimPrefix(u.Path, "/")
				if path == "" {
					return output.Error("invite link missing Cid")
				}
				inviteCid = path

				inviteFileKey = u.Fragment
				if inviteFileKey == "" {
					return output.Error("invite link missing key (should be after #)")
				}

				info, err := core.ViewSpaceInvite(inviteCid, inviteFileKey)
				if err != nil {
					return output.Error("failed to view invite: %w", err)
				}

				output.Info("Joining space '%s' created by %s...", info.SpaceName, info.CreatorName)
				spaceId = info.SpaceId
			} else {
				return output.Error("invalid invite link format, expected: https://invite.any.coop/{cid}#{key}")
			}

			if err := core.JoinSpace(networkId, spaceId, inviteCid, inviteFileKey); err != nil {
				return output.Error("failed to join space: %w", err)
			}

			output.Success("Successfully sent join request to space '%s'", spaceId)
			return nil
		},
	}

	cmd.Flags().StringVar(&networkId, "network", "", "Network `id` to join")
	cmd.Flags().StringVar(&inviteCid, "invite-cid", "", "Invite `cid` (extracted from invite link if not provided)")
	cmd.Flags().StringVar(&inviteFileKey, "invite-key", "", "Invite file `key` (extracted from invite link if not provided)")

	return cmd
}
