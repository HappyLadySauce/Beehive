// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package gateway

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/services/gateway/internal/logic/gateway"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Health check endpoint for gateway service.
func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := gateway.NewHealthLogic(r.Context(), svcCtx)
		resp, err := l.Health()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
