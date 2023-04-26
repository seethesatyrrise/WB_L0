package rest

import (
	"encoding/json"
	"fmt"
	"http-nats-psql/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	responseKey = "response"
)

func PublishError(c *gin.Context, err error, code int) {
	if err != nil {
		_ = c.Error(err)
	}

	c.JSON(code, gin.H{
		"error": err.Error(),
		"code":  code,
	})
}

func PublishData(c *gin.Context, data *models.Order) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Data(http.StatusOK, "application/json", dataJSON)
}
