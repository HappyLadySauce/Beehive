// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsHealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpsHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsHealthLogic {
	return &OpsHealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpsHealthLogic) OpsHealth() (resp *types.AdminEmptyResp, err error) {
	// 占位：可汇总各服务健康状态
	return &types.AdminEmptyResp{Code: 0, Message: "ok"}, nil
}
