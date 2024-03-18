package operations

import (
	"fulli9/shared"
	"math"
)

func UserFaves(user shared.User, faves [9]float32, workout shared.Workout) map[string]float32 {
	ret := user.ExerFavoriteRates

	for i, round := range workout.Exercises {
		for _, id := range round.ExerciseIDs {
			if val, ok := ret[id]; ok {
				modifier := val + (faves[i]-3.0)/9.0
				ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
			} else {
				modifier := 1 + (faves[i]-3.0)/9.0
				ret[id] = float32(math.Max(0.333, math.Min(2.25, float64(modifier))))
			}
		}
	}

	return ret
}
