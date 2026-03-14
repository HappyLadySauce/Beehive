// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsQueuesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpsQueuesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsQueuesLogic {
	return &OpsQueuesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpsQueuesLogic) OpsQueues() (resp *types.AdminEmptyResp, err error) {
	// 占位：运维接口，见 docs/API/admin-http-api.md
	return &types.AdminEmptyResp{Code: 0, Message: "ok"}, nil
}
