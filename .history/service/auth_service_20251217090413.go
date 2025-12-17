package service

import (
    "context"
    "errors"

    "uas_pbe/model"
    "uas_pbe/repository/postgres"
    "uas_pbe/middleware"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

/*
    AuthService
    Menangani logic register & login
*/
type AuthService interface {
    Register(ctx context.Context, name, email, password string) error
    Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
    repo postgres.UserRepo
}

/*
    Constructor
*/
func NewAuthService(userRepo postgres.UserRepo) AuthService {
    return &authService{repo: userRepo}
}

/*
    REGISTER USER
    - cek apakah email sudah terpakai
    - hash password
    - set role = student
*/
func (s *authService) Register(ctx context.Context, name, email, password string) error {

    // cek email apakah sudah ada
    existing, _ := s.repo.FindByEmail(ctx, email)
    if existing != nil {
        return errors.New("email sudah terdaftar")
    }

    // hash password
    hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

    // buat struct user baru
    u := model.User{
        ID:       uuid.NewString(),
        Name:     name,
        Email:    email,
        Password: string(hashed),
        Role:     "student", // default mahasiswa
    }

    return s.repo.Create(ctx, u)
}

/*
    LOGIN USER
    - ambil user berdasarkan email
    - cek password
    - generate JWT pakai helper middleware.GenerateJWT
*/
func (s *authService) Login(ctx context.Context, email, password string) (string, error) {

    // Ambil user
    u, err := s.repo.FindByEmail(ctx, email)
    if err != nil || u == nil {
        return "", errors.New("email atau password salah")
    }

    // Cek password
    if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
        return "", errors.New("email atau password salah")
    }

    // Generate JWT
    token, err := middleware.GenerateJWT(u.ID, u.Role)
    if err != nil {
        return "", err
    }

    return token, nil
}
type AuthService interface {
	Register(ctx context.Context, name, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)

	// Ambil profile user dari token
	GetProfile(ctx context.Context, userID string) (*model.User, error)
}
// GetProfile
// Mengambil data user berdasarkan user_id dari JWT
func (s *authService) GetProfile(ctx context.Context, userID string) (*model.User, error) {

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Password tidak dikirim ke response
	user.Password = ""

	return user, nil
}
