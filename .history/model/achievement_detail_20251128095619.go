package model

type AchievementDetail struct {
	RefID       string `bson:"ref_id"`        // ID referensi dari PostgreSQL
	Description string `bson:"description"`
	IsDeleted   bool   `bson:"is_deleted"`
}
