package service

import (
	"context"
	"testing"

	"uas_pbe/model"

	"github.com/stretchr/testify/assert"
)

//
// ================================
// MOCK REPOSITORY
// ================================
// Digunakan agar unit test TIDAK
// bergantung pada PostgreSQL & MongoDB
//

// -------------------------------
// Mock AchievementReferenceRepo
// -------------------------------
type mockAchievementRepo struct{}

// Create (dipakai oleh Create Achievement)
func (m *mockAchievementRepo) Create(
	ctx context.Context,
	ref model.AchievementReference,
) error {
	return nil
}

// CountByStudent (dipakai dashboard mahasiswa)
func (m *mockAchievementRepo) CountByStudent(
	ctx context.Context,
	studentID string,
) (map[string]int, error) {
	return map[string]int{
		"draft":     2,
		"submitted": 1,
		"approved":  3,
		"rejected":  0,
	}, nil
}

// Method lain tidak perlu dibuat
// selama tidak dipanggil oleh service
// yang sedang ditest

// -------------------------------
// Mock AchievementDetailRepo
// -------------------------------
type mockDetailRepo struct{}

// Create (dipakai oleh Create Achievement)
func (m *mockDetailRepo) Create(
	ctx context.Context,
	detail model.AchievementDetail,
) (string, error) {
	return "mock-mongo-id", nil
}

//
// ================================
// UNIT TEST
// ================================
//

// --------------------------------
// Test 1: Create Achievement
// --------------------------------
func TestCreateAchievement(t *testing.T) {
	achRepo := &mockAchievementRepo{}
	detailRepo := &mockDetailRepo{}

	service := NewAchievementService(achRepo, detailRepo)

	id, err := service.Create(
		context.Background(),
		"student-123",
		model.AchievementDetail{
			Description: "Juara 1 Lomba Coding",
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

// --------------------------------
// Test 2: Dashboard Mahasiswa
// --------------------------------
func TestGetStudentSummary(t *testing.T) {
	achRepo := &mockAchievementRepo{}
	detailRepo := &mockDetailRepo{}

	service := NewAchievementService(achRepo, detailRepo)

	result, err := service.GetStudentSummary(
		context.Background(),
		"student-123",
	)

	assert.NoError(t, err)
	assert.Equal(t, 2, result["draft"])
	assert.Equal(t, 1, result["submitted"])
	assert.Equal(t, 3, result["approved"])
	assert.Equal(t, 0, result["rejected"])
}
