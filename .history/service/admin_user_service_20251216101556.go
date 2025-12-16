// path: internal/service/admin_user_service.go  OR service/admin_user_service.go (sesuaikan projekmu)
package service

import (
    "context"
    "errors"

    "uas_pbe/model"
    pg "uas_pbe/repository/postgres"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

/*
    AdminUserService
    Menangani CRUD user yang dipakai Admin (FR-009)
*/
type AdminUserService struct {
    userRepo pg.UserRepo
}

func NewAdminUserService(r pg.UserRepo) *AdminUserService {
    return &AdminUserService{userRepo: r}
}

// Create user (Admin creates user and sets role)
// password plain is hashed here
func (s *AdminUserService) Create(ctx context.Context, name, email, password, role string) (*model.User, error) {
    existing, _ := s.userRepo.FindByEmail(ctx, email)
    if existing != nil {
        return nil, errors.New("email sudah terdaftar")
    }
    hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    u := model.User{
        ID:       uuid.NewString(),
        Name:     name,
        Email:    email,
        Password: string(hashed),
        Role:     role,
    }
    if err := s.userRepo.Create(ctx, u); err != nil {
        return nil, err
    }
    return &u, nil
}

func (s *AdminUserService) List(ctx context.Context) ([]model.User, error) {
    return s.userRepo.List(ctx)
}

func (s *AdminUserService) GetByID(ctx context.Context, id string) (*model.User, error) {
    return s.userRepo.GetByID(ctx, id)
}

func (s *AdminUserService) Update(ctx context.Context, id, name, email, password, role string) error {
    u, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if u == nil {
        return errors.New("user tidak ditemukan")
    }
    // if password provided -> hash it, else keep existing hash
    passHash := u.Password
    if password != "" {
        h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        passHash = string(h)
    }
    u.Name = name
    u.Email = email
    u.Password = passHash
    u.Role = role
    return s.userRepo.Update(ctx, *u)
}

func (s *AdminUserService) Delete(ctx context.Context, id string) error {
    return s.userRepo.Delete(ctx, id)
}

func (s *AdminUserService) UpdateRole(ctx context.Context, id, role string) error {
    return s.userRepo.UpdateRole(ctx, id, role)
}
