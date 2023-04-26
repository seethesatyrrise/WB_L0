package rest

import (
	"encoding/json"
	"http-nats-psql/internal/models"
	"http-nats-psql/internal/utils"
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
		utils.Logger.Error(err.Error())
		return
	}
	c.Data(http.StatusOK, "application/json", dataJSON)
}

func PublishDataBytes(c *gin.Context, data []byte) {
	c.Data(http.StatusOK, "application/json", data)
}
