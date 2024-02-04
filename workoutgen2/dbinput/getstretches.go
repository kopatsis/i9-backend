package dbinput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStretchesDB(database *mongo.Database) map[string][]shared.Stretch {
	collection := database.Collection("stretch")

	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(context.Background())

	var allstretches map[string][]shared.Stretch
	dynamics := []shared.Stretch{}
	statics := []shared.Stretch{}
	for cursor.Next(context.TODO()) {
		var str shared.Stretch
		if err := cursor.Decode(&str); err != nil {
			fmt.Println(err)
			return nil
		}
		if str.Status == "Dynamic" {
			dynamics = append(dynamics, str)
		} else {
			statics = append(statics, str)
		}
	}

	allstretches["Dynamic"] = dynamics
	allstretches["Static"] = statics

	return allstretches
}
