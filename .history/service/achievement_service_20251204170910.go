package service

import (
	"context"
	"errors"
	"time"

	"uas_pbe/model"
	"uas_pbe/repository/mongo"
	"uas_pbe/repository/postgres"
)

/*
	AchievementService
	Menangani semua bisnis proses prestasi:
	- Create draft
	- Update draft
	- Soft delete (status = deleted)
	- Submit (status = submitted)
	- Approve / Reject (lecture only)
	- Get detail by ID
*/
type AchievementService struct {
	achRepo     postgres.AchievementReferenceRepo
	detailRepo  mongo.AchievementDetailRepo
}

func NewAchievementService(
	achRepo postgres.AchievementReferenceRepo,
	detailRepo mongo.AchievementDetailRepo,
) *AchievementService {
	return &AchievementService{
		achRepo:    achRepo,
		detailRepo: detailRepo,
	}
}

/*
	1. CREATE ACHIEVEMENT (DRAFT)
	Flow:
	- Insert detail ke MongoDB
	- Insert reference ke PostgreSQL (status = draft)
*/
func (s *AchievementService) Create(ctx context.Context, studentID string, req model.AchievementDetail) (string, error) {

	// 1. Simpan detail di MongoDB
	mongoID, err := s.detailRepo.Create(ctx, req)
	if err != nil {
		return "", err
	}

	// 2. Buat reference di PostgreSQL
	ref := model.AchievementReference{
		ID:                 .NewUUID(),
		StudentID:          studentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
		CreatedAt:          time.Now(),
	}

	err = s.achRepo.Create(ctx, ref)
	if err != nil {
		return "", err
	}

	return ref.ID, nil
}

/*
	2. GET DETAIL BY ID
	Flow:
	- Ambil reference dari PostgreSQL
	- Ambil detail dari MongoDB
*/
func (s *AchievementService) GetByID(ctx context.Context, id string) (interface{}, error) {

	ref, err := s.achRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Ambil detail Mongo
	detail, err := s.detailRepo.GetByID(ctx, ref.MongoAchievementID)
	if err != nil {
		return nil, err
	}

	// Gabungkan jadi satu response
	return gin.H{
		"id":        ref.ID,
		"student_id": ref.StudentID,
		"status":     ref.Status,
		"detail":     detail,
	}, nil
}

/*
	3. UPDATE (Hanya boleh ketika status = draft)
*/
func (s *AchievementService) Update(ctx context.Context, id string, req model.AchievementDetail) error {

	ref, err := s.achRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if ref.Status != "draft" {
		return errors.New("hanya draft yang dapat diupdate")
	}

	return s.detailRepo.Update(ctx, ref.MongoAchievementID, req)
}

/*
	4. DELETE (Soft delete → status = deleted)
*/
func (s *AchievementService) Delete(ctx context.Context, id string) error {

	ref, err := s.achRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if ref.Status != "draft" {
		return errors.New("hanya draft yang dapat dihapus")
	}

	// Update detail di Mongo → is_deleted = true
	err = s.detailRepo.SoftDelete(ctx, ref.MongoAchievementID)
	if err != nil {
		return err
	}

	// Update status di PostgreSQL → deleted
	return s.achRepo.UpdateStatus(ctx, id, "deleted", "")
}

/*
	5. SUBMIT (status: draft → submitted)
*/
func (s *AchievementService) Submit(ctx context.Context, id string, studentID string) error {

	ref, err := s.achRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Hanya pemilik prestasi
	if ref.StudentID != studentID {
		return errors.New("tidak dapat submit prestasi milik orang lain")
	}

	if ref.Status != "draft" {
		return errors.New("hanya draft yang dapat disubmit")
	}

	return s.achRepo.UpdateStatus(ctx, id, "submitted", "")
}

/*
	6. APPROVE PRESTASI → dosen
*/
func (s *AchievementService) Approve(ctx context.Context, id string, lectureID string) error {

	ref, err := s.achRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if ref.Status != "submitted" {
		return errors.New("prestasi belum disubmit")
	}

	return s.achRepo.Verify(ctx, id, lectureID, "verified")
}

/*
	7. REJECT
*/
func (s *AchievementService) Reject(ctx context.Context, id string, lectureID string, note string) error {

	ref, err := s.achRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if ref.Status != "submitted" {
		return errors.New("prestasi belum disubmit")
	}

	return s.achRepo.Reject(ctx, id, lectureID, note)
}

/*
	List By Student
*/
func (s *AchievementService) ListByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	return s.achRepo.ListByStudent(ctx, studentID)
}
