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
					modifier := val + (float32(faves[i])-3.0)/9.0
					ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
				} else {
					modifier := 1 + (float32(faves[i])-3.0)/9.0
					ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
				}
			}
		}
	}

	return ret
}
