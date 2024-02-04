package dbinput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetExersDB(database *mongo.Database) map[string]shared.Exercise {
	collection := database.Collection("exercise")

	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(context.Background())

	var allexer map[string]shared.Exercise

	for cursor.Next(context.TODO()) {
		var exer shared.Exercise
		if err := cursor.Decode(&exer); err != nil {
			fmt.Println(err)
			return nil
		}
		allexer[exer.ID.Hex()] = exer
	}

	return allexer
}
