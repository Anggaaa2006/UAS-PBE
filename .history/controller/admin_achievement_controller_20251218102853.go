package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"uas_pbe/service"
	"uas_pbe/utils"
)

/*
	AdminAchievementController
	Controller khusus ADMIN untuk melihat semua prestasi
*/
type AdminAchievementController struct {
	svc *service.AdminAchievementService
}

/*
	Constructor
*/
func NewAdminAchievementController(
	s *service.AdminAchievementService,
) *AdminAchievementController {
	return &AdminAchievementController{svc: s}
}

/*
	GET /admin/achievements

	Fitur:
	- Filter status
	- Filter student_id
	- Pagination (page, limit)

	KETERKAITAN SRS:
	- FR-010 View All Achievements
*/
func (c *AdminAchievementController) ListAll(ctx *gin.Context) {

	// =========================
	// Ambil query parameter
	// =========================
	params := map[string]string{
		"status":     ctx.Query("status"),
		"student_id": ctx.Query("student_id"),
		"page":       ctx.DefaultQuery("page", "1"),
		"limit":      ctx.DefaultQuery("limit", "10"),
	}

	// =========================
	// Panggil service
	// =========================
	data, total, err := c.svc.ListAll(
		ctx.Request.Context(),
		params,
	)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// =========================
	// Response
	// =========================
	utils.SuccessData(ctx, gin.H{
		"total": total,
		"data":  data,
	})
}
