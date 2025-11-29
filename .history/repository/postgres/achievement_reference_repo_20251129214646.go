package postgres

import (
	"context"
	"database/sql"
	"errors"
	"uas_pbe/model"
)

/*
	AchievementReferenceRepo
	Interface repository prestasi utama yang disimpan di PostgreSQL.
*/
type AchievementReferenceRepo interface {

	// Create data prestasi baru
	Create(ctx context.Context, studentID, title, categoryID string) (string, error)

	// Ambil data prestasi berdasarkan ID
	GetByID(ctx context.Context, id string) (*model.AchievementReference, error)

	// Ambil semua prestasi mahasiswa tertentu
	GetByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error)

	// Update judul, kategori
	Update(ctx context.Context, id, title, categoryID string) error

	// Update status (draft, submitted, approved, rejected, deleted)
	UpdateStatus(ctx context.Context, id, status string) error

	// Soft delete → status jadi "deleted"
	SoftDelete(ctx context.Context, id string) error

	// Ambil semua data prestasi
	List(ctx context.Context) ([]model.AchievementReference, error)
}

/*
	Struct implementasi repository
*/
type achievementReferenceRepo struct{ db *sql.DB }

/*
	Constructor untuk membuat repository baru
*/
func NewAchievementReferenceRepo(db *sql.DB) AchievementReferenceRepo {
	return &achievementReferenceRepo{db}
}

/*
	Create
	Membuat data prestasi baru.
	Status awal: draft
*/
func (r *achievementReferenceRepo) Create(ctx context.Context, studentID, title, categoryID string) (string, error) {
	var id string
	query := `
		INSERT INTO achievement_references (student_id, title, category_id, status)
		VALUES ($1, $2, $3, 'draft')
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, studentID, title, categoryID).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

/*
	GetByID
	Mengambil data prestasi berdasarkan ID
*/
func (r *achievementReferenceRepo) GetByID(ctx context.Context, id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, title, category_id, status 
		FROM achievement_references
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	var a model.AchievementReference

	if err := row.Scan(&a.ID, &a.StudentID, &a.Title, &a.CategoryID, &a.Status); err != nil {
		return nil, errors.New("data tidak ditemukan")
	}

	return &a, nil
}

/*
	GetByStudentID
	Mengambil semua prestasi oleh 1 mahasiswa
*/
func (r *achievementReferenceRepo) GetByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, title, category_id, status
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference

	for rows.Next() {
		var a model.AchievementReference
		rows.Scan(&a.ID, &a.StudentID, &a.Title, &a.CategoryID, &a.Status)
		list = append(list, a)
	}

	return list, nil
}

/*
	Update
	Update judul & kategori prestasi
*/
func (r *achievementReferenceRepo) Update(ctx context.Context, id, title, categoryID string) error {
	query := `
		UPDATE achievement_references
		SET title = $1,
		    category_id = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, title, categoryID, id)
	return err
}

/*
	UpdateStatus
	Mengubah status prestasi (draft/submitted/approved/rejected/deleted)
*/
func (r *achievementReferenceRepo) UpdateStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE achievement_references
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

/*
	SoftDelete
	Soft delete → status = "deleted"
*/
func (r *achievementReferenceRepo) SoftDelete(ctx context.Context, id string) error {
	query := `
		UPDATE achievement_references
		SET status = 'deleted'
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

/*
	List
	Mengambil semua data prestasi (semua mahasiswa)
*/
func (r *achievementReferenceRepo) List(ctx context.Context) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, title, category_id, status
		FROM achievement_references
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference

	for rows.Next() {
		var a model.AchievementReference
		rows.Scan(&a.ID, &a.StudentID, &a.Title, &a.CategoryID, &a.Status)
		list = append(list, a)
	}

	return list, nil
}
