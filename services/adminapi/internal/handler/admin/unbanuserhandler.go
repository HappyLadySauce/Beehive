// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/logic/admin"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UnbanUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnbanReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewUnbanUserLogic(r.Context(), svcCtx)
		resp, err := l.UnbanUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
