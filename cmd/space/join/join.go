package join

import (
	"fmt"
	"github.com/anyproto/anytype-cli/core/config"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
)

func NewJoinCmd() *cobra.Command {
	var (
		networkID     string
		inviteCID     string
		inviteFileKey string
	)

	cmd := &cobra.Command{
		Use:   "join <invite-link>",
		Short: "Join a space",
		Long:  "Join a space using an invite link (https://invite.any.coop/...)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]
			var spaceID string

			if networkID == "" {
				networkID = config.AnytypeNetworkAddress
			}

			if strings.HasPrefix(input, "https://invite.any.coop/") {
				u, err := url.Parse(input)
				if err != nil {
					return fmt.Errorf("invalid invite link: %w", err)
				}

				path := strings.TrimPrefix(u.Path, "/")
				if path == "" {
					return fmt.Errorf("invite link missing CID")
				}
				inviteCID = path

				inviteFileKey = u.Fragment
				if inviteFileKey == "" {
					return fmt.Errorf("invite link missing key (should be after #)")
				}

				info, err := core.ViewSpaceInvite(inviteCID, inviteFileKey)
				if err != nil {
					return fmt.Errorf("failed to view invite: %w", err)
				}

				fmt.Printf("Joining space '%s' created by %s...\n", info.SpaceName, info.CreatorName)
				spaceID = info.SpaceID
			} else {
				return fmt.Errorf("invalid invite link format, expected: https://invite.any.coop/{cid}#{key}")
			}

			if err := core.JoinSpace(networkID, spaceID, inviteCID, inviteFileKey); err != nil {
				return fmt.Errorf("failed to join space: %w", err)
			}

			fmt.Printf("Successfully sent join request to space '%s'\n", spaceID)
			return nil
		},
	}

	cmd.Flags().StringVar(&networkID, "network", "", "Network ID (optional, defaults to Anytype network address)")
	cmd.Flags().StringVar(&inviteCID, "invite-cid", "", "Invite CID (optional, extracted from invite link if provided)")
	cmd.Flags().StringVar(&inviteFileKey, "invite-key", "", "Invite file key (optional, extracted from invite link if provided)")

	return cmd
}
