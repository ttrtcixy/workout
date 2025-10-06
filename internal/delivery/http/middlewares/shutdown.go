package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShutdownMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Context().Err() != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "server is shutting down",
			})
			return
		}

		c.Next()
	}
}
