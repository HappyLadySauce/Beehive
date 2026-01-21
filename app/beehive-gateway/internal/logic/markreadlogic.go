// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标记消息已读
func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkReadLogic) MarkRead(req *types.MarkReadReq) (resp *types.MarkReadResp, err error) {
	// todo: add your logic here and delete this line

	return
}
