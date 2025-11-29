package router

import (
	"github.com/gin-gonic/gin"

	"uas_pbe/controller"
	"uas_pbe/middleware"
)

// NewRouter membuat dan mengembalikan instance Gin Engine
// Di sini kita mendefinisikan semua endpoint API dan middleware yang digunakan
func NewRouter(
	auth *controller.AuthController,
	ach *controller.AchievementController,
	user *controller.UserController,
	student *controller.StudentController,
	report *controller.ReportController,
) *gin.Engine {

	// Membuat instance router Gin baru
	r := gin.Default()

	// Endpoint login (public route, tidak pakai middleware)
	r.POST("/auth/login", auth.Login)

	// Semua route dalam group ini akan menggunakan JWT middleware
	api := r.Group("/api/v1")
	api.Use(middleware.JWT()) // Protect all routes with JWT

	// ==============================
	// ACHIEVEMENT ROUTES
	// ==============================

	// GET /achievements  -> ambil semua prestasi mahasiswa
	api.GET("/achievements", ach.List)

	// POST /achievements -> mahasiswa membuat prestasi baru (status draft)
	api.POST("/achievements", ach.Create)

	// PUT /achievements/:id -> update prestasi (hanya bisa jika status draft)
	api.PUT("/achievements/:id", ach.Update)

	// DELETE /achievements/:id -> soft delete prestasi (status jadi "deleted")
	api.DELETE("/achievements/:id", ach.SoftDelete)

	// POST /achievements/:id/submit -> mahasiswa submit prestasi (status: submitted)
	api.POST("/achievements/:id/submit", ach.Submit)

	// POST /achievements/:id/verify -> dosen verifikasi prestasi (status: verified)
	api.POST("/achievements/:id/verify", ach.Verify)

	// POST /achievements/:id/reject -> dosen menolak prestasi (status: rejected)
	api.POST("/achievements/:id/reject", ach.Reject)

	return r
}
