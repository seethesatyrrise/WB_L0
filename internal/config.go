package internal

import (
	"github.com/caarlos0/env"
)

type DBConfig struct {
	PGAddress  string `env:"PG_ADDRESS" envDefault:":5432"`
	PGUser     string `env:"PG_USER" envDefault:"postgres"`
	PGPassword string `env:"PG_PASSWORD" envDefault:"postgres"`
	PGDatabase string `env:"PG_DATABASE" envDefault:"natsDB"`

	ServerPort string `env:"SERVER_PORT" envDefault:"3300"`
}

func GetDBConfig() (*DBConfig, error) {
	cfg := new(DBConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type JSConfig struct {
	URL string `default:"nats.DefaultURL"`
}

func GetJSConfig() *JSConfig {
	return new(JSConfig)
}
