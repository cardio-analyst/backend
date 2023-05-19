package config

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	databaseURLEnvKey = "DATABASE_URL"
	rmqUserEnvKey     = "RMQ_USER"
	rmqPasswordEnvKey = "RMQ_PASS"
)

type Config struct {
	Analytics AnalyticsConfig `yaml:"analytics"`
	Postgres  PostgresConfig  `yaml:"postgres"`
	RabbitMQ  RabbitMQConfig  `yaml:"rabbitmq"`
}

type AnalyticsConfig struct {
	GRPCAddress string `yaml:"grpc_address"`
}

type PostgresConfig struct {
	URI string `yaml:"uri"`
}

type RabbitMQConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`

	FeedbackQueue     RabbitMQQueueConfig `yaml:"feedback"`
	RegistrationQueue RabbitMQQueueConfig `yaml:"registration"`
}

type RabbitMQQueueConfig struct {
	Exchange   string `yaml:"exchange"`
	RoutingKey string `yaml:"routing_key"`
	Queue      string `yaml:"queue"`
}

func Load(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return Config{}, err
	}

	cfg.loadFromEnv()

	return cfg, nil
}

func (c *Config) loadFromEnv() {
	// if dsn was set at the environment
	if dsnFromEnv, exists := os.LookupEnv(databaseURLEnvKey); exists {
		c.Postgres.URI = dsnFromEnv
		log.Debug("database URI was set from environment")
	}

	// if RMQ user was set at the environment
	if rmqUser, exists := os.LookupEnv(rmqUserEnvKey); exists {
		c.RabbitMQ.User = rmqUser
		log.Debug("RabbitMQ user was set from environment")
	}

	// if RMQ password was set at the environment
	if rmqPassword, exists := os.LookupEnv(rmqPasswordEnvKey); exists {
		c.RabbitMQ.Password = rmqPassword
		log.Debug("RabbitMQ password was set from environment")
	}
}
