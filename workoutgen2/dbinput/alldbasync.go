package dbinput

import (
	"fulli9/shared"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, userID string) (shared.User, map[string][]shared.Stretch, map[string]shared.Exercise, []shared.Workout, shared.TypeMatrix) {
	var user shared.User
	var stretches map[string][]shared.Stretch
	var exercises map[string]shared.Exercise
	var pastWOs []shared.Workout
	var typeMatrix shared.TypeMatrix

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		user = GetUserDB(database, userID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		stretches = GetStretchesDB(database)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		exercises = GetExersDB(database)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pastWOs = GetPastWOsDB(database, userID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		typeMatrix = GetMatrix(database)
	}()
	wg.Wait()

	return user, stretches, exercises, pastWOs, typeMatrix
}
