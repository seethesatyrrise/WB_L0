package rest

import (
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

func PublishData(c *gin.Context, data interface{}) {
	c.Set(responseKey, data)

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

//
//func GetResponse(c *gin.Context) (interface{}, bool) {
//	return c.Get(responseKey)
//}
//
//func IsHTTPStatusSuccess(code int) bool {
//	return code >= http.StatusOK && code < http.StatusMultipleChoices
//}
