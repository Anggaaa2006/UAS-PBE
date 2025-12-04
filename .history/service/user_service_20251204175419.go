package service

import (
    "context"

    "uas_pbe/model"
    pg "uas_pbe/repository/postgres"
)

type UserService struct {
    repo pg.UserRepo
}

func NewUserService(r pg.UserRepo) *UserService {
    return &UserService{repo: r}
}

/*
    List semua user
    Service hanya meneruskan ke repository
*/
func (s *UserService) List(ctx context.Context) ([]model.User, error) {
    return s.repo.List(ctx)
}
