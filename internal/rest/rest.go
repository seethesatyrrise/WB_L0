package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/storage"
	"net/http"
)

type rest struct {
	storage *storage.Storage
	db      *database.DB
}

func NewRest(storage *storage.Storage, db *database.DB) *rest {
	return &rest{storage: storage, db: db}
}

func (r *rest) Register(api *gin.RouterGroup) {
	route := api.Group("/orders")
	{
		route.GET(":orderID", r.getOrderByID)
	}
}

func (r *rest) getOrderByID(c *gin.Context) {
	//ctx := c.Request.Context()
	orderID := c.Param("orderID")
	if orderID == "" {
		PublishError(c, errors.New("empty orderID"), http.StatusBadRequest)
		return
	}

	order, err := r.storage.GetOrderByID(orderID)
	if err == nil {
		PublishData(c, order)
		return
	}

	order, err = r.db.GetOrderByID(orderID)
	if err != nil {
		PublishError(c, err, http.StatusInternalServerError)
		return
	}

	PublishData(c, order)
}
