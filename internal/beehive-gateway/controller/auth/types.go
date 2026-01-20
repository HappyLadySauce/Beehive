package auth

import (
	authpb "github.com/HappyLadySauce/Beehive/pkg/api/proto/auth/v1"
)

// Handler 认证处理器
type Handler struct {
	authClient authpb.AuthServiceClient
}

// NewHandler 创建认证处理器
func NewHandler(authClient authpb.AuthServiceClient) *Handler {
	return &Handler{
		authClient: authClient,
	}
}
