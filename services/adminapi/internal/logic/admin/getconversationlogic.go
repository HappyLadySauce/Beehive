// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"
	"github.com/HappyLadySauce/Beehive/services/conversation/conversationservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationLogic {
	return &GetConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationLogic) GetConversation(req *types.GetConversationReq) (resp *types.GetConversationResp, err error) {
	if req.Id == "" {
		return &types.GetConversationResp{Code: 2001, Message: "参数错误"}, nil
	}
	rpcResp, err := l.svcCtx.ConversationSvc.GetConversation(l.ctx, &conversationservice.GetConversationRequest{Id: req.Id})
	if err != nil {
		return &types.GetConversationResp{Code: 5000, Message: err.Error()}, nil
	}
	if rpcResp == nil || rpcResp.Conversation == nil {
		return &types.GetConversationResp{Code: 3001, Message: "会话不存在"}, nil
	}
	c := rpcResp.Conversation
	return &types.GetConversationResp{
		Code:    0,
		Message: "ok",
		Data: types.GetConversationData{
			Id:           c.Id,
			Type:         c.Type,
			Name:         c.Name,
			MemberCount:  int(c.MemberCount),
			CreatedAt:    formatUnixTime(c.CreatedAt),
			LastActiveAt: formatUnixTime(c.LastActiveAt),
		},
	}, nil
}
