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

	notAllowed := []string{"Half Squats", "Double Pulsing Squats", "Triple Pulsing Squats", "Double Pulsing Lunges", "Triple Pulsing Lunges", "High Jumps", "Step High Jumps"}

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
			} else if i > 3 && exercise.PushupType == "Knee" {
				continue
			} else if exercise.PushupType == "Wall" {
				continue
			} else if slices.Contains(notAllowed, exercise.Name) {
				continue
			} else if i < 4 && exercise.PushupType == "Regular" {
				continue
			}

			if exercise.UnderCombos && !(i > 5 && exercise.CardioRating < 2) {
				currentCombo = append(currentCombo, exercise.ID.Hex())
			}
			if exercise.MaxLevel >= level && exercise.CardioRating > 3 && !(i > 5 && exercise.CardioRating < 4) {
				currentNormal = append(currentNormal, exercise.ID.Hex())
			}
			if exercise.InSplits && !(i > 5 && exercise.CardioRating < 2) {
				currentSplit = append(currentSplit, exercise.ID.Hex())
			}
		}
		allowedNormal[i] = currentNormal
		allowedCombo[i] = currentCombo
		allowedSplit[i] = currentSplit
	}

	return allowedNormal, allowedCombo, allowedSplit
}
