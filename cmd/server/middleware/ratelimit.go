package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware implements rate limiting using Redis
func RateLimitMiddleware(cache *infrastructure.Cache, maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.ClientIP()
		if userID := c.GetString("userID"); userID != "" {
			identifier = fmt.Sprintf("user:%s", userID)
		}
		allowed, err := cache.CheckRateLimit(c.Request.Context(), identifier, maxRequests, window)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit check failed"})
			c.Abort()
			return
		}
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"message": fmt.Sprintf("maximum %d requests per %v allowed", maxRequests, window),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
