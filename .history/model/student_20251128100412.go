package model

type Student struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	NIM      string `json:"nim"`
	Major    string `json:"major"`
	Semester int    `json:"semester"`
}
