package database

import (
	"github.com/caarlos0/env"
)

type Config struct {
	PGAddress  string `env:"PG_ADDRESS" envDefault:":5432"`
	PGUser     string `env:"PG_USER" envDefault:"postgres"`
	PGPassword string `env:"PG_PASSWORD" envDefault:"postgres"`
	PGDatabase string `env:"PG_DATABASE" envDefault:"natsDB"`

	ServerPort string `env:"SERVER_PORT" envDefault:"3300"`
}

func GetConfig() (*Config, error) {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
