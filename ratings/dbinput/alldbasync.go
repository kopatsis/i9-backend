package dbinput

import (
	"fulli9/shared"
	"sync"

	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, userID string, id string) (shared.User, shared.Workout, int64, map[string]shared.Exercise, error) {
	// func AllInputsAsync(database *mongo.Database, userID string, id string) (shared.User, shared.Workout, int64, int64, map[string]shared.Exercise, error) {
	var user shared.User
	var workout shared.Workout
	var countWO int64
	// var countUser int64
	var exercises map[string]shared.Exercise

	var wg sync.WaitGroup

	// errChan := make(chan error, 5)
	errChan := make(chan error, 4)
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
		countWO, err = GetUserWorkoutCount(database, userID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		exercises, err = GetExersDB(database)
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
	return user, workout, countWO, exercises, errGroup
}
