package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-search/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-search/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchMessagesLogic {
	return &SearchMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 搜索消息
func (l *SearchMessagesLogic) SearchMessages(in *search.SearchMessagesRequest) (*search.SearchMessagesResponse, error) {
	// todo: add your logic here and delete this line

	return &search.SearchMessagesResponse{}, nil
}
