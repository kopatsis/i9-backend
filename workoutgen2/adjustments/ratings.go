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

	hoursSince := 100000
	balance := [3]float32{1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0}

	for _, workout := range pastWOs {
		if workout.Status == "Unstarted" {
			continue
		}
		adjustment := int(time.Since(workout.Date.Time()).Hours())
		for _, exercise := range workout.Exercises {
			for _, id := range exercise.ExerciseIDs {
				hourAdj := math.Min(float64(ret[id]), float64(ret[id]-float32(18-adjustment)))
				ret[id] = float32(math.Max(0.0, hourAdj))
			}
		}
		if adjustment < hoursSince {
			balance = workout.GeneralTypeVals
		}
	}

	for id, adjustment := range user.ExerFavoriteRates {
		if _, ok := ret[id]; ok {
			ret[id] = ret[id] * adjustment
		}
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

	genTypeToPos := map[string]int{"Legs": 0, "Core": 1, "Push": 2}
	for _, exer := range exers {
		for _, gtype := range exer.GeneralType {
			ret[exer.ID.Hex()] *= 1.0 + (4.0 * (1.0/3.0 - balance[genTypeToPos[gtype]]))
		}
	}

	return ret
}

func AdjustBurpeeRatings(user shared.User, exers map[string]shared.Exercise, ratings map[string]float32) {
	for id, rating := range ratings {
		exer := exers[id]

		if user.PushupSetting == "Wall" && (exer.Name == "Step Burpees" || exer.Name == "Non-Pushup Burpees") {
			ratings[id] = rating * 2.25
		} else if user.PushupSetting == "Knee" && (exer.Name == "Step Burpees" || exer.Name == "Non-Pushup Burpees") {
			ratings[id] = rating * 1.5
		}
	}
}
