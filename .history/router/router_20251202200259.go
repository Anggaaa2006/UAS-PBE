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
	
		// student → create, submit
	ach.POST("", middleware.JWTAuth(), middleware.RoleMiddleware("student"), ctrl.Create)
	ach.POST("/:id/submit", middleware.JWTAuth(), middleware.RoleMiddleware("student"), ctrl.Submit)

	// lecture/dosen → approve/reject
	ach.POST("/:id/approve", middleware.JWTAuth(), middleware.RoleMiddleware("lecture", "dosen"), ctrl.Approve)
	ach.POST("/:id/reject", middleware.JWTAuth(), middleware.RoleMiddleware("lecture", "dosen"), ctrl.Reject)

	// semua role boleh
	ach.GET("/:id", middleware.JWTAuth(), ctrl.GetByID)

	// List prestasi berdasarkan mahasiswa
	r.GET("/students/:id/achievements", ctrl.ListByStudent)
}