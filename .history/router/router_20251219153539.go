package router

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	adminAchCtrl *controller.AdminAchievementController,
	adminUserCtrl *controller.AdminUserController,
	dashboardCtrl *controller.DashboardController, // ✅ TAMBAHAN
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
		auth.POST("/refresh", authCtrl.Refresh)
		auth.POST("/logout", authCtrl.Logout)
	}

	// ================================
	// ACHIEVEMENT
	// ================================
	ach := r.Group("/achievements", middleware.JWTMiddleware())
	{
		// Mahasiswa
		ach.POST("", middleware.RoleStudent(), achCtrl.Create)
		ach.PUT("/:id", middleware.RoleStudent(), achCtrl.Update)
		ach.POST("/:id/submit", middleware.RoleStudent(), achCtrl.Submit)
		ach.DELETE("/:id", middleware.RoleStudent(), achCtrl.Delete)

		// Dosen
		ach.POST("/:id/approve", middleware.RoleLecturer(), achCtrl.Approve)
		ach.POST("/:id/reject", middleware.RoleLecturer(), achCtrl.Reject)

		// Umum
		ach.GET("", achCtrl.List)
		ach.GET("/:id", achCtrl.GetByID)
		ach.GET("/:id/history", achCtrl.History)
		ach.POST("/:id/attachments", achCtrl.UploadAttachment)
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
	// ADMIN
	// ================================
	admin := r.Group(
		"/admin",
		middleware.JWTMiddleware(),
		middleware.RoleAdmin(),
	)
	{
		// ADMIN USER (FR-009)
		admin.POST("/users", adminUserCtrl.Create)
		admin.GET("/users", adminUserCtrl.List)
		admin.GET("/users/:id", adminUserCtrl.GetByID)
		admin.PUT("/users/:id", adminUserCtrl.Update)
		admin.DELETE("/users/:id", adminUserCtrl.Delete)
		admin.PUT("/users/:id/role", adminUserCtrl.UpdateRole)

		// ADMIN ACHIEVEMENT (FR-010)
		admin.GET("/achievements", adminAchCtrl.ListAll)
	}

	// ================================
	// DASHBOARD (ROLE BASED) 26–28
	// ================================
	dashboard := r.Group("/dashboard", middleware.JWTMiddleware())
	{
		// 26. Dashboard Mahasiswa
		dashboard.GET(
			"/student",
			middleware.RoleStudent(),
			dashboardCtrl.Student,
		)

		// 27. Dashboard Dosen
		dashboard.GET(
			"/lecturer",
			middleware.RoleLecturer(),
			dashboardCtrl.Lecturer,
		)

		// 28. Dashboard Admin
		dashboard.GET(
			"/admin",
			middleware.RoleAdmin(),
			dashboardCtrl.Admin,
		)
		// ================================
// SWAGGER DOCUMENTATION
// ================================
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	}
}
