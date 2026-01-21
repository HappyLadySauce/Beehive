package beehiveUser

import (
	"context"
	"fmt"

	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-user/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-user/options"
	"github.com/HappyLadySauce/Beehive/internal/beehive-user/store"
)

// Run 运行 User 服务器
func Run(ctx context.Context, opts *options.Options) error {
	// 1. 初始化配置
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

	// 3. 创建服务器实例
	server := NewUserServer(cfg, dbStore)

	// 4. 准备运行服务器
	if err := server.PrepareRun(); err != nil {
		// PrepareRun 失败时清理已创建的资源
		server.Shutdown()
		return fmt.Errorf("failed to prepare server: %w", err)
	}

	// 5. 运行服务器（阻塞直到收到停止信号）
	// Run 方法会在所有返回路径上调用 Shutdown()，确保资源被清理
	return server.Run(ctx)
}
