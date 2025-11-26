package router

import (
	"github.com/gin-gonic/gin"

	"UAS PBE/controller"
	"UAS PBE/middleware"
)

func NewRouter(
	auth *controller.AuthController,
	ach *controller.AchievementController,
	user *controller.UserController,
	student *controller.StudentController,
	report *controller.ReportController,
) *gin.Engine {

	r := gin.Default()

	r.POST("/auth/login", auth.Login)

	api := r.Group("/api/v1")
	api.Use(middleware.JWT())

	api.GET("/achievements", ach.List)
	api.POST("/achievements", ach.Create)
	api.PUT("/achievements/:id", ach.Update)
	api.DELETE("/achievements/:id", ach.SoftDelete)
	api.POST("/achievements/:id/submit", ach.Submit)
	api.POST("/achievements/:id/verify", ach.Verify)
	api.POST("/achievements/:id/reject", ach.Reject)

	return r
}
