package app

import (
	"context"
	"io"
	"os"

	pbauth "github.com/cardio-analyst/backend/pkg/api/proto/auth"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cardio-analyst/backend/internal/gateway/adapters/auth"
	"github.com/cardio-analyst/backend/internal/gateway/adapters/http"
	"github.com/cardio-analyst/backend/internal/gateway/adapters/migrator"
	"github.com/cardio-analyst/backend/internal/gateway/adapters/postgres"
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

type App struct {
	config  config.Config
	server  *http.Server
	closers []io.Closer
}

func New(configPath string) *App {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config data: %v", err)
	}

	var postgresMigrator *migrator.PostgresMigrator
	postgresMigrator, err = migrator.NewPostgresMigrator(cfg.Postgres.URI)
	if err != nil {
		log.Fatalf("failed to initialize postgres migrator: %v", err)
	}
	if err = postgresMigrator.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err = postgresMigrator.Close(); err != nil {
		log.Warnf("failed to close migrator: %v", err)
	}

	storage, err := postgres.NewStorage(cfg.Postgres.URI)
	if err != nil {
		log.Fatalf("failed to create postgres storage: %v", err)
	}

	smtpClient, err := smtp.NewClient(cfg.Gateway.SMTP.Host, cfg.Gateway.SMTP.Port, cfg.Gateway.SMTP.Username, cfg.Gateway.SMTP.Password)
	if err != nil {
		log.Fatalf("failed to create SMTP client: %v", err)
	}

	authGRPCConn, err := grpc.Dial(cfg.Services.Auth.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth gRPC server: %v", err)
	}
	authGRPCClient := pbauth.NewAuthServiceClient(authGRPCConn)

	authClient := auth.NewClient(authGRPCClient)

	services := service.NewServices(cfg, storage, smtpClient, authClient)

	srv := http.NewServer(services)

	return &App{
		config:  cfg,
		server:  srv,
		closers: []io.Closer{srv, storage, smtpClient, authGRPCConn},
	}
}

func (a *App) Start() {
	log.Info("the app is running")

	var group errgroup.Group
	group.Go(func() error {
		listenAddress := a.config.Gateway.HTTPAddress
		return a.server.Start(listenAddress)
	})

	if err := group.Wait(); err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}

func (a *App) Stop(_ context.Context) {
	for _, closer := range a.closers {
		if err := closer.Close(); err != nil {
			log.Warnf("failed to stop the closer: %T: %v", closer, err)
		}
	}

	log.Info("the app has been stopped")
}
