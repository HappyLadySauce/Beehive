// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/logic/admin"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OpsQueuesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewOpsQueuesLogic(r.Context(), svcCtx)
		resp, err := l.OpsQueues()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
