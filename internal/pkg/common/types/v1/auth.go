package v1

// RefreshTokenRequest Token 刷新请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse Token 刷新响应
type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// RevokeTokenRequest Token 撤销请求
type RevokeTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RevokeTokenResponse Token 撤销响应
type RevokeTokenResponse struct {
	Success bool `json:"success"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status string `json:"status"`
}
