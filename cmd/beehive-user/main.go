package main

import (
	"context"
	"os"

	"k8s.io/component-base/cli"

	"github.com/HappyLadySauce/Beehive/internal/beehive-user"
)

const (
	basename = "BeehiveUser"
)

func main() {
	ctx := context.TODO()
	cmd := beehiveUser.NewAPICommand(basename, ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}