// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package gateway

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/services/gateway/internal/logic/gateway"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// WebSocket entrypoint. See docs/API/websocket-client-api.md for JSON message envelope {type, tid, payload, error}.
func WsEntryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := gateway.NewWsEntryLogic(r.Context(), svcCtx)
		err := l.WsEntry()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
