package rest

import (
	"errors"
	"http-nats-psql/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type rest struct {
	service pkg.Repo
}

func NewRest(service pkg.Repo) *rest {
	return &rest{service}
}

func (r *rest) Register(api *gin.RouterGroup) {
	route := api.Group("/orders")
	{
		route.GET(":orderID", r.getOrderByID)
		route.GET("echo", func(c *gin.Context) {
			PublishData(c, "echo")
		})
	}
}

func (r *rest) getOrderByID(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("orderID")
	if orderID == "" {
		PublishError(c, errors.New("empty orderID"), http.StatusBadRequest)
		return
	}

	order, err := r.service.GetOrderByID(ctx, orderID)
	if err != nil {
		PublishError(c, err, http.StatusInternalServerError)
		return
	}

	PublishData(c, order)
}