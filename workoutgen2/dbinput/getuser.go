package dbinput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserDB(database *mongo.Database, username string) shared.User {
	collection := database.Collection("user")
	filter := bson.D{{Key: "username", Value: username}}

	// Define options for the query.
	optionsUser := options.FindOne()

	// Perform the query.
	var result shared.User
	err := collection.FindOne(context.Background(), filter, optionsUser).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return shared.User{}
	}

	return result

}
