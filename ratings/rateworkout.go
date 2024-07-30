package ratings

import (
	"fulli9/ratings/dbinput"
	"fulli9/ratings/dboutput"
	"fulli9/ratings/operations"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func RateWorkout(userID string, ratings, favorites [9]int, fullRating, fullFave int, onlyWorkout bool, completedRounds int, workoutID string, database *mongo.Database, boltDB *bbolt.DB) error {

	user, workout, countWO, exercises, err := dbinput.AllInputsAsync(database, boltDB, userID, workoutID)
	if err != nil {
		return err
	}

	// if countWO == 0 {
	// 	fmt.Println("No workouts for user")
	// 	return errors.New("no workouts")
	// }

	oldLevel := user.Level
	user.Level = operations.NewLevel(user, fullRating, workout.Difficulty, completedRounds, countWO)

	user.WORatedCt++
	user.DisplayLevel += int(4 + user.Level - oldLevel)

	user.ExerModifications, user.TypeModifications, user.RoundEndurance, user.TimeEndurance = operations.NewUserMods(user, ratings, workout, exercises, countWO, fullRating, onlyWorkout)
	user.ExerFavoriteRates = operations.UserFaves(user, favorites, workout, onlyWorkout)

	databaseRating := operations.CreateDatabaseRating(ratings, favorites, fullRating, fullFave, onlyWorkout, workout)

	if err := dboutput.SaveDBAllAsync(user, ratings, favorites, fullRating, fullFave, onlyWorkout, workout, databaseRating, database); err != nil {
		return err
	}

	return nil
}
