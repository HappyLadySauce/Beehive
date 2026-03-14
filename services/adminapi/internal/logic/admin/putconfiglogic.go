// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConfigLogic {
	return &PutConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutConfigLogic) PutConfig(req *types.PutConfigReq) (resp *types.AdminEmptyResp, err error) {
	if req.Key == "" {
		return &types.AdminEmptyResp{Code: 2001, Message: "参数错误"}, nil
	}
	// 占位：后续可写 etcd
	return &types.AdminEmptyResp{Code: 0, Message: "ok"}, nil
}
