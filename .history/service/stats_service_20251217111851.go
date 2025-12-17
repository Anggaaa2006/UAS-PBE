package service

import (
	"context"

	"uas_pbe/repository/postgres"
)

/*
	StatsService
	Service untuk menampilkan statistik prestasi

	KETERKAITAN SRS:
	- FR-011: Statistics
*/
type StatsService interface {
	GetAchievementStats(ctx context.Context) (map[string]int, error)

	// ✅ reports/student/:id
	GetStudentStats(ctx context.Context, studentID string) (map[string]int, error)
}

/*
	Implementasi StatsService
*/
type statsService struct {
	refRepo postgres.AchievementReferenceRepo
}

/*
	Constructor StatsService
*/
func NewStatsService(
	refRepo postgres.AchievementReferenceRepo,
) StatsService {
	return &statsService{
		refRepo: refRepo,
	}
}

/*
	GetAchievementStats
*/
func (s *statsService) GetAchievementStats(
	ctx context.Context,
) (map[string]int, error) {

	return s.refRepo.CountByStatus(ctx)
}

/*
	GetStudentStats
	Statistik prestasi untuk 1 mahasiswa

	Output contoh:
	{
		"draft": 2,
		"submitted": 3,
		"verified": 1
	}
*/
func (s *statsService) GetStudentStats(
	ctx context.Context,
	studentID string,
) (map[string]int, error) {

	// Ambil semua prestasi mahasiswa
	list, err := s.refRepo.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Hitung per status
	result := make(map[string]int)
	for _, a := range list {
		result[a.Status]++
	}

	return result, nil
}
