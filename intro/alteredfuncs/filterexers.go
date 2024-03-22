package alteredfuncs

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

func FilterExers(allExercises map[string]shared.Exercise, user shared.User, levelSteps []float32) ([9][]string, [9][]string, [9][]string) {

	allowedNormal := [9][]string{}
	allowedCombo := [9][]string{}
	allowedSplit := [9][]string{}

	for i, level := range levelSteps {
		currentNormal := []string{}
		currentCombo := []string{}
		currentSplit := []string{}
		for _, exercise := range allExercises {
			if user.PlyoTolerance == 0 && exercise.PlyoRating > 0 {
				continue
			} else if user.PlyoTolerance == 1 && exercise.PlyoRating > 1 {
				continue
			} else if user.PlyoTolerance == 2 && exercise.PlyoRating > 3 {
				continue
			} else if exercise.MinLevel > level {
				continue
			} else if slices.Contains(user.BannedExercises, exercise.ID.Hex()) {
				continue
			} else if intersects(user.BannedParts, exercise.BodyParts) {
				continue
			}
			if exercise.UnderCombos && exercise.PushupType != "Wall" {
				currentCombo = append(currentCombo, exercise.ID.Hex())
			}
			if exercise.MaxLevel >= level {
				currentNormal = append(currentNormal, exercise.ID.Hex())
			}
			if exercise.InSplits {
				currentSplit = append(currentSplit, exercise.ID.Hex())
			}
		}
		allowedNormal[i] = currentNormal
		allowedCombo[i] = currentCombo
		allowedSplit[i] = currentSplit
	}

	return allowedNormal, allowedCombo, allowedSplit
}
