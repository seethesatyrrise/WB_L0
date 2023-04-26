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
	var res []byte
	var mod orderModel.OrderSql
	queryRes, err := r.db.QueryOne(&mod,
		`select * from select_data(?);`,
		orderID)
	fmt.Println(mod)
	if err != nil {
		panic(err)
	}

	if queryRes.RowsReturned() == 0 {
		return nil, fmt.Errorf("order with id \"%s\" not found in bd", orderID)
	}

	if err != nil {
		return nil, err
	}

	order := &orderModel.Order{}
	err = json.Unmarshal(res, order)
	if err != nil {
		panic(err)
	}

	return order, nil
}
