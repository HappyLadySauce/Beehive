package beehiveAuth

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/client"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/options"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/service"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/store"
	"github.com/HappyLadySauce/Beehive/internal/pkg/registry"
	pb "github.com/HappyLadySauce/Beehive/pkg/api/proto/auth/v1"
	"github.com/google/uuid"
)

func Run(ctx context.Context, opts *options.Options) error {
	// 1. 创建配置
	cfg, err := config.CreateConfigFromOptions(opts)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// 2. 创建 Redis 连接
	klog.Info("Creating Redis connection...")
	redisStore, err := store.NewStore(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}
	defer func() {
		if err := redisStore.Close(); err != nil {
			klog.Errorf("Failed to close Redis connection: %v", err)
		}
	}()

	// 3. 创建 User Service gRPC 客户端
	klog.Info("Creating User Service gRPC client...")
	userClient, err := client.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize User Service client: %w", err)
	}
	defer func() {
		if err := userClient.Close(); err != nil {
			klog.Errorf("Failed to close User Service client: %v", err)
		}
	}()

	// 4. 创建 Auth Service 实例
	authService := service.NewService(cfg, redisStore, userClient)

	// 5. 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.Grpc.MaxMsgSize),
		grpc.MaxSendMsgSize(cfg.Grpc.MaxMsgSize),
	)

	// 6. 注册 Auth Service
	pb.RegisterAuthServiceServer(grpcServer, authService)

	// 7. 启动 gRPC 服务器
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.BindAddress, cfg.Grpc.BindPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	defer func() {
		if err := lis.Close(); err != nil {
			klog.Errorf("Failed to close listener: %v", err)
		}
	}()

	klog.Infof("Starting Auth Service gRPC server on %s", addr)

	// 8. 注册服务到 etcd（如果配置了 etcd）
	var serviceRegistry *registry.Registry
	var instanceID string
	if len(cfg.Etcd.Endpoints) > 0 && cfg.Etcd.Endpoints[0] != "" {
		serviceRegistry, err = registry.NewRegistry(
			cfg.Etcd.Endpoints,
			cfg.Etcd.DialTimeout,
			cfg.Etcd.Username,
			cfg.Etcd.Password,
			cfg.Etcd.Prefix,
		)
		if err != nil {
			klog.Warningf("Failed to create etcd registry, service registration disabled: %v", err)
		} else {
			// Ensure registry is closed even if registration fails
			defer func() {
				if serviceRegistry != nil {
					serviceRegistry.Close()
				}
			}()

			// 生成实例 ID
			instanceID = fmt.Sprintf("%s-%s", "auth", uuid.New().String()[:8])

			// 获取实际监听地址
			host := cfg.Grpc.BindAddress
			if host == "0.0.0.0" {
				host = "localhost" // 注册时使用 localhost，实际部署时应使用实际 IP
			}

			serviceInfo := &registry.ServiceInfo{
				ServiceName: "beehive-auth",
				Address:     host,
				Port:        cfg.Grpc.BindPort,
				InstanceID:  instanceID,
				Metadata:    make(map[string]string),
			}

			if err := serviceRegistry.Register(serviceInfo, 30); err != nil {
				klog.Warningf("Failed to register service to etcd: %v", err)
			} else {
				klog.Infof("Service registered to etcd: %s", instanceID)
				// Set up defer for deregistration only if registration succeeded
				defer func() {
					if err := serviceRegistry.Deregister("beehive-auth", instanceID); err != nil {
						klog.Errorf("Failed to deregister service from etcd: %v", err)
					} else {
						klog.Info("Service deregistered from etcd")
					}
				}()
			}
		}
	}

	// 9. 优雅关闭处理
	errChan := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	// 10. 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		klog.Errorf("gRPC server error: %v", err)
		grpcServer.GracefulStop()
		return err
	case sig := <-sigChan:
		klog.Infof("Received signal %v, shutting down gracefully...", sig)
		grpcServer.GracefulStop()
		klog.Info("Auth Service stopped successfully")
		return nil
	case <-ctx.Done():
		klog.Info("Context cancelled, shutting down...")
		grpcServer.GracefulStop()
		return nil
	}
}
