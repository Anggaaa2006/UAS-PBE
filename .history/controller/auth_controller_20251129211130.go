package controller

import (
	"uas_pbe/internal/service"
	"UAS PBE/utils"
	"github.com/gin-gonic/gin"
)

// AuthController menangani endpoint authentication:
// - POST /auth/login
// Controller hanya menerima request dan mengirim respose.
// Logika bisnis ada di AuthService.
type AuthController struct {
	service *service.AuthService
}

// Constructor (dipanggil di main.go)
func NewAuthController(s *service.AuthService) *AuthController {
	return &AuthController{service: s}
}

// Login handler
// Flow:
// 1. Ambil email & password dari body
// 2. Service akan melakukan pengecekan user + generate JWT
// 3. Return token
func (c *AuthController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid request body")
		return
	}

	token, err := c.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.SuccessData(ctx, gin.H{
		"token": token,
	})
}
