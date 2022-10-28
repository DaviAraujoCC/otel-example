package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MONGO_URI = os.Getenv("MONGO_URI")

func PingDB() error {
	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	defer client.Disconnect(context.Background())

	return nil
}