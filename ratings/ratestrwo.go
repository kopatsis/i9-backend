package ratings

import (
	"fulli9/ratings/dbinput"
	"fulli9/ratings/dboutput"
	"fulli9/ratings/operations"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func RateStrWorkout(userID string, favorites []int, fullFave int, onlyWO bool, workoutID string, database *mongo.Database, boltDB *bbolt.DB) error {
	user, workout, err := dbinput.AllStrInputsAsync(database, boltDB, userID, workoutID)
	if err != nil {
		return err
	}

	rating := operations.CreateStrDatabaseRating(favorites, fullFave, onlyWO, workout)

	user.StrFavoriteRates = operations.UserStrFaves(user, favorites, fullFave, workout, onlyWO)
	user.StrWORatedCt++

	if err := dboutput.SaveStrDBAllAsync(user, favorites, fullFave, onlyWO, workout, rating, database); err != nil {
		return err
	}

	return nil

}
