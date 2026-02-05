package main

import (
	"context"
	"os"

	"k8s.io/component-base/cli"

	beehiveAPI "github.com/HappyLadySauce/Beehive/internal/api"
)

const (
	basename = "BeehiveAPI"
)

// @title           Beehive API Service API
// @version         1.0
// @description     Beehive API 微服务 API 文档
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
	cmd := beehiveAPI.NewAPICommand(basename, ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
