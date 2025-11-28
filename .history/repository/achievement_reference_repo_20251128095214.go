package repository

import (
	"context"
	"database/sql"
	"errors"
	"UAS PBE/model"
)

type AchievementReferenceRepo interface {
	Create(ctx context.Context, studentID, title, categoryID string) (string, error)
	GetByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error)
	GetByID(ctx context.Context, id string) (*model.AchievementReference, error)
	Update(ctx context.Context, id string, title string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Reject(ctx context.Context, id string, reason string) error
	CountByStatus(ctx context.Context) (interface{}, error)
}

type achievementReferenceRepo struct{ db *sql.DB }

func NewAchievementReferenceRepo(db *sql.DB) AchievementReferenceRepo {
	return &achievementReferenceRepo{db}
}
