// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/logic"
	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 下载文件
func DownloadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewDownloadFileLogic(r.Context(), svcCtx)
		err := l.DownloadFile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
