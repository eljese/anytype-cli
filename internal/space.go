package internal

import (
	"fmt"
	"time"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func ApproveJoinRequest(token, spaceID, identity string, permissions model.ParticipantPermissions) error {
	client, err := GetGRPCClient()
	if err != nil {
		return err
	}
	ctx, cancel := ClientContextWithAuthTimeout(token, 5*time.Second)
	defer cancel()

	req := &pb.RpcSpaceRequestApproveRequest{
		SpaceId:     spaceID,
		Identity:    identity,
		Permissions: permissions,
	}
	_, err = client.SpaceRequestApprove(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to approve join request: %w", err)
	}
	return nil
}
