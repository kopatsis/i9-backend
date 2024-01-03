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

	var matrix datatypes.TypeMatrix
	err := collection.FindOne(context.Background(), filter).Decode(&matrix)
	if err != nil {
		fmt.Println(err)
		return datatypes.TypeMatrix{}
	}

	return matrix
}
