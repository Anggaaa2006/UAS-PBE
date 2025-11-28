package model

type AchievementReference struct {
	ID         string `json:"id"`
	StudentID  string `json:"student_id"`
	Title      string `json:"title"`
	CategoryID string `json:"category_id"`
	Status     string `json:"status"`
}
