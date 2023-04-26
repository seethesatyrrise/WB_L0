package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg"
	"http-nats-psql/pkg/orderModel"
)

type repo struct {
	db *pg.DB
}

func NewRepo(db *pg.DB) *repo {
	return &repo{db}
}

func (r *repo) GetOrderByID(ctx context.Context, orderID string) (*orderModel.Order, error) {
	//var res []byte
	var mod orderModel.OrderSql
	_, err := r.db.Query(&mod,
		`select * from select_data(?);`,
		orderID)
	if err != nil {
		panic(err)
	}

	if mod.Data == "" {
		return nil, fmt.Errorf("order with id \"%s\" not found in bd", orderID)
	}

	if err != nil {
		return nil, err
	}

	order := &orderModel.Order{}
	err = json.Unmarshal([]byte(mod.Data), order)
	if err != nil {
		panic(err)
	}

	return order, nil
}
