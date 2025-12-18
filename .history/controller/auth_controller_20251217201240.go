package controller

import (
	"uas_pbe/service"
	"uas_pbe/utils"

	"github.com/gin-gonic/gin"
)

/*
	AuthController
	Menangani endpoint authentication seperti:
	- POST /auth/register
	- POST /auth/login
	- GET  /auth/profile

	Controller hanya menerima request HTTP,
	memanggil Service untuk logic bisnis,
	dan mengembalikan response JSON.
*/
type AuthController struct {
	svc service.AuthService
}

/*
	NewAuthController
	Constructor untuk inisialisasi AuthController.
	Panggil ini di main.go
*/
func NewAuthController(s service.AuthService) *AuthController {
	return &AuthController{svc: s}
}

/*
	Struct request Register
*/
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
	POST /auth/register

	Flow:
	1. Ambil body JSON
	2. Panggil AuthService.Register
	3. Return message sukses / error
*/
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid request body")
		return
	}

	if err := c.svc.Register(ctx, req.Name, req.Email, req.Password); err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "akun berhasil didaftarkan")
}

/*
	Struct request Login
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
	POST /auth/login

	Flow:
	1. Ambil email & password
	2. Service cek user + password
	3. Generate JWT
	4. Return token
*/
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid request body")
		return
	}

	token, err := c.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		utils.Error(ctx, 401, err.Error())
		return
	}

	utils.SuccessData(ctx, gin.H{
		"token": token,
	})
}

/*
	GET /auth/profile
	(Protected - pakai JWT Middleware)

	Flow:
	1. Middleware JWT decode token
	2. User ID disimpan di context
	3. Ambil profile user dari Service
	4. Return data profile
*/
func (c *AuthController) Profile(ctx *gin.Context) {

	// Ambil user_id dari context (di-set oleh JWT middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Error(ctx, 401, "unauthorized")
		return
	}

	// Ambil data profile dari service
	profile, err := c.svc.GetProfile(ctx, userID.(string))
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.SuccessData(ctx, profile)
}
/*
	POST /auth/refresh
*/
func (c *AuthController) Refresh(ctx *gin.Context) {

	// Ambil dari middleware JWT
	userID := ctx.GetString("user_id")
	role := ctx.GetString("role")

	token, err := c.service.RefreshToken(
		ctx.Request.Context(),
		userID,
		role,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token": token,
	})
}

/*
	POST /auth/logout
*/
func (c *AuthController) Logout(ctx *gin.Context) {

	err := c.service.Logout(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "logout berhasil",
	})
}
