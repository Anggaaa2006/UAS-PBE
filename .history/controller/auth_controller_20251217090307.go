package controller

import (
	"uas_pbe/service"   // pakai service baru (bukan internal/service)
	"uas_pbe/utils"

	"github.com/gin-gonic/gin"
)

/*
	AuthController
	Menangani endpoint authentication seperti:
	- POST /auth/register
	- POST /auth/login

	Controller hanya menerima request dari HTTP,
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
	Struct untuk menerima body request register
*/
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
	POST /auth/register

	Flow:
	1. Ambil body JSON (name, email, password)
	2. Lempar ke AuthService.Register()
	3. Jika email duplikat → error
	4. Jika sukses → return message sukses
*/
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest

	// Validasi format JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid request body")
		return
	}

	// Proses register via Service
	if err := c.svc.Register(ctx, req.Name, req.Email, req.Password); err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	utils.Success(ctx, "akun berhasil didaftarkan")
}

/*
	Struct untuk menerima body request login
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
	POST /auth/login

	Flow:
	1. Ambil email & password
	2. Service.Login akan cek user + cek password
	3. Jika cocok → generate JWT Token
	4. Response token ke client
*/
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest

	// Validasi JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid request body")
		return
	}

	// Login, service akan cek password + generate token
	token, err := c.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		utils.Error(ctx, 401, err.Error())
		return
	}

	utils.SuccessData(ctx, gin.H{
		"token": token,
	})
}
// Profile
// GET /auth/profile
// Mengambil data user dari token JWT
func (c *AuthController) Profile(ctx *gin.Context) {

	// Ambil user_id dari JWT (di-set oleh middleware)
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.Error(ctx, 401, "unauthorized")
		return
	}

	// Ambil data user dari service
	user, err := c.service.GetProfile(ctx, userID)
	if err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}

	// Return data user (tanpa password)
	utils.SuccessData(ctx, user)
}
