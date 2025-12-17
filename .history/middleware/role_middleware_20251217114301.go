package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	GetRoleFromContext
	Mengambil role dari JWT middleware yang sebelumnya
	menyimpan role ke context.
*/
func getRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}

/*
	Middleware Role Student
	Hanya role: student
*/
func RoleStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := getRoleFromContext(c)

		if role != "student" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "akses hanya untuk mahasiswa",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

/*
	Middleware Role Lecturer
	Hanya role: lecturer
*/
func RoleLecturer() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := getRoleFromContext(c)

		if role != "lecturer" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "akses hanya untuk dosen",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
/*
	RoleAdmin
	Middleware untuk memastikan user adalah ADMIN

	Digunakan oleh:
	- Endpoint /admin/*
*/
func RoleAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Ambil role dari context (di-set oleh JWTMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "role tidak ditemukan",
			})
			return
		}

		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "akses khusus admin",
			})
			return
		}

		c.Next()
	}
}