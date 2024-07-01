package ratings

import (
	"errors"
	"fmt"
	"fulli9/ratings/dbinput"
	"fulli9/ratings/dboutput"
	"fulli9/ratings/operations"

	"go.mongodb.org/mongo-driver/mongo"
)

func RateWorkout(userID string, ratings, favorites [9]int, fullRating, fullFave int, onlyWorkout bool, completedRounds int, workoutID string, database *mongo.Database) error {

	user, workout, countWO, exercises, err := dbinput.AllInputsAsync(database, userID, workoutID)
	if err != nil {
		return err
	}

	if countWO == 0 {
		fmt.Println("No workouts for user")
		return errors.New("no workouts")
	}

	user.Level = operations.NewLevel(user, fullRating, workout.Difficulty, completedRounds, countWO)

	user.ExerModifications, user.TypeModifications, user.RoundEndurance, user.TimeEndurance = operations.NewUserMods(user, ratings, workout, exercises, countWO, fullRating, onlyWorkout)
	user.ExerFavoriteRates = operations.UserFaves(user, favorites, workout, onlyWorkout)

	databaseRating := operations.CreateDatabaseRating(ratings, favorites, fullRating, fullFave, onlyWorkout, workout)

	err = dboutput.SaveDBAllAsync(user, ratings, favorites, fullRating, fullFave, onlyWorkout, workout, databaseRating, database)
	if err != nil {
		return err
	}

	return nil
}
