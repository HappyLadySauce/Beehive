package beehiveGateway

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/client"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/connection"
	authCtrl "github.com/HappyLadySauce/Beehive/internal/beehive-gateway/controller/auth"
	commonCtrl "github.com/HappyLadySauce/Beehive/internal/beehive-gateway/controller/common"
	userCtrl "github.com/HappyLadySauce/Beehive/internal/beehive-gateway/controller/user"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/websocket"
	"github.com/HappyLadySauce/Beehive/internal/pkg/middleware"
)

// installControllers 初始化并注册所有控制器
func installControllers(
	cfg *config.Config,
	grpcClient *client.Client,
	connMgr *connection.Manager,
) *gin.Engine {
	// 初始化控制器
	klog.Info("Initializing controllers...")
	authHandler := authCtrl.NewHandler(grpcClient.AuthService())
	userHandler := userCtrl.NewHandler(grpcClient.UserService())
	commonHandler := commonCtrl.NewHandler()

	// 创建 WebSocket Handler
	klog.Info("Setting up WebSocket handler...")
	wsHandler := websocket.NewHandler(cfg, grpcClient, connMgr)

	// 设置路由
	return setupRoutes(cfg, authHandler, userHandler, commonHandler, wsHandler.HandleConnection)
}

// setupRoutes 设置所有路由
func setupRoutes(
	cfg *config.Config,
	authHandler *authCtrl.Handler,
	userHandler *userCtrl.Handler,
	commonHandler *commonCtrl.Handler,
	wsHandler func(c *gin.Context),
) *gin.Engine {
	// 根据日志级别设置 Gin 模式
	if cfg != nil && cfg.Log != nil && cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 应用中间件（优化顺序：Recovery -> RequestID -> CORS -> Logger）
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.Cors())
	router.Use(gin.Logger())

	// WebSocket 路由
	if wsHandler != nil {
		router.GET("/ws", wsHandler)
	}

	// API v1 路由组
	apiV1 := router.Group("/api/v1")
	{
		// 认证相关路由
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.HandleLogin)
			authGroup.POST("/refresh", authHandler.HandleRefreshToken)
			authGroup.POST("/revoke", authHandler.HandleRevokeToken)
		}

		// 用户相关路由
		userGroup := apiV1.Group("/user")
		{
			userGroup.POST("/register", userHandler.HandleRegister)
		}
	}

	// 健康检查路由
	router.GET("/health", commonHandler.HandleHealth)
	router.GET("/ready", commonHandler.HandleReady)

	return router
}

