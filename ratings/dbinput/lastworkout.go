package dbinput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func GetPastWOsDB(database *mongo.Database, idStr string) shared.Workout {

	collection := database.Collection("workouts")

	// filterWO := bson.D{
	// 	{Key: "username", Value: username},
	// }

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(idStr); err == nil {
		id = oid
	} else {
		// Handle error here, sorry
		return shared.Workout{}
	}

	// optionsWO := options.FindOne().SetSort(bson.D{{Key: "date", Value: -1}})

	var workout shared.Workout
	if err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&workout); err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			fmt.Println("No workout for user in database")
		}
		return shared.Workout{}
	}

	return workout

}
