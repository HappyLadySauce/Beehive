// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"github.com/HappyLadySauce/Beehive/edge/internal/config"
	"github.com/HappyLadySauce/Beehive/edge/internal/handler"
	"github.com/HappyLadySauce/Beehive/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
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

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
