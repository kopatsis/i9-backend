package dbinput

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserWorkoutCount(database *mongo.Database, username string) int64 {
	collection := database.Collection("workout")

	filter := bson.D{{Key: "username", Value: username}}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return count
}
