package app

import (
	"context"
	"io"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "github.com/cardio-analyst/backend/api/proto/analytics"
	grpcserver "github.com/cardio-analyst/backend/internal/analytics/adapters/grpc"
	"github.com/cardio-analyst/backend/internal/analytics/adapters/migrator"
	"github.com/cardio-analyst/backend/internal/analytics/adapters/postgres"
	"github.com/cardio-analyst/backend/internal/analytics/config"
	domain "github.com/cardio-analyst/backend/internal/analytics/domain/service"
	"github.com/cardio-analyst/backend/internal/analytics/ports/service"
	"github.com/cardio-analyst/backend/internal/pkg/rabbitmq"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

type App struct {
	server   *grpcserver.Server
	services service.Services
	closers  []io.Closer
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

	feedbackClient := rabbitmq.NewClient(rabbitmq.ClientOptions{
		User:         cfg.RabbitMQ.User,
		Password:     cfg.RabbitMQ.Password,
		Host:         cfg.RabbitMQ.Host,
		Port:         cfg.RabbitMQ.Port,
		ExchangeName: cfg.RabbitMQ.FeedbackQueue.Exchange,
		RoutingKey:   cfg.RabbitMQ.FeedbackQueue.RoutingKey,
		QueueName:    cfg.RabbitMQ.FeedbackQueue.Queue,
	})
	if err = feedbackClient.Connect(); err != nil {
		log.Fatalf("connecting to feedback RabbitMQ client: %v", err)
	}

	registrationClient := rabbitmq.NewClient(rabbitmq.ClientOptions{
		User:         cfg.RabbitMQ.User,
		Password:     cfg.RabbitMQ.Password,
		Host:         cfg.RabbitMQ.Host,
		Port:         cfg.RabbitMQ.Port,
		ExchangeName: cfg.RabbitMQ.RegistrationQueue.Exchange,
		RoutingKey:   cfg.RabbitMQ.RegistrationQueue.RoutingKey,
		QueueName:    cfg.RabbitMQ.RegistrationQueue.Queue,
	})
	if err = registrationClient.Connect(); err != nil {
		log.Fatalf("connecting to registration RabbitMQ client: %v", err)
	}

	storage, err := postgres.NewStorage(cfg.Postgres.URI)
	if err != nil {
		log.Fatalf("failed to create postgres storage: %v", err)
	}

	services := domain.NewServices(storage, feedbackClient, registrationClient)

	grpcListener, err := net.Listen("tcp", cfg.Analytics.GRPCAddress)
	if err != nil {
		log.Fatalf("failed to listen tcp on %v: %v", cfg.Analytics.GRPCAddress, err)
	}

	grpcServer := grpc.NewServer()
	server := grpcserver.NewServer(grpcServer, grpcListener, services)
	pb.RegisterAnalyticsServiceServer(grpcServer, server)

	return &App{
		server:   server,
		services: services,
		closers:  []io.Closer{server, feedbackClient, registrationClient, storage},
	}
}

func (a *App) Start() {
	var group errgroup.Group

	group.Go(func() error {
		return a.services.Feedback().ListenToFeedbackMessages()
	})

	group.Go(func() error {
		return a.services.Statistics().ListenToRegistrationMessages()
	})

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
