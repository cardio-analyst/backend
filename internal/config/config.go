package config

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	dsnEnvKey  = "DATABASE_URL"
	portEnvKey = "PORT"
)

type Config struct {
	Adapters AdaptersConfig `yaml:"adapters"`
	Services ServicesConfig `yaml:"services"`
}

type AdaptersConfig struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

type ServicesConfig struct {
	Auth AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	SigningKey []byte        // created from signingKey
	TokenTTL   time.Duration // created from tokenTTLMinutes

	signingKey      string `yaml:"signing_key"`
	tokenTTLMinutes int    `yaml:"token_ttl_minutes"`
}

func Load(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	// if dsn was set at the environment
	if dsnFromEnv, exists := os.LookupEnv(dsnEnvKey); exists {
		cfg.Adapters.Postgres.DSN = dsnFromEnv
	}

	// if port was set at the environment
	if portFromEnv, exists := os.LookupEnv(portEnvKey); exists {
		var port int
		port, err = strconv.Atoi(portFromEnv)
		if err == nil {
			cfg.Adapters.HTTP.Port = port
		}
	}

	cfg.Services.Auth.SigningKey = []byte(cfg.Services.Auth.signingKey)
	cfg.Services.Auth.TokenTTL = time.Duration(cfg.Services.Auth.tokenTTLMinutes) * time.Minute

	return &cfg, nil
}
