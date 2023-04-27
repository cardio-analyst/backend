package config

import (
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	accessTokenSigningKeyEnvKey  = "ACCESS_TOKEN_SIGNING_KEY"
	refreshTokenSigningKeyEnvKey = "REFRESH_TOKEN_SIGNING_KEY"

	secretKeySigningKeyEnvKey = "SECRET_KEY_SIGNING_KEY"

	accessTokenTTLEnvKey  = "ACCESS_TOKEN_TTL_SEC"
	refreshTokenTTLEnvKey = "REFRESH_TOKEN_TTL_SEC"
)

type Config struct {
	Auth  AuthConfig  `yaml:"auth"`
	Mongo MongoConfig `yaml:"mongo"`
}

type AuthConfig struct {
	GRPCAddress         string      `yaml:"grpc_address"`
	AccessToken         TokenConfig `yaml:"access_token"`
	RefreshToken        TokenConfig `yaml:"refresh_token"`
	SecretKeySigningKey string      `yaml:"secret_key_signing_key"`
}

type TokenConfig struct {
	SigningKey  string `yaml:"signing_key"`
	TokenTTLSec int64  `yaml:"token_ttl_sec"`
}

type MongoConfig struct {
	URI    string `yaml:"uri"`
	DBName string `yaml:"db_name"`
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
	// if signing keys were set at the environment
	if signingKey, exists := os.LookupEnv(accessTokenSigningKeyEnvKey); exists {
		c.Auth.AccessToken.SigningKey = signingKey
		log.Debug("access token signing key was set from environment")
	}
	if signingKey, exists := os.LookupEnv(refreshTokenSigningKeyEnvKey); exists {
		c.Auth.RefreshToken.SigningKey = signingKey
		log.Debug("refresh token signing key was set from environment")
	}
	if signingKey, exists := os.LookupEnv(secretKeySigningKeyEnvKey); exists {
		c.Auth.SecretKeySigningKey = signingKey
		log.Debug("secret key signing key was set from environment")
	}

	// if tokens ttl were set at the environment
	if tokenTTL, exists := os.LookupEnv(accessTokenTTLEnvKey); exists {
		if ttl, err := strconv.ParseInt(tokenTTL, 10, 64); err == nil {
			log.Debugf("access token TTL was set from environment: %v sec", ttl)
			c.Auth.AccessToken.TokenTTLSec = ttl
		}
	}
	if tokenTTL, exists := os.LookupEnv(refreshTokenTTLEnvKey); exists {
		if ttl, err := strconv.ParseInt(tokenTTL, 10, 64); err == nil {
			log.Debugf("refresh token TTL was set from environment: %v sec", ttl)
			c.Auth.RefreshToken.TokenTTLSec = ttl
		}
	}
}
