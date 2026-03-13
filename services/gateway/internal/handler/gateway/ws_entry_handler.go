// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package gateway

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/services/gateway/internal/logic/gateway"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/svc"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WsEntryHandler 处理 /ws：完成 HTTP→WebSocket 升级，注册连接并交由 logic 维持读循环。
func WsEntryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		deviceID := r.URL.Query().Get("deviceId")
		c := svcCtx.Hub.Register(conn, deviceID)
		defer func() {
			svcCtx.Hub.Unregister(c.ConnID)
			_ = c.Close()
		}()
		l := gateway.NewWsEntryLogic(r.Context(), svcCtx)
		_ = l.ServeConn(c)
	}
}
