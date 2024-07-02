package alteredfuncs

import (
	"context"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWorkoutByID(database *mongo.Database, workoutID string) (shared.Workout, error) {
	collection := database.Collection("workout")

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(workoutID); err == nil {
		id = oid
	} else {
		return shared.Workout{}, err
	}

	filter := bson.D{{Key: "_id", Value: id}}

	// Perform the query.
	var result shared.Workout
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return shared.Workout{}, err
	}

	return result, nil
}

func GetStretchWorkoutByID(database *mongo.Database, workoutID string) (shared.StretchWorkout, error) {
	collection := database.Collection("stretchworkout")

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(workoutID); err == nil {
		id = oid
	} else {
		return shared.StretchWorkout{}, err
	}

	filter := bson.D{{Key: "_id", Value: id}}

	// Perform the query.
	var result shared.StretchWorkout
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return shared.StretchWorkout{}, err
	}

	return result, nil
}
