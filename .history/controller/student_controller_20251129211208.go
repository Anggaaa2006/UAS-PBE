package controller

import (
	"UAS PBE/internal/service"
	"uas_pbe/utils"
	"github.com/gin-gonic/gin"
)

type StudentController struct {
	service *service.StudentService
}

func NewStudentController(s *service.StudentService) *StudentController {
	return &StudentController{service: s}
}

// GET /students/profile
// Mendapatkan profile mahasiswa dari JWT
func (c *StudentController) Profile(ctx *gin.Context) {
	studentID := utils.GetUserID(ctx)

	data, err := c.service.Profile(ctx, studentID)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.SuccessData(ctx, data)
}
