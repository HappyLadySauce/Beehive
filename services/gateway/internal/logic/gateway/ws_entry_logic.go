// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package gateway

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/gateway/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type WsEntryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// WebSocket entrypoint. See docs/API/websocket-client-api.md for JSON message envelope {type, tid, payload, error}.
func NewWsEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WsEntryLogic {
	return &WsEntryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WsEntryLogic) WsEntry() error {
	// todo: add your logic here and delete this line

	return nil
}
