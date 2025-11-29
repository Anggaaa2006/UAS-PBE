package controller

import (
	"/internal/service"
	"UAS PBE/utils"
	"github.com/gin-gonic/gin"
)

// ReportController menangani laporan (misal total prestasi)
type ReportController struct {
	service *service.ReportService
}

func NewReportController(s *service.ReportService) *ReportController {
	return &ReportController{service: s}
}

// GET /reports/achievement
func (c *ReportController) AchievementReport(ctx *gin.Context) {
	data, err := c.service.AchievementReport(ctx)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}
	utils.SuccessData(ctx, data)
}
