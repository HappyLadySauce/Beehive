// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/HappyLadySauce/Beehive/edge/internal/config"
	"github.com/HappyLadySauce/Beehive/imrpc/imrpcclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	IMRPC imrpcclient.Imrpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.IMRPC)
	return &ServiceContext{
		Config: c,
		IMRPC: imrpcclient.NewImrpc(client),
	}
}
