package auth

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/HappyLadySauce/Beehive/internal/pkg/common/types/v1"
	authpb "github.com/HappyLadySauce/Beehive/pkg/api/proto/auth/v1"
	"github.com/HappyLadySauce/Beehive/pkg/core"
)

// HandleLogin 处理登录请求
// @Summary 用户登录
// @Description 使用用户 ID 和密码进行登录，返回访问令牌和刷新令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "登录请求参数"
// @Success 200 {object} v1.LoginResponse
// @Failure 400 {object} core.ErrResponse
// @Failure 500 {object} core.ErrResponse
// @Router /api/v1/auth/login [post]
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
		core.WriteResponse(c, err, nil)
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
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body v1.RefreshTokenRequest true "刷新令牌请求参数"
// @Success 200 {object} v1.RefreshTokenResponse
// @Failure 400 {object} core.ErrResponse
// @Failure 500 {object} core.ErrResponse
// @Router /api/v1/auth/refresh [post]
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
		core.WriteResponse(c, err, nil)
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
// @Summary 撤销访问令牌
// @Description 撤销指定访问令牌，使其立即失效
// @Tags auth
// @Accept json
// @Produce json
// @Param request body v1.RevokeTokenRequest true "撤销令牌请求参数"
// @Success 200 {object} v1.RevokeTokenResponse
// @Failure 400 {object} core.ErrResponse
// @Failure 500 {object} core.ErrResponse
// @Router /api/v1/auth/revoke [post]
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
		core.WriteResponse(c, err, nil)
		return
	}

	response := v1.RevokeTokenResponse{
		Success: authResp.Success,
	}

	core.WriteResponse(c, nil, response)
}
