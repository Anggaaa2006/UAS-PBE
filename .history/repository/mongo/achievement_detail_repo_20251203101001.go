package mongo

import (
	"context"
	"errors"

	"uas_pbe/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
	AchievementDetailRepo
	Repository untuk MongoDB.
	Menyimpan detail prestasi (deskripsi, bukti, lampiran)
*/
type AchievementDetailRepo interface {
	Create(ctx context.Context, detail model.AchievementDetail) (string, error)
	GetByID(ctx context.Context, id string) (*model.AchievementDetail, error)
	Update(ctx context.Context, id string, req model.AchievementDetail) error
	SoftDelete(ctx context.Context, id string) error
}

type achievementDetailRepo struct {
	col *mongo.Collection
}

func NewAchievementDetailRepo(db *mongo.Database) AchievementDetailRepo {
	return &achievementDetailRepo{
		col: db.Collection("achievement_details"),
	}
}

/*
	Create detail prestasi ke MongoDB
*/
func (r *achievementDetailRepo) Create(ctx context.Context, detail model.AchievementDetail) (string, error) {

	detail.ID = primitive.NewObjectID()
	detail.IsDeleted = false

	_, err := r.col.InsertOne(ctx, detail)
	if err != nil {
		return "", err
	}

	return detail.ID.Hex(), nil
}

/*
	GetByID detail prestasi berdasarkan MongoID
*/
func (r *achievementDetailRepo) GetByID(ctx context.Context, id string) (*model.AchievementDetail, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid mongo id")
	}

	var detail model.AchievementDetail
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&detail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("detail tidak ditemukan")
		}
		return nil, err
	}

	return &detail, nil
}

/*
	Update detail prestasi (hanya draft)
*/
func (r *achievementDetailRepo) Update(ctx context.Context, id string, req model.AchievementDetail) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid mongo id")
	}

	update := bson.M{
		"$set": bson.M{
			"description": req.Description,
		},
	}

	_, err = r.col.UpdateByID(ctx, objID, update)
	return err
}

/*
	SoftDelete
	Menandai detail sebagai "deleted" tanpa menghapus dokumen
*/
func (r *achievementDetailRepo) SoftDelete(ctx context.Context, id string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid mongo id")
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
		},
	}

	_, err = r.col.UpdateByID(ctx, objID, update)

	return err
}
