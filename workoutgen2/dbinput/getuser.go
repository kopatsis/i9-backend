package dbinput

import (
	"context"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserDB(database *mongo.Database, userID string) (shared.User, error) {
	collection := database.Collection("user")

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		id = oid
	} else {
		return shared.User{}, err
	}

	filter := bson.D{{Key: "_id", Value: id}}

	// Define options for the query.
	optionsUser := options.FindOne()

	// Perform the query.
	var result shared.User
	err := collection.FindOne(context.Background(), filter, optionsUser).Decode(&result)
	if err != nil {
		return shared.User{}, err
	}

	return result, nil

}
