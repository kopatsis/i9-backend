package operations

import (
	"fulli9/shared"
	"math"
)

func UserFaves(user shared.User, faves [9]int, fave int, workout shared.Workout, onlyWorkout bool) map[string]float32 {
	ret := user.ExerFavoriteRates

	for i, round := range workout.Exercises {
		for _, id := range round.ExerciseIDs {

			var currentFave int
			var modifier float32

			if !onlyWorkout && i < len(faves) {
				currentFave = faves[i]
				modifier = (float32(currentFave) - 5.0) / 15.0
			} else {
				currentFave = fave
				modifier = (float32(currentFave) - 5.0) / 25.0
			}

			if val, ok := ret[id]; ok {
				ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(val+modifier))))
			} else {
				ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(1+modifier))))
			}
		}
	}

	return ret
}

func UserStrFaves(user shared.User, faves []int, fave int, workout shared.StretchWorkout, onlyWorkout bool) map[string]float32 {
	ret := user.StrFavoriteRates

	for i, id := range workout.Dynamics {
		var currentFave int
		var modifier float32

		if !onlyWorkout && i < len(faves) {
			currentFave = faves[i]
			modifier = (float32(currentFave) - 5.0) / 15.0
		} else {
			currentFave = fave
			modifier = (float32(currentFave) - 5.0) / 25.0
		}

		if val, ok := ret[id]; ok {
			ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(val+modifier))))
		} else {
			ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(1+modifier))))
		}
	}

	for i, id := range workout.Statics {
		var currentFave int
		var modifier float32

		if !onlyWorkout && i < len(faves) {
			currentFave = faves[i]
			modifier = (float32(currentFave) - 5.0) / 15.0
		} else {
			currentFave = fave
			modifier = (float32(currentFave) - 5.0) / 25.0
		}

		if val, ok := ret[id]; ok {
			ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(val+modifier))))
		} else {
			ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(1+modifier))))
		}
	}

	return ret
}
