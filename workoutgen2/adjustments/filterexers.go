package adjustments

import (
	"fulli9/shared"
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

func FilterExers(allExercises map[string]shared.Exercise, user shared.User, adjlevel float32) ([]string, []string, []string) {

	allowedNormal := []string{}
	allowedCombo := []string{}
	allowedSplit := []string{}

	for _, exercise := range allExercises {
		if user.PlyoTolerance == 0 && exercise.PlyoRating > 0 {
			continue
		} else if user.PlyoTolerance == 1 && exercise.PlyoRating > 1 {
			continue
		} else if user.PlyoTolerance == 2 && exercise.PlyoRating > 3 {
			continue
		} else if exercise.MinLevel > adjlevel {
			continue
		} else if slices.Contains(user.BannedExercises, exercise.ID.Hex()) {
			continue
		} else if intersects(user.BannedParts, exercise.BodyParts) {
			continue
		} else if exercise.PushupType != "" && user.PushupSetting != exercise.PushupType {
			continue
		}
		if exercise.UnderCombos {
			allowedCombo = append(allowedCombo, exercise.ID.Hex())
		}
		if exercise.MaxLevel >= adjlevel {
			allowedNormal = append(allowedNormal, exercise.ID.Hex())
		}
		if exercise.InSplits {
			allowedSplit = append(allowedSplit, exercise.ID.Hex())
		}
	}

	return allowedNormal, allowedCombo, allowedSplit
}
