package dbinput

import (
	"fulli9/workoutgen2/datatypes"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, username string) (datatypes.User, map[string][]datatypes.Stretch, map[string]datatypes.Exercise, []datatypes.Workout, datatypes.TypeMatrix) {
	var user datatypes.User
	var stretches map[string][]datatypes.Stretch
	var exercises map[string]datatypes.Exercise
	var pastWOs []datatypes.Workout
	var typeMatrix datatypes.TypeMatrix

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		user = GetUserDB(database, username)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		stretches = GetStetchesDB(database)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		exercises = GetExersDB(database)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pastWOs = GetPastWOsDB(database, username)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		typeMatrix = GetMatrix(database)
	}()
	wg.Wait()

	return user, stretches, exercises, pastWOs, typeMatrix
}
