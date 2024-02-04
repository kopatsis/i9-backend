package adjustments

import (
	"fulli9/shared"
	"math"
	"time"
)

func ExerRatings(exers map[string]shared.Exercise, pastWOs []shared.Workout, user shared.User) map[string]float32 {
	ret := map[string]float32{}
	for _, exercise := range exers {
		ret[exercise.ID.Hex()] = exercise.StartQuality
	}

	for _, workout := range pastWOs {
		adjustment := int(time.Since(workout.Date.Time()).Hours())
		for _, exercise := range workout.Exercises {
			for _, id := range exercise.ExerciseIDs {
				ret[id] = float32(math.Max(0.0, float64(ret[id]-float32(11-adjustment))))
			}
		}
	}

	for id, adjustment := range user.ExerFavoriteRates {
		ret[id] = 1.25 * adjustment
	}

	if user.PlyoTolerance > 3 {
		var plyoPlus float32
		if user.PlyoTolerance == 4 {
			plyoPlus = 1.15
		} else {
			plyoPlus = 1.3
		}
		for id := range ret {
			exer := exers[id]
			if exer.PlyoRating == 3 {
				ret[id] = plyoPlus * 1.1 * ret[id]
			} else if exer.PlyoRating == 4 {
				ret[id] = plyoPlus * 1.2 * ret[id]
			}
		}
	}

	return ret
}
