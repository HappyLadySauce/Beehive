package beehiveGateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/client"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/connection"
)

// GatewayServer Gateway 服务器
type GatewayServer struct {
	cfg        *config.Config
	grpcClient *client.Client
	connMgr    *connection.Manager
	httpServer *http.Server
}

// NewGatewayServer 创建 Gateway 服务器实例
func NewGatewayServer(cfg *config.Config, grpcClient *client.Client, connMgr *connection.Manager) *GatewayServer {
	return &GatewayServer{
		cfg:        cfg,
		grpcClient: grpcClient,
		connMgr:    connMgr,
	}
}

// PrepareRun 准备运行服务器
func (s *GatewayServer) PrepareRun() error {
	// 初始化路由和控制器
	klog.Info("Setting up routes and controllers...")
	handler := installControllers(s.cfg, s.grpcClient, s.connMgr)

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", s.cfg.InsecureServing.BindAddress, s.cfg.InsecureServing.BindPort)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return nil
}

// Run 运行服务器（阻塞直到收到停止信号）
func (s *GatewayServer) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	// 启动 HTTP 服务器
	go func() {
		klog.Infof("Starting Gateway server on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server failed: %w", err)
		}
	}()

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
		klog.Info("Gateway server stopped successfully")
		return nil
	case <-ctx.Done():
		klog.Info("Context cancelled, shutting down...")
		s.Shutdown()
		return nil
	}
}

// Shutdown 优雅关闭服务器
func (s *GatewayServer) Shutdown() {
	// 关闭 HTTP 服务器
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			klog.Errorf("HTTP server forced to shutdown: %v", err)
		}
	}

	// 关闭 gRPC 客户端
	if s.grpcClient != nil {
		if err := s.grpcClient.Close(); err != nil {
			klog.Errorf("Failed to close gRPC clients: %v", err)
		}
	}
}

