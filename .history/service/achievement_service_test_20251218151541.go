package service

import (
	"context"
	"testing"

	"uas_pbe/model"

	"github.com/stretchr/testify/assert"
)
type mockAchievementRepo struct{}

func (m *mockAchievementRepo) Create(ctx context.Context, ref model.AchievementReference) error {
	return nil
}

func (m *mockAchievementRepo) CountByStudent(ctx context.Context, studentID string) (map[string]int, error) {
	return map[string]int{
		"draft":     2,
		"submitted": 1,
		"approved":  3,
		"rejected":  0,
	}, nil
}
type mockDetailRepo struct{}

func (m *mockDetailRepo) Create(ctx context.Context, d model.AchievementDetail) (string, error) {
	return "mongo123", nil
}
