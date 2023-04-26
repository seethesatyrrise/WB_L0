package app

import (
	"github.com/kelseyhightower/envconfig"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/jetstream"
	"http-nats-psql/internal/server"
)

type AppConfig struct {
	DBConfig     database.DBConfig   `envconfig:"PG"`
	JSConfig     jetstream.JSConfig  `envconfig:"JETSTREAM"`
	ServerConfig server.ServerConfig `envconfig:"SERVER"`
}

func newConfig() (*AppConfig, error) {
	cfg := new(AppConfig)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
