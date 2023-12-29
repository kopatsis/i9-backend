package stretches

import (
	"context"
	"fmt"
	"fulli9/workoutgen/datatypes"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StretchesFromDB(database *mongo.Database, level float64) ([]datatypes.Stretch, []datatypes.Stretch) {
	collection := database.Collection("stretch")

	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	defer cursor.Close(context.Background())

	var allstretch []datatypes.Stretch
	if err := cursor.All(context.Background(), &allstretch); err != nil {
		fmt.Println(err)
		return nil, nil
	}

	var allStatic []datatypes.Stretch
	var allDynamic []datatypes.Stretch

	for _, str := range allstretch {
		if float64(str.MinLevel) < level {
			continue
		}
		if str.Status == "Dynamic" {
			allDynamic = append(allDynamic, str)
		} else {
			allStatic = append(allStatic, str)
		}
	}

	return allStatic, allDynamic
}

func GetStretches(database *mongo.Database, allExercises map[string]datatypes.Exercise, timesMap map[string]float32, exersToDo [9]map[string]string, level float64) map[string][]string {
	ret := map[string][]string{}

	allStatic, allDynamic := StretchesFromDB(database, level)

	ret["static"] = []string{}
	for i := 0; i < int(timesMap["staticSets"]); i++ {
		ret["static"] = append(ret["static"], allStatic[int(rand.Float64()*float64(len(allStatic)))].ID.Hex())
	}

	ret["dynamic"] = []string{}
	for i := 0; i < int(timesMap["dynamicSets"]); i++ {
		ret["dynamic"] = append(ret["dynamic"], allDynamic[int(rand.Float64()*float64(len(allDynamic)))].ID.Hex())
	}

	return ret
}
