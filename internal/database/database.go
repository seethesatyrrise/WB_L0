package database

import (
	"fmt"
	"github.com/go-pg/pg"
	"http-nats-psql/internal"
	"http-nats-psql/internal/models"
)

type DB struct {
	Handler *pg.DB
}

func NewDatabaseConnection(c *internal.DBConfig) (*DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     c.PGAddress,
		User:     c.PGUser,
		Password: c.PGPassword,
		Database: c.PGDatabase,
	})
	return &DB{Handler: db}, nil
}

func (db *DB) InsertOrder(data *models.Order) error {
	_, err := db.Handler.Exec(
		`select * from insert_data(?);`,
		data)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetOrderByID(id string) (*models.Order, error) {
	var mod models.SelectData
	_, err := db.Handler.QueryOne(&mod,
		`select * from select_data(?);`,
		id)
	if err != nil {
		panic(err)
	}

	if mod.Data.OrderUid == "" {
		return nil, fmt.Errorf("order with id \"%s\" not found in bd", id)
	}

	return &mod.Data, nil
}

func (db *DB) GetAllData() []models.Restoration {
	var mod []models.Restoration
	queryRes, err := db.Handler.Query(&mod,
		`select * from select_all();`,
	)
	if err != nil {
		panic(err)
	}

	if queryRes.RowsReturned() == 0 {
		fmt.Println("nothing to restore")
		return nil
	}

	return mod
}
