package alteredfuncs

import (
	// "fulli9/workoutgen2/datatypes"
	"fulli9/shared"
	"fulli9/workoutgen2/dbinput"
	"sync"

	"github.com/hashicorp/go-multierror"
	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, boltDB *bbolt.DB, userID string, workoutID string) (shared.User, map[string]shared.Exercise, []shared.Workout, shared.TypeMatrix, shared.Workout, error) {
	var user shared.User
	var exercises map[string]shared.Exercise
	var pastWOs []shared.Workout
	var typeMatrix shared.TypeMatrix
	var workout shared.Workout

	var wg sync.WaitGroup

	errChan := make(chan error, 5)
	var errGroup *multierror.Error

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		user, err = dbinput.GetUserDB(database, userID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		workout, err = GetWorkoutByID(database, workoutID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		exercises, err = dbinput.GetExersDB(database, boltDB)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		pastWOs, err = dbinput.GetPastWOsDB(database, userID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		typeMatrix, err = shared.GetMatrixHelper(database, boltDB)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	hasErr := false
	for err := range errChan {
		if err != nil {
			errGroup = multierror.Append(errGroup, err)
			hasErr = true
		}
	}

	if !hasErr {
		return user, exercises, pastWOs, typeMatrix, workout, nil
	}
	return user, exercises, pastWOs, typeMatrix, workout, errGroup
}
