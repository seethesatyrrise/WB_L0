package storage

import (
	"context"
	"errors"
	"sync"
)

type Storage struct {
	data map[string][]byte
	mu   *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{data: make(map[string][]byte), mu: &sync.RWMutex{}}
}

func (storage *Storage) RestoreData(ctx context.Context, l loader) error {
	res, err := l.GetAllData(ctx)
	if err != nil {
		return err
	}

	storage.mu.Lock()
	defer storage.mu.Unlock()
	for _, order := range res {
		storage.data[order.ID] = order.Data
	}
	return nil
}

func (storage *Storage) InsertOrder(id string, data []byte) error {
	storage.mu.Lock()
	storage.data[id] = data
	storage.mu.Unlock()

	return nil
}

func (storage *Storage) GetOrderByID(id string) ([]byte, error) {
	storage.mu.RLock()
	data, ok := storage.data[id]
	storage.mu.RUnlock()
	if ok {
		return data, nil
	}

	return nil, errors.New("no data in storage")
}
