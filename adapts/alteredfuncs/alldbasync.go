package alteredfuncs

import (
	"fulli9/workoutgen2/datatypes"
	"fulli9/workoutgen2/dbinput"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, username string, workoutID string) (datatypes.User, map[string]datatypes.Exercise, []datatypes.Workout, datatypes.TypeMatrix, datatypes.Workout) {
	var user datatypes.User
	var exercises map[string]datatypes.Exercise
	var pastWOs []datatypes.Workout
	var typeMatrix datatypes.TypeMatrix
	var workout datatypes.Workout

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		user = dbinput.GetUserDB(database, username)
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
		pastWOs = dbinput.GetPastWOsDB(database, username)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		typeMatrix = dbinput.GetMatrix(database)
	}()
	wg.Wait()

	return user, exercises, pastWOs, typeMatrix, workout
}
