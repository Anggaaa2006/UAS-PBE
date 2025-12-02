package controller

import (
	"net/http"
	"uas_pbe/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc service.AuthService
}

func NewAuthController(svc service.AuthService) *AuthController {
	return &AuthController{svc}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
	POST /auth/register
*/
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format request salah"})
		return
	}

	err := c.svc.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "akun berhasil dibuat"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
	POST /auth/login
*/
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "format request salah"})
		return
	}

	token, err := c.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
		"token":   token,
	})
}
