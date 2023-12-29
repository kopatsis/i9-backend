package dbinput

import (
	"context"
	"fmt"
	"fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMatrix(database *mongo.Database) datatypes.TypeMatrix {
	collection := database.Collection("typematrix")

	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
		return datatypes.TypeMatrix{}
	}
	defer cursor.Close(context.Background())

	var matrix datatypes.TypeMatrix
	if err = cursor.All(context.Background(), &matrix); err != nil {
		fmt.Println(err)
		return datatypes.TypeMatrix{}
	}

	return matrix
}
