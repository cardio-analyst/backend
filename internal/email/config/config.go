package config

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	smtpPasswordEnvKey = "SMTP_PASSWORD"
	rmqUserEnvKey      = "RMQ_USER"
	rmqPasswordEnvKey  = "RMQ_PASS"
)

type Config struct {
	Email    EmailConfig    `yaml:"email"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

type EmailConfig struct {
	SMTP SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RabbitMQConfig struct {
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
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
	// if smtp password was set at the environment
	if smtpPassword, exists := os.LookupEnv(smtpPasswordEnvKey); exists {
		c.Email.SMTP.Password = smtpPassword
		log.Debug("SMTP password was set from environment")
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
