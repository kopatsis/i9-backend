package adjustments

import (
	"fulli9/shared"
	"math"
	"time"
)

func ExerRatings(diff int, exers map[string]shared.Exercise, pastWOs []shared.Workout, user shared.User) map[string]float32 {
	ret := map[string]float32{}
	for _, exercise := range exers {
		ret[exercise.ID.Hex()] = exercise.StartQuality
	}

	if diff == 1 {
		for id, val := range ret {
			if exers[id].CardioRating > 2.85 {
				ret[id] = (val * 0.667)
			}
		}
	} else if diff > 4 {
		for id, val := range ret {
			if exers[id].CardioRating > 4 {
				ret[id] = (val * 1.45)
			}
		}
	}
	if diff == 6 {
		for id, val := range ret {
			if exers[id].CardioRating > 4.75 {
				ret[id] = (val * 1.45)
			}
		}
	}

	balance := [3]float32{1.0 / 2.0, 1.0 / 4.0, 1.0 / 4.0}

	for i, workout := range pastWOs {
		if workout.Status != "Rated" {
			continue
		}
		adjustment := int(time.Since(workout.Date.Time()).Hours())
		for _, exercise := range workout.Exercises {
			for _, id := range exercise.ExerciseIDs {
				hourAdj := math.Min(float64(ret[id]), float64(ret[id]-float32(18-adjustment)))
				ret[id] = float32(math.Max(0.0, hourAdj))
			}
		}
		weight := 1 / float32(math.Pow(2, float64(i+1)))

		balance[0] = weight*workout.GeneralTypeVals[0] + (1-weight)*balance[0]
		balance[1] = weight*workout.GeneralTypeVals[1] + (1-weight)*balance[1]
		balance[2] = weight*workout.GeneralTypeVals[2] + (1-weight)*balance[2]
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
			if gtype == "Legs" {
				ret[exer.ID.Hex()] *= 1.0 + (2.5 * (1.0/2.0 - balance[genTypeToPos[gtype]]))
			} else {
				ret[exer.ID.Hex()] *= 1.0 + (2.5 * (1.0/4.0 - balance[genTypeToPos[gtype]]))
			}

		}
	}

	return ret
}

func AdjustBurpeeRatings(user shared.User, exers map[string]shared.Exercise, ratings map[string]float32) {
	for id, rating := range ratings {
		exer := exers[id]

		if user.PushupSetting == "Wall" && exer.Parent == "Burpees" && exer.PushupType == "" {
			ratings[id] = rating * 1.667
		} else if user.PushupSetting == "Knee" && exer.Parent == "Burpees" && exer.PushupType == "Knee" {
			ratings[id] = rating * 1.333
		}
	}
}
