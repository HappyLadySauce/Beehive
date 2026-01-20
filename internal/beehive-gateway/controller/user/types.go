package user

import (
	userpb "github.com/HappyLadySauce/Beehive/pkg/api/proto/user/v1"
)

// Handler 用户处理器
type Handler struct {
	userClient userpb.UserServiceClient
}

// NewHandler 创建用户处理器
func NewHandler(userClient userpb.UserServiceClient) *Handler {
	return &Handler{
		userClient: userClient,
	}
}
