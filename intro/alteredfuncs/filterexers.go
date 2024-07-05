package alteredfuncs

import (
	"fulli9/shared"
)

func FilterExers(allExercises map[string]shared.Exercise, user shared.User, levelSteps []float32) ([9][]string, [9][]string, [9][]string) {

	allowedNormal := [9][]string{}
	allowedCombo := [9][]string{}
	allowedSplit := [9][]string{}

	for i := range levelSteps {
		currentNormal := []string{}
		currentCombo := []string{}
		currentSplit := []string{}
		for _, exercise := range allExercises {

			if i == 0 {
				if exercise.IntroGroup == 1 || exercise.IntroGroup == 10 {
					currentNormal = append(currentNormal, exercise.ID.Hex())
				}
			} else if i == 7 {
				if exercise.IntroGroup == 8 {
					currentCombo = append(currentCombo, exercise.ID.Hex())
				}
			} else if i == 8 {
				if exercise.IntroGroup == 9 || exercise.IntroGroup == 10 {
					currentSplit = append(currentSplit, exercise.ID.Hex())
				}
			} else {
				if exercise.IntroGroup == i+1 {
					currentNormal = append(currentNormal, exercise.ID.Hex())
				}
			}

		}
		allowedNormal[i] = currentNormal
		allowedCombo[i] = currentCombo
		allowedSplit[i] = currentSplit
	}

	return allowedNormal, allowedCombo, allowedSplit
}
