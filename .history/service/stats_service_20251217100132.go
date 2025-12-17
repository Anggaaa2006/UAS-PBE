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
    Mengambil statistik jumlah prestasi berdasarkan status

    Contoh output:
    {
        "draft": 5,
        "submitted": 10,
        "verified": 7,
        "rejected": 2
    }

    Digunakan untuk:
    - Dashboard Admin
    - Monitoring Dosen
*/
func (s *statsService) GetAchievementStats(
    ctx context.Context,
) (map[string]int, error) {

    /*
        Karena repository belum punya fungsi statistik,
        maka kita buat query langsung di service
        (masih sesuai arsitektur untuk UAS)
    */

    rows, err := s.refRepo.
        (interface {
            CountByStatus(context.Context) (map[string]int, error)
        }).
        CountByStatus(ctx)

    if err != nil {
        return nil, err
    }

    return rows, nil
}
type StatsService interface {
    GetAchievementStats(ctx context.Context) (map[string]int, error)

    // FR-011: Statistik prestasi per mahasiswa
    GetStudentStats(ctx context.Context, studentID string) (map[string]interface{}, error)
}
