package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cardio-analyst/backend/internal/analytics/app"
)

const (
	defaultConfigPath = "../../configs/analytics/config.yaml"

	shutdownTimeout = 10 * time.Second
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", defaultConfigPath, "path to config file")

	flag.Parse()

	// graceful shutdown

	application := app.New(configPath)
	go application.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdown()

	application.Stop(ctx)
}
