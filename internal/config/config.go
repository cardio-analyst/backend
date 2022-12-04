package config

import (
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	dsnEnvKey  = "DATABASE_URL"
	portEnvKey = "PORT"

	accessTokenSigningKeyEnvKey  = "ACCESS_TOKEN_SIGNING_KEY"
	refreshTokenSigningKeyEnvKey = "REFRESH_TOKEN_SIGNING_KEY"

	accessTokenTTLEnvKey  = "ACCESS_TOKEN_TTL_SEC"
	refreshTokenTTLEnvKey = "REFRESH_TOKEN_TTL_SEC"

	smtpPasswordEnvKey = "SMTP_PASSWORD"
)

type Config struct {
	Adapters AdaptersConfig `yaml:"adapters"`
	Services ServicesConfig `yaml:"services"`
}

type AdaptersConfig struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
	SMTP     SMTPConfig     `yaml:"smtp"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServicesConfig struct {
	Auth            AuthConfig            `yaml:"auth"`
	Recommendations RecommendationsConfig `yaml:"recommendations"`
}

type AuthConfig struct {
	AccessToken  TokenConfig `yaml:"access_token"`
	RefreshToken TokenConfig `yaml:"refresh_token"`
}

type TokenConfig struct {
	SigningKey  string `yaml:"signing_key"`
	TokenTTLSec int    `yaml:"token_ttl_sec"`
}

type RecommendationsConfig struct {
	HealthyEating    RecommendationConfig `yaml:"healthy_eating"`
	Smoking          RecommendationConfig `yaml:"smoking"`
	Lifestyle        RecommendationConfig `yaml:"lifestyle"`
	BMI              RecommendationConfig `yaml:"bmi"`
	CholesterolLevel RecommendationConfig `yaml:"cholesterol_level"`
	SBPLevel         RecommendationConfig `yaml:"sbp_level"`
}

type RecommendationConfig struct {
	What string `yaml:"what"`
	Why  string `yaml:"why"`
	How  string `yaml:"how"`
}

func Load(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	cfg.loadFromEnv()

	return &cfg, nil
}

func (c *Config) loadFromEnv() {
	// if dsn was set at the environment
	if dsnFromEnv, exists := os.LookupEnv(dsnEnvKey); exists {
		c.Adapters.Postgres.DSN = dsnFromEnv
	}

	// if port was set at the environment
	if portFromEnv, exists := os.LookupEnv(portEnvKey); exists {
		port, err := strconv.Atoi(portFromEnv)
		if err == nil {
			c.Adapters.HTTP.Port = port
		}
	}

	// if signing keys were set at the environment
	if signingKey, exists := os.LookupEnv(accessTokenSigningKeyEnvKey); exists {
		c.Services.Auth.AccessToken.SigningKey = signingKey
	}
	if signingKey, exists := os.LookupEnv(refreshTokenSigningKeyEnvKey); exists {
		c.Services.Auth.RefreshToken.SigningKey = signingKey
	}

	// if tokens ttl were set at the environment
	if tokenTTL, exists := os.LookupEnv(accessTokenTTLEnvKey); exists {
		ttl, err := strconv.Atoi(tokenTTL)
		if err == nil {
			log.Debugf("access token TTL was set from environment: %v sec", ttl)
			c.Services.Auth.AccessToken.TokenTTLSec = ttl
		}
	}
	if tokenTTL, exists := os.LookupEnv(refreshTokenTTLEnvKey); exists {
		ttl, err := strconv.Atoi(tokenTTL)
		if err == nil {
			log.Debugf("refresh token TTL was set from environment: %v sec", ttl)
			c.Services.Auth.RefreshToken.TokenTTLSec = ttl
		}
	}

	// if smtp password was set at the environment
	if smtpPassword, exists := os.LookupEnv(smtpPasswordEnvKey); exists {
		c.Adapters.SMTP.Password = smtpPassword
	}
}
