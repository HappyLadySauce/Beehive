package main

import (
	"flag"
	"fmt"

	"github.com/HappyLadySauce/Beehive/app/beehive-file/file"
	"github.com/HappyLadySauce/Beehive/app/beehive-file/internal/config"
	"github.com/HappyLadySauce/Beehive/app/beehive-file/internal/server"
	"github.com/HappyLadySauce/Beehive/app/beehive-file/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/file.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		file.RegisterFileServiceServer(grpcServer, server.NewFileServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
