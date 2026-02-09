package server

import (
	"github.com/HappyLadySauce/Beehive/edge/internal/svc"
	"github.com/HappyLadySauce/Beehive/common/socket"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCPServer struct {
	svcCtx *svc.ServiceContext
	Server *socket.Server
}
