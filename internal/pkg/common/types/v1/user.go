package v1

import authpb "github.com/HappyLadySauce/Beehive/api/proto/auth/v1"

// LoginRequest 登录请求
type LoginRequest struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string           `json:"token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresAt    int64            `json:"expires_at"`
	User         *authpb.UserInfo `json:"user"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Nickname    string `json:"nickname" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	ID string `json:"id"`
}
