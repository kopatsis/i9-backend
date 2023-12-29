package dbinput

import (
	"context"
	"fmt"
	"fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserDB(database *mongo.Database, username string) datatypes.User {
	collection := database.Collection("user")
	filter := bson.D{{Key: "username", Value: username}}

	// Define options for the query.
	optionsUser := options.FindOne()

	// Perform the query.
	var result datatypes.User
	err := collection.FindOne(context.Background(), filter, optionsUser).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return datatypes.User{}
	}

	return result

}
