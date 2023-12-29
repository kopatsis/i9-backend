package ratings

import (
	"errors"
	"fmt"
	"fulli9/ratings/dbinput"
	"fulli9/ratings/dboutput"
	"fulli9/ratings/operations"
	"fulli9/ratings/userinput"
	"fulli9/workoutgen2/dbhandler"
)

func RateWorkout(username string) error {

	client, database, err := dbhandler.ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
		return err
	}
	defer dbhandler.DisConnectDB(client)

	// user := dbinput.GetUserDB(database, username)
	// workout := dbinput.GetPastWOsDB(database, username)
	// countWO := dbinput.GetUserWorkoutCount(database, username)
	// countUser := dbinput.GetUserCount(database)
	// exercises := dbinput.GetExersDB(database)

	user, workout, countWO, countUser, exercises := dbinput.AllInputsAsync(database, username)

	if countWO == 0 {
		fmt.Println("No workouts for user")
		return errors.New("no workouts")
	}

	ratings := userinput.GetUserRatings(workout, exercises)

	newlevel := operations.NewLevel(user, ratings, workout.Difficulty, countWO)

	userExMod, userTypeMod := operations.NewUserMods(user, ratings, workout, exercises, countWO)

	exerFacts := operations.NewExerciseFactorialVars(ratings, workout, exercises, countUser)

	err = dboutput.SaveUser(user, newlevel, userExMod, userTypeMod, database)
	if err != nil {
		fmt.Printf("Error saving user modifications, try again %s\n", err.Error())
		return err
	}

	err = dboutput.SaveModifiedExercises(exerFacts, database)
	if err != nil {
		fmt.Printf("Error saving exercise modifications, try again %s\n", err.Error())
		return err
	}

	err = dboutput.SaveWorkout(ratings, workout, database)
	if err != nil {
		fmt.Printf("Error saving workout modifications, try again %s\n", err.Error())
		return err
	}

	return nil
}
