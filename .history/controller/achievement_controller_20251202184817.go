package controller

import (
	"net/http"
	"uas_pbe/service"

	"github.com/gin-gonic/gin"
)

/*
	AchievementController
	Interface controller untuk endpoint prestasi.
*/
type AchievementController interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	ListByStudent(c *gin.Context)
	Submit(c *gin.Context)
	Approve(c *gin.Context)
	Reject(c *gin.Context)
}

/*
	Struct implementasi controller
	memegang reference ke service
*/
type achievementController struct {
	svc service.AchievementService
}

/*
	NewAchievementController
	Constructor
*/
func NewAchievementController(svc service.AchievementService) AchievementController {
	return &achievementController{svc}
}

/*
	REQUEST DTO
	Struct untuk menerima JSON body dari client
*/
type CreateAchievementRequest struct {
	StudentID  string `json:"student_id"`
	Title      string `json:"title"`
	CategoryID string `json:"category_id"`
	Description string `json:"description"`
}

type UpdateAchievementRequest struct {
	Title       string `json:"title"`
	CategoryID  string `json:"category_id"`
	Description string `json:"description"`
}

/*
	POST /achievements
	Membuat prestasi baru
*/
func (ctrl *achievementController) Create(c *gin.Context) {
	var req CreateAchievementRequest

	// Bind request body ke struct req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format request tidak valid"})
		return
	}

	// Panggil service untuk membuat prestasi
	data, err := ctrl.svc.Create(
		c.Request.Context(),
		req.StudentID,
		req.Title,
		req.CategoryID,
		req.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "berhasil membuat prestasi",
		"data":    data,
	})
}

/*
	PUT /achievements/:id
	Update prestasi
*/
func (ctrl *achievementController) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAchievementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format request tidak valid"})
		return
	}

	err := ctrl.svc.Update(
		c.Request.Context(),
		id,
		req.Title,
		req.CategoryID,
		req.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "berhasil mengupdate prestasi"})
}

/*
	DELETE /achievements/:id
	Soft delete prestasi
*/
func (ctrl *achievementController) Delete(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.svc.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "berhasil menghapus prestasi"})
}

/*
	GET /achievements/:id
	Get detail prestasi (Postgres + Mongo)
*/
func (ctrl *achievementController) GetByID(c *gin.Context) {
	id := c.Param("id")

	ref, detail, err := ctrl.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reference": ref,
		"detail":    detail,
	})
}

/*
	GET /students/:id/achievements
	List prestasi milik satu mahasiswa
*/
func (ctrl *achievementController) ListByStudent(c *gin.Context) {
	studentID := c.Param("id")

	list, err := ctrl.svc.ListStudentAchievement(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

/*
	POST /achievements/:id/submit
	Mengubah status menjadi 'submitted'
*/
func (ctrl *achievementController) Submit(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.svc.Submit(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "berhasil submit prestasi"})
}

/*
	POST /achievements/:id/approve
	Mengubah status menjadi 'approved'
*/
func (ctrl *achievementController) Approve(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.svc.Approve(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "prestasi disetujui"})
}

/*
	POST /achievements/:id/reject
	Mengubah status menjadi 'rejected'
*/
func (ctrl *achievementController) Reject(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.svc.Reject(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "prestasi ditolak"})
}
