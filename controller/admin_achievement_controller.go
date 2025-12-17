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
*/
func (c *AdminAchievementController) ListAll(ctx *gin.Context) {

	data, err := c.svc.ListAll(ctx.Request.Context())
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessData(ctx, gin.H{
		"data": data,
	})
}
