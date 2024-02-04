package dbinput

import (
	"fulli9/workoutgen2/datatypes"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, username string, id string) (datatypes.User, datatypes.Workout, int64, int64, map[string]datatypes.Exercise) {
	var user datatypes.User
	var workout datatypes.Workout
	var countWO int64
	var countUser int64
	var exercises map[string]datatypes.Exercise

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		user = GetUserDB(database, username)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		workout = GetPastWOsDB(database, id)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		countWO = GetUserWorkoutCount(database, username)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		countUser = GetUserCount(database)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		exercises = GetExersDB(database)
	}()
	wg.Wait()

	return user, workout, countWO, countUser, exercises
}
