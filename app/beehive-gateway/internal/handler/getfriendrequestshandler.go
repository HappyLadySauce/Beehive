// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/logic"
	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取好友申请列表
func GetFriendRequestsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetFriendRequestsLogic(r.Context(), svcCtx)
		resp, err := l.GetFriendRequests()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
