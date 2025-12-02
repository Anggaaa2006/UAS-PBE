package service

import (
	"context"
	"errors"
	"uas_pbe/model"

	mongoRepo "uas_pbe/repository/mongo"
	postgresRepo "uas_pbe/repository/postgres"
)

/*
	AchievementService
	Interface yang mendefinisikan semua operasi bisnis terkait prestasi.
*/
type AchievementService interface {
	Create(ctx context.Context, studentID, title, categoryID, description string) (*model.AchievementReference, error)
	Update(ctx context.Context, id, title, categoryID, description string) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.AchievementReference, *model.AchievementDetail, error)
	ListStudentAchievement(ctx context.Context, studentID string) ([]model.AchievementReference, error)
	Submit(ctx context.Context, id string) error
	Approve(ctx context.Context, id string) error
	Reject(ctx context.Context, id string) error
}

/*
	Struct implementasi service
*/
type achievementService struct {
	refRepo    postgresRepo.AchievementReferenceRepo
	detailRepo mongoRepo.AchievementDetailRepo
}

/*
	NewAchievementService
	Constructor service prestasi
*/
func NewAchievementService(
	refRepo postgresRepo.AchievementReferenceRepo,
	detailRepo mongoRepo.AchievementDetailRepo,
) AchievementService {
	return &achievementService{
		refRepo:    refRepo,
		detailRepo: detailRepo,
	}
}

/*
	Create
	Business logic membuat prestasi:
	1. Insert ke PostgreSQL (judul, kategori, student)
	2. Insert detail ke MongoDB (desc, files, dll)
*/
func (s *achievementService) Create(ctx context.Context, studentID, title, categoryID, description string) (*model.AchievementReference, error) {

	// 1. Insert metadata ke Postgres
	refID, err := s.refRepo.Create(ctx, studentID, title, categoryID)
	if err != nil {
		return nil, err
	}

	// 2. Insert detail ke MongoDB
	if err := s.detailRepo.Create(ctx, refID, description); err != nil {
		return nil, err
	}

	// 3. Ambil kembali data untuk dikirim ke controller
	refData, err := s.refRepo.GetByID(ctx, refID)
	if err != nil {
		return nil, err
	}

	return refData, nil
}

/*
	Update
	Business logic update prestasi:
	1. Update title & category di PostgreSQL
	2. Update description di MongoDB
*/
func (s *achievementService) Update(ctx context.Context, id, title, categoryID, description string) error {

	// update metadata
	if err := s.refRepo.Update(ctx, id, title, categoryID); err != nil {
		return err
	}

	// update detail
	if err := s.detailRepo.Update(ctx, id, description); err != nil {
		return err
	}

	return nil
}

/*
	Delete
	Soft delete:
	1. PostgreSQL → status = deleted
	2. Mongo → is_deleted = true
*/
func (s *achievementService) Delete(ctx context.Context, id string) error {

	// soft delete postgres
	if err := s.refRepo.SoftDelete(ctx, id); err != nil {
		return err
	}

	// soft delete mongo
	if err := s.detailRepo.MarkDeleted(ctx, id); err != nil {
		return err
	}

	return nil
}

/*
	GetByID
	Mengambil 2 data sekaligus:
	1. Data metadata (Postgres)
	2. Data detail (Mongo)
*/
func (s *achievementService) GetByID(ctx context.Context, id string) (*model.AchievementReference, *model.AchievementDetail, error) {

	ref, err := s.refRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, errors.New("data prestasi tidak ditemukan")
	}

	detail, err := s.detailRepo.GetByRefID(ctx, id)
	if err != nil {
		return ref, nil, errors.New("detail prestasi tidak ditemukan")
	}

	return ref, detail, nil
}

/*
	ListStudentAchievement
	Mengambil semua prestasi untuk 1 mahasiswa
*/
func (s *achievementService) ListStudentAchievement(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	return s.refRepo.GetByStudentID(ctx, studentID)
}

/*
	Submit
	Mengubah status prestasi → submitted
*/
func (s *achievementService) Submit(ctx context.Context, id string) error {
	return s.refRepo.UpdateStatus(ctx, id, "submitted")
}

/*
	Approve
	Mengubah status prestasi → approved
*/
func (s *achievementService) Approve(ctx context.Context, id string) error {
	return s.refRepo.UpdateStatus(ctx, id, "approved")
}

/*
	Reject
	Mengubah status prestasi → rejected
*/
func (s *achievementService) Reject(ctx context.Context, id string) error {
	return s.refRepo.UpdateStatus(ctx, id, "rejected")
}
