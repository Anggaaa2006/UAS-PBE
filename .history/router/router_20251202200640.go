package router

import (
	"uas_pbe/controller"
	"uas_pbe/middleware"

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
		// ============================
		// STUDENT ONLY
		// ============================

		// Buat prestasi
		ach.POST("",
			middleware.JWTAuth(),                   // wajib login
			middleware.RoleMiddleware("student"),   // hanya student
			ctrl.Create,
		)

		// Submit prestasi
		ach.POST("/:id/submit",
			middleware.JWTAuth(),
			middleware.RoleMiddleware("student"),
			ctrl.Submit,
		)

		// ============================
		// DOSEN / LECTURE ONLY
		// ============================

		// Approve prestasi
		ach.POST("/:id/approve",
			middleware.JWTAuth(),
			middleware.RoleMiddleware("dosen", "lecture"),
			ctrl.Approve,
		)

		// Reject prestasi
		ach.POST("/:id/reject",
			middleware.JWTAuth(),
			middleware.RoleMiddleware("dosen", "lecture"),
			ctrl.Reject,
		)

		// ============================
		// PUBLIC (TAPI HARUS LOGIN)
		// ============================

		// Ambil prestasi berdasarkan ID
		ach.GET("/:id",
			middleware.JWTAuth(), // semua yang login boleh
			ctrl.GetByID,
		)

		// Update prestasi
		ach.PUT("/:id",
			middleware.JWTAuth(), // semua role boleh update? nanti bisa diatur
			ctrl.Update,
		)

		// Soft delete prestasi
		ach.DELETE("/:id",
			middleware.JWTAuth(),
			ctrl.Delete,
		)
	}

	// ===== List prestasi berdasarkan mahasiswa =====
	// Semua role boleh asal login
	r.GET("/students/:id/achievements",
		middleware.JWTAuth(),
		ctrl.ListByStudent,
	)
}
