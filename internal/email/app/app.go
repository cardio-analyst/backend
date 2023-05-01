package app

import (
	"context"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/cardio-analyst/backend/internal/email/adapters/smtp"
	"github.com/cardio-analyst/backend/internal/email/config"
	domain "github.com/cardio-analyst/backend/internal/email/domain/service"
	"github.com/cardio-analyst/backend/internal/email/ports/service"
	"github.com/cardio-analyst/backend/internal/pkg/rabbitmq"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

type App struct {
	emailService service.Email
	closers      []io.Closer
}

func New(configPath string) *App {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config data: %v", err)
	}

	smtpClient := smtp.NewClient(cfg.Email.SMTP.Host, cfg.Email.SMTP.Port, cfg.Email.SMTP.Username, cfg.Email.SMTP.Password)
	if err = smtpClient.Connect(); err != nil {
		log.Fatalf("connecting to SMTP server: %v", err)
	}

	rabbitmqClient := rabbitmq.NewClient(rabbitmq.ClientOptions{
		User:         cfg.RabbitMQ.User,
		Password:     cfg.RabbitMQ.Password,
		Host:         cfg.RabbitMQ.Host,
		Port:         cfg.RabbitMQ.Port,
		ExchangeName: cfg.RabbitMQ.Exchange,
		RoutingKey:   cfg.RabbitMQ.RoutingKey,
		QueueName:    cfg.RabbitMQ.Queue,
	})
	if err = rabbitmqClient.Connect(); err != nil {
		log.Fatalf("connecting to RabbitMQ: %v", err)
	}

	emailService := domain.NewEmailService(smtpClient, rabbitmqClient)

	return &App{
		emailService: emailService,
		closers:      []io.Closer{rabbitmqClient, smtpClient},
	}
}

func (a *App) Start() {
	log.Info("the app is running")

	var group errgroup.Group
	group.Go(func() error {
		return a.emailService.ListenToEmailMessages()
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
