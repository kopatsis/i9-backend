package ratings

import (
	"fmt"
	"fulli9/ratings/dbinput"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func RateStrWorkout(userID string, favorites []int, fullFave int, onlyWO bool, workoutID string, database *mongo.Database, boltDB *bbolt.DB) error {
	user, workout, strs, err := dbinput.AllStrInputsAsync(database, boltDB, userID, workoutID)
	if err != nil {
		return err
	}

	fmt.Println(user, workout, strs)

	return nil

}
