package beehiveAuth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/client"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/service"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/store"
	"github.com/HappyLadySauce/Beehive/internal/pkg/registry"
	pb "github.com/HappyLadySauce/Beehive/api/proto/auth/v1"
	"github.com/google/uuid"
)

// AuthServer Auth 微服务服务器
type AuthServer struct {
	cfg             *config.Config
	redisStore      *store.Store
	userClient      *client.Client
	authService     *service.Service
	grpcServer      *grpc.Server
	grpcListener    net.Listener
	httpServer      *http.Server
	serviceRegistry *registry.Registry
	instanceID      string
}

// NewAuthServer 创建 Auth 服务器实例
func NewAuthServer(cfg *config.Config, redisStore *store.Store, userClient *client.Client) *AuthServer {
	return &AuthServer{
		cfg:        cfg,
		redisStore: redisStore,
		userClient: userClient,
	}
}

// PrepareRun 准备运行服务器
func (s *AuthServer) PrepareRun() error {
	// 1. 创建 Auth Service 实例
	klog.Info("Creating Auth Service instance...")
	s.authService = service.NewService(s.cfg, s.redisStore, s.userClient)

	// 2. 创建 gRPC 服务器
	klog.Info("Creating gRPC server...")
	s.grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(s.cfg.Grpc.MaxMsgSize),
		grpc.MaxSendMsgSize(s.cfg.Grpc.MaxMsgSize),
	)

	// 3. 注册 Auth Service
	pb.RegisterAuthServiceServer(s.grpcServer, s.authService)

	// 4. 启动 gRPC 监听器
	addr := fmt.Sprintf("%s:%d", s.cfg.Grpc.BindAddress, s.cfg.Grpc.BindPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.grpcListener = lis
	klog.Infof("gRPC listener created on %s", addr)

	// 5. 注册服务到 etcd（如果配置了 etcd）
	if len(s.cfg.Etcd.Endpoints) > 0 && s.cfg.Etcd.Endpoints[0] != "" {
		serviceRegistry, err := registry.NewRegistry(
			s.cfg.Etcd.Endpoints,
			s.cfg.Etcd.DialTimeout,
			s.cfg.Etcd.Username,
			s.cfg.Etcd.Password,
			s.cfg.Etcd.Prefix,
		)
		if err != nil {
			klog.Warningf("Failed to create etcd registry, service registration disabled: %v", err)
		} else {
			s.serviceRegistry = serviceRegistry

			// 生成实例 ID
			s.instanceID = fmt.Sprintf("%s-%s", "auth", uuid.New().String()[:8])

			// 获取实际监听地址
			host := s.cfg.Grpc.BindAddress
			if host == "0.0.0.0" {
				host = "localhost" // 注册时使用 localhost，实际部署时应使用实际 IP
			}

			serviceInfo := &registry.ServiceInfo{
				ServiceName: "beehive-auth",
				Address:     host,
				Port:        s.cfg.Grpc.BindPort,
				InstanceID:  s.instanceID,
				Metadata:    make(map[string]string),
			}

			if err := serviceRegistry.Register(serviceInfo, 30); err != nil {
				klog.Warningf("Failed to register service to etcd: %v", err)
			} else {
				klog.Infof("Service registered to etcd: %s", s.instanceID)
			}
		}
	}

	// 6. 准备 HTTP 服务器
	if s.cfg.InsecureServing != nil && s.cfg.InsecureServing.BindPort != 0 {
		klog.Info("Preparing HTTP server...")
		handler := installRoutes(s.cfg)
		httpAddr := fmt.Sprintf("%s:%d", s.cfg.InsecureServing.BindAddress, s.cfg.InsecureServing.BindPort)
		s.httpServer = &http.Server{
			Addr:         httpAddr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}
	}

	return nil
}

// Run 运行服务器（阻塞直到收到停止信号）
func (s *AuthServer) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	// 启动 gRPC 服务器
	go func() {
		klog.Infof("Starting Auth Service gRPC server on %s", s.grpcListener.Addr())
		if err := s.grpcServer.Serve(s.grpcListener); err != nil {
			errChan <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	// 启动 HTTP 服务器
	if s.httpServer != nil {
		go func() {
			klog.Infof("Starting Auth HTTP server on %s", s.httpServer.Addr)
			if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- fmt.Errorf("HTTP server failed: %w", err)
			}
		}()
	}

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		klog.Errorf("Server error: %v", err)
		s.Shutdown()
		return err
	case sig := <-sigChan:
		klog.Infof("Received signal %v, shutting down gracefully...", sig)
		s.Shutdown()
		klog.Info("Auth Service stopped successfully")
		return nil
	case <-ctx.Done():
		klog.Info("Context cancelled, shutting down...")
		s.Shutdown()
		return nil
	}
}

// Shutdown 优雅关闭服务器
func (s *AuthServer) Shutdown() {
	// 关闭 gRPC 服务器
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	// 关闭 gRPC 监听器
	if s.grpcListener != nil {
		if err := s.grpcListener.Close(); err != nil {
			klog.Errorf("Failed to close gRPC listener: %v", err)
		}
	}

	// 关闭 HTTP 服务器
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			klog.Errorf("HTTP server forced to shutdown: %v", err)
		}
	}

	// 注销 etcd 服务
	if s.serviceRegistry != nil && s.instanceID != "" {
		if err := s.serviceRegistry.Deregister("beehive-auth", s.instanceID); err != nil {
			klog.Errorf("Failed to deregister service from etcd: %v", err)
		} else {
			klog.Info("Service deregistered from etcd")
		}
	}

	// 关闭 etcd 注册中心连接
	if s.serviceRegistry != nil {
		s.serviceRegistry.Close()
	}

	// 关闭 User Service gRPC 客户端
	if s.userClient != nil {
		if err := s.userClient.Close(); err != nil {
			klog.Errorf("Failed to close User Service client: %v", err)
		}
	}

	// 关闭 Redis 连接
	if s.redisStore != nil {
		if err := s.redisStore.Close(); err != nil {
			klog.Errorf("Failed to close Redis connection: %v", err)
		}
	}
}
