package inputs

import (
	"context"
	"fmt"
	"fulli9/workoutgen/datatypes"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDBInputs(database *mongo.Database, username string) (datatypes.User, []datatypes.Workout) {
	collection := database.Collection("UserWO")
	filter := bson.D{{Key: "username", Value: username}}

	// Define options for the query.
	optionsUser := options.FindOne()

	// Perform the query.
	var result datatypes.User
	err := collection.FindOne(context.Background(), filter, optionsUser).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return result, nil
	}

	collection = database.Collection("workouts")

	tenDaysAgo := primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -10))

	filterWO := bson.D{
		{Key: "created", Value: bson.D{{Key: "$gt", Value: tenDaysAgo}}},
		{Key: "username", Value: result.ID.Hex()},
	}

	optionsWO := options.Find().SetSort(bson.D{{Key: "created", Value: -1}}).SetLimit(7)

	cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
	if err != nil {
		fmt.Println(err)
		return result, nil
	}
	defer cursor.Close(context.Background())

	var pastWorkouts []datatypes.Workout
	err = cursor.All(context.Background(), &pastWorkouts)
	if err != nil {
		fmt.Println(err)
		return result, nil
	}

	return result, pastWorkouts

}
