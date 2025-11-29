package service

import (
	"context"
	"errors"
	"UAS PBE/internal/repository/postgres"
	"uas_pbe/utils"
)

type AuthService struct {
	userRepo postgres.UserRepo   // akses ke tabel users
	roleRepo postgres.RoleRepo   // untuk ambil role user
	cfg      *utils.Config       // untuk JWT secret
}

// Constructor
func NewAuthService(u postgres.UserRepo, r postgres.RoleRepo, c *utils.Config) *AuthService {
	return &AuthService{
		userRepo: u,
		roleRepo: r,
		cfg:      c,
	}
}

// Login melakukan:
// 1. cek user by email
// 2. validasi password
// 3. ambil role
// 4. generate JWT
func (s *AuthService) Login(ctx context.Context, email, pass string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("email tidak ditemukan")
	}

	if !utils.CheckPassword(pass, user.Password) {
		return "", errors.New("password salah")
	}

	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return "", errors.New("gagal mengambil role user")
	}

	token, err := utils.GenerateJWT(user.ID.String(), role.Name, s.cfg.JWTSecret)
	if err != nil {
		return "", errors.New("gagal generate token")
	}

	return token, nil
}
