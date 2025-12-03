package controller

import (
	"net/http"

	"uas_pbe/model"
	"uas_pbe/service"
	"uas_pbe/utils"

	"github.com/gin-gonic/gin"
)

// AchievementController menangani request prestasi
type AchievementController struct {
	svc *service.AchievementService
}

// Constructor
func NewAchievementController(s *service.AchievementService) *AchievementController {
	return &AchievementController{svc: s}
}

/*
	POST /achievements
	Student membuat prestasi (status = draft)
*/
func (c *AchievementController) Create(ctx *gin.Context) {
	studentID := utils.GetUserID(ctx)

	var req model.AchievementDetail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := c.svc.Create(ctx, studentID, req)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.SuccessData(ctx, gin.H{
		"id":      id,
		"message": "prestasi berhasil dibuat",
	})
}

/*
	GET /achievements/:id
	Mengambil detail prestasi
*/
func (c *AchievementController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	data, err := c.svc.GetByID(ctx, id)
	if err != nil {
		utils.Error(ctx, 404, err.Error())
		return
	}

	utils.SuccessData(ctx, data)
}

/*
	PUT /achievements/:id
	Update prestasi (hanya draft)
*/
func (c *AchievementController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req model.AchievementDetail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	err := c.svc.Update(ctx, id, req)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil diperbarui")
}

/*
	DELETE /achievements/:id
	Soft delete → status menjadi 'deleted'
*/
func (c *AchievementController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.svc.Delete(ctx, id)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil dihapus")
}

/*
	POST /achievements/:id/submit
	Student mengirim prestasi → status 'submitted'
*/
func (c *AchievementController) Submit(ctx *gin.Context) {
	id := ctx.Param("id")
	studentID := utils.GetUserID(ctx)

	err := c.svc.Submit(ctx, id, studentID)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil disubmit")
}

/*
	POST /achievements/:id/approve
	Dosen meng-approve → status 'verified'
*/
func (c *AchievementController) Approve(ctx *gin.Context) {
	id := ctx.Param("id")
	lectureID := utils.GetUserID(ctx)

	err := c.svc.Approve(ctx, id, lectureID)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil diverifikasi")
}

/*
	POST /achievements/:id/reject
	Dosen menolak prestasi → status 'rejected'
*/
func (c *AchievementController) Reject(ctx *gin.Context) {
	id := ctx.Param("id")
	lectureID := utils.GetUserID(ctx)

	var req struct {
		Note string `json:"note"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid request body")
		return
	}

	err := c.svc.Reject(ctx, id, lectureID, req.Note)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil ditolak")
}

/*
	GET /students/:id/achievements
	List prestasi berdasarkan student
*/
func (c *AchievementController) ListByStudent(ctx *gin.Context) {
	studentID := ctx.Param("id")

	list, err := c.svc.ListByStudent(ctx, studentID)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.SuccessData(ctx, list)
}
