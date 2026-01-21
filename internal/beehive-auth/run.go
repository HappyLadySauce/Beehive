package beehiveAuth

import (
	"context"
	"fmt"

	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/client"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/options"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/store"
)

// Run 运行 Auth 服务器
func Run(ctx context.Context, opts *options.Options) error {
	// 1. 初始化配置
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

	// 3. 创建 User Service gRPC 客户端
	klog.Info("Creating User Service gRPC client...")
	userClient, err := client.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize User Service client: %w", err)
	}

	// 4. 创建服务器实例
	server := NewAuthServer(cfg, redisStore, userClient)

	// 5. 准备运行服务器
	if err := server.PrepareRun(); err != nil {
		// PrepareRun 失败时清理已创建的资源
		server.Shutdown()
		return fmt.Errorf("failed to prepare server: %w", err)
	}

	// 6. 运行服务器（阻塞直到收到停止信号）
	// Run 方法会在所有返回路径上调用 Shutdown()，确保资源被清理
	return server.Run(ctx)
}
