package ratings

import (
	"context"
	"fmt"
	"fulli9/ratings/dbinput"

	// "fulli9/workoutgen2/datatypes"
	// "fulli9/workoutgen2/dbhandler"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RateIntroWorkout(username string, roundEnd float32, database *mongo.Database) error {

	// client, database, err := dbhandler.ConnectDB()
	// if err != nil {
	// 	fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
	// 	return err
	// }
	// defer dbhandler.DisConnectDB(client)

	user := dbinput.GetUserDB(database, username)
	levelSteps := []float32{50, 125, 225, 375, 575, 825, 1125, 1475, 1975}

	var userlevel float32
	if roundEnd < 9 {
		completed := int(roundEnd)
		userlevel = levelSteps[completed] + (roundEnd-float32(completed))*levelSteps[completed+1]
	} else {
		userlevel = levelSteps[8]
	}

	collection := database.Collection("user")

	filter := bson.M{"_id": user.ID}

	update := bson.M{"$set": bson.M{"level": userlevel}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
