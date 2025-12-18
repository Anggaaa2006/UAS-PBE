package postgres

import (
    "context"
    "database/sql"
    "errors"

    "uas_pbe/model"
)

/*
    AchievementReferenceRepo
    Repository PostgreSQL untuk tabel achievement_references

    Digunakan untuk:
    - Menyimpan referensi prestasi mahasiswa
    - Mengelola status prestasi (draft, submitted, verified, rejected)
    - Digunakan oleh mahasiswa, dosen, dan admin

    KETERKAITAN SRS:
    - FR-010: View All Achievements
*/
type AchievementReferenceRepo interface {
    Create(ctx context.Context, ref model.AchievementReference) error
    GetByID(ctx context.Context, id string) (*model.AchievementReference, error)
    UpdateStatus(ctx context.Context, id string, status string, note string) error
    Verify(ctx context.Context, id string, lecturerID string, note string) error
    Reject(ctx context.Context, id string, lecturerID string, note string) error
    ListByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error)

    // FR-010: Admin melihat semua prestasi
    ListAll(ctx context.Context) ([]model.AchievementReference, error)
    // FR-011: Statistik jumlah prestasi berdasarkan status
    CountByStatus(ctx context.Context) (map[string]int, error)
    // ✅ TAMBAHAN – FR History
	GetHistory(ctx context.Context, id string) ([]map[string]interface{}, error)
    // Statistik prestasi per mahasiswa
    CountByStudent(ctx context.Context, studentID string) (map[string]int, error)

}
/*
    Implementasi repository
*/
type achievementReferenceRepo struct {
    db *sql.DB
}

/*
    Constructor AchievementReferenceRepo
*/
func NewAchievementReferenceRepo(db *sql.DB) AchievementReferenceRepo {
    return &achievementReferenceRepo{db}
}

/*
    Create
    Digunakan mahasiswa untuk membuat prestasi baru
    Status awal: draft

    Alur SRS:
    - Mahasiswa input prestasi
    - Sistem menyimpan referensi ke PostgreSQL
*/
func (r *achievementReferenceRepo) Create(ctx context.Context, ref model.AchievementReference) error {

    query := `
        INSERT INTO achievement_references
        (id, student_id, mongo_achievement_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
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
    Mengambil satu prestasi berdasarkan ID

    Digunakan oleh:
    - Mahasiswa (lihat detail)
    - Dosen (verifikasi)
    - Admin (monitoring)

    FR-010: View All Achievements
*/
func (r *achievementReferenceRepo) GetByID(ctx context.Context, id string) (*model.AchievementReference, error) {

    query := `
        SELECT
            id,
            student_id,
            mongo_achievement_id,
            status,
            submitted_at,
            verified_at,
            verified_by,
            rejection_note,
            created_at,
            updated_at
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
    Digunakan untuk:
    - Submit prestasi (draft → submitted)
    - Soft delete (status → deleted)

    Dilakukan oleh mahasiswa
*/
func (r *achievementReferenceRepo) UpdateStatus(
    ctx context.Context,
    id string,
    status string,
    note string,
) error {

    query := `
        UPDATE achievement_references
        SET
            status = $1,
            rejection_note = $2,
            updated_at = NOW()
        WHERE id = $3
    `

    _, err := r.db.ExecContext(ctx, query, status, note, id)
    return err
}

/*
    Verify
    Digunakan dosen untuk menyetujui prestasi mahasiswa
    Status berubah menjadi: verified

    FR-010: Admin/Dosen melihat prestasi
*/
func (r *achievementReferenceRepo) Verify(
    ctx context.Context,
    id string,
    lecturerID string,
    note string,
) error {

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

    _, err := r.db.ExecContext(ctx, query, lecturerID, note, id)
    return err
}

/*
    Reject
    Digunakan dosen untuk menolak prestasi
    Wajib menyertakan catatan penolakan
*/
func (r *achievementReferenceRepo) Reject(
    ctx context.Context,
    id string,
    lecturerID string,
    note string,
) error {

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

    _, err := r.db.ExecContext(ctx, query, note, lecturerID, id)
    return err
}

/*
    ListByStudent
    Menampilkan seluruh prestasi milik satu mahasiswa

    Digunakan oleh:
    - Mahasiswa (riwayat prestasi)
    - Dosen wali
    - Admin

    FR-010: View All Achievements
*/
func (r *achievementReferenceRepo) ListByStudent(
    ctx context.Context,
    studentID string,
) ([]model.AchievementReference, error) {

    query := `
        SELECT
            id,
            student_id,
            mongo_achievement_id,
            status,
            submitted_at,
            verified_at,
            verified_by,
            rejection_note,
            created_at,
            updated_at
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
/*
    CountByStatus
    Menghitung jumlah prestasi berdasarkan status

    Digunakan untuk:
    - FR-011 Statistik
    - Dashboard Admin / Dosen

    Contoh hasil:
    {
        "draft": 3,
        "submitted": 5,
        "verified": 2,
        "rejected": 1
    }
*/
func (r *achievementReferenceRepo) CountByStatus(
    ctx context.Context,
) (map[string]int, error) {

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

    // Map untuk menyimpan hasil statistik
    result := make(map[string]int)

    for rows.Next() {
        var status string
        var total int

        err := rows.Scan(&status, &total)
        if err != nil {
            return nil, err
        }

        result[status] = total
    }

    return result, nil
}
/*
    ListAll
    Mengambil seluruh prestasi (untuk Admin)

    KETERKAITAN SRS:
    - FR-010: View All Achievements
*/
func (r *achievementReferenceRepo) ListAll(ctx context.Context) ([]model.AchievementReference, error) {

    query := `
        SELECT 
            id, student_id, mongo_achievement_id, status,
            submitted_at, verified_at, verified_by,
            rejection_note, created_at, updated_at
        FROM achievement_references
        ORDER BY created_at DESC
    `

    rows, err := r.db.QueryContext(ctx, query)
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
/*
    CountByStudent
    Menghitung jumlah prestasi mahasiswa berdasarkan status

    Digunakan untuk:
    - FR-011 laporan per mahasiswa
*/
func (r *achievementReferenceRepo) CountByStudent(
    ctx context.Context,
    studentID string,
) (map[string]int, error) {

    query := `
        SELECT status, COUNT(*)
        FROM achievement_references
        WHERE student_id = $1
        GROUP BY status
    `

    rows, err := r.db.QueryContext(ctx, query, studentID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[string]int)

    for rows.Next() {
        var status string
        var total int

        if err := rows.Scan(&status, &total); err != nil {
            return nil, err
        }
        result[status] = total
    }

    return result, nil
}
/*
	GetHistory
	Mengambil riwayat perubahan status prestasi
*/
func (r *achievementReferenceRepo) GetHistory(
	ctx context.Context,
	id string,
) ([]map[string]interface{}, error) {

	query := `
		SELECT status, updated_at
		FROM achievement_references
		WHERE id = $1
		ORDER BY updated_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []map[string]interface{}

	for rows.Next() {
		var status string
		var time string

		rows.Scan(&status, &time)

		history = append(history, map[string]interface{}{
			"status": status,
			"time":   time,
		})
	}

	return history, nil
}
/*
	CountByStudent
	Menghitung jumlah prestasi mahasiswa berdasarkan status
*/
func (r *achievementReferenceRepo) CountByStudent(
	ctx context.Context,
	studentID string,
) (map[string]int, error) {

	query := `
		SELECT status, COUNT(*) 
		FROM achievement_references
		WHERE student_id = $1
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)

	for rows.Next() {
		var status string
		var total int
		if err := rows.Scan(&status, &total); err != nil {
			return nil, err
		}
		result[status] = total
	}

	return result, nil
}
