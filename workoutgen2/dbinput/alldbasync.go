package dbinput

import (
	"fulli9/shared"
	"sync"

	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, userID string) (shared.User, map[string][]shared.Stretch, map[string]shared.Exercise, []shared.Workout, shared.TypeMatrix, error) {
	var user shared.User
	var stretches map[string][]shared.Stretch
	var exercises map[string]shared.Exercise
	var pastWOs []shared.Workout
	var typeMatrix shared.TypeMatrix

	var wg sync.WaitGroup

	errChan := make(chan error, 5)
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
		stretches, err = GetStretchesDB(database)
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		pastWOs, err = GetPastWOsDB(database, userID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		typeMatrix, err = GetMatrix(database)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		errGroup = multierror.Append(errGroup, err)
	}

	return user, stretches, exercises, pastWOs, typeMatrix, errGroup
}
