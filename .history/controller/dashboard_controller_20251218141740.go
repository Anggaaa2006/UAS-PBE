package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"uas_pbe/service"
	"uas_pbe/utils"
)

/*
	DashboardController
	Menangani dashboard berdasarkan role user
*/
type DashboardController struct {
	achSvc   *service.AchievementService
	statsSvc service.StatsService
}

func NewDashboardController(
	achSvc *service.AchievementService,
	statsSvc service.StatsService,
) *DashboardController {
	return &DashboardController{
		achSvc:   achSvc,
		statsSvc: statsSvc,
	}
}

/*
	GET /dashboard/student
	Dashboard mahasiswa
*/
func (c *DashboardController) Student(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	data, err := c.achSvc.GetStudentSummary(
		ctx.Request.Context(), // ✅ FIX
		userID,
	)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessData(ctx, data)
}

/*
	GET /dashboard/lecturer
	Dashboard dosen wali
*/
func (c *DashboardController) Lecturer(ctx *gin.Context) {
	data, err := c.achSvc.GetLecturerSummary(
		ctx.Request.Context(), // ✅ FIX
	)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessData(ctx, data)
}

/*
	GET /dashboard/admin
	Dashboard admin
*/
func (c *DashboardController) Admin(ctx *gin.Context) {
	data, err := c.statsSvc.GetAchievementStats(
		ctx.Request.Context(), // ✅ FIX
	)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessData(ctx, data)
}
