// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/HappyLadySauce/Beehive/common/libnet"
	"github.com/HappyLadySauce/Beehive/common/socket"
	"github.com/HappyLadySauce/Beehive/common/socketio"
	"github.com/HappyLadySauce/Beehive/edge/internal/config"
	"github.com/HappyLadySauce/Beehive/edge/internal/server"
	"github.com/HappyLadySauce/Beehive/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"golang.org/x/net/websocket"
)

var configFile = flag.String("f", "etc/edge-api.yaml", "the config file")

func main() {
	flag.Parse()

	var err error
	var c config.Config
	conf.MustLoad(*configFile, &c)
	srvCtx := svc.NewServiceContext(c)

	// 禁用统计日志
	logx.DisableStat()

	tcpServer := server.NewTCPServer(srvCtx)
	wsServer := server.NewWSServer(srvCtx)
	protocol := libnet.NewBeehiveProtocol()

	tcpServer.Server, err = socket.NewServe(c.Name, c.TCPListenOn, protocol, c.SendChanSize)
	if err != nil {
		panic(err)
	}
	wsServer.Server, err = socketio.NewServe(c.Name, c.WSListenOn, protocol, c.SendChanSize)
	if err != nil {
		panic(err)
	}
	http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		conn.PayloadType = websocket.BinaryFrame
		wsServer.HandleRequest(conn)
	}))

	go wsServer.Start()
	go tcpServer.HandleRequest()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
}
