package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse response standar untuk sukses
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}

// ErrorResponse response standar untuk error
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}
