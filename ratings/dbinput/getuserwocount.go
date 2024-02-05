package dbinput

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserWorkoutCount(database *mongo.Database, userID string) (int64, error) {
	collection := database.Collection("workout")

	filter := bson.D{{Key: "userid", Value: userID}}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return -1, err
	}
	return count, nil
}
