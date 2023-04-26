package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"http-nats-psql/pkg/orderModel"
)
import "github.com/go-pg/pg"

type storageModel struct {
	Data string `sql:"select_all"`
}

type Storage struct {
	Data map[string]orderModel.Order
	DB   *pg.DB
}

func NewStorage(db *pg.DB) *Storage {
	return &Storage{Data: make(map[string]orderModel.Order), DB: db}
}

func (storage *Storage) RestoreData(db *pg.DB) {
	var res []storageModel
	queryRes, err := db.Query(&res,
		`select select_all();`,
	)
	if err != nil {
		panic(err)
	}

	if queryRes.RowsReturned() == 0 {
		fmt.Println("nothing to restore")
		return
	}

	parsedData := &orderModel.Order{}
	for _, data := range res {
		err := json.Unmarshal([]byte(data.Data), parsedData)
		if err != nil {
			panic(err)
		}
		storage.Data[parsedData.OrderUid] = *parsedData
	}
}

func (storage *Storage) InsertOrder(order *orderModel.Order) error {
	storage.Data[order.OrderUid] = *order
	
	_, err := storage.DB.Exec(
		`select insert_data(?);`,
		order)

	if err != nil {
		return err
	}

	return nil
}

func (storage *Storage) GetOrderByID(ctx context.Context, orderID string) (*orderModel.Order, error) {
	data, ok := storage.Data[orderID]
	if !ok {
		return nil, fmt.Errorf("order with id \"%s\" not found in bd", orderID)
	}
	return &data, nil
}
