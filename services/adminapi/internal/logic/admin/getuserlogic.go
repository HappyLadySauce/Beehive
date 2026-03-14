// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"
	"github.com/HappyLadySauce/Beehive/services/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.GetUserResp, err error) {
	if req.Id == "" {
		return &types.GetUserResp{Code: 2001, Message: "参数错误"}, nil
	}
	rpcResp, err := l.svcCtx.UserSvc.GetUser(l.ctx, &userservice.GetUserRequest{Id: req.Id})
	if err != nil {
		return &types.GetUserResp{Code: 5000, Message: err.Error()}, nil
	}
	if rpcResp == nil || rpcResp.User == nil {
		return &types.GetUserResp{Code: 3001, Message: "用户不存在"}, nil
	}
	u := rpcResp.User
	return &types.GetUserResp{
		Code:    0,
		Message: "ok",
		Data: types.GetUserData{
			Id:          u.Id,
			Nickname:    u.Nickname,
			Email:       "",
			Status:      "active",
			CreatedAt:   "",
			LastLoginAt: "",
			Profile:     types.UserProfile{AvatarUrl: u.AvatarUrl, Bio: u.Bio},
		},
	}, nil
}
