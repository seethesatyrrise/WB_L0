package storage

import (
	"context"
	"http-nats-psql/internal/models"
)

//go:generate mockgen -source $GOFILE -package $GOPACKAGE -destination mocks_test.go -mock_names=loader=MockLoader

type loader interface {
	GetAllData(ctx context.Context) ([]models.Restoration, error)
}
