package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Requested-With") != "XMLHttpRequest" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid X-Requested-With header",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
