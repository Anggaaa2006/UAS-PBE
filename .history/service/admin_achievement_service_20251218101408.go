package service

import (
    "context"

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
    ListAll
    Mengambil seluruh prestasi mahasiswa (Admin)

    Alur:
    1. Ambil data reference dari PostgreSQL
    2. Ambil detail prestasi dari MongoDB
    3. Gabungkan menjadi response lengkap

    KETERKAITAN SRS:
    - FR-010: View All Achievements
*/
func (s *AdminAchievementService) ListAll(
    ctx context.Context,
) ([]model.AchievementFullResponse, error) {

    refs, err := s.refRepo.ListAll(ctx)
    if err != nil {
        return nil, err
    }

    var result []model.AchievementFullResponse

    for _, ref := range refs {
        detail, _ := s.detailRepo.GetByID(ctx, ref.MongoAchievementID)

        result = append(result, model.AchievementFullResponse{
            ID:        ref.ID,
            StudentID: ref.StudentID,
            Status:    ref.Status,
            Detail:    detail,
        })
    }

    return result, nil
}
/*
	ListAll
	Admin melihat semua prestasi + filter + pagination
*/
func (s *AdminAchievementService) ListAll(
	ctx context.Context,
	params map[string]string,
) ([]model.AchievementRefe, int, error) {

	status := params["status"]
	studentID := params["student_id"]

	page := 1
	limit := 10

	// pagination sederhana
	if p, err := strconv.Atoi(params["page"]); err == nil {
		page = p
	}
	if l, err := strconv.Atoi(params["limit"]); err == nil {
		limit = l
	}

	return s.refRepo.FindAll(ctx, status, studentID, page, limit)
}
