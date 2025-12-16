package model

import "time"

/*
	AchievementReference
	Model utama untuk tabel PostgreSQL (prestasi)
*/
type AchievementReference struct {
	ID                 string     `db:"id"`
	StudentID          string     `db:"student_id"`
	MongoAchievementID string     `db:"mongo_achievement_id"`
	Status             string     `db:"status"`

	// Untuk proses SUBMIT oleh mahasiswa
	SubmittedAt *time.Time `db:"submitted_at"` // nullable

	// Untuk proses APPROVE / REJECT oleh dosen
	VerifiedAt    *time.Time `db:"verified_at"`     // nullable
	VerifiedBy    *string    `db:"verified_by"`     // nullable
	RejectionNote *string    `db:"rejection_note"`  // nullable

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
