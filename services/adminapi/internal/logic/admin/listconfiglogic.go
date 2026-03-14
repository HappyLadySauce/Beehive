// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConfigLogic {
	return &ListConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListConfigLogic) ListConfig(req *types.ListConfigReq) (resp *types.ListConfigResp, err error) {
	// 占位：后续可对接 etcd 或配置服务
	return &types.ListConfigResp{Code: 0, Message: "ok", Data: types.ListConfigData{Items: []types.ConfigItem{}}}, nil
}
