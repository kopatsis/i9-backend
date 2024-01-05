package dboutput

import (
	"context"
	"fmt"
	"fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveWorkout(ratings [9]float32, workout datatypes.Workout, database *mongo.Database) error {

	newExercises := [9]datatypes.WorkoutRound{}
	for i, round := range workout.Exercises {
		newRound := round
		newRound.Rating = ratings[i]
		newExercises[i] = newRound
	}
	newWorkout := workout
	newWorkout.Exercises = newExercises
	newWorkout.Status = "Rated"

	collection := database.Collection("workout")

	filter := bson.M{"_id": workout.ID}

	update := bson.M{"$set": newWorkout}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
