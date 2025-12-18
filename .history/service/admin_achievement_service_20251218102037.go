package service

import (
	"context"
	"strconv"

	"uas_pbe/model"
	mg "uas_pbe/repository/mongo"
	pg "uas_pbe/repository/postgres"
)

/*
	AdminAchievementService
	Service untuk ADMIN melihat seluruh prestasi mahasiswa

	KETERKAITAN SRS:
	- FR-010: View All Achievements
*/
type AdminAchievementService struct {
	refRepo    pg.AchievementReferenceRepo
	detailRepo mg.AchievementDetailRepo
}

/*
	Constructor
*/
func NewAdminAchievementService(
	ref pg.AchievementReferenceRepo,
	detail mg.AchievementDetailRepo,
) *AdminAchievementService {
	return &AdminAchievementService{
		refRepo:    ref,
		detailRepo: detail,
	}
}

/*
	⚠️ VERSI LAMA (DINONAKTIFKAN)
	Tetap disimpan untuk dokumentasi UAS
*/
// func (s *AdminAchievementService) ListAll(
// 	ctx context.Context,
// ) ([]model.AchievementFullResponse, error) {
// 	return nil, nil
// }

/*
	ListAll (AKTIF)
	Admin melihat semua prestasi + filter + pagination
*/
func (s *AdminAchievementService) ListAll(
	ctx context.Context,
	params map[string]string,
) ([]model.AchievementFullResponse, int, error) {

	refs, err := s.refRepo.ListAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	status := params["status"]
	studentID := params["student_id"]

	page := 1
	limit := 10

	if p, err := strconv.Atoi(params["page"]); err == nil {
		page = p
	}
	if l, err := strconv.Atoi(params["limit"]); err == nil {
		limit = l
	}

	// ===== FILTER =====
	var filtered []model.AchievementReference
	for _, r := range refs {
		if status != "" && r.Status != status {
			continue
		}
		if studentID != "" && r.StudentID != studentID {
			continue
		}
		filtered = append(filtered, r)
	}

	total := len(filtered)

	// ===== PAGINATION =====
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		return []model.AchievementFullResponse{}, total, nil
	}
	if end > total {
		end = total
	}

	var result []model.AchievementFullResponse
	for _, ref := range filtered[start:end] {
		detail, _ := s.detailRepo.GetByID(ctx, ref.MongoAchievementID)

		result = append(result, model.AchievementFullResponse{
			ID:        ref.ID,
			StudentID: ref.StudentID,
			Status:    ref.Status,
			Detail:    detail,
		})
	}

	return result, total, nil
}
