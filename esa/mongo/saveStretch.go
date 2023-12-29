package mongo

import (
	"context"
	"fmt"
	"fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveStretch(collection *mongo.Collection, stretches []datatypes.Stretch) map[string]string {

	ret := map[string]string{}
	for _, stretch := range stretches {
		insertResult, err := collection.InsertOne(context.Background(), stretch)
		if err != nil {
			fmt.Println(err)
		}
		id := insertResult.InsertedID.(primitive.ObjectID).Hex()
		ret[id] = stretch.Name
	}
	return ret

}
