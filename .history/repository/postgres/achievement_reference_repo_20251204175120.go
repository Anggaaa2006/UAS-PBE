package postgres

import (
    "context"
    "database/sql"
    "errors"

    "uas_pbe/model"
)

/*
    AchievementReferenceRepo
    Interface PostgreSQL untuk reference prestasi
*/
type AchievementReferenceRepo interface {
    Create(ctx context.Context, ref model.AchievementReference) error
    GetByID(ctx context.Context, id string) (*model.AchievementReference, error)
    UpdateStatus(ctx context.Context, id string, status string, rejectionNote string) error
    Verify(ctx context.Context, id string, lectureID string, note string) error
    Reject(ctx context.Context, id string, lectureID string, note string) error
    ListByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error)
}

type achievementReferenceRepo struct {
    db *sql.DB
}

func NewAchievementReferenceRepo(db *sql.DB) AchievementReferenceRepo {
    return &achievementReferenceRepo{db}
}

/*
    Create (status = draft)
*/
func (r *achievementReferenceRepo) Create(ctx context.Context, ref model.AchievementReference) error {

    query := `
        INSERT INTO achievement_references 
        (id, student_id, mongo_achievement_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW());
    `

    _, err := r.db.ExecContext(
        ctx,
        query,
        ref.ID,
        ref.StudentID,
        ref.MongoAchievementID,
        ref.Status,
    )

    return err
}

/*
    GetByID
*/
func (r *achievementReferenceRepo) GetByID(ctx context.Context, id string) (*model.AchievementReference, error) {

    query := `
        SELECT 
            id, student_id, mongo_achievement_id, status,
            submitted_at, verified_at, verified_by,
            rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE id = $1
    `

    row := r.db.QueryRowContext(ctx, query, id)

    var ref model.AchievementReference

    err := row.Scan(
        &ref.ID,
        &ref.StudentID,
        &ref.MongoAchievementID,
        &ref.Status,
        &ref.SubmittedAt,
        &ref.VerifiedAt,
        &ref.VerifiedBy,
        &ref.RejectionNote,
        &ref.CreatedAt,
        &ref.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("prestasi tidak ditemukan")
        }
        return nil, err
    }

    return &ref, nil
}

/*
    UpdateStatus
    - Untuk delete → deleted
    - Untuk submit → submitted
*/
func (r *achievementReferenceRepo) UpdateStatus(ctx context.Context, id string, status string, rejectionNote string) error {

    query := `
        UPDATE achievement_references
        SET status = $1, rejection_note = $2, updated_at = NOW()
        WHERE id = $3
    `
    _, err := r.db.ExecContext(ctx, query, status, rejectionNote, id)
    return err
}

/*
    Verify
    Approve prestasi → status = verified
    NOTE: service mengirim note = "" (kosong)
*/
func (r *achievementReferenceRepo) Verify(ctx context.Context, id string, lectureID string, note string) error {

    query := `
        UPDATE achievement_references
        SET 
            status = 'verified',
            verified_by = $1,
            rejection_note = $2,
            verified_at = NOW(),
            updated_at = NOW()
        WHERE id = $3
    `

    _, err := r.db.ExecContext(ctx, query, lectureID, note, id)
    return err
}

/*
    Reject
*/
func (r *achievementReferenceRepo) Reject(ctx context.Context, id string, lectureID string, note string) error {

    query := `
        UPDATE achievement_references
        SET 
            status = 'rejected',
            rejection_note = $1,
            verified_by = $2,
            verified_at = NOW(),
            updated_at = NOW()
        WHERE id = $3
    `

    _, err := r.db.ExecContext(ctx, query, note, lectureID, id)
    return err
}

/*
    List By Student
*/
func (r *achievementReferenceRepo) ListByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error) {

    query := `
        SELECT 
            id, student_id, mongo_achievement_id, status,
            submitted_at, verified_at, verified_by,
            rejection_note, created_at, updated_at
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
        var ref model.AchievementReference

        err := rows.Scan(
            &ref.ID,
            &ref.StudentID,
            &ref.MongoAchievementID,
            &ref.Status,
            &ref.SubmittedAt,
            &ref.VerifiedAt,
            &ref.VerifiedBy,
            &ref.RejectionNote,
            &ref.CreatedAt,
            &ref.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }

        list = append(list, ref)
    }

    return list, nil
}
