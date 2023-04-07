package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	postgresDSNEnvKey  = "DATABASE_URL"
	smtpPasswordEnvKey = "SMTP_PASSWORD"
)

type Config struct {
	Gateway         GatewayConfig         `yaml:"gateway"`
	Postgres        PostgresConfig        `yaml:"postgres"`
	Recommendations RecommendationsConfig `yaml:"recommendations"`
	Services        ServicesConfig        `yaml:"services"`
}

type GatewayConfig struct {
	HTTPAddress string     `yaml:"http_address"`
	SMTP        SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
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
	Auth ServiceConfig `yaml:"auth"`
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
	if dsnFromEnv, exists := os.LookupEnv(postgresDSNEnvKey); exists {
		c.Postgres.DSN = dsnFromEnv
	}

	// if smtp password was set at the environment
	if smtpPassword, exists := os.LookupEnv(smtpPasswordEnvKey); exists {
		c.Gateway.SMTP.Password = smtpPassword
	}
}
