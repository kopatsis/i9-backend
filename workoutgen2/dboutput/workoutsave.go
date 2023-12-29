package dboutput

import (
	"context"
	"fmt"
	"fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/mongo"
)

func SaveNewWorkout(database *mongo.Database, workout datatypes.Workout) error {
	collection := database.Collection("workout")
	_, err := collection.InsertOne(context.Background(), workout)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
