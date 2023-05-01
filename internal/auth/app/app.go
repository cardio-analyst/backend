package app

import (
	"context"
	"io"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "github.com/cardio-analyst/backend/api/proto/auth"
	grpcserver "github.com/cardio-analyst/backend/internal/auth/adapters/grpc"
	"github.com/cardio-analyst/backend/internal/auth/adapters/mongo"
	"github.com/cardio-analyst/backend/internal/auth/config"
	"github.com/cardio-analyst/backend/internal/auth/domain/service"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

type App struct {
	server  *grpcserver.Server
	closers []io.Closer
}

func New(configPath string) *App {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config data: %v", err)
	}

	storage, err := mongo.NewStorage(cfg.Mongo.URI, cfg.Mongo.DBName)
	if err != nil {
		log.Fatalf("failed to create mongodb storage: %v", err)
	}

	services := service.NewServices(cfg.Auth, storage)

	grpcListener, err := net.Listen("tcp", cfg.Auth.GRPCAddress)
	if err != nil {
		log.Fatalf("failed to listen tcp on %v: %v", cfg.Auth.GRPCAddress, err)
	}

	grpcServer := grpc.NewServer()
	server := grpcserver.NewServer(grpcServer, grpcListener, services)
	pb.RegisterAuthServiceServer(grpcServer, server)

	return &App{
		server:  server,
		closers: []io.Closer{server, storage},
	}
}

func (a *App) Start() {
	var group errgroup.Group
	group.Go(func() error {
		return a.server.Serve()
	})

	log.Info("the app is running")

	if err := group.Wait(); err != nil {
		log.Fatalf("app: %v", err)
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
