package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func AutoapproveTask(ctx context.Context, spaceID, role string) error {
	var permissions model.ParticipantPermissions
	switch role {
	case "Editor":
		permissions = model.ParticipantPermissions_Writer
	case "Viewer":
		fallthrough
	default:
		permissions = model.ParticipantPermissions_Reader
	}

	token, err := core.GetStoredToken()
	if err != nil || token == "" {
		return fmt.Errorf("failed to get stored token; are you logged in?")
	}

	er, err := core.ListenForEvents(token)
	if err != nil || er == nil {
		return fmt.Errorf("failed to start event listener: %w", err)
	}

	// Optionally, monitor the server status.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				status, err := core.IsGRPCServerRunning()
				if err != nil || !status {
					return
				}
			}
		}
	}()

	// Main loop: poll for join request events and approve them.
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			joinReq, err := core.WaitForJoinRequestEvent(er, spaceID)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if err := core.ApproveJoinRequest(joinReq.SpaceId, joinReq.Identity, permissions); err != nil {
				fmt.Printf("Failed to approve join request: %v\n", err)
			} else {
				fmt.Printf("Successfully approved join request for identity %s\n", joinReq.Identity)
			}
		}
	}
}
