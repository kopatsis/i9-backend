package alteredfuncs

import "fulli9/shared"

func ModifyExisting(workout shared.Workout, reps [9][]float32, pairs [9][]bool, genRatings [3]float32, cardioRatings [9]float32, cardioRating float32) shared.Workout {

	for i := range workout.Exercises {
		workout.Exercises[i].Pairs = pairs[i]
		workout.Exercises[i].Reps = reps[i]
	}

	workout.CardioRating = cardioRating
	workout.CardioRatings = cardioRatings
	workout.GeneralTypeVals = genRatings

	workout.Status = "Unstarted"
	workout.PausedTime = 0

	return workout
}
