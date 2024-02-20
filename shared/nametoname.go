package shared

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthIDtoMongoID(database *mongo.Database, authID string) (string, error) {
	collection := database.Collection("user")

	var user User
	err := collection.FindOne(context.Background(), bson.M{"username": authID}).Decode(&user)
	if err != nil {
		return "", err
	}

	return user.ID.Hex(), nil
}
