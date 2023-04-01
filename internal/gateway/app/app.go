package app

import (
	"context"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/cardio-analyst/backend/internal/gateway/adapters/http"
	"github.com/cardio-analyst/backend/internal/gateway/adapters/postgres"
	"github.com/cardio-analyst/backend/internal/gateway/adapters/postgres_migrator"
	"github.com/cardio-analyst/backend/internal/gateway/adapters/smtp"
	"github.com/cardio-analyst/backend/internal/gateway/config"
	"github.com/cardio-analyst/backend/internal/gateway/domain/service"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

type app struct {
	config  config.Config
	server  *http.Server
	closers []io.Closer
}

func NewApp(configPath string) *app {
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

	storage, err := postgres.NewStorage(cfg.Adapters.Postgres)
	if err != nil {
		log.Fatalf("failed to create postgres storage: %v", err)
	}

	smtpClient, err := smtp.NewClient(cfg.Adapters.SMTP)
	if err != nil {
		log.Fatalf("failed to create SMTP client: %v", err)
	}

	services := service.NewServices(cfg.Services, storage, smtpClient)

	srv := http.NewServer(services)

	return &app{
		config:  *cfg,
		server:  srv,
		closers: []io.Closer{srv, storage, smtpClient},
	}
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

func (a *app) Stop(_ context.Context) {
	for _, closer := range a.closers {
		if err := closer.Close(); err != nil {
			log.Warnf("failed to stop the closer: %T: %v", closer, err)
		}
	}

	log.Info("the app has been stopped")
}
