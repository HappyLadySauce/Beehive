package common

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/HappyLadySauce/Beehive/internal/pkg/common/types/v1"
	"github.com/HappyLadySauce/Beehive/pkg/core"
)

// HandleHealth 处理健康检查请求
func (h *Handler) HandleHealth(c *gin.Context) {
	response := v1.HealthResponse{
		Status: "ok",
	}
	core.WriteResponse(c, nil, response)
}

// HandleReady 处理就绪检查请求
func (h *Handler) HandleReady(c *gin.Context) {
	// 可以在这里检查依赖服务的连接状态
	response := v1.HealthResponse{
		Status: "ready",
	}
	core.WriteResponse(c, nil, response)
}
