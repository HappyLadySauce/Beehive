package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-search/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-search/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMessageIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageIndexLogic {
	return &DeleteMessageIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除消息索引
func (l *DeleteMessageIndexLogic) DeleteMessageIndex(in *search.DeleteMessageIndexRequest) (*search.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &search.CommonResponse{}, nil
}
