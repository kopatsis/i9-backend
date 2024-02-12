package dbinput

import (
	"context"
	"fulli9/shared"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPastWOsDB(database *mongo.Database, userID string) ([]shared.Workout, error) {

	collection := database.Collection("workout")

	tenDaysAgo := primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -10))

	filterWO := bson.D{
		{Key: "date", Value: bson.D{{Key: "$gt", Value: tenDaysAgo}}},
		{Key: "userid", Value: userID},
	}

	optionsWO := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).SetLimit(7)

	cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var pastWorkouts []shared.Workout
	err = cursor.All(context.Background(), &pastWorkouts)
	if err != nil {
		return nil, err
	}

	return pastWorkouts, nil

}
