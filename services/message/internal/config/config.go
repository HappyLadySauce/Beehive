package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// PostgresDSN 例如 postgres://user:password@127.0.0.1:5432/beehive?sslmode=disable
	PostgresDSN string `json:",optional"`

	// RabbitMQ 可选；配置后 PostMessage 成功会发布 message.created 事件
	RabbitMQURL      string `json:",optional"` // amqp://guest:guest@127.0.0.1:5672/
	RabbitMQExchange string `json:",optional"`  // im.events
	RabbitMQRouteKey string `json:",optional"` // message.created
}
