package model

import "time"

/*
	AchievementReference
	Model utama untuk tabel PostgreSQL (prestasi)
*/
type AchievementReference struct {
	ID                 string    `db:"id"`
	StudentID          string    `db:"student_id"`
	Category           string    `db:"category"`
	Status             string    `db:"status"`
	MongoAchievementID string    `db:"mongo_achievement_id"`

	// Untuk proses APPROVE / REJECT oleh dosen
	LecturerID    string    `db:"lecturer_id"`   // dosen yang menangani prestasi
	VerifiedAt    time.Time `db:"verified_at"`   // waktu approve
	VerifiedBy    string    `db:"verified_by"`   // dosen yang approve
	RejectionNote string    `db:"rejection_note"`// catatan alasan ditolak

	// Untuk proses SUBMIT oleh mahasiswa
	SubmittedAt time.Time `db:"submitted_at"`

	// Info metadata
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
