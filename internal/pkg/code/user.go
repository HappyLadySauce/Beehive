package code

// 用户相关错误码 (10410-10419)
const (
	ErrUserAlreadyExists  = iota + 10410 // 用户已存在
	ErrUserNotFound                      // 用户不存在
	ErrEmailAlreadyExists                // 邮箱已被使用
	ErrInvalidEmail                      // 无效的邮箱格式
	ErrWeakPassword                      // 密码强度不足
)
