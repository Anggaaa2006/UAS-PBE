package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	RoleMiddleware
	Middleware ini dipakai untuk membatasi akses berdasarkan role.

	Contoh:
	routes.Use(RoleMiddleware("student"))
	routes.Use(RoleMiddleware("lecture"))
*/
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Ambil role yang sudah disimpan di JWT middleware
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in token"})
			c.Abort()
			return
		}

		userRole := roleValue.(string)

		// Cek apakah role user termasuk ke role yang diizinkan
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				// role valid → lanjut
				c.Next()
				return
			}
		}

		// Role tidak cocok → blok
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access denied: insufficient permissions",
		})
		c.Abort()
	}
}
