package beehiveUser

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-user/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-user/options"
	"github.com/HappyLadySauce/Beehive/internal/beehive-user/service"
	"github.com/HappyLadySauce/Beehive/internal/beehive-user/store"
	pb "github.com/HappyLadySauce/Beehive/pkg/api/proto/user/v1"
)

func Run(ctx context.Context, opts *options.Options) error {
	// 1. 创建配置
	cfg, err := config.CreateConfigFromOptions(opts)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// 2. 创建数据库连接
	klog.Info("Creating database connection...")
	dbStore, err := store.NewStore(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		if err := dbStore.Close(); err != nil {
			klog.Errorf("Failed to close database connection: %v", err)
		}
	}()

	// 3. 创建 User Service 实例
	userService := service.NewService(dbStore)

	// 4. 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.Grpc.MaxMsgSize),
		grpc.MaxSendMsgSize(cfg.Grpc.MaxMsgSize),
	)

	// 5. 注册 User Service
	pb.RegisterUserServiceServer(grpcServer, userService)

	// 6. 启动 gRPC 服务器
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

	klog.Infof("Starting User Service gRPC server on %s", addr)

	// 7. 优雅关闭处理
	errChan := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	// 8. 等待中断信号
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
		klog.Info("User Service stopped successfully")
		return nil
	case <-ctx.Done():
		klog.Info("Context cancelled, shutting down...")
		grpcServer.GracefulStop()
		return nil
	}
}
