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

func SaveStrWorkout(faves []int, fullFave int, onlyWO bool, workout shared.StretchWorkout, database *mongo.Database) error {

	if workout.AvgFaves == -1 {
		workout.AvgFaves = float32(fullFave)
	} else {
		workout.AvgFaves = (((workout.AvgFaves) * float32(workout.RatedCount)) + float32(fullFave)) / float32(workout.RatedCount+1)
	}

	workout.RatedCount++
	workout.Status = "Rated"

	now := primitive.NewDateTimeFromTime(time.Now())
	workout.RatedDates = append(workout.RatedDates, now)
	workout.DateToFaves[now] = fullFave

	collection := database.Collection("stretchworkout")

	filter := bson.M{"_id": workout.ID}

	update := bson.M{"$set": workout}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
