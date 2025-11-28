package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware mengecek apakah user punya role tertentu
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in token"})
			c.Abort()
			return
		}

		if role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
