// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbanUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbanUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbanUserLogic {
	return &UnbanUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbanUserLogic) UnbanUser(req *types.UnbanReq) (resp *types.AdminEmptyResp, err error) {
	if req.Id == "" {
		return &types.AdminEmptyResp{Code: 2001, Message: "参数错误"}, nil
	}
	// 占位：解封后续可对接 Auth/User
	return &types.AdminEmptyResp{Code: 0, Message: "ok"}, nil
}
