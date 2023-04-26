package database

import (
	"github.com/go-pg/pg"
)

func NewDatabaseConnection(c *Config) (*pg.DB, error) {
	return pg.Connect(&pg.Options{
		Addr:     c.PGAddress,
		User:     c.PGUser,
		Password: c.PGPassword,
		Database: c.PGDatabase,
	}), nil

}
