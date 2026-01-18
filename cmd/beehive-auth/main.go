package main

import (
	"context"
	"os"

	"k8s.io/component-base/cli"

	"github.com/HappyLadySauce/Beehive/internal/beehive-auth"
)

const (
	basename = "BeehiveAuth"
)

func main() {
	ctx := context.TODO()
	cmd := beehiveAuth.NewAPICommand(basename, ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}