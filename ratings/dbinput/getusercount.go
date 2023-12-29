package dbinput

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserCount(database *mongo.Database) int64 {
	collection := database.Collection("user")

	filter := bson.D{}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return count
}
