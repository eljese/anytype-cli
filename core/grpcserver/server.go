//go:build !nogrpcserver
// +build !nogrpcserver

package grpcserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/anytype-heart/core/api"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/anyproto/anytype-cli/core/config"
	cli_core "github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-heart/core"
	"github.com/anyproto/anytype-heart/core/event"
	"github.com/anyproto/anytype-heart/metrics"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/util/grpcprocess"
)

var log = logging.Logger("anytype-heart")

const grpcWebStartedMessagePrefix = "gRPC Web proxy started at: "

type Server struct {
	mw           *core.Middleware
	grpcServer   *grpc.Server
	webServer    *http.Server
	proxyServer  *http.Server
	grpcListener net.Listener
	webListener  net.Listener
	proxyListener net.Listener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(grpcAddr, grpcWebAddr, proxyAddr, internalAPIAddr string) error {
	app.StartWarningAfter = time.Second * 5

	if os.Getenv("ANYTYPE_LOG_LEVEL") == "" {
		os.Setenv("ANYTYPE_LOG_LEVEL", "ERROR")
	}

	metrics.Service.InitWithKeys(metrics.DefaultInHouseKey)

	log.Info("Starting anytype-heart...")
	s.mw = core.New()
	s.mw.SetEventSender(event.NewGrpcSender())

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

	s.proxyListener, err = net.Listen("tcp", proxyAddr)
	if err != nil {
		s.grpcListener.Close()
		s.webListener.Close()
		return fmt.Errorf("failed to listen on %s: %w", proxyAddr, err)
	}

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

	if os.Getenv("ANYTYPE_GRPC_NO_DEBUG_TIMEOUT") != "1" {
		unaryInterceptors = append(unaryInterceptors, metrics.LongMethodsInterceptor)
	}

	unaryInterceptors = append(unaryInterceptors, grpcprocess.ProcessInfoInterceptor(
		"/anytype.ClientCommands/AccountLocalLinkNewChallenge",
	))

	s.grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(20*1024*1024),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
	)

	service.RegisterClientCommandsServer(s.grpcServer, s.mw)

	if metrics.Enabled {
		grpc_prometheus.EnableHandlingTimeHistogram()
	}

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

	// Setup Proxy Server
	targetUrl, _ := url.Parse("http://" + internalAPIAddr)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	proxyMux := http.NewServeMux()
	proxyMux.HandleFunc("/v1/file", s.handleFileUpload)
	proxyMux.Handle("/", proxy)

	s.proxyServer = &http.Server{
		Handler:           proxyMux,
		ReadHeaderTimeout: 30 * time.Second,
	}

	go func() {
		log.Infof("Starting gRPC server on %s", s.grpcListener.Addr())
		if err := s.grpcServer.Serve(s.grpcListener); err != nil {
			log.Errorf("gRPC server error: %v", err)
		}
	}()

	go func() {
		fmt.Printf("%s%s\n", grpcWebStartedMessagePrefix, s.webListener.Addr())
		if err := s.webServer.Serve(s.webListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("gRPC-Web server error: %v", err)
		}
	}()

	go func() {
		log.Infof("Starting API Proxy server on %s -> %s", s.proxyListener.Addr(), internalAPIAddr)
		if err := s.proxyServer.Serve(s.proxyListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("API Proxy server error: %v", err)
		}
	}()

	api.SetMiddlewareParams(s.mw)

	return nil
}

func (s *Server) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 100MB max
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	spaceId := r.FormValue("space_id")
	if spaceId == "" {
		http.Error(w, "space_id is required", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create temp file
	tempFile, err := os.CreateTemp("", "anytype-upload-*"+filepath.Ext(handler.Filename))
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Connect to gRPC
	// Use stored session token (sudo mode for file upload)
	token, _, err := cli_core.GetStoredSessionToken()
	if err != nil {
		http.Error(w, "Failed to get stored session token", http.StatusInternalServerError)
		return
	}

	conn, err := grpc.NewClient(config.DefaultGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := service.NewClientCommandsClient(conn)

	ctx := context.Background()
	if token != "" {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("token", token))
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := &pb.RpcFileUploadRequest{
		SpaceId:   spaceId,
		LocalPath: tempFile.Name(),
	}

	resp, err := client.FileUpload(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("gRPC error: %v", err), http.StatusInternalServerError)
		return
	}

	if resp.Error != nil && resp.Error.Code != 0 {
		tokenPreview := "empty"
		if len(token) > 5 {
			tokenPreview = token[:5] + "..."
		}
		http.Error(w, fmt.Sprintf("Upload error: %s (Token: %s, Space: %s)", resp.Error.Description, tokenPreview, spaceId), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"object_id": resp.ObjectId,
	})
}

func (s *Server) Stop() error {
	log.Info("Shutting down servers...")

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	if s.webServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.webServer.Shutdown(ctx); err != nil {
			log.Errorf("HTTP server shutdown error: %v", err)
		}
	}

	if s.proxyServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.proxyServer.Shutdown(ctx); err != nil {
			log.Errorf("Proxy server shutdown error: %v", err)
		}
	}

	if s.mw != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.mw.AppShutdown(ctx, &pb.RpcAppShutdownRequest{})
	}

	log.Info("Servers stopped")
	return nil
}
