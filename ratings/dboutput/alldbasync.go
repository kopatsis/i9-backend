package dboutput

import (
	"fulli9/shared"
	"sync"

	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveDBAllAsync(user shared.User, ratings, faves [9]int, fullRating, fullFave int, onlyWorkout bool, workout shared.Workout, databaseRating shared.StoredRating, database *mongo.Database) error {

	var wg sync.WaitGroup

	errChan := make(chan error, 3)
	var errGroup *multierror.Error

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := SaveUser(user, database)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := SaveWorkout(ratings, faves, fullRating, fullFave, onlyWorkout, workout, database)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := SaveRating(databaseRating, database)
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

	if hasErr {
		return errGroup
	}
	return nil

}
