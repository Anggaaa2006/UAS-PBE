package model

type AchievementFullResponse struct {
    ID        string             `json:"id"`
    StudentID string             `json:"student_id"`
    Status    string             `json:"status"`
    Detail    *AchievementDetail  `json:"detail"`
}
