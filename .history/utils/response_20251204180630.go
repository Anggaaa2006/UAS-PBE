package utils

import (
	"github.com/gin-gonic/gin"
)

/*
	Error response standar
	Contoh output:
	{
		"success": false,
		"message": "gagal mengambil data"
	}
*/
func Error(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

/*
	Success (pakai message saja)
	Contoh output:
	{
		"success": true,
		"message": "berhasil"
	}
*/
func Success(ctx *gin.Context, message string) {
	ctx.JSON(200, gin.H{
		"success": true,
		"message": message,
	})
}

/*
	SuccessData (pakai data)
	Contoh output:
	{
		"success": true,
		"data": { ... }
	}
*/
func SuccessData(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}
