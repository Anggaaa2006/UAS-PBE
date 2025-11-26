package controller

import (
	"UAS PBE/internal/service"
	"UAS PBE/utils"
	"github.com/gin-gonic/gin"
)

// UserController untuk data user admin/dosen/mahasiswa
type UserController struct {
	service *service.UserService
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{service: s}
}

// GET /users
func (c *UserController) List(ctx *gin.Context) {
	data, err := c.service.List(ctx)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}
	utils.SuccessData(ctx, data)
}
