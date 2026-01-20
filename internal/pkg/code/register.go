package code

import "net/http"

func init() {
	// 通用错误码
	register(ErrSuccess, http.StatusOK, "success")
	register(ErrBind, http.StatusBadRequest, "invalid request")

	// 认证相关错误码 (10401-10409)
	register(ErrInvalidCredentials, http.StatusUnauthorized, "invalid username or password")
	register(ErrTokenExpired, http.StatusUnauthorized, "token has expired")
	register(ErrTokenInvalid, http.StatusUnauthorized, "invalid token")
	register(ErrTokenRevoked, http.StatusUnauthorized, "token has been revoked")
	register(ErrUnauthorized, http.StatusUnauthorized, "unauthorized")

	// 用户相关错误码 (10410-10419)
	register(ErrUserAlreadyExists, http.StatusConflict, "user already exists")
	register(ErrUserNotFound, http.StatusNotFound, "user not found")
	register(ErrEmailAlreadyExists, http.StatusConflict, "email already in use")
	register(ErrInvalidEmail, http.StatusBadRequest, "invalid email format")
	register(ErrWeakPassword, http.StatusBadRequest, "password is too weak")
}
