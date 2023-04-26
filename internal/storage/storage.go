package storage

import (
	"errors"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/models"
	"sync"
)

type Storage struct {
	data map[string]models.Order
	mu   *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{data: make(map[string]models.Order), mu: &sync.RWMutex{}}
}

func (storage *Storage) RestoreData(db *database.DB) {
	res := db.GetAllData()

	for _, order := range res {
		storage.data[order.ID] = order.Data
	}
}

func (storage *Storage) InsertOrder(data *models.Order) error {
	storage.mu.Lock()
	storage.data[data.OrderUid] = *data
	storage.mu.Unlock()

	return nil
}

func (storage *Storage) GetOrderByID(id string) (*models.Order, error) {
	storage.mu.RLock()
	data, ok := storage.data[id]
	storage.mu.RUnlock()
	if ok {
		return &data, nil
	}

	return nil, errors.New("no data in storage")
}
