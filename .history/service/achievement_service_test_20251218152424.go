package service

import (
	"context"
	"testing"

	"uas_pbe/model"

	"github.com/stretchr/testify/assert"
)

//
// ================================
// MOCK POSTGRES REPO
// ================================
type mockAchievementRepo struct{}

// ===== methods WAJIB agar implement interface =====

func (m *mockAchievementRepo) Create(
	ctx context.Context,
	ref model.AchievementReference,
) error {
	return nil
}

func (m *mockAchievementRepo) GetByID(
	ctx context.Context,
	id string,
) (*model.AchievementReference, error) {
	return &model.AchievementReference{
		ID:                 id,
		StudentID:          "student-123",
		MongoAchievementID: "mongo-id",
		Status:             "draft",
	}, nil
}

func (m *mockAchievementRepo) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
	note string,
) error {
	return nil
}

func (m *mockAchievementRepo) Verify(
	ctx context.Context,
	id string,
	lecturerID string,
	note string,
) error {
	return nil
}

func (m *mockAchievementRepo) Reject(
	ctx context.Context,
	id string,
	lecturerID string,
	note string,
) error {
	return nil
}

func (m *mockAchievementRepo) ListByStudent(
	ctx context.Context,
	studentID string,
) ([]model.AchievementReference, error) {
	return []model.AchievementReference{}, nil
}

func (m *mockAchievementRepo) ListAll(
	ctx context.Context,
) ([]model.AchievementReference, error) {
	return []model.AchievementReference{}, nil
}

func (m *mockAchievementRepo) GetHistory(
	ctx context.Context,
	id string,
) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

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

func (m *mockAchievementRepo) CountByStatus(
	ctx context.Context,
) (map[string]int, error) {
	return map[string]int{
		"submitted": 2,
		"approved":  5,
		"rejected":  1,
	}, nil
}

//
// ================================
// MOCK MONGO REPO
// ================================
type mockDetailRepo struct{}

func (m *mockDetailRepo) Create(
	ctx context.Context,
	detail model.AchievementDetail,
) (string, error) {
	return "mock-mongo-id", nil
}

func (m *mockDetailRepo) GetByID(
	ctx context.Context,
	id string,
) (*model.AchievementDetail, error) {
	return &model.AchievementDetail{
		Description: "Prestasi Mock",
	}, nil
}

func (m *mockDetailRepo) Update(
	ctx context.Context,
	id string,
	req model.AchievementDetail,
) error {
	return nil
}

func (m *mockDetailRepo) SoftDelete(
	ctx context.Context,
	id string,
) error {
	return nil
}

//
// ================================
// UNIT TEST
// ================================

// --------------------------------
// Test Create Achievement
// --------------------------------
func TestCreateAchievement(t *testing.T) {
	achRepo := &mockAchievementRepo{}
	detailRepo := &mockDetailRepo{}

	svc := NewAchievementService(achRepo, detailRepo)

	id, err := svc.Create(
		context.Background(),
		"student-123",
		model.AchievementDetail{
			Description: "Juara 1 Nasional",
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

// --------------------------------
// Test Dashboard Mahasiswa
// --------------------------------
func TestGetStudentSummary(t *testing.T) {
	achRepo := &mockAchievementRepo{}
	detailRepo := &mockDetailRepo{}

	svc := NewAchievementService(achRepo, detailRepo)

	result, err := svc.GetStudentSummary(
		context.Background(),
		"student-123",
	)

	assert.NoError(t, err)
	assert.Equal(t, 2, result["draft"])
	assert.Equal(t, 1, result["submitted"])
	assert.Equal(t, 3, result["approved"])
	assert.Equal(t, 0, result["rejected"])
}
