package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"uas_pbe/service"
	"uas_pbe/utils"
)

/*
	AchievementAttachmentController
	Menangani upload & lihat attachment prestasi
*/
type AchievementAttachmentController struct {
	svc *service.AchievementAttachmentService
}

func NewAchievementAttachmentController(
	s *service.AchievementAttachmentService,
) *AchievementAttachmentController {
	return &AchievementAttachmentController{svc: s}
}

/*
	POST /achievements/:id/attachments

	Upload file bukti prestasi
*/
func (c *AchievementAttachmentController) Upload(ctx *gin.Context) {
	achievementID := ctx.Param("id")

	file, err := ctx.FormFile("file")
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "file wajib diupload")
		return
	}

	if err := c.svc.Upload(ctx, achievementID, file); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, "attachment berhasil diupload")
}

/*
	GET /achievements/:id/attachments

	Lihat daftar attachment prestasi
*/
func (c *AchievementAttachmentController) List(ctx *gin.Context) {
	achievementID := ctx.Param("id")

	data, err := c.svc.List(ctx, achievementID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessData(ctx, data)
}
