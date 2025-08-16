package core

import (
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestClientContextWithAuth(t *testing.T) {
	testToken := "test-token-123"
	ctx := ClientContextWithAuth(testToken)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("Context should have metadata")
	}

	tokenHeaders := md.Get("token")
	if len(tokenHeaders) == 0 {
		t.Fatal("Context should have token header")
	}

	if tokenHeaders[0] != testToken {
		t.Errorf("Token header = %v, want %v", tokenHeaders[0], testToken)
	}
}

func TestClientContextWithAuthTimeout(t *testing.T) {
	testToken := "test-token-456"
	timeout := 5 * time.Second

	ctx, cancel := ClientContextWithAuthTimeout(testToken, timeout)
	defer cancel()

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("Context should have metadata")
	}

	tokenHeaders := md.Get("token")
	if len(tokenHeaders) == 0 {
		t.Fatal("Context should have token header")
	}

	if tokenHeaders[0] != testToken {
		t.Errorf("Token header = %v, want %v", tokenHeaders[0], testToken)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("Context should have deadline")
	}

	until := time.Until(deadline)
	if until > timeout || until < 0 {
		t.Errorf("Context deadline not set correctly: %v", until)
	}
}
