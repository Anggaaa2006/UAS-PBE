package model

import "time"

/*
	AchievementReference
	Menyimpan data prestasi di PostgreSQL

	- ID               → PK di PostgreSQL
	- StudentID        → user_id mahasiswa
	- Category         → kategori prestasi
	- Status           → draft / submitted / verified / rejected / deleted
	- MongoAchievementID → ID detail prestasi di Mongo
	- LecturerID       → dosen yang approve/reject
	- SubmittedAt      → waktu prestasi dikirim mahasiswa
	- CreatedAt        → waktu dibuat
	- UpdatedAt        → waktu diperbarui
*/
type AchievementReference struct {
	ID                 string    `db:"id"`
	StudentID          string    `db:"student_id"`
	Category           string    `db:"category"`
	Status             string    `db:"status"`
	MongoAchievementID string    `db:"mongo_achievement_id"`
	LecturerID         string    `db:"lecturer_id"`
	SubmittedAt        time.Time `db:"submitted_at"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
