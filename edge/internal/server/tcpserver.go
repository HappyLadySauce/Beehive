package server

import (
	"github.com/HappyLadySauce/Beehive/edge/internal/svc"
	// "github.com/HappyLadySauce/Beehive/pkg/socket"
)


type TCPServer struct {
	svcCtx *svc.ServiceContext
	Server *socket.Server
}