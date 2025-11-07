package middleware

import (
	"net/http"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/config"
	"github.com/gin-gonic/gin"
)

func NoRouteMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		wildcard := ctx.Request.URL.Path
		ctx.JSON(http.StatusNotFound, config.Response{
			Success: false,
			Status:  "error",
			Data:    gin.H{"message": "please check your url", "path": wildcard},
		})

		ctx.Abort()
	}
}
