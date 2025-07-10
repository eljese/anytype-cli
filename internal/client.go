package internal

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultTimeout = 5 * time.Second

var (
	clientInstance service.ClientCommandsClient
	grpcConn       *grpc.ClientConn
	once           sync.Once
	initErr        error
)

// GetGRPCClient initializes (if needed) and returns the shared gRPC client
func GetGRPCClient() (service.ClientCommandsClient, error) {
	// Ensure we only initialize once (singleton)
	once.Do(func() {
		var err error
		grpcConn, err = grpc.NewClient("dns:///127.0.0.1:31007", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			initErr = fmt.Errorf("failed to connect to gRPC server: %w", err)
			return
		}
		clientInstance = service.NewClientCommandsClient(grpcConn)
	})

	if initErr != nil {
		return nil, initErr
	}
	return clientInstance, nil
}

// CloseGRPCConnection ensures the connection is properly closed
func CloseGRPCConnection() {
	if grpcConn != nil {
		_ = grpcConn.Close()
	}
}

// IsGRPCServerRunning checks if the gRPC server is reachable
func IsGRPCServerRunning() (bool, error) {
	err := GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.AppGetVersion(ctx, &pb.RpcAppGetVersionRequest{})
		return err
	})

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func ClientContextWithAuth(token string) context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("token", token))
}

// ClientContextWithAuthTimeout creates a context with both authentication and timeout
func ClientContextWithAuthTimeout(token string, timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := ClientContextWithAuth(token)
	return context.WithTimeout(ctx, timeout)
}

// GRPCCall is a helper that reduces boilerplate for gRPC calls
// It gets the client, token, creates context with timeout, and executes the function
func GRPCCall(fn func(ctx context.Context, client service.ClientCommandsClient) error) error {
	client, err := GetGRPCClient()
	if err != nil {
		return fmt.Errorf("error connecting to gRPC server: %w", err)
	}

	token, err := GetStoredToken()
	if err != nil {
		return fmt.Errorf("failed to get stored token: %w", err)
	}

	ctx, cancel := ClientContextWithAuthTimeout(token, defaultTimeout)
	defer cancel()

	return fn(ctx, client)
}

// GRPCCallNoAuth is like GRPCCall but without authentication
func GRPCCallNoAuth(fn func(ctx context.Context, client service.ClientCommandsClient) error) error {
	client, err := GetGRPCClient()
	if err != nil {
		return fmt.Errorf("error connecting to gRPC server: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return fn(ctx, client)
}
