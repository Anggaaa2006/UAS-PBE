package controller

import (
	"uas_pbe/internal/service"
	"uas_/utils"
	"github.com/gin-gonic/gin"
)

// AchievementController menangani semua endpoint prestasi mahasiswa
// Router → Controller → Service
type AchievementController struct {
	service *service.AchievementService
}

func NewAchievementController(s *service.AchievementService) *AchievementController {
	return &AchievementController{service: s}
}

// ==============================
// GET /achievements
// Ambil semua prestasi (filter by student dari JWT)
// ==============================
func (c *AchievementController) List(ctx *gin.Context) {
	userID := utils.GetUserID(ctx)

	data, err := c.service.List(ctx, userID)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.SuccessData(ctx, data)
}

// ==============================
// POST /achievements
// Mahasiswa membuat prestasi baru (status = draft)
// ==============================
func (c *AchievementController) Create(ctx *gin.Context) {
	userID := utils.GetUserID(ctx)

	var req service.AchievementCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid body")
		return
	}

	err := c.service.Create(ctx, userID, req)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil dibuat (draft)")
}

// ==============================
// PUT /achievements/:id
// Update data draft
// ==============================
func (c *AchievementController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := utils.GetUserID(ctx)

	var req service.AchievementUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid body")
		return
	}

	err := c.service.Update(ctx, userID, id, req)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil diperbarui")
}

// ==============================
// DELETE /achievements/:id
// Soft delete → status "deleted"
// ==============================
func (c *AchievementController) SoftDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := utils.GetUserID(ctx)

	err := c.service.SoftDelete(ctx, userID, id)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil dihapus (soft delete)")
}

// ==============================
// SUBMIT
// ==============================
func (c *AchievementController) Submit(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := utils.GetUserID(ctx)

	err := c.service.Submit(ctx, userID, id)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil dikirim (submitted)")
}

// ==============================
// VERIFY
// ==============================
func (c *AchievementController) Verify(ctx *gin.Context) {
	id := ctx.Param("id")
	verifierID := utils.GetUserID(ctx)

	err := c.service.Verify(ctx, verifierID, id)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil diverifikasi (verified)")
}

// ==============================
// REJECT
// ==============================
func (c *AchievementController) Reject(ctx *gin.Context) {
	id := ctx.Param("id")
	verifierID := utils.GetUserID(ctx)

	var req struct {
		Reason string `json:"reason"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid reason")
		return
	}

	err := c.service.Reject(ctx, verifierID, id, req.Reason)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "prestasi berhasil ditolak (rejected)")
}
