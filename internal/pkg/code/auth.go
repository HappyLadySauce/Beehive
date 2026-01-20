package code

// 认证相关错误码 (10401-10409)
const (
	ErrInvalidCredentials = iota + 10401 // 无效的用户名或密码
	ErrTokenExpired                      // Token 已过期
	ErrTokenInvalid                      // Token 无效
	ErrTokenRevoked                      // Token 已被撤销
	ErrUnauthorized                      // 未授权
)
