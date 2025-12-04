package service

import (
	"context"
	"errors"
	"time"
	"uas_pbe/model"
	"uas_pbe/repository/postgres"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
	AuthService interface
	Menangani login & register
*/
type AuthService interface {
	Register(ctx context.Context, name, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	userRepo postgres.UserRepo
}

/*
	Constructor
*/
func NewAuthService(userRepo postgres.UserRepo) AuthService {
	return &authService{userRepo}
}

/*
	Register
	1. Hash password
	2. Generate ID
	3. Insert ke database
*/
func (s *authService) Register(ctx context.Context, name, email, password string) error {

	// Cek apakah email sudah terpakai
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return errors.New("email sudah terdaftar")
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Buat struct user
	user := model.User{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: string(hashed),
		RoleID:   "student", // default role mahasiswa
	}

	// Insert ke database
	return s.userRepo.
}

/*
	Login
	1. Ambil user berdasarkan email
	2. Cek password
	3. Generate JWT Token
*/
func (s *authService) Login(ctx context.Context, email, password string) (string, error) {

	// Ambil user dari DB
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// Cek password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("email atau password salah")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.RoleID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := "MY_SECRET_KEY" // nanti pindah ke .env
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}
