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
	adminAchCtrl *controller.AdminAchievementController, // ← pastikan ada
	adminUserCtrl *controller.AdminUserController
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
		auth.GET("/profile", authCtrl.Profile)

		// ✅ TAMBAHAN
		auth.POST("/refresh", authCtrl.Refresh)
		auth.POST("/logout", authCtrl.Logout)
	}

	// ================================
	// ACHIEVEMENT (Protected)
	// ================================
	ach := r.Group("/achievements", middleware.JWTMiddleware())
	{
		ach.POST("", middleware.RoleStudent(), achCtrl.Create)
		ach.PUT("/:id", middleware.RoleStudent(), achCtrl.Update)
		ach.POST("/:id/submit", middleware.RoleStudent(), achCtrl.Submit)
		ach.DELETE("/:id", middleware.RoleStudent(), achCtrl.Delete)

		ach.POST("/:id/approve", middleware.RoleLecturer(), achCtrl.Approve)
		ach.POST("/:id/reject", middleware.RoleLecturer(), achCtrl.Reject)

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
	// STATISTIK GLOBAL (FR-011)
	// ================================
	stats := r.Group("/stats", middleware.JWTMiddleware())
	{
		stats.GET("/achievements", statsCtrl.GetAchievementStats)
	}

	// ================================
	// REPORTS PER MAHASISWA (FR-011)
	// ================================
	reports := r.Group("/reports", middleware.JWTMiddleware())
	{
		reports.GET("/student/:id", statsCtrl.GetStudentStats)
	}

	// ================================
	// ADMIN – VIEW ALL ACHIEVEMENTS (FR-010)
	// ================================
	admin := r.Group(
		"/admin",
		middleware.JWTMiddleware(),
		middleware.RoleAdmin(),
	)
	{
		// ADMIN – USERS (FR-009)
    	admin.POST("/users", adminUserCtrl.Create)
    	admin.GET("/users", adminUserCtrl.List)
    	admin.GET("/users/:id", adminUserCtrl.GetByID)
    	admin.PUT("/users/:id", adminUserCtrl.Update)
    	admin.DELETE("/users/:id", adminUserCtrl.Delete)
    	admin.PUT("/users/:id/role", adminUserCtrl.UpdateRole)
		// ADMIN – ACHIEVEMENTS (FR-010)
		admin.GET("/achievements", adminAchCtrl.ListAll)
	}
}
