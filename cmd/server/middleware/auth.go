package middleware

import (
	"net/http"
	"strings"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtService *infrastructure.JWTService, cache *infrastructure.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrTokenMissing.Error()})
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
		refresh := c.Request.URL.Path == "/api/v1/auth/refresh"

		claims, err := jwtService.ValidateToken(token, refresh)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrTokenInvalid.Error()})
			c.Abort()
			return
		}

		prefix := "token:"
		if refresh {
			prefix = "refresh:"
		}
		key := prefix + claims.UserID.String()
		stored, err := cache.Get(c, key)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrTokenExpired.Error()})
			c.Abort()
			return
		}
		if stored != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrTokenExpired.Error()})
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
