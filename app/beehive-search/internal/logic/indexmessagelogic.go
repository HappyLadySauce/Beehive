package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-search/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-search/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIndexMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexMessageLogic {
	return &IndexMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 索引消息
func (l *IndexMessageLogic) IndexMessage(in *search.IndexMessageRequest) (*search.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &search.CommonResponse{}, nil
}
