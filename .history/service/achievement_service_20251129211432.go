package service

import (
	"context"
	"errors"

	mg "uas_pbe/internal/repository/mongo"
	pg "uas_pbe/internal/repository/postgres"
)

// ======== REQUEST STRUCTS (DTO) ========
type AchievementCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id"`
}

type AchievementUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AchievementService: menangani seluruh logika prestasi
type AchievementService struct {
	mongoRepo mg.AchievementRepo          // detail prestasi (file, deskripsi panjang)
	refRepo   pg.AchievementReferenceRepo // metadata prestasi + status (postgres)
}

func NewAchievementService(m mg.AchievementRepo, r pg.AchievementReferenceRepo) *AchievementService {
	return &AchievementService{
		mongoRepo: m,
		refRepo:   r,
	}
}

// ====================
// LIST PRESTASI
// ====================
func (s *AchievementService) List(ctx context.Context, studentID string) (interface{}, error) {
	return s.refRepo.GetByStudent(ctx, studentID)
}

// ====================
// CREATE (status: draft)
// ====================
func (s *AchievementService) Create(ctx context.Context, studentID string, req AchievementCreateRequest) error {
	// simpan metadata di postgres (status draft)
	refID, err := s.refRepo.Create(ctx, studentID, req.Title, req.CategoryID)
	if err != nil {
		return err
	}

	// simpan detail di mongo
	return s.mongoRepo.Create(ctx, refID, req.Description)
}

// ====================
// UPDATE (hanya draft)
// ====================
func (s *AchievementService) Update(ctx context.Context, studentID, id string, req AchievementUpdateRequest) error {
	// cek status: harus draft
	ref, err := s.refRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ref.Status != "draft" {
		return errors.New("hanya prestasi draft yang bisa diupdate")
	}

	// update postgres
	err = s.refRepo.Update(ctx, id, req.Title)
	if err != nil {
		return err
	}

	// update mongo
	return s.mongoRepo.Update(ctx, id, req.Description)
}

// ====================
// SOFT DELETE (status: deleted)
// ====================
func (s *AchievementService) SoftDelete(ctx context.Context, studentID, id string) error {
	// ubah status postgres → deleted
	err := s.refRepo.UpdateStatus(ctx, id, "deleted")
	if err != nil {
		return err
	}

	// tandai is_deleted di mongo
	return s.mongoRepo.MarkDeleted(ctx, id)
}

// ====================
// SUBMIT (draft → submitted)
// ====================
func (s *AchievementService) Submit(ctx context.Context, studentID, id string) error {
	ref, err := s.refRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ref.Status != "draft" {
		return errors.New("hanya draft yang bisa disubmit")
	}

	return s.refRepo.UpdateStatus(ctx, id, "submitted")
}

// ====================
// VERIFY (submitted → verified)
// ====================
func (s *AchievementService) Verify(ctx context.Context, verifier, id string) error {
	ref, err := s.refRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ref.Status != "submitted" {
		return errors.New("hanya prestasi submitted yang bisa diverifikasi")
	}

	return s.refRepo.UpdateStatus(ctx, id, "verified")
}

// ====================
// REJECT (submitted → rejected)
// ====================
func (s *AchievementService) Reject(ctx context.Context, verifier, id, reason string) error {
	ref, err := s.refRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ref.Status != "submitted" {
		return errors.New("hanya prestasi submitted yang bisa ditolak")
	}

	return s.refRepo.Reject(ctx, id, reason)
}
