//go:build !nogrpcserver
// +build !nogrpcserver

package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/anytype-heart/core/api"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"github.com/anyproto/anytype-heart/core"
	"github.com/anyproto/anytype-heart/core/event"
	"github.com/anyproto/anytype-heart/metrics"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/util/grpcprocess"
)

var log = logging.Logger("anytype-embedded-server")

type Server struct {
	mw           *core.Middleware
	grpcServer   *grpc.Server
	webServer    *http.Server
	grpcListener net.Listener
	webListener  net.Listener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(grpcAddr, grpcWebAddr string) error {
	// Initialize the app start warning
	app.StartWarningAfter = time.Second * 5

	// Set up logging to stdout
	if os.Getenv("ANYTYPE_LOG_LEVEL") == "" {
		os.Setenv("ANYTYPE_LOG_LEVEL", "ERROR")
	}

	// Initialize metrics
	metrics.Service.InitWithKeys(metrics.DefaultInHouseKey)

	log.Info("Starting Anytype server...")

	// Create middleware
	s.mw = core.New()
	s.mw.SetEventSender(event.NewGrpcSender())

	// Create listeners
	var err error
	s.grpcListener, err = net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcAddr, err)
	}

	s.webListener, err = net.Listen("tcp", grpcWebAddr)
	if err != nil {
		s.grpcListener.Close()
		return fmt.Errorf("failed to listen on %s: %w", grpcWebAddr, err)
	}

	// Setup interceptors
	var unaryInterceptors []grpc.UnaryServerInterceptor

	if metrics.Enabled {
		unaryInterceptors = append(unaryInterceptors, grpc_prometheus.UnaryServerInterceptor)
	}

	unaryInterceptors = append(unaryInterceptors, metrics.UnaryTraceInterceptor)
	unaryInterceptors = append(unaryInterceptors, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = s.mw.Authorize(ctx, req, info, handler)
		if err != nil {
			log.Errorf("authorize: %s", err)
		}
		return
	})

	// Add debug timeout interceptor if not disabled
	if os.Getenv("ANYTYPE_GRPC_NO_DEBUG_TIMEOUT") != "1" {
		unaryInterceptors = append(unaryInterceptors, metrics.LongMethodsInterceptor)
	}

	// Add process info interceptor
	unaryInterceptors = append(unaryInterceptors, grpcprocess.ProcessInfoInterceptor(
		"/anytype.ClientCommands/AccountLocalLinkNewChallenge",
	))

	// Create gRPC server
	s.grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(20*1024*1024),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
	)

	// Register service
	service.RegisterClientCommandsServer(s.grpcServer, s.mw)

	if metrics.Enabled {
		grpc_prometheus.EnableHandlingTimeHistogram()
	}

	// Create gRPC-Web proxy
	webrpc := grpcweb.WrapServer(
		s.grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool { return true }),
	)

	s.webServer = &http.Server{
		Handler:           webrpc,
		ReadHeaderTimeout: 30 * time.Second,
	}

	// Start servers in goroutines
	go func() {
		log.Infof("Starting gRPC server on %s", s.grpcListener.Addr())
		if err := s.grpcServer.Serve(s.grpcListener); err != nil {
			log.Errorf("gRPC server error: %v", err)
		}
	}()

	go func() {
		// Print the required message for JS client compatibility
		fmt.Printf("gRPC Web proxy started at: %s\n", s.webListener.Addr())
		if err := s.webServer.Serve(s.webListener); err != nil && err != http.ErrServerClosed {
			log.Errorf("gRPC-Web server error: %v", err)
		}
	}()

	api.SetMiddlewareParams(s.mw)

	// Give servers a moment to start
	time.Sleep(100 * time.Millisecond)

	return nil
}

func (s *Server) Stop() error {
	log.Info("Shutting down servers...")

	// Gracefully stop gRPC server
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	// Shutdown HTTP server
	if s.webServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.webServer.Shutdown(ctx); err != nil {
			log.Errorf("HTTP server shutdown error: %v", err)
		}
	}

	// Stop middleware - using AppShutdown as per anytype-heart implementation
	if s.mw != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.mw.AppShutdown(ctx, &pb.RpcAppShutdownRequest{})
	}

	log.Info("Servers stopped")
	return nil
}
