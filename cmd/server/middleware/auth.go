package middleware

import (
	"net/http"
	"strings"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtService *infrastructure.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}
		token := parts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrTokenInvalid.Error()})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID.String())
		c.Set("userEmail", claims.Email)
		c.Set("isAdmin", claims.IsAdmin)
		c.Next()
	}
}

// AdminMiddleware checks if user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": domain.ErrForbidden.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
