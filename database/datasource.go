package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client       *mongo.Client
	registerColl *mongo.Collection
}

func Connect(mongoURI string, collection string, mongoName string) (*MongoClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		//log.Warn().Err(err).Msg("unable to connect to mongo database")
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.Connect(ctx)
	return &MongoClient{
		client:       client,
		registerColl: client.Database(mongoName).Collection(collection),
	}, nil
}

func (db *MongoClient) Disconnect() {
	db.client.Disconnect(context.Background())
}
