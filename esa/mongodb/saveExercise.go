package mongodb

import (
	"context"
	"fmt"
	"fulli9/workoutgen2/datatypes"

	// "time"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveExercise(collection *mongo.Collection, exercises []datatypes.Exercise) map[string]string {

	ret := map[string]string{}
	for _, exer := range exercises {
		insertResult, err := collection.InsertOne(context.Background(), exer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		id := insertResult.InsertedID.(primitive.ObjectID).Hex()
		ret[id] = exer.Name
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// cursor, err := collection.Find(ctx, bson.M{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }
	// defer cursor.Close(ctx)

	// nameToID := map[string]string{}

	// for cursor.Next(ctx) {
	// 	var exer datatypes.RetExercise
	// 	if err := cursor.Decode(&exer); err != nil {
	// 		fmt.Println(err)
	// 		return nil
	// 	}
	// 	nameToID[exer.Name] = exer.ID.Hex()
	// }

	// if err := cursor.Err(); err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }

	// for name, compatMap := range compats {
	// 	saveID, err := primitive.ObjectIDFromHex(nameToID[name])
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return nil
	// 	}
	// 	saveMap := map[string][2]float32{}

	// 	for compat_name, vals := range compatMap {
	// 		saveMap[nameToID[compat_name]] = vals
	// 	}

	// 	filter := bson.M{"_id": saveID}

	// 	update := bson.M{"$set": bson.M{"compatibles": saveMap}}

	// 	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
	// 		fmt.Println(err)
	// 		return nil
	// 	}
	// }

	return ret

}
