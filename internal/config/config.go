package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

const dsnEnvKey = "DSN"

type Config struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
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
	dsnFromEnv, exists := os.LookupEnv(dsnEnvKey)
	if exists {
		cfg.Postgres.DSN = dsnFromEnv
	}

	return &cfg, nil
}
