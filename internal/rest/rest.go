package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"html/template"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/storage"
	"http-nats-psql/internal/utils"
	"net/http"
)

type Rest struct {
	storage       *storage.Storage
	db            *database.DB
	OrderTemplate *template.Template
}

func NewRest(storage *storage.Storage, db *database.DB) *Rest {
	return &Rest{storage: storage, db: db, OrderTemplate: utils.OrderTemplate()}
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
		PublishOrder(c, order, r.OrderTemplate)
		return
	}
	utils.Logger.Error(err.Error())

	data, err := r.db.GetOrderByID(ctx, orderID)
	if err != nil {
		PublishError(c, err, http.StatusInternalServerError)
		return
	}

	PublishOrder(c, data, r.OrderTemplate)
}
