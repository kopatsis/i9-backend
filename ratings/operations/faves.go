package operations

import (
	"fulli9/shared"
	"math"
)

func UserFaves(user shared.User, faves [9]int, workout shared.Workout, onlyWorkout bool) map[string]float32 {
	ret := user.ExerFavoriteRates

	if !onlyWorkout {
		for i, round := range workout.Exercises {
			for _, id := range round.ExerciseIDs {
				if val, ok := ret[id]; ok {
					modifier := val + (float32(faves[i])-5.0)/20.0
					ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
				} else {
					modifier := 1 + (float32(faves[i])-5.0)/20.0
					ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
				}
			}
		}
	}

	return ret
}

func UserStrFaves(user shared.User, faves []int, fave int, workout shared.StretchWorkout, onlyWorkout bool) map[string]float32 {
	ret := user.StrFavoriteRates

	for i, id := range workout.Dynamics {
		var currentFave int
		if !onlyWorkout && i < len(faves) {
			currentFave = faves[i]
		} else {
			currentFave = fave
		}

		if val, ok := ret[id]; ok {
			modifier := val + (float32(currentFave)-5.0)/20.0
			ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
		} else {
			modifier := 1 + (float32(currentFave)-5.0)/20.0
			ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
		}
	}

	for i, id := range workout.Statics {
		var currentFave int
		if !onlyWorkout && i < len(faves) {
			currentFave = faves[i]
		} else {
			currentFave = fave
		}

		if val, ok := ret[id]; ok {
			modifier := val + (float32(currentFave)-5.0)/20.0
			ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
		} else {
			modifier := 1 + (float32(currentFave)-5.0)/20.0
			ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
		}
	}

	return ret
}
