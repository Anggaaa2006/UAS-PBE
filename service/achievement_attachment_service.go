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
