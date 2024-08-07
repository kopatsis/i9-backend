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

func RateIntroWorkout(userID string, roundEnd float32, database *mongo.Database) error {

	user, err := dbinput.GetUserDB(database, userID)
	if err != nil {
		return err
	}

	levelSteps := []float32{50, 100, 200, 350, 550, 800, 1100, 1500, 2000}

	var userlevel float32
	if roundEnd < 9 {
		completed := int(roundEnd)
		userlevel = 3 * (levelSteps[completed] + (roundEnd-float32(completed))*levelSteps[completed+1]) / 4
	} else {
		userlevel = levelSteps[9]
	}

	pushupSetting := "Wall"
	if roundEnd > 3.99 {
		pushupSetting = "Regular"
	} else if roundEnd > 2.49 {
		pushupSetting = "Knee"
	}

	collection := database.Collection("user")

	filter := bson.M{"_id": user.ID}

	dispLevel := user.DisplayLevel
	if dispLevel == 0 {
		dispLevel = int(userlevel)
	} else {
		dispLevel += 5
	}

	update := bson.M{"$set": bson.M{"level": userlevel, "pushsetting": pushupSetting, "assessed": true, "displevel": dispLevel}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
