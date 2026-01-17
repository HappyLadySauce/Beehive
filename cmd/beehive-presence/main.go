package main

import (
	"context"
	"os"

	"k8s.io/component-base/cli"

	"github.com/HappyLadySauce/Beehive/internal/beehive-presence"
)

const (
	basename = "BeehivePresence"
)

func main() {
	ctx := context.TODO()
	cmd := beehivePresence.NewAPICommand(basename, ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}