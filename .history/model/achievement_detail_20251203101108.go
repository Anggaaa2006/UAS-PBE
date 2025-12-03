package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AchievementDetail struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	RefID       string             `bson:"ref_id"`
	Description string             `bson:"description"`
	IsDeleted   bool               `bson:"is_deleted"`
}
