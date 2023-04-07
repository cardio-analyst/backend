package app

import (
	"context"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	grpcserver "github.com/cardio-analyst/backend/internal/auth/adapters/grpc"
	"github.com/cardio-analyst/backend/internal/auth/adapters/mongo"
	"github.com/cardio-analyst/backend/internal/auth/config"
	"github.com/cardio-analyst/backend/internal/auth/domain/service"
	pb "github.com/cardio-analyst/backend/pkg/api/proto/auth"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

type App struct {
	storage *mongo.Storage
	server  *grpcserver.Server
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
		storage: storage,
		server:  server,
	}
}

func (a *App) Start() {
	var group errgroup.Group
	group.Go(func() error {
		return a.server.Serve()
	})

	log.Info("the app is running")

	if err := group.Wait(); err != nil {
		log.Fatalf("grpc server: %v", err)
	}
}

func (a *App) Stop(ctx context.Context) {
	if err := a.server.Close(); err != nil {
		log.Warnf("failed to stop grpc server: %v", err)
	}
	if err := a.storage.Close(ctx); err != nil {
		log.Warnf("failed to close mongodb connection: %v", err)
	}
	log.Info("the app has been stopped")
}
