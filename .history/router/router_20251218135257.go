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
	adminUserCtrl *controller.AdminUserController,
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
		ach.GET("", achCtrl.List)
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

// ================================
// DASHBOARD (ROLE BASED)
// ================================
// Dashboard digunakan untuk menampilkan ringkasan data
// berdasarkan peran (role) user yang login
//
// 26. Mahasiswa  → ringkasan prestasi milik sendiri
// 27. Dosen      → ringkasan status prestasi mahasiswa
// 28. Admin      → statistik global sistem
dashboard := r.Group(
	"/dashboard",
	middleware.JWTMiddleware(), // wajib login (JWT)
)
{
	// --------------------------------
	// 26. DASHBOARD MAHASISWA
	// --------------------------------
	// Endpoint:
	// GET /dashboard/student
	//
	// Hak Akses:
	// - Hanya MAHASISWA
	//
	// Fungsi:
	// - Menampilkan jumlah prestasi:
	//   draft, submitted, approved, rejected
	dashboard.GET(
		"/student",
		middleware.RoleStudent(),
		dashboardCtrl.Student,
	)

	// --------------------------------
	// 27. DASHBOARD DOSEN
	// --------------------------------
	// Endpoint:
	// GET /dashboard/lecturer
	//
	// Hak Akses:
	// - Hanya DOSEN
	//
	// Fungsi:
	// - Menampilkan jumlah prestasi mahasiswa
	//   yang perlu ditinjau (submitted, approved, rejected)
	dashboard.GET(
		"/lecturer",
		middleware.RoleLecturer(),
		dashboardCtrl.Lecturer,
	)

	// --------------------------------
	// 28. DASHBOARD ADMIN
	// --------------------------------
	// Endpoint:
	// GET /dashboard/admin
	//
	// Hak Akses:
	// - Hanya ADMIN
	//
	// Fungsi:
	// - Menampilkan statistik global sistem:
	//   total mahasiswa, total prestasi, status prestasi
	dashboard.GET(
		"/admin",
		middleware.RoleAdmin(),
		dashboardCtrl.Admin,
	)
}
)