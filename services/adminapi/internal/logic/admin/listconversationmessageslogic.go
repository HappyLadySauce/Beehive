// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"
	"github.com/HappyLadySauce/Beehive/services/message/messageservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConversationMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListConversationMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConversationMessagesLogic {
	return &ListConversationMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListConversationMessagesLogic) ListConversationMessages(req *types.ListMessagesReq) (resp *types.ListMessagesResp, err error) {
	if req.Id == "" {
		return &types.ListMessagesResp{Code: 2001, Message: "参数错误"}, nil
	}
	limit := int32(req.PageSize)
	if limit <= 0 {
		limit = 50
	}
	rpcResp, err := l.svcCtx.MessageSvc.GetHistory(l.ctx, &messageservice.GetHistoryRequest{
		ConversationId: req.Id,
		BeforeTime:     0,
		Limit:          limit,
	})
	if err != nil {
		return &types.ListMessagesResp{Code: 5000, Message: err.Error()}, nil
	}
	items := make([]types.MessageItem, 0)
	if rpcResp != nil {
		for _, m := range rpcResp.Items {
			body := types.MessageBody{}
			if m.Body != nil {
				body.Type = m.Body.Type
				body.Text = m.Body.Text
			}
			items = append(items, types.MessageItem{
				ServerMsgId:     m.ServerMsgId,
				ConversationId:  m.ConversationId,
				FromUserId:      m.FromUserId,
				Body:            body,
				ServerTime:      formatUnixTime(m.ServerTime),
			})
		}
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	// GetHistory RPC 未返回 total 总数，此处用本页条数；完整分页需 MessageService 提供 Count 或返回 Total
	total := len(items)
	return &types.ListMessagesResp{
		Code:    0,
		Message: "ok",
		Data: types.ListMessagesData{
			Items:    items,
			Page:     page,
			PageSize: int(limit),
			Total:    total,
		},
	}, nil
}
