package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// PostgreSQL 连接 DSN，例如：
	// postgres://user:password@127.0.0.1:5432/beehive?sslmode=disable
	PostgresDSN string `json:",optional"`

	// Redis 连接配置，用于 user profile 缓存。
	RedisAddr     string `json:",optional"`
	RedisPassword string `json:",optional"`
	RedisDB       int    `json:",optional"`

	// 用户 Profile 缓存 TTL（秒），<=0 时使用默认值。
	UserProfileTTLSeconds int `json:",optional"`
}
