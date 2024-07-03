package dbinput

import (
	"fmt"
	"fulli9/shared"
	"sync"

	"github.com/hashicorp/go-multierror"
	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, boltDB *bbolt.DB, userID string) (shared.User, map[string][]shared.Stretch, map[string]shared.Exercise, []shared.Workout, shared.TypeMatrix, error) {
	var user shared.User
	stretches := map[string][]shared.Stretch{}
	exercises := map[string]shared.Exercise{}
	pastWOs := []shared.Workout{}
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
			fmt.Println(err)
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		stretches, err = GetStretchesDB(database, boltDB)
		if err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		exercises, err = GetExersDB(database)
		if err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		pastWOs, err = GetPastWOsDB(database, userID)
		if err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		typeMatrix, err = GetMatrix(database)
		if err != nil {
			fmt.Println(err)
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

	if hasErr {
		return user, stretches, exercises, pastWOs, typeMatrix, errGroup
	}

	return user, stretches, exercises, pastWOs, typeMatrix, nil
}
