package service

import (
	"context"
	mg "UAS PBE/repository/mongo"
	pg "UAS PBE/repository/postgres"
)

type ReportService struct {
	mongoRepo mg.AchievementRepo
	refRepo   pg.AchievementReferenceRepo
}

func NewReportService(m mg.AchievementRepo, r pg.AchievementReferenceRepo) *ReportService {
	return &ReportService{
		mongoRepo: m,
		refRepo:   r,
	}
}

// Contoh laporan: jumlah prestasi per status
func (s *ReportService) AchievementReport(ctx context.Context) (interface{}, error) {
	return s.refRepo.CountByStatus(ctx)
}
