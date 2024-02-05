package alteredfuncs

import (
	// "fulli9/workoutgen2/datatypes"
	"fulli9/shared"
	"fulli9/workoutgen2/dbinput"
	"sync"

	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, userID string, workoutID string) (shared.User, map[string]shared.Exercise, []shared.Workout, shared.TypeMatrix, shared.Workout, error) {
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
		exercises, err = dbinput.GetExersDB(database)
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
		typeMatrix, err = dbinput.GetMatrix(database)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		errGroup = multierror.Append(errGroup, err)
	}

	return user, exercises, pastWOs, typeMatrix, workout, errGroup
}
