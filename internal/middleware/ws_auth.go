package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"dnd-simulator/internal/auth"
)

// WebSocketAuthMiddleware handles authentication for WebSocket connections
// It checks both Authorization header and query parameter for the token
func WebSocketAuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// First, try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				// No Bearer prefix found, try query parameter
				tokenString = ""
			}
		}

		// If no token in header, try query parameter (for WebSocket connections)
		if tokenString == "" {
			tokenString = c.Query("token")
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
				c.Abort()
				return
			}
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}