// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"
	"github.com/HappyLadySauce/Beehive/services/presence/presenceservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKickUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickUserLogic {
	return &KickUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KickUserLogic) KickUser(req *types.KickReq) (resp *types.AdminEmptyResp, err error) {
	if req.Id == "" {
		return &types.AdminEmptyResp{Code: 2001, Message: "参数错误"}, nil
	}
	if len(req.SessionIds) == 0 {
		sessResp, err := l.svcCtx.PresenceSvc.GetOnlineSessions(l.ctx, &presenceservice.GetOnlineSessionsRequest{UserId: req.Id})
		if err != nil {
			return &types.AdminEmptyResp{Code: 5000, Message: err.Error()}, nil
		}
		if sessResp != nil {
			for _, s := range sessResp.Sessions {
				_, _ = l.svcCtx.PresenceSvc.UnregisterSession(l.ctx, &presenceservice.UnregisterSessionRequest{UserId: req.Id, ConnId: s.ConnId})
			}
		}
	} else {
		for _, connId := range req.SessionIds {
			_, _ = l.svcCtx.PresenceSvc.UnregisterSession(l.ctx, &presenceservice.UnregisterSessionRequest{UserId: req.Id, ConnId: connId})
		}
	}
	return &types.AdminEmptyResp{Code: 0, Message: "ok"}, nil
}
