package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/authservice"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// AuthMiddleware 从 Authorization: Bearer <token> 取 token，调用 AuthService.ValidateToken，
// 成功则将 userId 写入 context（AdminUserIDKey），失败返回 { code: 1001 }。
func AuthMiddleware(svcCtx *svc.ServiceContext) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				httpx.WriteJson(w, http.StatusOK, map[string]interface{}{"code": CodeUnauth, "message": "未认证", "data": nil})
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			resp, err := svcCtx.AuthSvc.ValidateToken(r.Context(), &authservice.ValidateTokenRequest{AccessToken: token})
			if err != nil {
				httpx.WriteJson(w, http.StatusOK, map[string]interface{}{"code": CodeInternal, "message": err.Error(), "data": nil})
				return
			}
			if resp == nil || !resp.Valid || resp.UserId == "" {
				httpx.WriteJson(w, http.StatusOK, map[string]interface{}{"code": CodeUnauth, "message": "未认证", "data": nil})
				return
			}
			ctx := context.WithValue(r.Context(), AdminUserIDKey, resp.UserId)
			next(w, r.WithContext(ctx))
		}
	}
}
