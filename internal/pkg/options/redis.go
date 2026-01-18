package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// RedisOptions Redis 配置选项
type RedisOptions struct {
	Addr     string `json:"addr" mapstructure:"addr"`
	Password string `json:"password" mapstructure:"password"`
	DB       int    `json:"db" mapstructure:"db"`
}

// NewRedisOptions 创建 Redis 配置选项
func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

// Validate 验证 Redis 配置
func (o *RedisOptions) Validate() []error {
	var errs []error

	if o.Addr == "" {
		errs = append(errs, fmt.Errorf("redis.addr cannot be empty"))
	}

	if o.DB < 0 {
		errs = append(errs, fmt.Errorf("redis.db cannot be negative"))
	}

	return errs
}

// AddFlags 添加命令行标志
func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Addr, "redis.addr", o.Addr, "Redis server address")
	fs.StringVar(&o.Password, "redis.password", o.Password, "Redis server password")
	fs.IntVar(&o.DB, "redis.db", o.DB, "Redis database number")
}
