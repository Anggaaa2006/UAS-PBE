package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	NewMongoConn
	Membuat koneksi MongoDB
*/
func NewMongoConn(ctx context.Context, uri string) (*mongo.Client, error) {

	// Membuat client MongoDB dengan URI
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Timeout 10 detik
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect ke server MongoDB
	if err := client.Connect(ctx2); err != nil {
		return nil, err
	}

	return client, nil
}
