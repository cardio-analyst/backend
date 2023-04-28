package app

import (
	"context"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/cardio-analyst/backend/internal/analytics/adapters/migrator"
	"github.com/cardio-analyst/backend/internal/analytics/adapters/postgres"
	"github.com/cardio-analyst/backend/internal/analytics/adapters/rabbitmq"
	"github.com/cardio-analyst/backend/internal/analytics/config"
	"github.com/cardio-analyst/backend/internal/analytics/domain/service"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

type App struct {
	rabbitMQClient *rabbitmq.Client
	closers        []io.Closer
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

	services := service.NewServices(storage)

	rabbitmqClient := rabbitmq.NewClient(rabbitmq.ClientOptions{
		User:            cfg.RabbitMQ.User,
		Password:        cfg.RabbitMQ.Password,
		Host:            cfg.RabbitMQ.Host,
		Port:            cfg.RabbitMQ.Port,
		ExchangeName:    cfg.RabbitMQ.Exchange,
		RoutingKey:      cfg.RabbitMQ.RoutingKey,
		QueueName:       cfg.RabbitMQ.Queue,
		MessagesHandler: services.Feedback().MessagesHandler(),
	})
	if err = rabbitmqClient.Connect(); err != nil {
		log.Fatalf("initializing RabbitMQ client: %v", err)
	}

	return &App{
		rabbitMQClient: rabbitmqClient,
		closers:        []io.Closer{rabbitmqClient, storage},
	}
}

func (a *App) Start() {
	log.Info("the app is running")

	var group errgroup.Group
	group.Go(func() error {
		return a.rabbitMQClient.Consume()
	})

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
