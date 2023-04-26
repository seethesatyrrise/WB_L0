package pkg

import (
	"context"
	"http-nats-psql/pkg/orderModel"
)

type Repo interface {
	GetOrderByID(ctx context.Context, orderID string) (*orderModel.Order, error)
}
