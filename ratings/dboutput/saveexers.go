package dboutput

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveModifiedExercises(exerFacts map[string][3]float32, database *mongo.Database) error {

	collection := database.Collection("exercise")

	var bulkOps []mongo.WriteModel

	for id, vars := range exerFacts {
		formattedID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		filter := bson.M{"_id": formattedID}
		update := bson.M{"$set": bson.M{"repvars": vars}}
		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		bulkOps = append(bulkOps, updateModel)
	}

	_, err := collection.BulkWrite(context.Background(), bulkOps)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
