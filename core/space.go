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

func JoinSpace(networkID, spaceID, inviteCID, inviteFileKey string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceJoinRequest{
			NetworkId:     networkID,
			SpaceId:       spaceID,
			InviteCid:     inviteCID,
			InviteFileKey: inviteFileKey,
		}
		_, err := client.SpaceJoin(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to join space: %w", err)
		}
		return nil
	})
}

func LeaveSpace(spaceID string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceDeleteRequest{
			SpaceId: spaceID,
		}
		_, err := client.SpaceDelete(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to leave space: %w", err)
		}
		return nil
	})
}

type SpaceInviteInfo struct {
	SpaceID           string
	SpaceName         string
	SpaceIconCID      string
	CreatorName       string
	IsGuestUserInvite bool
	InviteType        model.InviteType
}

func ViewSpaceInvite(inviteCID, inviteFileKey string) (*SpaceInviteInfo, error) {
	var info *SpaceInviteInfo
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceInviteViewRequest{
			InviteCid:     inviteCID,
			InviteFileKey: inviteFileKey,
		}
		resp, err := client.SpaceInviteView(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to view space invite: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcSpaceInviteViewResponseError_NULL {
			return fmt.Errorf("space invite view error: %s", resp.Error.Description)
		}
		info = &SpaceInviteInfo{
			SpaceID:           resp.SpaceId,
			SpaceName:         resp.SpaceName,
			SpaceIconCID:      resp.SpaceIconCid,
			CreatorName:       resp.CreatorName,
			IsGuestUserInvite: resp.IsGuestUserInvite,
			InviteType:        resp.InviteType,
		}
		return nil
	})
	return info, err
}
