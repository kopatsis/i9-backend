package dbinput

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserCount(database *mongo.Database) (int64, error) {
	collection := database.Collection("user")

	filter := bson.D{}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return -1, err
	}
	return count, nil
}
