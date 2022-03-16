package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	users  *mongo.Collection
	auths  *mongo.Collection
	creds  *mongo.Collection
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
		client: client,
		users:  client.Database(mongoName).Collection("users"),
		auths:  client.Database(mongoName).Collection("auths"),
		creds:  client.Database(mongoName).Collection("creds"),
	}, nil
}

func (db *MongoClient) Disconnect() {
	db.client.Disconnect(context.Background())
}
