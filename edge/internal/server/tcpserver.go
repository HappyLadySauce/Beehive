package server

import (
	"github.com/HappyLadySauce/Beehive/edge/internal/svc"
	"github.com/HappyLadySauce/Beehive/common/socket"
	"github.com/HappyLadySauce/Beehive/edge/client"

	"github.com/zeromicro/go-zero/core/logx"
)

// TCP服务器
type TCPServer struct {
	svcCtx *svc.ServiceContext
	Server *socket.Server
}

// 创建TCP服务器
func NewTCPServer(svcCtx *svc.ServiceContext) *TCPServer {
	return &TCPServer{svcCtx: svcCtx}
}

// 处理请求
func (srv *TCPServer) HandleRequest() {
	for {
		session, err := srv.Server.Accept()
		if err != nil {
			panic(err)
		}
		cli := client.NewClient(srv.Server.Manager, session, srv.svcCtx.IMRPC)
		go srv.SessionLoop(cli)
	}
}

// 会话循环
func (srv *TCPServer) SessionLoop(client *client.Client) {
	message, err := client.Receive()
	if err != nil {
		logx.Errorf("[SessionLoop] client.Receive error: %v", err)
		_ = client.Close()
		return
	}

	// 登录
	err = client.Login(message)
	if err != nil {
		logx.Errorf("[SessionLoop] client.Login error: %v", err)
		_ = client.Close()
		return
	}

	// 心跳检测
	go client.Heartbeat()

	// 处理消息
	for {
		message, err := client.Receive()
		if err != nil {
			logx.Errorf("[SessionLoop] client.Receive error: %v", err)
			_ = client.Close()
			return
		}

		err = client.HandlPackage(message)
		if err != nil {
			logx.Errorf("[SessionLoop] client.HandleMessage error: %v", err)
			_ = client.Close()
			return
		}
	}
}

