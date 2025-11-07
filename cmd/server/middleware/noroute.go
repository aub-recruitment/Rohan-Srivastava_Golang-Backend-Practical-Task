package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRouteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		wildcard := c.Request.URL.Path
		c.JSON(http.StatusNotFound, gin.H{"message": "please check your url", "path": wildcard})

		c.Abort()
	}
}
