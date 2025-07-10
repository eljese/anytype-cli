package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/cheggaaa/mb/v3"
)

type EventReceiver struct {
	queue  *mb.MB[*pb.EventMessage]
	stream service.ClientCommands_ListenSessionEventsClient
	cancel context.CancelFunc
	once   sync.Once
}

var (
	eventReceiverInstance *EventReceiver
	erOnce                sync.Once
	erInitErr             error
)

// ListenForEvents ensures a single EventReceiver instance is used.
func ListenForEvents(token string) (*EventReceiver, error) {
	erOnce.Do(func() {
		eventReceiverInstance, erInitErr = startListeningForEvents(token)
	})
	if erInitErr != nil {
		return nil, erInitErr
	}
	return eventReceiverInstance, nil
}

func startListeningForEvents(token string) (*EventReceiver, error) {
	client, err := GetGRPCClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get gRPC client: %w", err)
	}

	// Create authenticated context with cancel capability for the stream
	authCtx := ClientContextWithAuth(token)
	ctx, cancel := context.WithCancel(authCtx)

	stream, err := client.ListenSessionEvents(ctx, &pb.StreamRequest{
		Token: token,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start event stream: %w", err)
	}

	er := &EventReceiver{
		queue:  mb.New[*pb.EventMessage](0),
		stream: stream,
		cancel: cancel,
	}
	go er.receiveLoop(ctx)

	return er, nil
}

// receiveLoop continuously receives events from the stream
func (er *EventReceiver) receiveLoop(ctx context.Context) {
	defer er.cancel()

	for {
		event, err := er.stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("🔄 Event stream ended")
			return
		}
		if err != nil {
			if ctx.Err() != nil {
				// Context cancelled, clean shutdown
				return
			}
			fmt.Printf("✗ Event stream error: %v\n", err)
			return
		}

		for _, msg := range event.Messages {
			if err := er.queue.Add(ctx, msg); err != nil {
				if ctx.Err() != nil {
					return
				}
				fmt.Printf("✗ Failed to add event to queue: %v\n", err)
			}
		}
	}
}

// WaitForAccountID waits for an accountShow event and returns the account ID
func WaitForAccountID(er *EventReceiver) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a condition that filters for accountShow events
	cond := er.queue.NewCond().WithFilter(func(msg *pb.EventMessage) bool {
		return msg.GetAccountShow() != nil && msg.GetAccountShow().GetAccount() != nil
	})

	msg, err := cond.WaitOne(ctx)
	if err != nil {
		return "", fmt.Errorf("timeout waiting for account ID: %w", err)
	}

	return msg.GetAccountShow().GetAccount().Id, nil
}

// WaitForJoinRequestEvent waits for a join request for the specified space
func WaitForJoinRequestEvent(er *EventReceiver, spaceID string) (*model.NotificationRequestToJoin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create a condition that filters for join request events
	cond := er.queue.NewCond().WithFilter(func(msg *pb.EventMessage) bool {
		if ns := msg.GetNotificationSend(); ns != nil && ns.Notification != nil {
			if req := ns.Notification.GetRequestToJoin(); req != nil {
				return req.SpaceId == spaceID
			}
		}
		return false
	})

	msg, err := cond.WaitOne(ctx)
	if err != nil {
		return nil, fmt.Errorf("timeout waiting for join request: %w", err)
	}

	return msg.GetNotificationSend().Notification.GetRequestToJoin(), nil
}

// WaitOne waits for any single event with optional timeout
func (er *EventReceiver) WaitOne(ctx context.Context) (*pb.EventMessage, error) {
	return er.queue.WaitOne(ctx)
}

// WaitForEvent waits for an event matching the predicate
func (er *EventReceiver) WaitForEvent(ctx context.Context, predicate func(*pb.EventMessage) bool) (*pb.EventMessage, error) {
	cond := er.queue.NewCond().WithFilter(predicate)
	return cond.WaitOne(ctx)
}

// Close stops the event receiver
func (er *EventReceiver) Close() {
	er.once.Do(func() {
		er.cancel()
		_ = er.queue.Close()
	})
}

// CloseEventReceiver closes the global event receiver instance if it exists
func CloseEventReceiver() {
	if eventReceiverInstance != nil {
		eventReceiverInstance.Close()
	}
}
