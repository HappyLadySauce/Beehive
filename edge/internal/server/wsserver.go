package server

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive/edge/internal/svc"
	"github.com/HappyLadySauce/Beehive/edge/client"
	"github.com/HappyLadySauce/Beehive/common/socketio"
	"golang.org/x/net/websocket"

	"github.com/zeromicro/go-zero/core/logx"
)

type WSServer struct {
	svcCtx *svc.ServiceContext
	Server *socketio.Server
}

func NewWSServer(svcCtx *svc.ServiceContext) *WSServer {
	return &WSServer{svcCtx: svcCtx}
}

func (ws *WSServer) Start() {
	err := http.ListenAndServe(ws.Server.Address, nil)
	if err != nil {
		panic(err)
	}
}

func (ws *WSServer) HandleRequest(conn *websocket.Conn) {
	session, err := ws.Server.Accept(conn)
	if err != nil {
		panic(err)
	}
	cli := client.NewClient(ws.Server.Manager, session, ws.svcCtx.IMRPC)
	go ws.SessionLoop(cli)
}

func (ws *WSServer) SessionLoop(client *client.Client) {
	message, err := client.Receive()
	if err != nil {
		logx.Errorf("[SessionLoop] client.Receive error: %v", err)
		_ = client.Close()
		return
	}
	err = client.Login(message)
	if err != nil {
		logx.Errorf("[SessionLoop] client.Login error: %v", err)
		_ = client.Close()
		return
	}

	for {
		message, err = client.Receive()
		if err != nil {
			logx.Errorf("[SessionLoop] client.Receive error: %v", err)
			_ = client.Close()
			return
		}
		err = client.HandlPackage(message)
		if err != nil {
			logx.Errorf("[SessionLoop] client.HandlePackage error: %v", err)
			_ = client.Close()
			return
		}
	}
}