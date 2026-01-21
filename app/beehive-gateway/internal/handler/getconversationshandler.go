// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/logic"
	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取会话列表
func GetConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetConversationsLogic(r.Context(), svcCtx)
		resp, err := l.GetConversations()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
