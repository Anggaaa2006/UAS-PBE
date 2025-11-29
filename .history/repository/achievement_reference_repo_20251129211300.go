package repository

import (
	"context"
	"database/sql"
	"errors"
	"uas_pbe/model"
)

// Interface untuk service layer
type AchievementReferenceRepo interface {
	Create(ctx context.Context, studentID, title, categoryID string) (string, error)
	GetByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error)
	GetByID(ctx context.Context, id string) (*model.AchievementReference, error)
	Update(ctx context.Context, id string, title string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Reject(ctx context.Context, id string, reason string) error
	CountByStatus(ctx context.Context) (interface{}, error)
}

type achievementReferenceRepo struct {
	db *sql.DB
}

// Constructor
func NewAchievementReferenceRepo(db *sql.DB) AchievementReferenceRepo {
	return &achievementReferenceRepo{db}
}

// =====================================================
// CREATE PRESTASI (metadata di PostgreSQL)
// =====================================================
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

// =====================================================
// LIST PRESTASI BY STUDENT
// =====================================================
func (r *achievementReferenceRepo) GetByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, title, category_id, status
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
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

// =====================================================
// GET DETAIL METADATA BY ID
// =====================================================
func (r *achievementReferenceRepo) GetByID(ctx context.Context, id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, title, category_id, status
		FROM achievement_references
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var a model.AchievementReference
	if err := row.Scan(&a.ID, &a.StudentID, &a.Title, &a.CategoryID, &a.Status); err != nil {
		return nil, errors.New("prestasi tidak ditemukan")
	}

	return &a, nil
}

// =====================================================
// UPDATE TITLE (HANYA UNTUK DRAFT)
// =====================================================
func (r *achievementReferenceRepo) Update(ctx context.Context, id string, title string) error {
	query := `
		UPDATE achievement_references
		SET title = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, title, id)
	return err
}

// =====================================================
// UPDATE STATUS (draft → submitted → verified/rejected/deleted)
// =====================================================
func (r *achievementReferenceRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// =====================================================
// REJECT PRESTASI (+ ALASAN PENOLAKAN)
// =====================================================
func (r *achievementReferenceRepo) Reject(ctx context.Context, id string, reason string) error {
	query := `
		UPDATE achievement_references
		SET status = 'rejected',
		    reject_reason = $1,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, reason, id)
	return err
}

// =====================================================
// REPORT: HITUNG SETIAP STATUS
// =====================================================
func (r *achievementReferenceRepo) CountByStatus(ctx context.Context) (interface{}, error) {
	query := `
		SELECT status, COUNT(*)
		FROM achievement_references
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type Result struct {
		Status string `json:"status"`
		Total  int    `json:"total"`
	}

	var results []Result
	for rows.Next() {
		var res Result
		rows.Scan(&res.Status, &res.Total)
		results = append(results, res)
	}
	return results, nil
}
