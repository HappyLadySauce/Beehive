// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/logic"
	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户登录
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
