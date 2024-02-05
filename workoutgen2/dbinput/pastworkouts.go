package dbinput

import (
	"context"
	"fmt"
	"fulli9/shared"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPastWOsDB(database *mongo.Database, userID string) []shared.Workout {

	collection := database.Collection("workouts")

	tenDaysAgo := primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -10))

	filterWO := bson.D{
		{Key: "date", Value: bson.D{{Key: "$gt", Value: tenDaysAgo}}},
		{Key: "userid", Value: userID},
	}

	optionsWO := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).SetLimit(7)

	cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(context.Background())

	var pastWorkouts []shared.Workout
	err = cursor.All(context.Background(), &pastWorkouts)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return pastWorkouts

}
