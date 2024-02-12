package dboutput

import (
	"context"
	"fulli9/shared"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveNewWorkout(database *mongo.Database, workout shared.Workout) (primitive.ObjectID, error) {
	collection := database.Collection("workout")
	result, err := collection.InsertOne(context.Background(), workout)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
