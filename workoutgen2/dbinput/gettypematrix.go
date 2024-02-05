package dbinput

import (
	"context"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMatrix(database *mongo.Database) (shared.TypeMatrix, error) {
	collection := database.Collection("typematrix")

	filter := bson.D{}

	var matrix shared.TypeMatrix
	err := collection.FindOne(context.Background(), filter).Decode(&matrix)
	if err != nil {
		return shared.TypeMatrix{}, err
	}

	return matrix, nil
}
