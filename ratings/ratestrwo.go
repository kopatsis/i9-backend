package ratings

import (
	"fmt"
	"fulli9/ratings/dbinput"
	"fulli9/ratings/operations"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func RateStrWorkout(userID string, favorites []int, fullFave int, onlyWO bool, workoutID string, database *mongo.Database, boltDB *bbolt.DB) error {
	user, workout, strs, err := dbinput.AllStrInputsAsync(database, boltDB, userID, workoutID)
	if err != nil {
		return err
	}

	rating := operations.CreateStrDatabaseRating(favorites, fullFave, onlyWO, workout)

	fmt.Println(user, strs, rating)

	return nil

}
