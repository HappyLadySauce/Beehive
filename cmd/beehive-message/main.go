package main

import (
	"context"
	"os"

	"k8s.io/component-base/cli"

	"github.com/HappyLadySauce/Beehive/internal/beehive-message"
)

const (
	basename = "BeehiveMessage"
)

func main() {
	ctx := context.TODO()
	cmd := beehiveMessage.NewAPICommand(basename, ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}