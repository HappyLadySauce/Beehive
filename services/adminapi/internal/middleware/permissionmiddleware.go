package middleware

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/authservice"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// PermissionMiddleware 在认证之后检查当前用户是否具备指定权限，调用 AuthService.CheckPermission；
// 若 context 中无 userId（AdminUserIDKey）返回 1001，若无权限或 RPC 错误返回 1003。
func PermissionMiddleware(svcCtx *svc.ServiceContext, permission string) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userId, _ := r.Context().Value(AdminUserIDKey).(string)
			if userId == "" {
				httpx.WriteJson(w, http.StatusOK, map[string]interface{}{"code": CodeUnauth, "message": "未认证", "data": nil})
				return
			}
			resp, err := svcCtx.AuthSvc.CheckPermission(r.Context(), &authservice.CheckPermissionRequest{
				UserId:     userId,
				Permission: permission,
			})
			if err != nil {
				httpx.WriteJson(w, http.StatusOK, map[string]interface{}{"code": CodeForbidden, "message": "权限不足", "data": nil})
				return
			}
			if resp == nil || !resp.Allowed {
				httpx.WriteJson(w, http.StatusOK, map[string]interface{}{"code": CodeForbidden, "message": "权限不足", "data": nil})
				return
			}
			next(w, r)
		}
	}
}
