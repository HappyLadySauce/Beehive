package beehiveGateway

import (
	"context"
	"fmt"

	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/client"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/connection"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/options"
)

// run 运行 Gateway 服务器
func run(ctx context.Context, opts *options.Options) error {
	// 1. 初始化配置
	cfg, err := config.CreateConfigFromOptions(opts)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// 2. 初始化 gRPC 客户端
	klog.Info("Creating gRPC clients...")
	grpcClient, err := client.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize gRPC clients: %w", err)
	}

	// 3. 创建 Connection Manager
	klog.Info("Creating connection manager...")
	connMgr := connection.NewManager()

	// 4. 创建服务器实例
	server := NewGatewayServer(cfg, grpcClient, connMgr)

	// 5. 准备运行服务器
	if err := server.PrepareRun(); err != nil {
		return fmt.Errorf("failed to prepare server: %w", err)
	}

	// 6. 运行服务器（阻塞直到收到停止信号）
	return server.Run(ctx)
}
