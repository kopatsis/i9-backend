package dbinput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMatrix(database *mongo.Database) shared.TypeMatrix {
	collection := database.Collection("typematrix")

	filter := bson.D{}

	var matrix shared.TypeMatrix
	err := collection.FindOne(context.Background(), filter).Decode(&matrix)
	if err != nil {
		fmt.Println(err)
		return shared.TypeMatrix{}
	}

	return matrix
}
