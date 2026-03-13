// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package gateway

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Health check endpoint for gateway service.
func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLogic {
	return &HealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthLogic) Health() (resp *types.HealthResp, err error) {
	// todo: add your logic here and delete this line

	return
}
