package dboutput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/mongo"
)

func SaveStretchWorkout(database *mongo.Database, workout shared.StretchWorkout) error {
	collection := database.Collection("stretchworkout")
	_, err := collection.InsertOne(context.Background(), workout)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
