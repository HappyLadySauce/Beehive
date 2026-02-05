// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/imapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/imapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImapiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImapiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImapiLogic {
	return &ImapiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImapiLogic) Imapi(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
