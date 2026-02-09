// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/edge/internal/svc"
	"github.com/HappyLadySauce/Beehive/edge/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EdgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEdgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EdgeLogic {
	return &EdgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EdgeLogic) Edge(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
