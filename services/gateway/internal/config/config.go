// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	// GatewayID 用于多实例部署时区分本实例，Presence 注册会话时会使用。
	GatewayID string `json:",optional"`
}
