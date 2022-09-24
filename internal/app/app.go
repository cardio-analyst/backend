package app

import (
	"context"
	"fmt"
	"io"

	"github.com/cardio-analyst/backend/internal/adapters/postgres"
	"github.com/cardio-analyst/backend/internal/domain/users"
	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"

	"github.com/cardio-analyst/backend/internal/adapters/http"
	"github.com/cardio-analyst/backend/internal/config"
)

type app struct {
	server  *http.Server
	config  *config.Config
	closers []io.Closer
}

func NewApp(appCtx context.Context, configPath string) *app {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config data: %v", err)
	}

	database, err := postgres.NewDatabase(appCtx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to establish database connection: %v", err)
	}

	userService := users.NewUserService(database)

	srv := http.NewServer(userService)

	return &app{
		server:  srv,
		config:  cfg,
		closers: []io.Closer{srv, database},
	}
}

func (a *app) Start() {
	log.Info("the app is running")

	var group errgroup.Group
	group.Go(func() error {
		listenAddress := fmt.Sprintf(":%d", a.config.HTTP.Port)
		return a.server.Start(listenAddress)
	})

	if err := group.Wait(); err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}

func (a *app) Stop() {
	for _, closer := range a.closers {
		if err := closer.Close(); err != nil {
			log.Errorf("failed to stop the closer: %T: %v", closer, err)
		}
	}

	log.Info("the app has been stopped")
}
