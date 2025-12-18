package service

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "uas_pbe/model"
    "uas_pbe/repository/mongo"
    "uas_pbe/repository/postgres"
)

/*
    AchievementService
    - Create
    - Update
    - Delete (soft delete)
    - Submit
    - Approve
    - Reject
    - GetByID
*/
type AchievementService struct {
    achRepo    postgres.AchievementReferenceRepo
    detailRepo mongo.AchievementDetailRepo
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
    1. CREATE (DRAFT)
*/
func (s *AchievementService) Create(ctx context.Context, studentID string, req model.AchievementDetail) (string, error) {

    // Simpan detail di Mongo
    mongoID, err := s.detailRepo.Create(ctx, req)
    if err != nil {
        return "", err
    }

    // Simpan reference di Postgre
    ref := model.AchievementReference{
        ID:                 uuid.NewString(),
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
    2. GET BY ID (Mengambil reference + detail)
    RETURN: gabungan data
    NOTE:
    - Service TIDAK boleh return gin.H → controller yang bentuk response
*/
func (s *AchievementService) GetByID(ctx context.Context, id string) (*model.AchievementFullResponse, error) {

    ref, err := s.achRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    detail, err := s.detailRepo.GetByID(ctx, ref.MongoAchievementID)
    if err != nil {
        return nil, err
    }

    // Data digabung dalam struct response model
    return &model.AchievementFullResponse{
        ID:        ref.ID,
        StudentID: ref.StudentID,
        Status:    ref.Status,
        Detail:    detail,
    }, nil
}

/*
    3. UPDATE (hanya draft)
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
    4. DELETE (soft delete)
*/
func (s *AchievementService) Delete(ctx context.Context, id string) error {

    ref, err := s.achRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    if ref.Status != "draft" {
        return errors.New("hanya draft yang dapat dihapus")
    }

    // Soft delete di Mongo
    err = s.detailRepo.SoftDelete(ctx, ref.MongoAchievementID)
    if err != nil {
        return err
    }

    // Status → deleted
    return s.achRepo.UpdateStatus(ctx, id, "deleted", "")
}

/*
    5. SUBMIT
*/
func (s *AchievementService) Submit(ctx context.Context, id string, studentID string) error {

    ref, err := s.achRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    if ref.StudentID != studentID {
        return errors.New("tidak dapat submit prestasi milik orang lain")
    }

    if ref.Status != "draft" {
        return errors.New("hanya draft yang dapat disubmit")
    }

    return s.achRepo.UpdateStatus(ctx, id, "submitted", "")
}

/*
    6. APPROVE
*/
func (s *AchievementService) Approve(ctx context.Context, id string, lecturerID string) error {

    ref, err := s.achRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    if ref.Status != "submitted" {
        return errors.New("prestasi belum disubmit")
    }

    // VERIFY = status approved
    return s.achRepo.Verify(ctx, id, lecturerID, "")
}

/*
    7. REJECT
*/
func (s *AchievementService) Reject(ctx context.Context, id string, lecturerID string, note string) error {

    ref, err := s.achRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    if ref.Status != "submitted" {
        return errors.New("prestasi belum disubmit")
    }

    return s.achRepo.Reject(ctx, id, lecturerID, note)
}

/*
    8. LIST BY STUDENT
*/
func (s *AchievementService) ListByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
    return s.achRepo.ListByStudent(ctx, studentID)
}
/*
	ListByRole
	Menentukan data prestasi berdasarkan role
*/
func (s *achievementService) ListByRole(
	ctx context.Context,
	userID string,
	role string,
) ([]model.AchievementReference, error) {

	switch role {

	case "student":
		return s.refRepo.ListByStudent(ctx, userID)

	case "lecturer":
		// sementara dosen bisa lihat semua
		// (idealnya filter mahasiswa bimbingan)
		return s.refRepo.ListAll(ctx)

	case "admin":
		return s.refRepo.ListAll(ctx)

	default:
		return nil, errors.New("role tidak valid")
	}
}
func (s *achievementService) GetHistory(
	ctx context.Context,
	id string,
) ([]map[string]interface{}, error) {

	return s.refRepo.GetHistory(ctx, id)
}
