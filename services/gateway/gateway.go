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
	// 不在此处 defer MessageSendLimit.Close()：server.Stop() 返回时仍有 in-flight 的 WebSocket 请求可能调用 Allow()，
	// 若先关闭 Redis 会造成竞态。进程退出时由 OS 回收连接；需显式关闭时应在优雅退出流程中先停止接收请求并等待请求排空后再关闭。
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
