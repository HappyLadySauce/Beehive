package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/imrpc/imrpc"
	"github.com/HappyLadySauce/Beehive/imrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostMessageLogic {
	return &PostMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostMessageLogic) PostMessage(in *imrpc.PostMsg) (*imrpc.PostReponse, error) {
	// todo: add your logic here and delete this line

	return &imrpc.PostReponse{}, nil
}
