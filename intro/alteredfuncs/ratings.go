package alteredfuncs

import (
	"fulli9/shared"
	"slices"
)

func ExerRatings(exers map[string]shared.Exercise, user shared.User) map[string]float32 {
	ret := map[string]float32{}

	for _, exercise := range exers {
		if slices.Contains(exercise.GeneralType, "Core") {
			ret[exercise.ID.Hex()] = 1.5 * exercise.StartQuality
		} else if slices.Contains(exercise.GeneralType, "Push") {
			ret[exercise.ID.Hex()] = 1.125 * exercise.StartQuality
		} else {
			ret[exercise.ID.Hex()] = exercise.StartQuality
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

	return ret
}
