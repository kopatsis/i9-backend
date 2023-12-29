package dboutput

import (
	"context"
	"fmt"
	"fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveUser(user datatypes.User, newlevel float32, userExMod, userTypeMod map[string]float32, database *mongo.Database) error {
	collection := database.Collection("user")

	filter := bson.M{"_id": user.ID}

	update := bson.M{"$set": bson.M{"level": newlevel, "exermods": userExMod, "typemods": userTypeMod}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
