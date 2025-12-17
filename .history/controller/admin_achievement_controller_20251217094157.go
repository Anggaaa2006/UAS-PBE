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
func NewAdminAchievementController(s *service.AdminAchievementService) *AdminAchievementController {
	return &AdminAchievementController{svc: s}
}

/*
	GET /admin/achievements

	Fungsi:
	- Menampilkan semua prestasi mahasiswa
	- Bisa difilter (status, student_id)
	- Bisa pagination
*/
func (c *AdminAchievementController) ListAll(ctx *gin.Context) {

	// Ambil query parameter
	params := map[string]string{
		"status":     ctx.Query("status"),
		"student_id": ctx.Query("student_id"),
		"page":       ctx.Query("page"),
		"limit":      ctx.Query("limit"),
	}

	data, total, err := c.svc.ListAll(ctx, params)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessData(ctx, gin.H{
		"total": total,
		"data":  data,
	})
}
