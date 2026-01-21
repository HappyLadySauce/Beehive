package beehiveUser

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	commonCtrl "github.com/HappyLadySauce/Beehive/internal/beehive-gateway/controller/common"
	"github.com/HappyLadySauce/Beehive/internal/beehive-user/config"
	"github.com/HappyLadySauce/Beehive/internal/pkg/middleware"
)

// installRoutes 构建 User 微服务的 HTTP 路由（健康检查 / Swagger）。
func installRoutes(cfg *config.Config) *gin.Engine {
	// 根据日志级别设置 Gin 模式
	if cfg != nil && cfg.Log != nil && cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 中间件顺序：Recovery -> RequestID -> CORS -> Swagger -> Logger
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Cors())
	router.Use(middleware.Swagger())
	router.Use(gin.Logger())

	// 健康检查路由（/healthz, /readyz）
	commonHandler := commonCtrl.NewHandler()
	router.GET("/healthz", commonHandler.HandleHealth)
	router.GET("/readyz", commonHandler.HandleReady)

	klog.Info("User HTTP routes installed: /healthz, /readyz, /swagger")

	// 预留后续在 User 微服务暴露更多 HTTP 接口的挂载点。

	return router
}
