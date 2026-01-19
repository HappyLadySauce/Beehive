package main

import (
	"context"
	"os"

	"k8s.io/component-base/cli"

	"github.com/HappyLadySauce/Beehive/internal/beehive-search"
)

const (
	basename = "BeehiveSearch"
)

func main() {
	ctx := context.TODO()
	cmd := beehiveSearch.NewAPICommand(basename, ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}