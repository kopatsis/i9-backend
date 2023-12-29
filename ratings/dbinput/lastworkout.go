package dbinput

import (
	"context"
	"fmt"
	"fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPastWOsDB(database *mongo.Database, username string) datatypes.Workout {

	collection := database.Collection("workouts")

	filterWO := bson.D{
		{Key: "username", Value: username},
	}

	optionsWO := options.FindOne().SetSort(bson.D{{Key: "date", Value: -1}})

	var workout datatypes.Workout
	if err := collection.FindOne(context.Background(), filterWO, optionsWO).Decode(&workout); err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			fmt.Println("No workout for user in database")
		}
		return datatypes.Workout{}
	}

	return workout

}
