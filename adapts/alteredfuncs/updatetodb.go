package alteredfuncs

import (
	"context"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveUpdatedWorkout(database *mongo.Database, workout shared.Workout) error {
	collection := database.Collection("workout")
	filter := bson.M{"_id": workout.ID}
	_, err := collection.ReplaceOne(context.TODO(), filter, workout)
	if err != nil {
		return err
	}
	return nil
}

func SaveUpdatedStretchWorkout(database *mongo.Database, workout shared.StretchWorkout) error {
	collection := database.Collection("stretchworkout")
	filter := bson.M{"_id": workout.ID}
	_, err := collection.ReplaceOne(context.TODO(), filter, workout)
	if err != nil {
		return err
	}
	return nil
}
