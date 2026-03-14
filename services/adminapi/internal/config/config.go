// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	AuthRpc         zrpc.RpcClientConf `json:",optional"`
	UserRpc         zrpc.RpcClientConf `json:",optional"`
	PresenceRpc     zrpc.RpcClientConf `json:",optional"`
	MessageRpc      zrpc.RpcClientConf `json:",optional"`
	ConversationRpc zrpc.RpcClientConf `json:",optional"`
}
