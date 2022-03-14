package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type RecordNameObj struct {
	Name string `json:"name"`
	// TimeStamp time.Time
}

func (db *MongoClient) StoreRecord(recordObj RecordNameObj) error {
	collection := db.registerColl
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, recordObj)
	if err != nil || res == nil {
		log.Print("\nunable to insert entry into DB in database package\n")
		return err
	}
	return nil
}

// check if name is available, if available return true
func (db *MongoClient) CheckName(name string) (bool, error) {
	collection := db.registerColl
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return false, err
	}

	if result == 0 {
		return true, nil
	}
	return false, nil
}
