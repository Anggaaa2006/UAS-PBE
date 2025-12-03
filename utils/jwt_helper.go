package utils

import (
	"github.com/gin-gonic/gin"
)

/*
	GetUserID
	Mengambil user_id dari JWT yang sudah disimpan
	oleh JWTMiddleware ke dalam context.
*/
func GetUserID(ctx *gin.Context) string {
	id, exists := ctx.Get("user_id")
	if !exists {
		return ""
	}
	return id.(string)
}

/*
	GetUserRole
	Mengambil role user dari JWT
	(role: student, lecturer, admin).
*/
func GetUserRole(ctx *gin.Context) string {
	role, exists := ctx.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}
