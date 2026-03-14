// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/HappyLadySauce/Beehive/services/auth/authservice"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/ws"
	"github.com/HappyLadySauce/Beehive/services/presence/presenceservice"
	"github.com/HappyLadySauce/Beehive/services/user/userservice"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Hub    *ws.Hub

	AuthSvc     authservice.AuthService
	PresenceSvc presenceservice.PresenceService
	UserSvc     userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	authCli := zrpc.MustNewClient(c.AuthRpcConf)
	presenceCli := zrpc.MustNewClient(c.PresenceRpcConf)
	userCli := zrpc.MustNewClient(c.UserRpcConf)

	return &ServiceContext{
		Config:      c,
		Hub:         ws.NewHub(c.GatewayID),
		AuthSvc:     authservice.NewAuthService(authCli),
		PresenceSvc: presenceservice.NewPresenceService(presenceCli),
		UserSvc:     userservice.NewUserService(userCli),
	}
}