package middleware

import (
	"net/http"
	"strings"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/config"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/utils"
	"github.com/gin-gonic/gin"
)

// CustomResponseMiddleware formats your JSON responses uniformly
func CustomResponseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// Check if response is already written
		if ctx.IsAborted() {
			return
		}

		// Decide your custom response structure
		if len(ctx.Errors) > 0 {
			// Collect all error messages from ctx.Errors
			var errorMessages []string
			for _, ginErr := range ctx.Errors {
				errorMessages = append(errorMessages, ginErr.Error())
			}

			statusCode := http.StatusInternalServerError
			if ctx.Writer.Status() != http.StatusOK {
				statusCode = ctx.Writer.Status()
			}

			// Aggregate all messages into one error message
			message := config.APIError{
				Code:    statusCode,
				Message: strings.Join(errorMessages, "; "),
			}

			// Now check if custom error data exists in context and override handling
			if val, exists := ctx.Get(utils.ErrorContextKey); exists {
				switch v := val.(type) {
				case string:
					message = config.APIError{
						Code:    statusCode,
						Message: v,
					}
				case config.APIError:
					message = v
				}
			}

			ctx.JSON(statusCode, config.Response{
				Success: false,
				Status:  "error",
				Data:    message,
				Error:   ctx.Errors[0],
			})
			ctx.Abort()
			return
		}

		// For successful responses
		var responseData interface{}
		if v, exists := ctx.Get(utils.ResponseDataKey); exists {
			responseData = v
		}
		ctx.JSON(ctx.Writer.Status(), config.Response{
			Success: true,
			Status:  "success",
			Data:    responseData,
		})
	}
}
