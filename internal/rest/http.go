package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"html/template"
	"http-nats-psql/internal/models"
	"http-nats-psql/internal/utils"
)

const (
	responseKey = "response"
)

func PublishError(c *gin.Context, err error, code int) {
	if err != nil {
		_ = c.Error(err)
	}

	c.Data(code, "text/html", []byte(err.Error()))
}

func PublishOrder(c *gin.Context, data []byte, orderTemplate *template.Template) {
	model := &models.Order{}
	err := json.Unmarshal(data, model)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	orderTemplate.Execute(c.Writer, *model)
}
