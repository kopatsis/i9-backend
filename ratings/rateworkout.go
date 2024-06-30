package ratings

import (
	"errors"
	"fmt"
	"fulli9/ratings/dbinput"
	"fulli9/ratings/dboutput"
	"fulli9/ratings/operations"

	"go.mongodb.org/mongo-driver/mongo"
)

func RateWorkout(userID string, ratings [9]float32, favorites [9]float32, fullRating float32, fullFave float32, onlyWorkout bool, workoutID string, database *mongo.Database) error {

	user, workout, countWO, countUser, exercises, err := dbinput.AllInputsAsync(database, userID, workoutID)
	if err != nil {
		return err
	}

	if countWO == 0 {
		fmt.Println("No workouts for user")
		return errors.New("no workouts")
	}

	newlevel, average := operations.NewLevel(user, ratings, workout.Difficulty, countWO)

	ratings = operations.AdjustRatings(ratings, workout)

	userExMod, userTypeMod, roundEndur, timeEndur := operations.NewUserMods(user, ratings, workout, exercises, countWO, average)
	userFavMod := operations.UserFaves(user, favorites, workout)

	exerFacts := operations.NewExerciseFactorialVars(ratings, workout, exercises, countUser)

	userErr, exerErr, workoutErr := dboutput.SaveDBAllAsync(user, newlevel, userExMod, userTypeMod, roundEndur, timeEndur, ratings, workout, exerFacts, userFavMod, database)

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
