package user

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/HappyLadySauce/Beehive/internal/pkg/common/types/v1"
	userpb "github.com/HappyLadySauce/Beehive/pkg/api/proto/user/v1"
	"github.com/HappyLadySauce/Beehive/pkg/core"
)

// HandleRegister 处理注册请求
// @Summary 用户注册
// @Description 创建一个新的用户账号
// @Tags user
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "注册请求参数"
// @Success 200 {object} v1.RegisterResponse
// @Failure 400 {object} core.ErrResponse
// @Failure 500 {object} core.ErrResponse
// @Router /api/v1/user/register [post]
// HandleRegister 处理注册请求
func (h *Handler) HandleRegister(c *gin.Context) {
	var req v1.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponseBindErr(c, err, nil)
		return
	}

	// 调用 User Service
	userReq := &userpb.RegisterRequest{
		Nickname:    req.Nickname,
		Avatar:      req.Avatar,
		Email:       req.Email,
		Description: req.Description,
		Password:    req.Password,
	}

	userResp, err := h.userClient.Register(c.Request.Context(), userReq)
	if err != nil {
		core.WriteResponseBindErr(c, err, nil)
		return
	}

	response := v1.RegisterResponse{
		ID: userResp.Id,
	}

	core.WriteResponse(c, nil, response)
}
