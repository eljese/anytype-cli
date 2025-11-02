package core

import (
	"context"
	"fmt"

	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

func JoinSpace(networkId, spaceId, inviteCId, inviteFileKey string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceJoinRequest{
			NetworkId:     networkId,
			SpaceId:       spaceId,
			InviteCid:     inviteCId,
			InviteFileKey: inviteFileKey,
		}
		_, err := client.SpaceJoin(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to join space: %w", err)
		}
		return nil
	})
}

func LeaveSpace(spaceId string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceDeleteRequest{
			SpaceId: spaceId,
		}
		_, err := client.SpaceDelete(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to leave space: %w", err)
		}
		return nil
	})
}

func ViewSpaceInvite(inviteCId, inviteFileKey string) (*pb.RpcSpaceInviteViewResponse, error) {
	var resp *pb.RpcSpaceInviteViewResponse
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceInviteViewRequest{
			InviteCid:     inviteCId,
			InviteFileKey: inviteFileKey,
		}
		var err error
		resp, err = client.SpaceInviteView(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to view space invite: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcSpaceInviteViewResponseError_NULL {
			return fmt.Errorf("space invite view error: %s", resp.Error.Description)
		}
		return nil
	})
	return resp, err
}

type SpaceListItem struct {
	SpaceId string
	Name    string
	Status  model.SpaceStatus
}

// ListSpaces returns a list of all available spaces
func ListSpaces() ([]SpaceListItem, error) {
	techSpaceId, err := config.GetTechSpaceIdFromConfig()
	if err != nil {
		return nil, fmt.Errorf("tech space Id not found in config - please login first: %w", err)
	}

	var spaces []SpaceListItem
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcObjectSearchRequest{
			SpaceId: techSpaceId,
			Filters: []*model.BlockContentDataviewFilter{
				{
					RelationKey: bundle.RelationKeyResolvedLayout.String(),
					Condition:   model.BlockContentDataviewFilter_Equal,
					Value:       pbtypes.Int64(int64(model.ObjectType_spaceView)),
				},
				{
					RelationKey: bundle.RelationKeySpaceLocalStatus.String(),
					Condition:   model.BlockContentDataviewFilter_In,
					Value:       pbtypes.IntList(int(model.SpaceStatus_Unknown), int(model.SpaceStatus_Ok)),
				},
				{
					RelationKey: bundle.RelationKeySpaceAccountStatus.String(),
					Condition:   model.BlockContentDataviewFilter_In,
					Value:       pbtypes.IntList(int(model.SpaceStatus_Unknown), int(model.SpaceStatus_SpaceActive)),
				},
			},
			Sorts: []*model.BlockContentDataviewSort{
				{
					RelationKey:    bundle.RelationKeySpaceOrder.String(),
					Type:           model.BlockContentDataviewSort_Asc,
					NoCollate:      true,
					EmptyPlacement: model.BlockContentDataviewSort_End,
				},
			},
			Keys: []string{
				bundle.RelationKeyTargetSpaceId.String(),
				bundle.RelationKeyName.String(),
				bundle.RelationKeySpaceLocalStatus.String(),
			},
		}

		resp, err := client.ObjectSearch(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to search spaces: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectSearchResponseError_NULL {
			return fmt.Errorf("object search error: %s", resp.Error.Description)
		}

		for _, record := range resp.Records {
			item := SpaceListItem{}

			// Get space Id
			if spaceIdVal := pbtypes.GetString(record, bundle.RelationKeyTargetSpaceId.String()); spaceIdVal != "" {
				item.SpaceId = spaceIdVal
			}

			// Get name
			if nameVal := pbtypes.GetString(record, bundle.RelationKeyName.String()); nameVal != "" {
				item.Name = nameVal
			}

			// Get status
			if statusVal := pbtypes.GetInt64(record, bundle.RelationKeySpaceLocalStatus.String()); statusVal != 0 {
				item.Status = model.SpaceStatus(statusVal)
			}

			if item.SpaceId != "" {
				spaces = append(spaces, item)
			}
		}

		return nil
	})

	return spaces, err
}
