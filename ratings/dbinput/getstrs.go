package dbinput

import (
	"context"
	"fmt"
	"fulli9/shared"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStrsDB(database *mongo.Database, boltDB *bbolt.DB) (map[string]shared.Stretch, error) {

	strs, err := shared.GetStretchHelper(database, boltDB)
	if err != nil {
		return nil, err
	}

	allstrs := map[string]shared.Stretch{}
	for _, str := range strs {
		allstrs[str.ID.Hex()] = str
	}

	return allstrs, nil
}

func GetPastStrWO(database *mongo.Database, idStr string) (shared.StretchWorkout, error) {

	collection := database.Collection("stretchworkout")

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(idStr); err == nil {
		id = oid
	} else {
		return shared.StretchWorkout{}, err
	}

	var workout shared.StretchWorkout
	if err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&workout); err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			fmt.Println("No workout for user in database")
			return shared.StretchWorkout{}, err
		} else {
			return shared.StretchWorkout{}, nil
		}
	}

	return workout, nil

}
