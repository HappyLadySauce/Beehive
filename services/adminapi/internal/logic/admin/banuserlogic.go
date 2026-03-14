// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BanUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBanUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanUserLogic {
	return &BanUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BanUserLogic) BanUser(req *types.BanReq) (resp *types.AdminEmptyResp, err error) {
	if req.Id == "" {
		return &types.AdminEmptyResp{Code: 2001, Message: "参数错误"}, nil
	}
	// 占位：封禁状态后续可对接 Auth 或 User 的 Ban RPC
	return &types.AdminEmptyResp{Code: 0, Message: "ok"}, nil
}
