package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// Redis 连接配置，用于会话存储（user session 索引与 session 详情）。
	RedisAddr     string `json:",optional"`
	RedisPassword string `json:",optional"`
	RedisDB       int    `json:",optional"`

	// 会话 TTL（秒），心跳刷新时续期；<=0 时使用默认 90。
	SessionTTLSeconds int `json:",optional"`
}
