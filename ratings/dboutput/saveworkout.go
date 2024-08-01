package dboutput

import (
	"context"
	"fmt"
	"fulli9/shared"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveWorkout(ratings, faves [9]int, fullRating, fullFave int, onlyWorkout bool, workout shared.Workout, database *mongo.Database) error {

	if workout.AvgRating == -1 {
		workout.AvgRating = float32(fullRating)
	} else {
		workout.AvgRating = (((workout.AvgRating) * float32(workout.RatedCount)) + float32(fullRating)) / float32(workout.RatedCount+1)
	}

	if workout.AvgFaves == -1 {
		workout.AvgFaves = float32(fullFave)
	} else {
		workout.AvgFaves = (((workout.AvgFaves) * float32(workout.RatedCount)) + float32(fullFave)) / float32(workout.RatedCount+1)
	}

	workout.RatedCount++
	workout.Status = "Rated"

	now := primitive.NewDateTimeFromTime(time.Now())
	workout.RatedDates = append(workout.RatedDates, now)
	workout.DateToRatings[now] = fullRating
	workout.DateToFaves[now] = fullFave

	if !onlyWorkout {
		for i, round := range workout.Exercises {
			if round.AvgRating == -1 {
				round.AvgRating = float32(ratings[i])
			} else {
				round.AvgRating = (((workout.AvgRating) * float32(workout.RatedCount)) + float32(ratings[i])) / float32(workout.RatedCount+1)
			}

			if round.AvgFaves == -1 {
				round.AvgFaves = float32(faves[i])
			} else {
				round.AvgFaves = (((round.AvgFaves) * float32(workout.RatedCount)) + float32(faves[i])) / float32(workout.RatedCount+1)
			}

			workout.Exercises[i] = round
		}
	}

	collection := database.Collection("workout")

	filter := bson.M{"_id": workout.ID}

	update := bson.M{"$set": workout}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
