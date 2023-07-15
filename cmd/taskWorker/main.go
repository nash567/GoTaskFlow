package main

import (
	"context"
	"flag"

	"github.com/GoTaskFlow/cmd/taskWorker/app"
)

const (
	defaultConfPath = "./cmd/taskWorker/local.yaml"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", defaultConfPath, "config file to load")
	flag.Parse()

	application := &app.Application{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application.Init(ctx, configFile)
	application.Start(ctx)

	w := application.RegisterWorkflow()
	application.RunWorker(ctx, w)

}
