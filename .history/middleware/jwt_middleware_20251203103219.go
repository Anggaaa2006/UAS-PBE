package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
	JWTMiddleware
	==========================
	Middleware ini digunakan untuk setiap endpoint
	yang membutuhkan login.
	
	Flow:
	1. Ambil token dari header Authorization (Bearer token)
	2. Validasi token pakai ValidateToken()
	3. Simpan user_id dan role ke context
	4. Lanjutkan request ke controller berikutnya
*/
func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Ambil header Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing Authorization header",
			})
			ctx.Abort()
			return
		}

		// Format header harus: Bearer token
		split := strings.Split(authHeader, " ")
		if len(split) != 2 || strings.ToLower(split[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			ctx.Abort()
			return
		}

		tokenString := split[1]

		// Validasi token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}

		// Simpan user_id dan role ke Gin context
		ctx.Set("user_id", claims["user_id"])
		ctx.Set("role", claims["role"])

		// Lanjutkan
		ctx.Next()
	}
}
git