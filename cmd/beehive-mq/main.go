package main

import (
	"context"
	"os"

	"k8s.io/component-base/cli"

	beehiveMQ "github.com/HappyLadySauce/Beehive/internal/mq"
)

const (
	basename = "BeehiveMQ"
)

// @title           Beehive MQ Service API
// @version         1.0
// @description     Beehive MQ 微服务 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

func main() {
	ctx := context.TODO()
	cmd := beehiveMQ.NewAPICommand(basename, ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
