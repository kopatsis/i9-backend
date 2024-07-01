package ratings

import (
	// "errors"
	// "fmt"
	// "fulli9/ratings/dbinput"
	// "fulli9/ratings/dboutput"
	// "fulli9/ratings/operations"

	"go.mongodb.org/mongo-driver/mongo"
)

func RateWorkout(userID string, ratings, favorites [9]int, fullRating, fullFave int, onlyWorkout bool, completedRounds int, workoutID string, database *mongo.Database) error {

	// user, workout, countWO, exercises, err := dbinput.AllInputsAsync(database, userID, workoutID)
	// if err != nil {
	// 	return err
	// }

	// if countWO == 0 {
	// 	fmt.Println("No workouts for user")
	// 	return errors.New("no workouts")
	// }

	// newlevel := operations.NewLevel(user, fullRating, workout.Difficulty, completedRounds, countWO)

	// userExMod, userTypeMod, roundEndur, timeEndur := operations.NewUserMods(user, ratings, workout, exercises, countWO, average)
	// userFavMod := operations.UserFaves(user, favorites, workout)

	// databaseRating := operations.CreateDatabaseRating(ratings, favorites, fullRating, fullFave, onlyWorkout, workout)

	// userErr, exerErr, workoutErr := dboutput.SaveDBAllAsync(user, newlevel, userExMod, userTypeMod, roundEndur, timeEndur, ratings, workout, userFavMod, database)

	// errortext := ""
	// if userErr != nil {
	// 	errortext += "Error with saving user: " + userErr.Error() + "\n"
	// }
	// if exerErr != nil {
	// 	errortext += "Error with saving exers: " + exerErr.Error() + "\n"
	// }
	// if workoutErr != nil {
	// 	errortext += "Error with saving workout: " + workoutErr.Error() + "\n"
	// }
	// if errortext != "" {
	// 	return errors.New(errortext)
	// }

	return nil
}
