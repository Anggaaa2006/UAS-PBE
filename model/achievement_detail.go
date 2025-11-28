package model

type AchievementDetail struct {
	RefID       string `bson:"ref_id" json:"ref_id"`
	Description string `bson:"description" json:"description"`
	IsDeleted   bool   `bson:"is_deleted" json:"is_deleted"`
}
