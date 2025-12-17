package router

import (
	"github.com/gin-gonic/gin"

	"uas_pbe/controller"
	"uas_pbe/middleware"
)

/*
	RegisterRoutes
	Fungsi utama untuk mendaftarkan semua endpoint.
	Dipanggil dari main.go
*/
func RegisterRoutes(
	r *gin.Engine,
	authCtrl *controller.AuthController,
	achCtrl *controller.AchievementController,
	statsCtrl *controller.StatsController,
) {

	// ================================
	// AUTH (Public)
	// ================================
	r.POST("/auth/login", authCtrl.Login)

	// ================================
	// AUTH (Protected)
	// ================================
	auth := r.Group("/auth", middleware.JWTMiddleware())
	{
		// Ambil profile user dari JWT
		auth.GET("/profile", authCtrl.Profile)
	}

	// ================================
	// ACHIEVEMENT (Protected)
	// ================================
	ach := r.Group("/achievements", middleware.JWTMiddleware())
	{
		// ================================
		// MAHASISWA
		// ================================
		ach.POST("", middleware.RoleStudent(), achCtrl.Create)
		ach.PUT("/:id", middleware.RoleStudent(), achCtrl.Update)
		ach.POST("/:id/submit", middleware.RoleStudent(), achCtrl.Submit)
		ach.DELETE("/:id", middleware.RoleStudent(), achCtrl.Delete)

		// ================================
		// DOSEN WALI
		// ================================
		ach.POST("/:id/approve", middleware.RoleLecturer(), achCtrl.Approve)
		ach.POST("/:id/reject", middleware.RoleLecturer(), achCtrl.Reject)

		// ================================
		// LIHAT DETAIL
		// ================================
		ach.GET("/:id", achCtrl.GetByID)
	}

	// ================================
	// LIST PRESTASI MAHASISWA
	// ================================
	r.GET(
		"/students/:id/achievements",
		middleware.JWTMiddleware(),
		achCtrl.ListByStudent,
	)

	// ================================
	// STATISTIK PRESTASI (FR-011)
	// ================================
	stats := r.Group("/stats", middleware.JWTMiddleware())
	{
		stats.GET("/achievements", statsCtrl.GetAchievementStats)
	}
}
	// ================================
	// REPORTS PER MAHASISWA (FR-011)
	// ================================
	reports := r.Group("/reports", middleware.JWTMiddleware())
	{
		reports.GET(
			"/student/:id",
			statsCtrl.GetStudentStats,
		)
	}
