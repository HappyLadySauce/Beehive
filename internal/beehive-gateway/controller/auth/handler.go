package auth

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/HappyLadySauce/Beehive/internal/pkg/common/types/v1"
	authpb "github.com/HappyLadySauce/Beehive/pkg/api/proto/auth/v1"
	"github.com/HappyLadySauce/Beehive/pkg/core"
)

// HandleLogin 处理登录请求
func (h *Handler) HandleLogin(c *gin.Context) {
	var req v1.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponseBindErr(c, err, nil)
		return
	}

	// 调用 Auth Service
	authReq := &authpb.LoginRequest{
		Id:       req.ID,
		Password: req.Password,
	}

	authResp, err := h.authClient.Login(c.Request.Context(), authReq)
	if err != nil {
		core.HandleGRPCError(c, err)
		return
	}

	response := v1.LoginResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    authResp.ExpiresAt,
		User:         authResp.User,
	}

	core.WriteResponse(c, nil, response)
}

// HandleRefreshToken 处理 Token 刷新请求
func (h *Handler) HandleRefreshToken(c *gin.Context) {
	var req v1.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponseBindErr(c, err, nil)
		return
	}

	// 调用 Auth Service
	authReq := &authpb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	authResp, err := h.authClient.RefreshToken(c.Request.Context(), authReq)
	if err != nil {
		core.HandleGRPCError(c, err)
		return
	}

	response := v1.RefreshTokenResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    authResp.ExpiresAt,
	}

	core.WriteResponse(c, nil, response)
}

// HandleRevokeToken 处理 Token 撤销请求
func (h *Handler) HandleRevokeToken(c *gin.Context) {
	var req v1.RevokeTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponseBindErr(c, err, nil)
		return
	}

	// 调用 Auth Service
	authReq := &authpb.RevokeTokenRequest{
		Token: req.Token,
	}

	authResp, err := h.authClient.RevokeToken(c.Request.Context(), authReq)
	if err != nil {
		core.HandleGRPCError(c, err)
		return
	}

	response := v1.RevokeTokenResponse{
		Success: authResp.Success,
	}

	core.WriteResponse(c, nil, response)
}
