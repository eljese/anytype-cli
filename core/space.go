package core

import (
	"context"
	"fmt"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func ApproveJoinRequest(spaceID, identity string, permissions model.ParticipantPermissions) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceRequestApproveRequest{
			SpaceId:     spaceID,
			Identity:    identity,
			Permissions: permissions,
		}
		_, err := client.SpaceRequestApprove(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to approve join request: %w", err)
		}
		return nil
	})
}
