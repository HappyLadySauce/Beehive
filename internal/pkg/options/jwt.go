package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// JWTOptions JWT 配置选项
type JWTOptions struct {
	Secret              string `json:"secret" mapstructure:"secret"`
	ExpireHours         int    `json:"expire-hours" mapstructure:"expire-hours"`
	RefreshExpireHours  int    `json:"refresh-expire-hours" mapstructure:"refresh-expire-hours"`
}

// NewJWTOptions 创建 JWT 配置选项
func NewJWTOptions() *JWTOptions {
	return &JWTOptions{
		Secret:             "your-secret-key-change-in-production",
		ExpireHours:        24,  // 24 小时
		RefreshExpireHours: 168, // 7 天
	}
}

// Validate 验证 JWT 配置
func (o *JWTOptions) Validate() []error {
	var errs []error

	if o.Secret == "" {
		errs = append(errs, fmt.Errorf("jwt.secret cannot be empty"))
	}

	if o.ExpireHours <= 0 {
		errs = append(errs, fmt.Errorf("jwt.expire-hours must be greater than 0"))
	}

	if o.RefreshExpireHours <= 0 {
		errs = append(errs, fmt.Errorf("jwt.refresh-expire-hours must be greater than 0"))
	}

	return errs
}

// AddFlags 添加命令行标志
func (o *JWTOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Secret, "jwt.secret", o.Secret, "JWT secret key for signing tokens")
	fs.IntVar(&o.ExpireHours, "jwt.expire-hours", o.ExpireHours, "Access token expiration time in hours")
	fs.IntVar(&o.RefreshExpireHours, "jwt.refresh-expire-hours", o.RefreshExpireHours, "Refresh token expiration time in hours")
}
