package dbinput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetPastWOsDB(database *mongo.Database, idStr string) (shared.Workout, error) {

	collection := database.Collection("workout")

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(idStr); err == nil {
		id = oid
	} else {
		return shared.Workout{}, err
	}

	var workout shared.Workout
	if err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&workout); err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			fmt.Println("No workout for user in database")
			return shared.Workout{}, err
		} else {
			return shared.Workout{}, nil
		}
	}

	return workout, nil

}
