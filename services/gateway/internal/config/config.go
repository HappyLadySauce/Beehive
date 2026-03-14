// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	// GatewayID 用于多实例部署时区分本实例，Presence 注册会话时会使用。
	GatewayID string `json:",optional"`
	// AllowedOrigins 允许的 WebSocket Origin 列表（如 https://app.example.com）。为空时仅允许同源（Origin 与 Host 一致或未携带）。
	AllowedOrigins []string `json:",optional"`

	// AuthRpcConf 为 AuthService 的 zrpc 客户端配置；可选，未配置时 auth.login/auth.tokenLogin/auth.logout 不可用。
	AuthRpcConf zrpc.RpcClientConf `json:",optional"`
	// PresenceRpcConf 为 PresenceService 的 zrpc 客户端配置；可选，未配置时 presence.ping 不可用。
	PresenceRpcConf zrpc.RpcClientConf `json:",optional"`
	// UserRpcConf 为 UserService 的 zrpc 客户端配置（获取用户资料等）；可选，未配置时 user.me 不可用。
	UserRpcConf zrpc.RpcClientConf `json:",optional"`
	// ConversationRpcConf 为 ConversationService 的 zrpc 客户端配置；可选，未配置时 conversation.list 不可用。
	ConversationRpcConf zrpc.RpcClientConf `json:",optional"`
	// MessageRpcConf 为 MessageService 的 zrpc 客户端配置；可选，未配置时 message.send / message.history 不可用。
	MessageRpcConf zrpc.RpcClientConf `json:",optional"`
}

// UserRpcConfigured 判断是否已配置 UserService（Etcd 或 Endpoints），未配置时 Gateway 可不依赖 User 服务启动。
func (c *Config) UserRpcConfigured() bool {
	return len(c.UserRpcConf.Endpoints) > 0 || c.UserRpcConf.Etcd.Key != ""
}

// ConversationRpcConfigured 判断是否已配置 ConversationService。
func (c *Config) ConversationRpcConfigured() bool {
	return len(c.ConversationRpcConf.Endpoints) > 0 || c.ConversationRpcConf.Etcd.Key != ""
}

// MessageRpcConfigured 判断是否已配置 MessageService。
func (c *Config) MessageRpcConfigured() bool {
	return len(c.MessageRpcConf.Endpoints) > 0 || c.MessageRpcConf.Etcd.Key != ""
}

func (c *Config) AuthRpcConfigured() bool {
	return len(c.AuthRpcConf.Endpoints) > 0 || c.AuthRpcConf.Etcd.Key != ""
}

func (c *Config) PresenceRpcConfigured() bool {
	return len(c.PresenceRpcConf.Endpoints) > 0 || c.PresenceRpcConf.Etcd.Key != ""
}