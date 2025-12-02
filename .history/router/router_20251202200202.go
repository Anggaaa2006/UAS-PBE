package router

import (
	"uas_pbe/controller"

	"github.com/gin-gonic/gin"
)

/*
	RegisterAchievementRoutes
	Fungsi ini dipanggil di main.go untuk mendaftarkan
	semua endpoint prestasi ke Gin router.
*/
func RegisterAchievementRoutes(r *gin.Engine, ctrl controller.AchievementController) {

	// Group khusus /achievements
	ach := r.Group("/achievements")

	{
		// Buat prestasi
		ach.POST("", ctrl.Create)

		// Ambil prestasi berdasarkan id
		ach.GET("/:id", ctrl.GetByID)

		// Update prestasi
		ach.PUT("/:id", ctrl.Update)

		// Soft delete prestasi
		ach.DELETE("/:id", ctrl.Delete)

		// Submit prestasi
		ach.POST("/:id/submit", ctrl.Submit)

		// Approve prestasi
		ach.POST("/:id/approve", ctrl.Approve)

		// Reject prestasi
		ach.POST("/:id/reject", ctrl.Reject)
	}
	

	// List prestasi berdasarkan mahasiswa
	r.GET("/students/:id/achievements", ctrl.ListByStudent)
}
