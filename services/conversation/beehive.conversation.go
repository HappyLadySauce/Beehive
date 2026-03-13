package main

import (
	"flag"
	"fmt"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/config"
	"github.com/HappyLadySauce/Beehive/services/conversation/internal/server"
	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/beehive.conversation.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterConversationServiceServer(grpcServer, server.NewConversationServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
