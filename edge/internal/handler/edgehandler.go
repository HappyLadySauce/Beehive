// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/edge/internal/logic"
	"github.com/HappyLadySauce/Beehive/edge/internal/svc"
	"github.com/HappyLadySauce/Beehive/edge/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func EdgeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewEdgeLogic(r.Context(), svcCtx)
		resp, err := l.Edge(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
