package router

import (
	"github.com/gin-gonic/gin"

	"uas_pbe/controller"
	"uas_pbe/middleware"
)

/*
	RegisterRoutes
	Fungsi utama untuk mendaftarkan semua endpoint.
	Ini dipanggil dari main.go
*/
func RegisterRoutes(
	r *gin.Engine,
	authCtrl *controller.AuthController,
	achCtrl *controller.AchievementController,
) {

	// ================================
	// AUTH (Public)
	// ================================
	r.POST("/auth/login", authCtrl.Login)

	// ================================
	// ACHIEVEMENT (Protected)
	// Semua endpoint harus lewat JWT middleware
	// ================================
	ach := r.Group("/achievements", middleware.JWTMiddleware())

	{
		// STUDENT ONLY
		ach.POST("", middleware.RoleStudent(), achCtrl.Create)
		ach.PUT("/:id", middleware.RoleStudent(), achCtrl.Update)
		ach.POST("/:id/submit", middleware.RoleStudent(), achCtrl.Submit)
		ach.DELETE("/:id", middleware.RoleStudent(), achCtrl.Delete)

		// LECTURER ONLY
		ach.POST("/:id/approve", middleware.RoleLecturer(), achCtrl.Approve)
		ach.POST("/:id/reject", middleware.RoleLecturer(), achCtrl.Reject)

		// Student OR Lecturer boleh lihat detail prestasi
		ach.GET("/:id", achCtrl.GetByID)
	}

	// ================================
	// LIST PRESTASI BERDASARKAN MAHASISWA
	// Endpoint terbuka untuk student/lecturer
	// ================================
	r.GET("/students/:id/achievements", middleware.JWTMiddleware(), achCtrl.ListByStudent)
}
