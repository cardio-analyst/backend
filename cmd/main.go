package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/cardio-analyst/backend/internal/app"
)

const defaultConfigPath = "./configs/config.yaml"

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", defaultConfigPath, "path to config file")

	flag.Parse()

	// graceful shutdown
	applicationCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	defer cancel()

	application := app.NewApp(applicationCtx, configPath)
	go application.Start()
	<-applicationCtx.Done()
	application.Stop()
}
