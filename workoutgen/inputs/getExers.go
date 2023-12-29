package inputs

import (
	"context"
	"fmt"
	"fulli9/workoutgen/datatypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllExers(database *mongo.Database) []datatypes.Exercise {
	collection := database.Collection("exercise")

	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(context.Background())

	var allexer []datatypes.Exercise
	if err := cursor.All(context.Background(), &allexer); err != nil {
		fmt.Println(err)
		return nil
	}

	return allexer
}
