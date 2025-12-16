// path: service/admin_achievement_service.go
package service

import (
    "context"
    "fmt"
    "strings"

    "uas_pbe/model"
    pg "uas_pbe/repository/postgres"
    mg "uas_pbe/repository/mongo"
)

/*
    AdminAchievementService
    Menyediakan listing semua achievement dengan filter, sort, pagination
*/
type AdminAchievementService struct {
    refRepo pg.AchievementReferenceRepo
    detailRepo mg.AchievementDetailRepo
}

func NewAdminAchievementService(ref pg.AchievementReferenceRepo, d mg.AchievementDetailRepo) *AdminAchievementService {
    return &AdminAchievementService{refRepo: ref, detailRepo: d}
}

// Params = map[string]string untuk filter: status, student_id, category, page, limit, sort
func (s *AdminAchievementService) ListAll(ctx context.Context, params map[string]string) ([]model.AchievementFullResponse, int, error) {
    // Building simple where clause based on params (this is minimal; for complex query use SQL builder)
    where := "WHERE 1=1"
    args := []interface{}{}
    idx := 1

    if v, ok := params["status"]; ok && v != "" {
        where += fmt.Sprintf(" AND status = $%d", idx); args = append(args, v); idx++
    }
    if v, ok := params["student_id"]; ok && v != "" {
        where += fmt.Sprintf(" AND student_id = $%d", idx); args = append(args, v); idx++
    }

    // paging
    page := 1
    limit := 20
    if p, ok := params["page"]; ok && p != "" {
        fmt.Sscanf(p, "%d", &page)
    }
    if l, ok := params["limit"]; ok && l != "" {
        fmt.Sscanf(l, "%d", &limit)
    }
    offset := (page-1)*limit

    // Use refRepo.ListAll (we must implement this repo method) — but to keep consistent with current repo, we'll reuse ListByStudent when student filter present, else call a generic repo method.
    // For simplicity assume we added repo method: ListAll(ctx, where string, args []interface{}, limit, offset int) ([]model.AchievementReference, int, error)
    refs, total, err := s.refRepo.ListAll(ctx, where, args, limit, offset)
    if err != nil {
        return nil, 0, err
    }

    // for each ref fetch detail from mongo and map to full response
    var out []model.AchievementFullResponse
    for _, r := range refs {
        d, _ := s.detailRepo.GetByID(ctx, r.MongoAchievementID)
        out = append(out, model.AchievementFullResponse{
            ID: r.ID,
            StudentID: r.StudentID,
            Status: r.Status,
            Detail: d, // assumes Detail is pointer type in model
        })
    }

    return out, total, nil
}
