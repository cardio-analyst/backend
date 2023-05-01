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
	Gateway         GatewayConfig         `yaml:"gateway"`
	Postgres        PostgresConfig        `yaml:"postgres"`
	RabbitMQ        RabbitMQConfig        `yaml:"rabbitmq"`
	Recommendations RecommendationsConfig `yaml:"recommendations"`
	Services        ServicesConfig        `yaml:"services"`
}

type GatewayConfig struct {
	HTTPAddress string `yaml:"http_address"`
}

type PostgresConfig struct {
	URI string `yaml:"uri"`
}

type RabbitMQConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`

	EmailsQueue   RabbitMQQueueConfig `yaml:"emails"`
	FeedbackQueue RabbitMQQueueConfig `yaml:"feedback"`
}

type RabbitMQQueueConfig struct {
	Exchange   string `yaml:"exchange"`
	RoutingKey string `yaml:"routing_key"`
	Queue      string `yaml:"queue"`
}

type RecommendationsConfig struct {
	HealthyEating    RecommendationConfig `yaml:"healthy_eating"`
	Smoking          RecommendationConfig `yaml:"smoking"`
	Lifestyle        RecommendationConfig `yaml:"lifestyle"`
	BMI              RecommendationConfig `yaml:"bmi"`
	CholesterolLevel RecommendationConfig `yaml:"cholesterol_level"`
	SBPLevel         RecommendationConfig `yaml:"sbp_level"`
	Risk             RecommendationConfig `yaml:"risk"`
}

type RecommendationConfig struct {
	What string `yaml:"what"`
	Why  string `yaml:"why"`
	How  string `yaml:"how"`
}

type ServicesConfig struct {
	Auth      ServiceConfig `yaml:"auth"`
	Analytics ServiceConfig `yaml:"analytics"`
}

type ServiceConfig struct {
	GRPCAddress string `yaml:"grpc_address"`
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
