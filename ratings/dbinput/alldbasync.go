package dbinput

import (
	"fulli9/shared"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func AllInputsAsync(database *mongo.Database, userID string, id string) (shared.User, shared.Workout, int64, int64, map[string]shared.Exercise) {
	var user shared.User
	var workout shared.Workout
	var countWO int64
	var countUser int64
	var exercises map[string]shared.Exercise

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		user = GetUserDB(database, userID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		workout = GetPastWOsDB(database, id)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		countWO = GetUserWorkoutCount(database, userID)
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
