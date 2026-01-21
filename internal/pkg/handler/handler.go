package handler

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/HappyLadySauce/Beehive/internal/pkg/common/types/v1"
	"github.com/HappyLadySauce/Beehive/pkg/core"
)


// HandleHealthz 处理健康检查请求
// @Summary 健康检查
// @Description 检查服务是否健康
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} v1.HealthResponse
// @Failure 500 {object} core.ErrResponse
// @Router /healthz [get]
func HandleHealthz(c *gin.Context) {
	response := v1.HealthResponse{
		Status: "ok",
	}
	core.WriteResponse(c, nil, response)
}

// HandleReadyz 处理就绪检查请求
// @Summary 就绪检查
// @Description 检查服务是否就绪
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} v1.HealthResponse
// @Failure 500 {object} core.ErrResponse
// @Router /readyz [get]
func HandleReadyz(c *gin.Context) {
	response := v1.HealthResponse{
		Status: "ready",
	}

	core.WriteResponse(c, nil, response)
}
