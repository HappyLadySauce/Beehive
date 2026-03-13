// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	// GatewayID 用于多实例部署时区分本实例，Presence 注册会话时会使用。
	GatewayID string `json:",optional"`
	// AllowedOrigins 允许的 WebSocket Origin 列表（如 https://app.example.com）。为空时仅允许同源（Origin 与 Host 一致或未携带）。
	AllowedOrigins []string `json:",optional"`
}
