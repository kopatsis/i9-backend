package dboutput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveUser(user shared.User, database *mongo.Database) error {
	collection := database.Collection("user")

	filter := bson.M{"_id": user.ID}

	update := bson.M{"$set": user}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
