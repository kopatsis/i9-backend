package dboutput

import (
	"fulli9/shared"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func SaveDBAllAsync(user shared.User, newlevel float32, userExMod, userTypeMod map[string]float32, roundEndur, timeEndur map[int]float32, ratings [9]float32, workout shared.Workout, exerFacts map[string][3]float32, database *mongo.Database) (error, error, error) {

	var userErr error
	var exerErr error
	var workoutErr error

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		userErr = SaveUser(user, newlevel, userExMod, userTypeMod, roundEndur, timeEndur, database)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		exerErr = SaveModifiedExercises(exerFacts, database)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		workoutErr = SaveWorkout(ratings, workout, database)
	}()

	// if userErr != nil {
	// 	fmt.Printf("Error saving user modifications, try again %s\n", userErr.Error())
	// 	return userErr
	// }

	// if exerErr != nil {
	// 	fmt.Printf("Error saving exercise modifications, try again %s\n", exerErr.Error())
	// 	return exerErr
	// }

	// if workoutErr != nil {
	// 	fmt.Printf("Error saving workout modifications, try again %s\n", workoutErr.Error())
	// 	return workoutErr
	// }

	return userErr, exerErr, workoutErr
}
