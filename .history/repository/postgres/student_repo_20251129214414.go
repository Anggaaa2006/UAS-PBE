package postgres

import (
	"context"
	"database/sql"
	"uas_pbe/model"
)

/*
	StudentRepo interface
	Mengatur query untuk tabel students
*/
type StudentRepo interface {
	GetByID(ctx context.Context, userID string) (*model.Student, error)
}

type studentRepo struct {
	db *sql.DB
}

/*
	NewStudentRepo
	Constructor student repo
*/
func NewStudentRepo(db *sql.DB) StudentRepo {
	return &studentRepo{db}
}

/*
	GetByID
	Mengambil data mahasiswa berdasarkan user_id
*/
func (r *studentRepo) GetByID(ctx context.Context, userID string) (*model.Student, error) {
	query := `
		SELECT id, user_id, nim, major, semester
		FROM students WHERE user_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, userID)

	var s model.Student
	if err := row.Scan(&s.ID, &s.UserID, &s.NIM, &s.Major, &s.Semester); err != nil {
		return nil, err
	}

	return &s, nil
}
