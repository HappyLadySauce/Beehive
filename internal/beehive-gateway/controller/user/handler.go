package user

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/HappyLadySauce/Beehive/internal/pkg/common/types/v1"
	userpb "github.com/HappyLadySauce/Beehive/pkg/api/proto/user/v1"
	"github.com/HappyLadySauce/Beehive/pkg/core"
)

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
		core.HandleGRPCError(c, err)
		return
	}

	response := v1.RegisterResponse{
		ID: userResp.Id,
	}

	core.WriteResponse(c, nil, response)
}
