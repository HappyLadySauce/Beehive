package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// PostgresDSN 例如 postgres://user:password@127.0.0.1:5432/beehive?sslmode=disable
	PostgresDSN string `json:",optional"`
}
