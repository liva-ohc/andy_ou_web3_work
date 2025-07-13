package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"code":    status,
		"message": message,
		"data":    data,
	})
}

func Success(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, "success", data)
}

func Fail(c *gin.Context, message string) {
	Response(c, http.StatusInternalServerError, message, nil)
}
