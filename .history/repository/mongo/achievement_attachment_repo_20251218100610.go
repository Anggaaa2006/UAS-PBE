package mongo

import (
	"context"
	"mime/multipart"
)

/*
	AchievementAttachmentRepo
	Menyimpan metadata file attachment
*/
type AchievementAttachmentRepo interface {
	Save(ctx context.Context, achievementID string, file *multipart.FileHeader) error
	List(ctx context.Context, achievementID string) ([]string, error)
}
