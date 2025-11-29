package mongo

import (
	"context"
	"errors"

	"uas_pbe/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
	AchievementDetailRepo interface
	Mengatur detail prestasi dalam MongoDB
*/
type AchievementDetailRepo interface {
	Create(ctx context.Context, refID string, description string) error
	Update(ctx context.Context, refID string, description string) error
	MarkDeleted(ctx context.Context, refID string) error
	GetByRefID(ctx context.Context, refID string) (*model.AchievementDetail, error)
}

/*
	Struct implementasi repository
*/
type achievementDetailRepo struct {
	col *mongo.Collection
}

/*
	NewAchievementDetailRepo
	Constructor repository Mongo
*/
func NewAchievementDetailRepo(db *mongo.Database) AchievementDetailRepo {
	return &achievementDetailRepo{
		col: db.Collection("achievement_details"), // nama collection
	}
}

/*
	Create
	Membuat detail prestasi baru di MongoDB
*/
func (r *achievementDetailRepo) Create(ctx context.Context, refID string, description string) error {
	doc := model.AchievementDetail{
		RefID:       refID,
		Description: description,
		IsDeleted:   false,
	}

	_, err := r.col.InsertOne(ctx, doc)
	return err
}

/*
	Update
	Mengupdate deskripsi prestasi
*/
func (r *achievementDetailRepo) Update(ctx context.Context, refID string, description string) error {
	filter := bson.M{"ref_id": refID}
	update := bson.M{"$set": bson.M{"description": description}}

	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

/*
	MarkDeleted
	Soft delete -> is_deleted = true
*/
func (r *achievementDetailRepo) MarkDeleted(ctx context.Context, refID string) error {
	filter := bson.M{"ref_id": refID}
	update := bson.M{"$set": bson.M{"is_deleted": true}}

	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

/*
	GetByRefID
	Mengambil detail prestasi berdasarkan ID referensi PostgreSQL
*/
func (r *achievementDetailRepo) GetByRefID(ctx context.Context, refID string) (*model.AchievementDetail, error) {
	filter := bson.M{"ref_id": refID}

	var result model.AchievementDetail
	err := r.col.FindOne(ctx, filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("detail prestasi tidak ditemukan")
	}
	return &result, err
}
