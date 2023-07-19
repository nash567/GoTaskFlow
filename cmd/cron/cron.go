package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/GoTaskFlow/cmd/cron/app"
)

const (
	defaultConfPath = "./local.yaml"
)

func main() {
	var configFile, migrationPath, seedDataPath string
	flag.StringVar(&configFile, "config", defaultConfPath, "config file to load")

	flag.Parse()
	application := &app.Application{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application.Init(ctx, configFile, migrationPath, seedDataPath)
	application.Start(ctx)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigterm
	// application.Stop(ctx)
	defer func(cancel context.CancelFunc) {
		cancel()
		os.Exit(0)
	}(cancel)

}
