package alteredfuncs

import (
	"fulli9/shared"
)

func ExerRatings(exers map[string]shared.Exercise, user shared.User) map[string]float32 {
	ret := map[string]float32{}

	for _, exercise := range exers {
		if exercise.Parent == "Pushups" {
			ret[exercise.ID.Hex()] = 0.6 * exercise.StartQuality
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
