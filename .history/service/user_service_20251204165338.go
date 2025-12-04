package service

import (
	"context"
	pg "uas_pbe/repository/postgres"
)

type UserService struct {
	repo pg.UserRepo
}

func NewUserService(r pg.UserRepo) *UserService {
	return &UserService{repo: r}
}

// List semua user
func (s *UserService) List(ctx context.Context) (interface{}, error) {
	return s.repo.List(ctx)
}
