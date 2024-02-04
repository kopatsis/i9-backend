package alteredfuncs

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWorkoutByID(database *mongo.Database, workoutID string) shared.Workout {
	collection := database.Collection("workout")

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(workoutID); err == nil {
		id = oid
	} else {
		// Handle error here, sorry
		return shared.Workout{}
	}

	filter := bson.D{{Key: "_id", Value: id}}

	// Perform the query.
	var result shared.Workout
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return shared.Workout{}
	}

	return result
}
