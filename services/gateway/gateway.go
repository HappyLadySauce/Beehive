// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/HappyLadySauce/Beehive/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/handler"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	if ctx.PushConsumer != nil {
		go ctx.PushConsumer.Run(context.Background())
		defer ctx.PushConsumer.Close()
	}
	if ctx.MessageSendLimit != nil {
		defer ctx.MessageSendLimit.Close()
	}
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
