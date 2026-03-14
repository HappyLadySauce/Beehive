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

type GetUserSessionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSessionsLogic {
	return &GetUserSessionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSessionsLogic) GetUserSessions(req *types.GetUserReq) (resp *types.GetUserSessionsResp, err error) {
	if req.Id == "" {
		return &types.GetUserSessionsResp{Code: 2001, Message: "参数错误"}, nil
	}
	rpcResp, err := l.svcCtx.PresenceSvc.GetOnlineSessions(l.ctx, &presenceservice.GetOnlineSessionsRequest{UserId: req.Id})
	if err != nil {
		return &types.GetUserSessionsResp{Code: 5000, Message: err.Error()}, nil
	}
	sessions := make([]types.SessionItem, 0)
	if rpcResp != nil {
		for _, s := range rpcResp.Sessions {
			sessions = append(sessions, types.SessionItem{
				GatewayId:  s.GatewayId,
				ConnId:     s.ConnId,
				DeviceId:   s.DeviceId,
				DeviceType: s.DeviceType,
				Ip:         "", // SessionInfo 未包含 Ip，需 proto 扩展
				LoginAt:    "", // SessionInfo 未包含 LoginAt，需 proto 扩展
				LastPingAt: formatUnixTime(s.LastPingAt),
			})
		}
	}
	return &types.GetUserSessionsResp{
		Code:    0,
		Message: "ok",
		Data:    types.GetUserSessionsData{Sessions: sessions},
	}, nil
}
