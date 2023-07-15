package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoTaskFlow/cmd/cronWorker/app"
)

const (
	defaultConfPath      = "./local.yaml"
	defaultMigrationPath = "./build/db/migrations/"
	defaultSeedDataPath  = "./build/db/seed/"
)

func main() {
	fmt.Println(time.Now().Format(time.RFC3339))
	var configFile, migrationPath, seedDataPath string
	flag.StringVar(&configFile, "config", defaultConfPath, "config file to load")
	flag.StringVar(&migrationPath, "migrations", defaultMigrationPath, "path to SQL migration directory")
	flag.StringVar(&seedDataPath, "seedData", defaultSeedDataPath, "path to SQL seed data directory")

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
