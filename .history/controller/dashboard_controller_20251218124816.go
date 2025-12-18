package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"uas_pbe/service"
)

/*
	DashboardController
	Endpoint dashboard berdasarkan role user

	SRS:
	- FR-012: Dashboard Mahasiswa
	- FR-013: Dashboard Dosen
	- FR-014: Dashboard Admin
*/
type DashboardController struct {
	achSvc   *service.AchievementService
	statsSvc *service.StatsService
}

func NewDashboardController(
	ach *service.AchievementService,
	stats *service.StatsService,
) *DashboardController {
	return &DashboardController{
		achSvc:   ach,
		statsSvc: stats,
	}
}

/*
	GET /dashboard/student
	Mahasiswa melihat ringkasan prestasinya
*/
func (c *DashboardController) Student(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	data, err := c.achSvc.GetStudentSummary(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

/*
	GET /dashboard/lecturer
	Dosen melihat status prestasi mahasiswa
*/
func (c *DashboardController) Lecturer(ctx *gin.Context) {
	data, err := c.achSvc.GetLecturerSummary(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

/*
	GET /dashboard/admin
	Admin melihat statistik global sistem
*/
func (c *DashboardController) Admin(ctx *gin.Context) {
	data, err := c.statsSvc.GetGlobalDashboard(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
