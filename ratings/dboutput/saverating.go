package dboutput

import (
	"context"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/mongo"
)

func SaveRating(databaseRating shared.StoredRating, database *mongo.Database) error {
	collection := database.Collection("ratings")
	_, err := collection.InsertOne(context.TODO(), databaseRating)
	return err
}
