package service

import (
	"context"
	pg "UAS PBE/internal/repository/postgres"
)

type StudentService struct {
	repo pg.StudentRepo
}

func NewStudentService(r pg.StudentRepo) *StudentService {
	return &StudentService{repo: r}
}

// Profile mahasiswa berdasarkan JWT
func (s *StudentService) Profile(ctx context.Context, studentID string) (interface{}, error) {
	return s.repo.GetByID(ctx, studentID)
}
