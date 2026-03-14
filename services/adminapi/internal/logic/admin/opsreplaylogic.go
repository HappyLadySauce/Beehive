// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsReplayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpsReplayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsReplayLogic {
	return &OpsReplayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpsReplayLogic) OpsReplay() (resp *types.AdminEmptyResp, err error) {
	// 占位：测试环境消息重放
	return &types.AdminEmptyResp{Code: 0, Message: "ok"}, nil
}
