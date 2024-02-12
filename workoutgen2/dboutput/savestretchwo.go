package dboutput

import (
	"context"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveStretchWorkout(database *mongo.Database, workout shared.StretchWorkout) (primitive.ObjectID, error) {
	collection := database.Collection("stretchworkout")
	result, err := collection.InsertOne(context.Background(), workout)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
