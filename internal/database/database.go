package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"http-nats-psql/internal/models"
	"http-nats-psql/internal/utils"
	"time"
)

type DB struct {
	Handler *pgx.Conn
}

func NewDatabaseConnection(c *DBConfig) (*DB, error) {
	db, err := pgx.Connect(pgx.ConnConfig{
		Host:     c.PGAddress,
		User:     c.PGUser,
		Password: c.PGPassword,
		Database: c.PGDatabase,
	})
	if err != nil {
		return nil, err
	}

	return &DB{Handler: db}, nil
}

func (db *DB) InsertOrder(ctx context.Context, data []byte) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var id string
	row := db.Handler.QueryRowEx(ctx,
		`select * from insert_data($1);`,
		nil, data)

	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return "", err
}

func (db *DB) GetOrderByID(ctx context.Context, id string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var data []byte
	row := db.Handler.QueryRowEx(ctx,
		`select * from select_data($1);`,
		nil, id)

	err := row.Scan(&data)
	if err != nil {
		return nil, err
	}
	if data == nil {
		utils.Logger.Info("rows with id: " + id + " not found")
		return nil, fmt.Errorf("rows with id: %s not found", id)
	}

	return data, nil
}

func (db *DB) GetAllData(ctx context.Context) ([]models.Restoration, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var mod []models.Restoration
	rows, err := db.Handler.QueryEx(ctx,
		`select * from select_all();`, nil,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := new(models.Restoration)
		if err := rows.Scan(&r.ID, &r.Data); err != nil {
			return nil, err
		}
		mod = append(mod, *r)
	}

	return mod, nil
}
