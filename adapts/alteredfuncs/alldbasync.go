package alteredfuncs

import (
	// "fulli9/workoutgen2/datatypes"
	"fulli9/shared"
	"fulli9/workoutgen2/dbinput"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, userID string, workoutID string) (shared.User, map[string]shared.Exercise, []shared.Workout, shared.TypeMatrix, shared.Workout) {
	var user shared.User
	var exercises map[string]shared.Exercise
	var pastWOs []shared.Workout
	var typeMatrix shared.TypeMatrix
	var workout shared.Workout

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		user = dbinput.GetUserDB(database, userID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		workout = GetWorkoutByID(database, workoutID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		exercises = dbinput.GetExersDB(database)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pastWOs = dbinput.GetPastWOsDB(database, userID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		typeMatrix = dbinput.GetMatrix(database)
	}()
	wg.Wait()

	return user, exercises, pastWOs, typeMatrix, workout
}
