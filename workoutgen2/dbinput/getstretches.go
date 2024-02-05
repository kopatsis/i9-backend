package dbinput

import (
	"context"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStretchesDB(database *mongo.Database) (map[string][]shared.Stretch, error) {
	collection := database.Collection("stretch")

	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var allstretches map[string][]shared.Stretch
	dynamics := []shared.Stretch{}
	statics := []shared.Stretch{}
	for cursor.Next(context.TODO()) {
		var str shared.Stretch
		if err := cursor.Decode(&str); err != nil {
			return nil, err
		}
		if str.Status == "Dynamic" {
			dynamics = append(dynamics, str)
		} else {
			statics = append(statics, str)
		}
	}

	allstretches["Dynamic"] = dynamics
	allstretches["Static"] = statics

	return allstretches, nil
}
