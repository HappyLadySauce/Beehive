package common

import "github.com/gin-gonic/gin"

type IHandler interface {
	HandleHealth(c *gin.Context)
	HandleReady(c *gin.Context)
}

// Handler 公共处理器
type Handler struct{}

// NewHandler 创建公共处理器
func NewHandler() *Handler {
	return &Handler{}
}