package app

import (
	"context"
	"fmt"
	"github.com/cardio-analyst/backend/internal/domain/disease"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/cardio-analyst/backend/internal/adapters/http"
	"github.com/cardio-analyst/backend/internal/adapters/postgres"
	"github.com/cardio-analyst/backend/internal/adapters/postgres_migrator"
	"github.com/cardio-analyst/backend/internal/config"
	"github.com/cardio-analyst/backend/internal/domain/auth"
	"github.com/cardio-analyst/backend/internal/domain/user"
)

type app struct {
	server  *http.Server
	config  *config.Config
	closers []io.Closer
}

func NewApp(appCtx context.Context, configPath string) *app {
	initializeLogger()

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config data: %v", err)
	}

	var migrator *postgres_migrator.PostgresMigrator
	migrator, err = postgres_migrator.NewPostgresMigrator(cfg.Adapters.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to initialize postgres migrator: %v", err)
	}
	if err = migrator.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err = migrator.Close(); err != nil {
		log.Warnf("failed to close migrator: %v", err)
	}

	database, err := postgres.NewDatabase(appCtx, cfg.Adapters.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to establish database connection: %v", err)
	}

	authService := auth.NewAuthService(cfg.Services.Auth, database, database)
	userService := user.NewUserService(database)
	diseaseService := disease.NewDiseaseService(database)

	srv := http.NewServer(authService, userService, diseaseService)

	return &app{
		server:  srv,
		config:  cfg,
		closers: []io.Closer{srv, database},
	}
}

func initializeLogger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

func (a *app) Start() {
	log.Info("the app is running")

	var group errgroup.Group
	group.Go(func() error {
		listenAddress := fmt.Sprintf(":%v", a.config.Adapters.HTTP.Port)
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
