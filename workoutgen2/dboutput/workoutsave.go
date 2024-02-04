package dboutput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/mongo"
)

func SaveNewWorkout(database *mongo.Database, workout shared.Workout) error {
	collection := database.Collection("workout")
	_, err := collection.InsertOne(context.Background(), workout)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
