package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoConn membuat koneksi ke MongoDB dan mengembalikan *mongo.Client
func NewMongoConn(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Timeout 10 detik
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = client.Connect(ctx2)
	if err != nil {
		return nil, err
	}

	return client, nil
}
