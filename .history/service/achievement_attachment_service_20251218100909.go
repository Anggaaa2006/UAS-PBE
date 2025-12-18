package service

import (
	"context"
	"mime/multipart"

	"uas_pbe/repository/mongo"
)

/*
	AchievementAttachmentService
	Logika upload & baca attachment
*/
type AchievementAttachmentService struct {
	repo mongo.AchievementAttachmentRepo
}

func NewAchievementAttachmentService(
	r mongo.AchievementAttachmentRepo,
) *AchievementAttachmentService {
	return &AchievementAttachmentService{repo: r}
}

/*
	Upload attachment
*/
func (s *AchievementAttachmentService) Upload(
	ctx context.Context,
	achievementID string,
	file *multipart.FileHeader,
) error {
	return s.repo.Save(ctx, achievementID, file)
}

/*
	List attachment
*/
func (s *AchievementAttachmentService) List(
	ctx context.Context,
	achievementID string,
) ([]string, error) {
	return s.repo.List(ctx, achievementID)
}
/*
	ListAll
	Admin melihat semua prestasi + filter + pagination
*/
func (s *AdminAchievementService) ListAll(
	ctx context.Context,
	params map[string]string,
) ([]model.Achievement, int, error) {

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
