package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// PostgreSQL 连接 DSN，例如：
	// postgres://user:password@127.0.0.1:5432/beehive?sslmode=disable
	PostgresDSN string

	// Redis 连接配置，用于登录态、黑名单或 RBAC 缓存（可选）。
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// AccessToken 与 RefreshToken 的有效期（秒），<=0 时使用默认值。
	AccessTokenTTLSeconds  int
	RefreshTokenTTLSeconds int
}
