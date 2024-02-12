package ratings

import (
	"errors"
	"fmt"
	"fulli9/ratings/dbinput"
	"fulli9/ratings/dboutput"
	"fulli9/ratings/operations"

	// "fulli9/ratings/userinput"
	// "fulli9/workoutgen2/dbhandler"
	"go.mongodb.org/mongo-driver/mongo"
)

func RateWorkout(userID string, ratings [9]float32, workoutID string, database *mongo.Database) error {

	// client, database, err := dbhandler.ConnectDB()
	// if err != nil {
	// 	fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
	// 	return err
	// }
	// defer dbhandler.DisConnectDB(client)

	// user := dbinput.GetUserDB(database, username)
	// workout := dbinput.GetPastWOsDB(database, username)
	// countWO := dbinput.GetUserWorkoutCount(database, username)
	// countUser := dbinput.GetUserCount(database)
	// exercises := dbinput.GetExersDB(database)

	user, workout, countWO, countUser, exercises, err := dbinput.AllInputsAsync(database, userID, workoutID)
	if err != nil {
		return err
	}

	if countWO == 0 {
		fmt.Println("No workouts for user")
		return errors.New("no workouts")
	}

	// ratings := userinput.GetUserRatings(workout, exercises)

	newlevel, average := operations.NewLevel(user, ratings, workout.Difficulty, countWO)

	userExMod, userTypeMod, roundEndur, timeEndur := operations.NewUserMods(user, ratings, workout, exercises, countWO, average)

	exerFacts := operations.NewExerciseFactorialVars(ratings, workout, exercises, countUser)

	// err := dboutput.SaveUser(user, newlevel, userExMod, userTypeMod, database)
	// if err != nil {
	// 	fmt.Printf("Error saving user modifications, try again %s\n", err.Error())
	// 	return err
	// }

	// err = dboutput.SaveModifiedExercises(exerFacts, database)
	// if err != nil {
	// 	fmt.Printf("Error saving exercise modifications, try again %s\n", err.Error())
	// 	return err
	// }

	// err = dboutput.SaveWorkout(ratings, workout, database)
	// if err != nil {
	// 	fmt.Printf("Error saving workout modifications, try again %s\n", err.Error())
	// 	return err
	// }

	userErr, exerErr, workoutErr := dboutput.SaveDBAllAsync(user, newlevel, userExMod, userTypeMod, roundEndur, timeEndur, ratings, workout, exerFacts, database)

	errortext := ""
	if userErr != nil {
		errortext += "Error with saving user: " + userErr.Error() + "\n"
	}
	if exerErr != nil {
		errortext += "Error with saving exers: " + exerErr.Error() + "\n"
	}
	if workoutErr != nil {
		errortext += "Error with saving workout: " + workoutErr.Error() + "\n"
	}
	if errortext != "" {
		return errors.New(errortext)
	}

	return nil
}
