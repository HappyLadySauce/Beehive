package beehiveAuth

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/config"
	"github.com/HappyLadySauce/Beehive/internal/pkg/handler"
	"github.com/HappyLadySauce/Beehive/internal/pkg/middleware"

	_ "github.com/HappyLadySauce/Beehive/internal/beehive-auth/api/swagger/docs"
)

// installRoutes 构建 Auth 微服务的 HTTP 路由（健康检查 / Swagger）。
func installRoutes(cfg *config.Config) *gin.Engine {
	// 根据日志级别设置 Gin 模式
	if cfg != nil && cfg.Log != nil && cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 中间件顺序：Recovery -> RequestID -> CORS -> Logger
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Cors())
	router.Use(gin.Logger())

	// 健康检查路由（/healthz, /readyz）
	router.GET("/healthz", handler.HandleHealthz)
	router.GET("/readyz", handler.HandleReadyz)

	// Swagger 文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	klog.Info("Auth HTTP routes installed: /healthz, /readyz, /swagger")

	// 预留后续在 Auth 微服务暴露更多 HTTP 接口的挂载点。

	return router
}
