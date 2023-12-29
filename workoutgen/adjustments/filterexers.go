package adjustments

import (
	"fulli9/workoutgen/datatypes"
	"slices"
)

func intersects(slice1 []int, slice2 []int) bool {
	for _, item1 := range slice1 {
		for _, item2 := range slice2 {
			if item1 == item2 {
				return true
			}
		}
	}
	return false
}

func FilterExers(allExercises []datatypes.Exercise, user datatypes.User, adjlevel float64) map[string]datatypes.Exercise {
	retExers := map[string]datatypes.Exercise{}
	for _, exercise := range allExercises {
		if exercise.PlyoRating > float32(user.PlyoTolerance) {
			continue
		} else if exercise.MaxLevel < float32(adjlevel) {
			continue
		} else if exercise.MinLevel > float32(adjlevel) {
			continue
		} else if slices.Contains(user.BannedExercises, exercise.ID.Hex()) {
			continue
		} else if intersects(user.BannedParts, exercise.BodyParts) {
			continue
		} else {
			retExers[exercise.ID.Hex()] = exercise
		}
	}
	return retExers
}
