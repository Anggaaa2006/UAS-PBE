package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	InitMongo
	Menghubungkan aplikasi ke MongoDB
*/
func InitMongo(cfg Config) (*mongo.Database, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURL))
	if err != nil {
		return nil, err
	}

	// Connect dengan timeout 10 detik
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.MongoDBName), nil
}
