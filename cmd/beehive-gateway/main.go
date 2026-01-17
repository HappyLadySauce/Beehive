package main

import (
	"context"
	"os"

	"k8s.io/component-base/cli"

	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway"
)

const (
	basename = "BeehiveGateway"
)

func main() {
	ctx := context.TODO()
	cmd := beehivegateway.NewAPICommand(basename, ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}