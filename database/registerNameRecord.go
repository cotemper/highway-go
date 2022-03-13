package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *Mongo) StoreRecord(creator string, name string) error {
	collection := db.registerColl
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newID := primitive.NewObjectID().Hex()

	newEntry := struct {
		ID             string
		Creator        string
		NameToRegister string
		TimeStamp      time.Time
	}{
		ID:             newID,
		Creator:        creator,
		NameToRegister: name,
		TimeStamp:      time.Now(),
	}

	res, err := collection.InsertOne(ctx, newEntry)
	if err != nil || res == nil {
		log.Print("\nunable to insert entry into DB in database package\n")
		return err
	}
	return nil
}
