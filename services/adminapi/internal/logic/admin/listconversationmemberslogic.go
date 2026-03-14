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

type ListConversationMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListConversationMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConversationMembersLogic {
	return &ListConversationMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListConversationMembersLogic) ListConversationMembers(req *types.GetConversationReq) (resp *types.ListMembersResp, err error) {
	if req.Id == "" {
		return &types.ListMembersResp{Code: 2001, Message: "参数错误"}, nil
	}
	rpcResp, err := l.svcCtx.ConversationSvc.ListMembers(l.ctx, &conversationservice.ListMembersRequest{ConversationId: req.Id})
	if err != nil {
		return &types.ListMembersResp{Code: 5000, Message: err.Error()}, nil
	}
	items := make([]types.MemberItem, 0)
	if rpcResp != nil {
		for _, m := range rpcResp.Items {
			items = append(items, types.MemberItem{
				UserId:   m.UserId,
				Role:     m.Role,
				JoinedAt: formatUnixTime(m.JoinedAt),
				Status:   m.Status,
			})
		}
	}
	return &types.ListMembersResp{
		Code:    0,
		Message: "ok",
		Data:    types.ListMembersData{Items: items},
	}, nil
}
