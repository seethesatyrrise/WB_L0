package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/storage"
	"net/http"
)

type Rest struct {
	storage *storage.Storage
	db      *database.DB
}

func NewRest(storage *storage.Storage, db *database.DB) *Rest {
	return &Rest{storage: storage, db: db}
}

func (r *Rest) Register(api *gin.RouterGroup) {
	route := api.Group("/orders")
	{
		route.GET(":orderID", r.getOrderByID)
	}
}

func (r *Rest) getOrderByID(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("orderID")
	if orderID == "" {
		PublishError(c, errors.New("empty orderID"), http.StatusBadRequest)
		return
	}

	order, err := r.storage.GetOrderByID(orderID)
	if err == nil {
		PublishDataBytes(c, order)
		return
	}

	data, err := r.db.GetOrderByID(ctx, orderID)
	if err != nil {
		PublishError(c, err, http.StatusInternalServerError)
		return
	}

	PublishDataBytes(c, data)
}
