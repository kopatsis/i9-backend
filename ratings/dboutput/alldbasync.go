package dboutput

import (
	"fulli9/shared"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func SaveDBAllAsync(user shared.User, newlevel float32, userExMod, userTypeMod map[string]float32, roundEndur, timeEndur map[int]float32, ratings [9]float32, workout shared.Workout, userFavMod map[string]float32, database *mongo.Database) (error, error, error) {

	var userErr error
	var exerErr error
	var workoutErr error

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		userErr = SaveUser(user, newlevel, userExMod, userTypeMod, roundEndur, timeEndur, userFavMod, database)
	}()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	exerErr = SaveAnalysisRatings(ratings, database)
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	workoutErr = SaveWorkout(ratings, workout, database)
	// }()

	wg.Wait()

	return userErr, exerErr, workoutErr
}
