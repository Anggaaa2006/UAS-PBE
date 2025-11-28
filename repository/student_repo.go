package repository

import (
	"context"
	"database/sql"
	"UAS PBE/model"
)

type StudentRepo interface {
	GetByID(ctx context.Context, id string) (*model.Student, error)
}

type studentRepo struct {
	db *sql.DB
}

func NewStudentRepo(db *sql.DB) StudentRepo {
	return &studentRepo{db}
}

// Ambil profile mahasiswa berdasarkan user_id
func (r *studentRepo) GetByID(ctx context.Context, id string) (*model.Student, error) {
	query := `
		SELECT id, user_id, nim, major, semester
		FROM students
		WHERE user_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s model.Student
	if err := row.Scan(&s.ID, &s.UserID, &s.NIM, &s.Major, &s.Semester); err != nil {
		return nil, err
	}
	return &s, nil
}
