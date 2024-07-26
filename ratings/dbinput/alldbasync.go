package dbinput

import (
	"fulli9/shared"
	"sync"

	"github.com/hashicorp/go-multierror"
	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, boltDB *bbolt.DB, userID string, id string) (shared.User, shared.Workout, int64, map[string]shared.Exercise, error) {
	var user shared.User
	var workout shared.Workout
	var countWO int64
	var exercises map[string]shared.Exercise

	var wg sync.WaitGroup

	errChan := make(chan error, 3)
	var errGroup *multierror.Error

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		user, err = GetUserDB(database, userID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		workout, err = GetPastWOsDB(database, id)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		exercises, err = GetExersDB(database, boltDB)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	hasErr := false
	for err := range errChan {
		if err != nil {
			hasErr = true
			errGroup = multierror.Append(errGroup, err)
		}
	}

	if !hasErr {
		return user, workout, countWO, exercises, nil
	}

	countWO = int64(user.WORatedCt) + 2

	return user, workout, countWO, exercises, errGroup
}

func AllStrInputsAsync(database *mongo.Database, boltDB *bbolt.DB, userID string, id string) (shared.User, shared.StretchWorkout, map[string]shared.Stretch, error) {
	var user shared.User
	var workout shared.StretchWorkout
	var strs map[string]shared.Stretch

	var wg sync.WaitGroup

	errChan := make(chan error, 3)
	var errGroup *multierror.Error

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		user, err = GetUserDB(database, userID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		workout, err = GetPastStrWO(database, id)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		strs, err = GetStrsDB(database, boltDB)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	hasErr := false
	for err := range errChan {
		if err != nil {
			hasErr = true
			errGroup = multierror.Append(errGroup, err)
		}
	}

	if !hasErr {
		return user, workout, strs, nil
	}

	return user, workout, strs, errGroup
}
